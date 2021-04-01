package ccaddrepo

import (
	"embed"
	"io/fs"
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
	reporterID, err := AddOnCodeClimate("", "")
	if assert.NoError(t, err) {
		assert.Equal(t, "", reporterID)
	}
}

func TestSetReporterIDSecret(t *testing.T) {
	err := SetReporterIDSecret("parro-it/examplerepo", "42", "ghtoken")
	if assert.NoError(t, err) {
		assert.Equal(t, nil, err)
	}
}
