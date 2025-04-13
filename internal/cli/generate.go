package cli

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/tmknom/actdocs/internal/format"
	"github.com/tmknom/actdocs/internal/parse"
	"github.com/tmknom/actdocs/internal/read"
)

func NewGenerateCommand(formatterConfig *format.FormatterConfig, io *IO) *cobra.Command {
	option := &GenerateOption{IO: io}
	return &cobra.Command{
		Use:   "generate",
		Short: "Generate documentation",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.SetPrefix(fmt.Sprintf("[%s] [%s] ", AppName, cmd.Name()))
			if len(args) > 0 {
				runner := NewGenerateRunner(args[0], formatterConfig, option)
				return runner.Run()
			}
			return cmd.Usage()
		},
	}
}

type GenerateRunner struct {
	source string
	*format.FormatterConfig
	*GenerateOption
}

func NewGenerateRunner(source string, config *format.FormatterConfig, option *GenerateOption) *GenerateRunner {
	return &GenerateRunner{
		source:          source,
		FormatterConfig: config,
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
	log.Printf("read: %s", r.source)

	factory := &parse.ParserFactory{Raw: yaml}
	parser, err := factory.Factory(r.FormatterConfig)
	if err != nil {
		return err
	}
	log.Printf("selected parser: %T", parser)

	content, err := parser.Parse(yaml)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(r.OutWriter, content)
	return err
}
