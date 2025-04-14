package cli

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/tmknom/actdocs/internal/conf"
	"github.com/tmknom/actdocs/internal/format"
	"github.com/tmknom/actdocs/internal/parse"
	"github.com/tmknom/actdocs/internal/read"
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
	log.Printf("read: %s", r.source)

	factory := &parse.ParserFactory{Raw: yaml}
	parser, err := factory.Factory(r.FormatterConfig, r.SortConfig)
	if err != nil {
		return err
	}
	log.Printf("selected parser: %T", parser)

	content, err := parser.ParseAST(yaml)
	if err != nil {
		return err
	}

	formatted := ""
	switch content.(type) {
	case *parse.ActionAST:
		formatter := format.NewActionFormatter(r.FormatterConfig)
		formatted = formatter.Format(content.(*parse.ActionAST))
	case *parse.WorkflowAST:
		formatter := format.NewWorkflowFormatter(r.FormatterConfig)
		formatted = formatter.Format(content.(*parse.WorkflowAST))
	default:
		return fmt.Errorf("unsupported AST type: %T", content)
	}

	_, err = fmt.Fprintln(r.OutWriter, formatted)
	return err
}
