package actdocs

import (
	"io"

	"github.com/spf13/cobra"
)

type App struct{}

func NewApp() *App {
	return &App{}
}

func (app *App) Run(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	rootCmd := &cobra.Command{
		Use:   "actdocs",
		Short: "Generate documentation from Custom Actions and Reusable Workflows",
	}
	rootCmd.SetArgs(args[1:])
	rootCmd.SetIn(stdin)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)

	config := NewTemplateConfig(rootCmd.OutOrStdout())
	rootCmd.PersistentFlags().StringVarP(&config.OutputFile, "output-file", "o", "", "file path to insert output into (default \"\")")
	rootCmd.PersistentFlags().BoolVar(&config.Stdout, "stdout", false, "write result to stdout")

	rootCmd.AddCommand(&cobra.Command{
		Use:   "workflow",
		Short: "Generate docs for Reusable Workflows",
		RunE: func(cmd *cobra.Command, args []string) error {
			return NewWorkflowCmd(config, args, cmd.InOrStdin(), cmd.OutOrStdout(), cmd.ErrOrStderr()).Run()
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "action",
		Short: "Generate docs for Custom Actions",
		RunE: func(cmd *cobra.Command, args []string) error {
			return NewActionCmd(config, args, cmd.InOrStdin(), cmd.OutOrStdout(), cmd.ErrOrStderr()).Run()
		},
	})

	return rootCmd.Execute()
}
