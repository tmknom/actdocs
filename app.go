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

	rootCmd.AddCommand(&cobra.Command{
		Use:   "workflow",
		Short: "Generate docs for Reusable Workflows",
		RunE: func(cmd *cobra.Command, args []string) error {
			return NewWorkflowCmd().Run(cmd, args)
		},
	})

	return rootCmd.Execute()
}
