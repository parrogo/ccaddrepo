// Package ccaddrepo export a function that adds a repository to CodeClimate
// and add its report id as a secret to GitHub repository.
package ccaddrepo

// Func answers
func Func() int {
	return 42
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
