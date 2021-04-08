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

func TestGetRepoID(t *testing.T) {
	cctoken, ok := os.LookupEnv("CC_TOKEN")
	if !ok {
		t.Fatalf("\nTO RUN THIS TEST, DECLARE A CC_TOKEN ENV VAR WITH A GITHUB API TOKEN\n")
	}
	cc := CodeClimate(cctoken)
	t.Run("check this repo ID", func(t *testing.T) {
		ID, err := cc.GetRepoID("parrogo/ccaddrepo")
		if assert.NoError(t, err) {
			assert.Equal(t, "606b6ca958f8666dcc00a693", ID)
		}
	})
}
func TestAddRepo(t *testing.T) {
	cctoken, ok := os.LookupEnv("CC_TOKEN")
	if !ok {
		t.Fatalf("\nTO RUN THIS TEST, DECLARE A CC_TOKEN ENV VAR WITH A GITHUB API TOKEN\n")
	}
	cc := CodeClimate(cctoken)
	t.Run("check centro", func(t *testing.T) {
		ID, err := cc.GetRepoID("parrogo/centro")
		if err == nil {
			err = cc.DeleteRepo(ID)
			if !assert.NoError(t, err) {
				return
			}
		}

		reporterID, err := cc.AddRepo("parrogo/centro")
		if !assert.NoError(t, err) {
			return
		}

		assert.Greater(t, len(reporterID), 10)
	})
}
func TestDeleteRepo(t *testing.T) {
	cctoken, ok := os.LookupEnv("CC_TOKEN")
	if !ok {
		t.Fatalf("\nTO RUN THIS TEST, DECLARE A CC_TOKEN ENV VAR WITH A GITHUB API TOKEN\n")
	}
	cc := CodeClimate(cctoken)
	t.Run("check centro", func(t *testing.T) {
		ID, err := cc.GetRepoID("parrogo/centro")
		if err != nil && err.Error() == "repository not found: parrogo/centro" {
			_, err := cc.AddRepo("parrogo/centro")
			if !assert.NoError(t, err) {
				return
			}
		}

		if !assert.NoError(t, err) {
			return
		}
		assert.NoError(t, cc.DeleteRepo(ID))
	})
}

func TestGetOwnOrgID(t *testing.T) {
	cctoken, ok := os.LookupEnv("CC_TOKEN")
	if !ok {
		t.Fatalf("\nTO RUN THIS TEST, DECLARE A CC_TOKEN ENV VAR WITH A GITHUB API TOKEN\n")
	}
	cc := CodeClimate(cctoken)

	t.Run("check wrong org name", func(t *testing.T) {
		ID, err := cc.GetOwnOrgID("nonexistentorg")
		if assert.Error(t, err) {
			assert.Equal(t, "", ID)
		}
	})
	t.Run("check parrogo ID", func(t *testing.T) {
		ID, err := cc.GetOwnOrgID("parrogo")
		if assert.NoError(t, err) {
			assert.Equal(t, "606ddb6085fc77593c001ec5", ID)
		}
	})

	t.Run("check parro-it ID", func(t *testing.T) {
		ID, err := cc.GetOwnOrgID("parro-it")
		if assert.NoError(t, err) {
			assert.Equal(t, "5610f3b56956806eff013739", ID)
		}
	})

}
