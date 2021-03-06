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
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	sodium "github.com/GoKillers/libsodium-go/cryptobox"
	"github.com/google/go-github/v34/github"
	"golang.org/x/oauth2"
)

// SecretsOptions contains all the options
// to setup Code Climate related secrets in a
// GH repo.
type SecretsOptions struct {
	RepoSlug   string
	GHToken    string
	ReporterID string
	BadgeID    string
	ID         string
}

// SetReporterIDSecret setup a secret on a github repository
// named CC_TEST_REPORTER_ID containing the specified reporterID.
// The function uses GitHub API to add the secret, so a ghtoken
// has to be specified.
func SetReporterIDSecret(options SecretsOptions) error {
	parts := strings.Split(options.RepoSlug, "/")
	owner := parts[0]
	repo := parts[1]

	ctx, client, err := githubAuth(options.GHToken)
	if err != nil {
		return fmt.Errorf("unable to authorize using env GITHUB_AUTH_TOKEN: %w", err)
	}

	if err := addRepoSecret(ctx, client, owner, repo, "CC_REPO_ID", options.ID); err != nil {
		return err
	}
	fmt.Println("CodeCLimate repo ID stored in secrets.CC_REPO_ID")

	if err := addRepoSecret(ctx, client, owner, repo, "CC_TEST_REPORTER_ID", options.ReporterID); err != nil {
		return err
	}
	fmt.Println("CodeCLimate test reporter ID stored in secrets.CC_TEST_REPORTER_ID")

	if err := addRepoSecret(ctx, client, owner, repo, "CC_BADGE", options.BadgeID); err != nil {
		return err
	}
	fmt.Println("CodeCLimate badge ID stored in secrets.CC_BADGE")

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
