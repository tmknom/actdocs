package main

import (
	"log"
	"os"

	"github.com/tmknom/actdocs/internal/cli"
)

// Specify explicitly in ldflags
// For full details, see Makefile and .goreleaser.yml
var (
	name    = ""
	version = ""
	commit  = ""
	date    = ""
)

func main() {
	app := cli.NewApp(name, version, commit, date)
	if err := app.Run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
