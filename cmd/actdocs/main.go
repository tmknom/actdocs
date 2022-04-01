package main

import (
	"log"
	"os"

	"github.com/tmknom/actdocs"
)

var (
	// Specify explicitly in ldflags
	name = ""

	// GoReleaser sets these values by default
	// https://goreleaser.com/customization/build/
	version = ""
	commit  = ""
	date    = ""
)

func main() {
	app := actdocs.NewApp(name, version, commit, date)
	if err := app.Run(os.Stdin, os.Stdout, os.Stderr); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
