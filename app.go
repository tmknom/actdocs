package actdocs

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type App struct {
	debug bool
}

func NewApp() *App {
	return &App{
		debug: false,
	}
}

func (a *App) Run(stdin io.Reader, stdout, stderr io.Writer) error {
	rootCmd := &cobra.Command{
		Use:   "actdocs",
		Short: "Generate documentation from Custom Actions and Reusable Workflows",
	}

	// setup log
	rootCmd.PersistentFlags().BoolVar(&a.debug, "debug", false, "enable debug log")
	cobra.OnInitialize(func() { a.setupLog() })

	// setup I/O
	rootCmd.SetIn(stdin)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)

	// setup global flags
	config := NewConfig(rootCmd.OutOrStdout())
	rootCmd.PersistentFlags().StringVarP(&config.OutputFile, "output-file", "o", "", "file path to insert output into (default \"\")")
	rootCmd.PersistentFlags().BoolVar(&config.SortByName, "sort-by-name", false, "sort items by name")

	// setup commands
	rootCmd.AddCommand(&cobra.Command{
		Use:   "generate",
		Short: "Generate docs",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.SetPrefix(fmt.Sprintf("[actdocs] [%s] ", cmd.Name()))
			log.Printf("start: command = %s, config = %#v", cmd.Name(), config)
			generateCmd := NewGenerateCmd(config, cmd.InOrStdin(), cmd.OutOrStdout(), cmd.ErrOrStderr())
			if len(args) > 0 {
				generateCmd.filename = args[0]
			}
			return generateCmd.Run()
		},
	})

	return rootCmd.Execute()
}

func (a *App) setupLog() {
	log.SetOutput(io.Discard)
	if a.debug {
		log.SetOutput(os.Stderr)
		log.SetPrefix("[actdocs] ")
	}
	log.Printf("start: %s", strings.Join(os.Args, " "))
}
