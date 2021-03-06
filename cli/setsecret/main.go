package main

import (
	"bufio"
	"flag"
	"fmt"
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
	badgeID    string
	ID         string
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
	flag.StringVar(&options.ID, "id", "", "CodeCLimate repo ID")
	flag.StringVar(&options.reporterID, "testrepid", "", "CodeCLimate repo test reporter ID")
	flag.StringVar(&options.badgeID, "badgeid", "", "CodeCLimate repo badge ID")

	flag.Parse()

	if options.version {
		fmt.Println(Version)
		os.Exit(0)
	}

	if options.reporterID == "" {
		stdinbuf := bufio.NewReader(os.Stdin)

		linebuf, err := stdinbuf.ReadBytes('\n')
		fatal(err)
		options.ID = strings.Trim(string(linebuf), " \t\r\n")

		linebuf, err = stdinbuf.ReadBytes('\n')
		fatal(err)
		options.reporterID = strings.Trim(string(linebuf), " \t\r\n")

		linebuf, err = stdinbuf.ReadBytes('\n')
		fatal(err)
		options.badgeID = strings.Trim(string(linebuf), " \t\r\n")
	}

	if options.ID == "" {
		usage("id flag not specified")
	}
	if options.badgeID == "" {
		usage("badgeid flag not specified")
	}
	if options.reporterID == "" {
		usage("testrepid flag not specified")
	}
	if options.repo == "" {
		usage("r flag not specified")
	}

	if options.token == "" {
		usage("t flag not specified")
	}

	//fmt.Println(options.repo, options.token, options.reporterID)
	fatal(ccaddrepo.SetReporterIDSecret(ccaddrepo.SecretsOptions{
		RepoSlug:   options.repo,
		GHToken:    options.token,
		ReporterID: options.reporterID,
		BadgeID:    options.badgeID,
		ID:         options.ID,
	}))

}
