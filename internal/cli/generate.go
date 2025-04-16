package cli

import (
	"fmt"
	"log"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/tmknom/actdocs/internal/action"
	"github.com/tmknom/actdocs/internal/conf"
	"github.com/tmknom/actdocs/internal/read"
	"github.com/tmknom/actdocs/internal/workflow"
)

func NewGenerateCommand(formatterConfig *conf.FormatterConfig, sortConfig *conf.SortConfig, io *IO) *cobra.Command {
	option := &GenerateOption{IO: io}
	return &cobra.Command{
		Use:   "generate",
		Short: "Generate documentation",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.SetPrefix(fmt.Sprintf("[%s] [%s] ", AppName, cmd.Name()))
			if len(args) > 0 {
				runner := NewGenerateRunner(args[0], formatterConfig, sortConfig, option)
				return runner.Run()
			}
			return cmd.Usage()
		},
	}
}

type GenerateRunner struct {
	source string
	*conf.FormatterConfig
	*conf.SortConfig
	*GenerateOption
}

func NewGenerateRunner(source string, formatter *conf.FormatterConfig, sort *conf.SortConfig, option *GenerateOption) *GenerateRunner {
	return &GenerateRunner{
		source:          source,
		FormatterConfig: formatter,
		SortConfig:      sort,
		GenerateOption:  option,
	}
}

type GenerateOption struct {
	*IO
}

func (r *GenerateRunner) Run() error {
	reader := &read.SourceReader{}
	yaml, err := reader.Read(r.source)
	if err != nil {
		return err
	}

	formatted, err := Generate(yaml, r.FormatterConfig, r.SortConfig)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(r.OutWriter, formatted)
	return err
}

func Generate(yaml []byte, formatter *conf.FormatterConfig, sort *conf.SortConfig) (string, error) {
	if regexp.MustCompile(ActionRegex).Match(yaml) {
		return action.Generate(yaml, formatter, sort)
	} else if regexp.MustCompile(WorkflowRegex).Match(yaml) {
		return workflow.Generate(yaml, formatter, sort)
	}
	return "", fmt.Errorf("not found parser: invalid YAML file")
}
