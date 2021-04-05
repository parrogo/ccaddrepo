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
)

const ccGetURL = "https://api.codeclimate.com/v1/repos"
const ccURL = "https://api.codeclimate.com/v1/github/repos"
const ccBodyFormat = `{"data":{"type": "repos","attributes": {"url": "https://github.com/%s"}}}`

// AddOnCodeClimate ask CodeClimate servers to add specified repo
// The requests uses CodeClimate API, cctoken is an API token
// that you can get here: https://codeclimate.com/profile/tokens
//
// The function return a string containing a CodeClimate TEST REPORTER ID
// or in case of failure, an error value.
func AddOnCodeClimate(githubRepo string, cctoken string) (string, error) {
	var body bytes.Buffer

	_, err := body.Write([]byte(fmt.Sprintf(ccBodyFormat, githubRepo)))
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", ccURL, &body)
	if err != nil {
		return "", err
	}

	req.Header.Add("Accept", "application/vnd.api+json")
	req.Header.Add("Authorization", "Token token="+cctoken)
	req.Header.Add("Content-Type", "application/vnd.api+json")

	var c http.Client
	res, err := c.Do(req)
	if err != nil {
		return "", err
	}

	fmt.Println(res.Status)
	// 201 Created
	defer res.Body.Close()

	resbuf, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return parse(resbuf)
}

func parse(resbuf []byte) (string, error) {
	fmt.Println(string(resbuf))
	var results commandResult
	err := json.Unmarshal(resbuf, &results)
	if err != nil {
		return "", err
	}
	fmt.Println(results)
	return results.Data.Attributes.TestReporterID, nil
}

type commandResult struct {
	Data struct {
		Attributes struct {
			TestReporterID string `json:"test_reporter_id"`
		}
	}
}

// SetReporterIDSecret setup a secret on a github repository
// named CC_TEST_REPORTER_ID containing the specified reporterID.
// The function uses GitHub API to add the secret, so a ghtoken
// has to be specified.
func SetReporterIDSecret(githubRepo string, reporterID string, ghtoken string) error {
	return nil
}

// echo '{"data":{"type": "repos","attributes": {"url": "https://github.com/'${GITHUB_REPOSITORY}'"}}}' > body.json
// curl -s \
//   -H "Accept: application/vnd.api+json" \
//   -H "Authorization: Token token=${CODE_CLIMATE_TOKEN}" \
//   -H "Content-Type: application/vnd.api+json" \
//   -d @body.json \
//   https://api.codeclimate.com/v1/github/repos > .codeclimate.json

//       - name: Test & publish code coverage
//         uses: paambaati/codeclimate-action@v2.7.5
//         env:
//           CC_TEST_REPORTER_ID: ${---{ secrets.CC_TEST_REPORTER_ID }}
//         with:
//           coverageCommand: go test -coverprofile=c.out ./...
//           prefix: github.com/{{.Author}}/{{.RepoName}}
