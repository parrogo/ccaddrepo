// Package ccaddrepo allows to add a GitHub repository to CodeClimate
// and setup its report id as a CC_TEST_REPORTER_ID secret in the repository.
//
// This two function works together in order to automate the setup of
// a GitHub repository.
//
// If you are using this package in GitHub Actions,
// you can easily publish coverages reports to CodeClimate
// using e.g. [paambaati/codeclimate-action](https://github.com/paambaati/codeclimate-action)
package ccaddrepo

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	sodium "github.com/GoKillers/libsodium-go/cryptobox"
	"github.com/google/go-github/v34/github"
	"golang.org/x/oauth2"
)

const ccBaseURL = "https://api.codeclimate.com/v1/"
const ccGetURL = "https://api.codeclimate.com/v1/repos"
const ccURL = "https://api.codeclimate.com/v1/github/repos"
const ccBodyFormat = `{"data":{"type": "repos","attributes": {"url": "https://github.com/%s"}}}`

// CodeClimate represent an authenticated
// session on Code Climate.
type CodeClimate string

func (cc CodeClimate) doRequest(method string, URL string, body io.Reader, response interface{}) error {
	req, err := http.NewRequest(method, ccBaseURL+URL, body)
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/vnd.api+json")
	req.Header.Add("Authorization", "Token token="+string(cc))
	req.Header.Add("Content-Type", "application/vnd.api+json")
	var c http.Client
	res, err := c.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if response == nil {
		return nil
	}

	resbuf, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	//fmt.Println(string(resbuf))
	return json.Unmarshal(resbuf, response)

}

// GetRepoID returns the ID of a repository
func (cc CodeClimate) GetRepoID(reposlug string) (string, error) {
	var response struct {
		Data []struct {
			ID string
		}
	}

	err := cc.doRequest("GET", "repos?github_slug="+reposlug, nil, &response)
	if err != nil {
		return "", err
	}
	if len(response.Data) == 0 {
		return "", fmt.Errorf("repository not found: %s", reposlug)
	}
	data := response.Data[0]
	return data.ID, nil

}

// DeleteRepo remove a repository from CodeClimate
func (cc CodeClimate) DeleteRepo(repoid string) error {

	err := cc.doRequest("DELETE", "repos/"+repoid, nil, nil)
	if err != nil {
		return err
	}

	return nil

}

// GetOwnOrgID returns the ID of an organization
func (cc CodeClimate) GetOwnOrgID(orgname string) (string, error) {
	var response struct {
		Data []struct {
			ID         string
			Attributes struct {
				Name string
			}
		}
	}

	err := cc.doRequest("GET", "orgs", nil, &response)
	if err != nil {
		return "", err
	}
	for _, data := range response.Data {
		if data.Attributes.Name == orgname {
			return data.ID, nil
		}
	}
	return "", fmt.Errorf("org ID not found in response data")
}

// AddRepo create a repository within an organization
// and return the reporter ID.
func (cc CodeClimate) AddRepo(reposlug string) (string, error) {
	var response struct {
		Data struct {
			Attributes struct {
				TestReporterID string `json:"test_reporter_id"`
			}
		}
	}

	parts := strings.Split(reposlug, "/")
	org := parts[0]

	orgID, err := cc.GetOwnOrgID(org)
	if err != nil {
		return "", err
	}

	URL := fmt.Sprintf("orgs/%s/repos", orgID)

	var body bytes.Buffer

	_, err = body.Write([]byte(fmt.Sprintf(ccBodyFormat, reposlug)))
	if err != nil {
		return "", err
	}

	err = cc.doRequest("POST", URL, &body, &response)
	if err != nil {
		return "", err
	}
	return response.Data.Attributes.TestReporterID, nil
}

// SetReporterIDSecret setup a secret on a github repository
// named CC_TEST_REPORTER_ID containing the specified reporterID.
// The function uses GitHub API to add the secret, so a ghtoken
// has to be specified.
func SetReporterIDSecret(githubRepo string, reporterID string, ghtoken string) error {
	parts := strings.Split(githubRepo, "/")
	owner := parts[0]
	repo := parts[1]

	ctx, client, err := githubAuth(ghtoken)
	if err != nil {
		return fmt.Errorf("unable to authorize using env GITHUB_AUTH_TOKEN: %w", err)
	}

	if err := addRepoSecret(ctx, client, owner, repo, "CC_TEST_REPORTER_ID", reporterID); err != nil {
		return err
	}

	return nil
}

// addRepoSecret will add a secret to a GitHub repo for use in GitHub Actions.
//
// Finally, the secretName and secretValue will determine the name of the secret added and it's corresponding value.
//
// The actual transmission of the secret value to GitHub using the api requires that the secret value is encrypted
// using the public key of the target repo. This encryption must be done using sodium.
//
// First, the public key of the repo is retrieved. The public key comes base64
// encoded, so it must be decoded prior to use in sodiumlib.
//
// Second, the secret value is converted into a slice of bytes.
//
// Third, the secret is encrypted with sodium.CryptoBoxSeal using the repo's decoded public key.
//
// Fourth, the encrypted secret is encoded as a base64 string to be used in a github.EncodedSecret type.
//
// Fifth, The other two properties of the github.EncodedSecret type are determined. The name of the secret to be added
// (string not base64), and the KeyID of the public key used to encrypt the secret.
// This can be retrieved via the public key's GetKeyID method.
//
// Finally, the github.EncodedSecret is passed into the GitHub client.Actions.CreateOrUpdateRepoSecret method to
// populate the secret in the GitHub repo.
func addRepoSecret(ctx context.Context, client *github.Client, owner string, repo, secretName string, secretValue string) error {
	publicKey, _, err := client.Actions.GetRepoPublicKey(ctx, owner, repo)
	if err != nil {
		return err
	}

	encryptedSecret, err := encryptSecretWithPublicKey(publicKey, secretName, secretValue)
	if err != nil {
		return err
	}

	if _, err := client.Actions.CreateOrUpdateRepoSecret(ctx, owner, repo, encryptedSecret); err != nil {
		return fmt.Errorf("Actions.CreateOrUpdateRepoSecret returned error: %v", err)
	}

	return nil
}

// githubAuth returns a GitHub client and context.
func githubAuth(token string) (context.Context, *github.Client, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	return ctx, client, nil
}

func encryptSecretWithPublicKey(publicKey *github.PublicKey, secretName string, secretValue string) (*github.EncryptedSecret, error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey.GetKey())
	if err != nil {
		return nil, fmt.Errorf("base64.StdEncoding.DecodeString was unable to decode public key: %v", err)
	}

	secretBytes := []byte(secretValue)
	encryptedBytes, exit := sodium.CryptoBoxSeal(secretBytes, decodedPublicKey)
	if exit != 0 {
		return nil, errors.New("sodium.CryptoBoxSeal exited with non zero exit code")
	}

	encryptedString := base64.StdEncoding.EncodeToString(encryptedBytes)
	keyID := publicKey.GetKeyID()
	encryptedSecret := &github.EncryptedSecret{
		Name:           secretName,
		KeyID:          keyID,
		EncryptedValue: encryptedString,
	}
	return encryptedSecret, nil
}
