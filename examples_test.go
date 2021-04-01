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
func ExampleAddOnCodeClimate() {
	reporterID, err := ccaddrepo.AddOnCodeClimate("", "")
	if err != nil {
		panic(err)
	}
	fmt.Println(reporterID)

}

func ExampleSetReporterIDSecret() {
	err := ccaddrepo.SetReporterIDSecret("parro-it/examplerepo", "42", "ghtoken")
	if err != nil {
		panic(err)
	}
}
