package main

import (
	"log"
	"os"

	"github.com/tmknom/actdocs"
)

func main() {
	app := actdocs.NewApp()
	if err := app.Run(os.Args, os.Stdin, os.Stdout, os.Stderr); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
