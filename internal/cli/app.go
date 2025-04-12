package cli

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tmknom/actdocs/internal/format"
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
	formatterConfig := format.DefaultFormatterConfig()
	rootCmd.PersistentFlags().StringVar(&formatterConfig.Format, "format", format.DefaultFormat, "output format [markdown json]")
	rootCmd.PersistentFlags().BoolVar(&formatterConfig.Omit, "omit", format.DefaultOmit, "omit for markdown if item not exists")
	rootCmd.PersistentFlags().BoolVarP(&formatterConfig.Sort, "sort", "s", format.DefaultSort, "sort items by name and required")
	rootCmd.PersistentFlags().BoolVar(&formatterConfig.SortByName, "sort-by-name", format.DefaultSortByName, "sort items by name")
	rootCmd.PersistentFlags().BoolVar(&formatterConfig.SortByRequired, "sort-by-required", format.DefaultSortByRequired, "sort items by required")

	// setup version option
	version := fmt.Sprintf("%s version %s", AppName, AppVersion)
	rootCmd.SetVersionTemplate(version)

	// setup commands
	rootCmd.AddCommand(a.newGenerateCommand(formatterConfig))
	rootCmd.AddCommand(a.newInjectCommand(formatterConfig))

	return rootCmd.Execute()
}

func (a *App) newGenerateCommand(formatterConfig *format.FormatterConfig) *cobra.Command {
	cfg := NewGeneratorConfig(formatterConfig)
	return &cobra.Command{
		Use:   "generate",
		Short: "Generate documentation",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.SetPrefix(fmt.Sprintf("[%s] [%s] ", AppName, cmd.Name()))
			log.Printf("start: command = %s, formatterConfig = %#v", cmd.Name(), cfg)
			if len(args) > 0 {
				runner := NewGenerator(cfg, a.IO, args[0])
				return runner.Run()
			}
			return cmd.Usage()
		},
	}
}

func (a *App) newInjectCommand(formatterConfig *format.FormatterConfig) *cobra.Command {
	cfg := NewInjectorConfig(formatterConfig)
	command := &cobra.Command{
		Use:   "inject",
		Short: "Inject generated documentation to existing file",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.SetPrefix(fmt.Sprintf("[%s] [%s] ", AppName, cmd.Name()))
			log.Printf("start: command = %s, formatterConfig = %#v", cmd.Name(), cfg)
			if len(args) > 0 {
				runner := NewInjector(cfg, a.IO, args[0])
				return runner.Run()
			}
			return cmd.Usage()
		},
	}

	command.PersistentFlags().StringVarP(&cfg.OutputFile, "file", "f", "", "file path to insert output into (default \"\")")
	command.PersistentFlags().BoolVar(&cfg.DryRun, "dry-run", false, "dry run")
	return command
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
