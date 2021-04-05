package ccaddrepo

import (
	"embed"
	"io/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed fixtures
var fixtureRootFS embed.FS
var fixtureFS, _ = fs.Sub(fixtureRootFS, "fixtures")

func TestPlaceHolder(t *testing.T) {
	content, err := fs.ReadFile(fixtureFS, "placeholder")
	if assert.NoError(t, err) {
		assert.Equal(t, "this is a placeholder", string(content))
	}
}

//func TestAddOnCodeClimate(t *testing.T) {
//	cctoken, ok := os.LookupEnv("CC_TOKEN")
//	if !ok {
//		t.Fatalf("\nTO RUN THIS TEST, DECLARE A CC_TOKEN ENV VAR WITH CODE CLIMATE TOKEN\n")
//	}
//
//	reporterID, err := AddOnCodeClimate("parrogo/ccaddrepo", cctoken)
//	if assert.NoError(t, err) && assert.NotEmpty(t, reporterID) {
//		assert.Greater(t, len(reporterID), 6)
//	}
//}

func TestParseResponse(t *testing.T) {
	response, err := fs.ReadFile(fixtureFS, "ccaddresponse.json")
	if !assert.NoError(t, err) {
		return
	}
	res, err := parse(response)
	assert.Equal(t, "fakefakefake", res)

}

func TestSetReporterIDSecret(t *testing.T) {
	cctoken, ok := os.LookupEnv("GH_WORKFLOW")
	if !ok {
		t.Fatalf("\nTO RUN THIS TEST, DECLARE A GH_WORKFLOW ENV VAR WITH A GITHUB API TOKEN\n")
	}

	err := SetReporterIDSecret("parrogo/ccaddrepo", "42", cctoken)
	if assert.NoError(t, err) {
		assert.Equal(t, nil, err)
	}
}
