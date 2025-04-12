package cli

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tmknom/actdocs/internal/config"
)

type App struct {
	*config.IO
	*config.Ldflags
	debug bool
}

func NewApp(name string, version string, commit string, date string) *App {
	return &App{
		Ldflags: config.NewLdflags(name, version, commit, date),
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
	a.IO = config.NewIO(rootCmd.InOrStdin(), rootCmd.OutOrStdout(), rootCmd.ErrOrStderr())

	// setup log
	rootCmd.PersistentFlags().BoolVar(&a.debug, "debug", false, "show debugging output")
	cobra.OnInitialize(func() { a.setupLog(args) })

	// setup global flags
	cfg := config.DefaultGlobalConfig()
	rootCmd.PersistentFlags().StringVar(&cfg.Format, "format", "markdown", "output format [markdown json]")
	rootCmd.PersistentFlags().BoolVar(&cfg.Omit, "omit", false, "omit for markdown if item not exists")
	rootCmd.PersistentFlags().BoolVarP(&cfg.Sort, "sort", "s", false, "sort items by name and required")
	rootCmd.PersistentFlags().BoolVar(&cfg.SortByName, "sort-by-name", false, "sort items by name")
	rootCmd.PersistentFlags().BoolVar(&cfg.SortByRequired, "sort-by-required", false, "sort items by required")

	// setup version option
	version := fmt.Sprintf("%s version %s (%s)", a.Name, a.Version, a.Date)
	rootCmd.SetVersionTemplate(version)

	// setup commands
	rootCmd.AddCommand(a.newGenerateCommand(cfg))
	rootCmd.AddCommand(a.newInjectCommand(cfg))

	return rootCmd.Execute()
}

func (a *App) newGenerateCommand(globalConfig *config.GlobalConfig) *cobra.Command {
	cfg := NewGeneratorConfig(globalConfig)
	return &cobra.Command{
		Use:   "generate",
		Short: "Generate documentation",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.SetPrefix(fmt.Sprintf("[%s] [%s] ", a.Name, cmd.Name()))
			log.Printf("start: command = %s, config = %#v", cmd.Name(), cfg)
			if len(args) > 0 {
				runner := NewGenerator(cfg, a.IO, args[0])
				return runner.Run()
			}
			return cmd.Usage()
		},
	}
}

func (a *App) newInjectCommand(globalConfig *config.GlobalConfig) *cobra.Command {
	cfg := NewInjectorConfig(globalConfig)
	command := &cobra.Command{
		Use:   "inject",
		Short: "Inject generated documentation to existing file",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.SetPrefix(fmt.Sprintf("[%s] [%s] ", a.Name, cmd.Name()))
			log.Printf("start: command = %s, config = %#v", cmd.Name(), cfg)
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
		log.SetPrefix(fmt.Sprintf("[%s] ", a.Name))
	}
	log.Printf("start: %s", strings.Join(os.Args, " "))
	log.Printf("args: %q", args)
	log.Printf("ldflags: %+v", a.Ldflags)
}

func (a *App) isDebug() bool {
	switch os.Getenv("ACTDOCS_DEBUG") {
	case "true", "1", "yes":
		return true
	default:
		return false
	}
}
