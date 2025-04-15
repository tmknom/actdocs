package cli

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/tmknom/actdocs/internal/conf"
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
	formatted, err := Orchestrate(r.source, r.FormatterConfig, r.SortConfig)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(r.OutWriter, formatted)
	return err
}
