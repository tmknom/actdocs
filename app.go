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
	*IO
	*Ldflags
	debug bool
}

func NewApp(name string, version string, commit string, date string) *App {
	return &App{
		Ldflags: NewLdflags(name, version, commit, date),
		debug:   false,
	}
}

func (a *App) Run(args []string, inReader io.Reader, outWriter, errWriter io.Writer) error {
	rootCmd := &cobra.Command{
		Use:     a.Name,
		Short:   "Generate documentation from Custom Actions and Reusable Workflows",
		Version: a.Version,
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
	config := DefaultGlobalConfig()
	rootCmd.PersistentFlags().StringVar(&config.Format, "format", DefaultFormat, "output format [markdown json]")
	rootCmd.PersistentFlags().BoolVarP(&config.Sort, "sort", "s", DefaultSort, "sort items by name and required")
	rootCmd.PersistentFlags().BoolVar(&config.SortByName, "sort-by-name", DefaultSortByName, "sort items by name")
	rootCmd.PersistentFlags().BoolVar(&config.SortByRequired, "sort-by-required", DefaultSortByRequired, "sort items by required")

	// setup version option
	version := fmt.Sprintf("%s version %s (%s)", a.Name, a.Version, a.Date)
	rootCmd.SetVersionTemplate(version)

	// setup commands
	rootCmd.AddCommand(a.newGenerateCommand(config))
	rootCmd.AddCommand(a.newInjectCommand(config))

	return rootCmd.Execute()
}

func (a *App) newGenerateCommand(globalConfig *GlobalConfig) *cobra.Command {
	config := NewGeneratorConfig(globalConfig)
	return &cobra.Command{
		Use:   "generate",
		Short: "Generate documentation",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.SetPrefix(fmt.Sprintf("[%s] [%s] ", a.Name, cmd.Name()))
			log.Printf("start: command = %s, config = %#v", cmd.Name(), config)
			runner := NewGenerator(config, a.IO, NewYamlFile(args))
			return runner.Run()
		},
	}
}

func (a *App) newInjectCommand(globalConfig *GlobalConfig) *cobra.Command {
	config := NewInjectorConfig(globalConfig)
	command := &cobra.Command{
		Use:   "inject",
		Short: "Inject generated documentation to existing file",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.SetPrefix(fmt.Sprintf("[%s] [%s] ", a.Name, cmd.Name()))
			log.Printf("start: command = %s, config = %#v", cmd.Name(), config)
			runner := NewInjector(config, a.IO, NewYamlFile(args))
			return runner.Run()
		},
	}

	command.PersistentFlags().StringVarP(&config.OutputFile, "file", "f", "", "file path to insert output into (default \"\")")
	command.PersistentFlags().BoolVar(&config.DryRun, "dry-run", false, "dry run")
	return command
}

func (a *App) setupLog(args []string) {
	log.SetOutput(io.Discard)
	if a.debug {
		log.SetOutput(os.Stderr)
		log.SetPrefix(fmt.Sprintf("[%s] ", a.Name))
	}
	log.Printf("start: %s", strings.Join(os.Args, " "))
	log.Printf("args: %q", args)
	log.Printf("ldflags: %+v", a.Ldflags)
}
