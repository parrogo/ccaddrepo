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

func TestAddOnCodeClimate(t *testing.T) {
	cctoken, ok := os.LookupEnv("CC_TOKEN")
	if !ok {
		t.Fatalf("\nTO RUN THIS TEST, DECLARE A CC_TOKEN ENV VAR WITH CODE CLIMATE TOKEN\n")
	}

	reporterID, err := AddOnCodeClimate("parrogo/ccaddrepo", cctoken)
	if assert.NoError(t, err) && assert.NotEmpty(t, reporterID) {
		assert.Greater(t, len(reporterID), 6)
	}
}

func TestParseResponse(t *testing.T) {
	response, err := fs.ReadFile(fixtureFS, "ccaddresponse.json")
	if !assert.NoError(t, err) {
		return
	}
	res, err := parse(response)
	assert.Equal(t, "fakefakefake", res)

}

func TestSetReporterIDSecret(t *testing.T) {
	err := SetReporterIDSecret("parro-it/examplerepo", "42", "ghtoken")
	if assert.NoError(t, err) {
		assert.Equal(t, nil, err)
	}
}
