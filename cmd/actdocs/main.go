package main

import (
	"log"
	"os"

	"github.com/tmknom/actdocs"
)

var (
	version  = ""
	revision = ""
)

func main() {
	app := actdocs.NewApp(version, revision)
	if err := app.Run(os.Stdin, os.Stdout, os.Stderr); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
