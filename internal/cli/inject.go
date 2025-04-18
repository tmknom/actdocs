package cli

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/tmknom/actdocs/internal/action"
	"github.com/tmknom/actdocs/internal/conf"
	"github.com/tmknom/actdocs/internal/workflow"
)

func NewInjectCommand(formatter *conf.FormatterConfig, sort *conf.SortConfig, io *IO) *cobra.Command {
	option := &InjectOption{IO: io}
	command := &cobra.Command{
		Use:   "inject",
		Short: "Inject generated documentation to existing file",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.SetPrefix(fmt.Sprintf("[%s] [%s] ", AppName, cmd.Name()))
			log.Printf("start: command = %s, option = %#v", cmd.Name(), option)
			if len(args) > 0 {
				runner := NewInjectRunner(args[0], formatter, sort, option)
				return runner.Run()
			}
			return cmd.Usage()
		},
	}

	command.PersistentFlags().StringVarP(&option.OutputFile, "file", "f", "", "file path to insert output into (default \"\")")
	command.PersistentFlags().BoolVar(&option.DryRun, "dry-run", false, "dry run")
	return command
}

type InjectRunner struct {
	source string
	*conf.FormatterConfig
	*conf.SortConfig
	*InjectOption
}

func NewInjectRunner(source string, formatter *conf.FormatterConfig, sort *conf.SortConfig, option *InjectOption) *InjectRunner {
	return &InjectRunner{
		source:          source,
		FormatterConfig: formatter,
		SortConfig:      sort,
		InjectOption:    option,
	}
}

type InjectOption struct {
	OutputFile string
	DryRun     bool
	*IO
}

func (r *InjectRunner) Run() error {
	reader := &SourceReader{}
	yaml, err := reader.Read(r.source)
	if err != nil {
		return err
	}

	dest, err := os.Open(r.OutputFile)
	if err != nil {
		return err
	}
	defer func(file *os.File) { err = file.Close() }(dest)

	result, err := Inject(yaml, dest, r.FormatterConfig, r.SortConfig)
	if err != nil {
		return err
	}

	if r.DryRun {
		_, err = fmt.Fprintf(r.OutWriter, result)
		return err
	}
	return os.WriteFile(r.OutputFile, []byte(result), 0644)
}

func Inject(yaml []byte, reader io.Reader, formatter *conf.FormatterConfig, sort *conf.SortConfig) (string, error) {
	if regexp.MustCompile(ActionRegex).Match(yaml) {
		return action.Inject(yaml, reader, formatter, sort)
	} else if regexp.MustCompile(WorkflowRegex).Match(yaml) {
		return workflow.Inject(yaml, reader, formatter, sort)
	}
	return "", fmt.Errorf("not found parser: invalid YAML file")
}
