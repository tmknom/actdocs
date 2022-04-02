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
	name    string
	version string
	commit  string
	date    string
	debug   bool
}

func NewApp(name string, version string, commit string, date string) *App {
	return &App{
		name:    name,
		version: version,
		commit:  commit,
		date:    date,
		debug:   false,
	}
}

func (a *App) Run(stdin io.Reader, stdout, stderr io.Writer) error {
	rootCmd := &cobra.Command{
		Use:     a.name,
		Short:   "Generate documentation from Custom Actions and Reusable Workflows",
		Version: a.version,
	}

	// setup log
	rootCmd.PersistentFlags().BoolVar(&a.debug, "debug", false, "show debugging output")
	cobra.OnInitialize(func() { a.setupLog() })

	// setup I/O
	rootCmd.SetIn(stdin)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)

	// setup global flags
	config := NewConfig(rootCmd.OutOrStdout())
	rootCmd.PersistentFlags().StringVar(&config.Format, "format", "markdown", "output format [markdown json]")
	rootCmd.PersistentFlags().BoolVarP(&config.Sort, "sort", "s", false, "sort items by name and required")
	rootCmd.PersistentFlags().BoolVar(&config.SortByName, "sort-by-name", false, "sort items by name")
	rootCmd.PersistentFlags().BoolVar(&config.SortByRequired, "sort-by-required", false, "sort items by required")

	// setup version option
	version := fmt.Sprintf("%s version %s (%s)", a.name, a.version, a.date)
	rootCmd.SetVersionTemplate(version)

	// setup commands
	rootCmd.AddCommand(a.newGenerateCommand(config))
	rootCmd.AddCommand(a.newInjectCommand(config))

	return rootCmd.Execute()
}

func (a *App) newGenerateCommand(config *Config) *cobra.Command {
	return &cobra.Command{
		Use:   "generate",
		Short: "Generate documentation",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.SetPrefix(fmt.Sprintf("[%s] [%s] ", a.name, cmd.Name()))
			log.Printf("start: command = %s, config = %#v", cmd.Name(), config)
			generateCmd := NewGenerateCmd(config, cmd.InOrStdin(), cmd.OutOrStdout(), cmd.ErrOrStderr())
			if len(args) > 0 {
				generateCmd.filename = args[0]
			}
			return generateCmd.Run()
		},
	}
}

func (a *App) newInjectCommand(config *Config) *cobra.Command {
	command := &cobra.Command{
		Use:   "inject",
		Short: "Inject generated documentation to existing file",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.SetPrefix(fmt.Sprintf("[%s] [%s] ", a.name, cmd.Name()))
			log.Printf("start: command = %s, config = %#v", cmd.Name(), config)
			injectCmd := NewInjectCmd(config, cmd.InOrStdin(), cmd.OutOrStdout(), cmd.ErrOrStderr())
			if len(args) > 0 {
				injectCmd.filename = args[0]
			}
			return injectCmd.Run()
		},
	}

	command.PersistentFlags().StringVarP(&config.OutputFile, "file", "f", "", "file path to insert output into (default \"\")")
	return command
}

func (a *App) setupLog() {
	log.SetOutput(io.Discard)
	if a.debug {
		log.SetOutput(os.Stderr)
		log.SetPrefix(fmt.Sprintf("[%s] ", a.name))
	}
	log.Printf("start: %s", strings.Join(os.Args, " "))
	log.Printf("name: %s, version: %s, date: %s, commit: %s", a.name, a.version, a.date, a.commit)
}
