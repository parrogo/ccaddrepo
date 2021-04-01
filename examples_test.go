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

// This example show how to use ccaddrepo.Func()
func ExampleFunc() {
	fmt.Println(ccaddrepo.Func())
	// Output: 42
}
