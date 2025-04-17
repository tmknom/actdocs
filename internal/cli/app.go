package cli

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tmknom/actdocs/internal/conf"
)

// AppName is the cli name (set by main.go)
var AppName string

// AppVersion is the current version (set by main.go)
var AppVersion string

type App struct {
	*IO
	debug bool
}

func NewApp(name string, version string, commit string, date string) *App {
	return &App{
		debug: false,
	}
}

func (a *App) Run(args []string, inReader io.Reader, outWriter, errWriter io.Writer) error {
	rootCmd := &cobra.Command{
		Use:     AppName,
		Short:   "Generate documentation from Custom Actions and Reusable Workflows",
		Version: AppVersion,
	}

	// override default settings
	rootCmd.SetArgs(args)
	rootCmd.SetIn(inReader)
	rootCmd.SetOut(outWriter)
	rootCmd.SetErr(errWriter)
	a.IO = NewIO(rootCmd.InOrStdin(), rootCmd.OutOrStdout(), rootCmd.ErrOrStderr())

	// setup log
	rootCmd.PersistentFlags().BoolVar(&a.debug, "debug", false, "show debugging output")
	cobra.OnInitialize(func() { a.setupLog(args) })

	// setup global flags
	formatterConfig := conf.DefaultFormatterConfig()
	sortConfig := conf.DefaultSortConfig()
	rootCmd.PersistentFlags().StringVar(&formatterConfig.Format, "format", conf.DefaultFormat, "output format [markdown json]")
	rootCmd.PersistentFlags().BoolVar(&formatterConfig.Omit, "omit", conf.DefaultOmit, "omit for markdown if item not exists")
	rootCmd.PersistentFlags().BoolVarP(&sortConfig.Sort, "sort", "s", conf.DefaultSort, "sort items by name and required")
	rootCmd.PersistentFlags().BoolVar(&sortConfig.SortByName, "sort-by-name", conf.DefaultSortByName, "sort items by name")
	rootCmd.PersistentFlags().BoolVar(&sortConfig.SortByRequired, "sort-by-required", conf.DefaultSortByRequired, "sort items by required")

	// setup version option
	version := fmt.Sprintf("%s version %s", AppName, AppVersion)
	rootCmd.SetVersionTemplate(version)

	// setup commands
	rootCmd.AddCommand(NewGenerateCommand(formatterConfig, sortConfig, a.IO))
	rootCmd.AddCommand(NewInjectCommand(formatterConfig, sortConfig, a.IO))

	return rootCmd.Execute()
}

func (a *App) setupLog(args []string) {
	log.SetOutput(io.Discard)
	if a.isDebug() || a.debug {
		log.SetOutput(os.Stderr)
		log.SetPrefix(fmt.Sprintf("[%s] ", AppName))
	}
	log.Printf("start: %s", strings.Join(os.Args, " "))
	log.Printf("args: %q", args)
}

func (a *App) isDebug() bool {
	switch os.Getenv("ACTDOCS_DEBUG") {
	case "true", "1", "yes":
		return true
	default:
		return false
	}
}

type IO struct {
	InReader  io.Reader
	OutWriter io.Writer
	ErrWriter io.Writer
}

func NewIO(inReader io.Reader, outWriter, errWriter io.Writer) *IO {
	return &IO{
		InReader:  inReader,
		OutWriter: outWriter,
		ErrWriter: errWriter,
	}
}

const (
	ActionRegex   = `(?m)^[\s]*runs:`
	WorkflowRegex = `(?m)^[\s]*workflow_call:`
)
