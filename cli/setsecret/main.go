package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/parrogo/ccaddrepo"
)

// Version of the command
var Version string = "development"

var options struct {
	version    bool
	repo       string
	token      string
	reporterID string
}

func fatal(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}

func usage(msg string) {
	fmt.Fprintf(os.Stderr, "Wrong usage: %s\n", msg)
	flag.Usage()
	os.Exit(1)
}

func main() {

	repoDefault := os.Getenv("GITHUB_REPOSITORY")
	tokenDefault := os.Getenv("GH_WORKFLOW")

	flag.BoolVar(&options.version, "v", false, "print version of the command to stdout")

	flag.StringVar(&options.repo, "r", repoDefault, "GitHub user/repo of the repository you want to add")
	flag.StringVar(&options.token, "t", tokenDefault, "GitHub API token")
	flag.StringVar(&options.reporterID, "id", "", "CodeCLimate reporter ID")

	flag.Parse()

	if options.version {
		fmt.Println(Version)
		os.Exit(0)
	}

	if options.reporterID == "" {
		stdin, err := io.ReadAll(os.Stdin)
		fatal(err)
		options.reporterID = strings.Trim(string(stdin), " \t\r\n")
	}

	if options.reporterID == "" {
		usage("id flag not specified")
	}

	if options.repo == "" {
		usage("r flag not specified")
	}

	if options.token == "" {
		usage("t flag not specified")
	}

	//fmt.Println(options.repo, options.token, options.reporterID)
	fatal(ccaddrepo.SetReporterIDSecret(options.repo, options.reporterID, options.token))

	fmt.Println("CodeCLimate reporter ID stored in secrets.CC_TEST_REPORTER_ID")
}
