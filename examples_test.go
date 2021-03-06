package ccaddrepo_test

import (
	"embed"
	"fmt"
	"io/fs"

	"github.com/parrogo/ccaddrepo"
)

//go:embed fixtures
var fixtureRootFS embed.FS
var fixtureFS, _ = fs.Sub(fixtureRootFS, "fixtures")

// This example show how to use ccaddrepo.AddOnCodeClimate()
func ExampleCodeClimate() {
	cc := ccaddrepo.CodeClimate("token")
	reporterID, err := cc.AddRepo("author/repo")
	if err != nil {
		panic(err)
	}
	fmt.Println(reporterID)

}

func ExampleSetReporterIDSecret() {
	err := ccaddrepo.SetReporterIDSecret(ccaddrepo.SecretsOptions{
		RepoSlug:   "",
		GHToken:    "",
		ReporterID: "",
		BadgeID:    "",
		ID:         "",
	})
	if err != nil {
		panic(err)
	}
}
