package main

import (
	"flag"
	"fmt"
	"os"
)

// Version of the command
var Version string = "development"

var options struct {
	version bool
	repo    string
	token   string
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
	os.Exit(1)
}

func usage(msg string) {
	fmt.Fprintf(os.Stderr, "Wrong usage: %s\n", msg)
	flag.Usage()
	os.Exit(1)
}

func main() {

	repoDefault := os.Getenv("GITHUB_REPOSITORY")
	tokenDefault := os.Getenv("CC_TOKEN")

	flag.BoolVar(&options.version, "v", false, "print version of the command to stdout")

	flag.StringVar(&options.repo, "r", repoDefault, "GitHub user/repo of the repository you want to add")
	flag.StringVar(&options.token, "t", tokenDefault, "CodeClimate API token")

	flag.Parse()

	if options.version {
		fmt.Println(Version)
		os.Exit(0)
	}

	if options.repo == "" {
		usage("r flag not specified")
	}

	if options.token == "" {
		usage("t flag not specified")
	}

	fmt.Println(options.repo, options.token)
	reporterID, err := "boh", error(nil) //ccaddrepo.AddOnCodeClimate(options.repo, options.token)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
	fmt.Println(reporterID)
}
