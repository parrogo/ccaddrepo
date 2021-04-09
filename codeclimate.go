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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const ccBaseURL = "https://api.codeclimate.com/v1/"
const ccGetURL = "https://api.codeclimate.com/v1/repos"
const ccURL = "https://api.codeclimate.com/v1/github/repos"
const ccBodyFormat = `{"data":{"type": "repos","attributes": {"url": "https://github.com/%s"}}}`

// CodeClimate represent an authenticated
// session on Code Climate.
type CodeClimate string

func (cc CodeClimate) doRequest(method string, URL string, body io.Reader, response interface{}) (string, error) {
	req, err := http.NewRequest(method, ccBaseURL+URL, body)
	if err != nil {
		return "", err
	}

	req.Header.Add("Accept", "application/vnd.api+json")
	req.Header.Add("Authorization", "Token token="+string(cc))
	req.Header.Add("Content-Type", "application/vnd.api+json")
	var c http.Client
	res, err := c.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	resbuf, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	if response == nil {
		return string(resbuf), nil
	}

	err = json.Unmarshal(resbuf, response)
	if err != nil {
		return "", err
	}

	return string(resbuf), nil

}

// GetRepoID returns the ID of a repository
func (cc CodeClimate) GetRepoID(reposlug string) (string, error) {
	var response struct {
		Data []struct {
			ID string
		}
	}

	rowdata, err := cc.doRequest("GET", "repos?github_slug="+reposlug, nil, &response)
	if err != nil {
		return "", err
	}
	if len(response.Data) == 0 {
		return "", fmt.Errorf("repository not found: %s\nRESPONSE: %s", reposlug, rowdata)
	}
	data := response.Data[0]
	return data.ID, nil

}

// DeleteRepo remove a repository from CodeClimate
func (cc CodeClimate) DeleteRepo(repoid string) error {

	_, err := cc.doRequest("DELETE", "repos/"+repoid, nil, nil)
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

	rowdata, err := cc.doRequest("GET", "orgs", nil, &response)
	if err != nil {
		return "", err
	}
	for _, data := range response.Data {
		if data.Attributes.Name == orgname {
			return data.ID, nil
		}
	}

	return "", fmt.Errorf("org ID `%s` not found in response data\nRESPONSE:\n%v", orgname, rowdata)
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

	_, err = cc.doRequest("POST", URL, &body, &response)
	if err != nil {
		return "", err
	}
	return response.Data.Attributes.TestReporterID, nil
}
