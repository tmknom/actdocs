package cli

import (
	"fmt"
	"log"

	"github.com/tmknom/actdocs/internal/action"
	"github.com/tmknom/actdocs/internal/conf"
	"github.com/tmknom/actdocs/internal/format"
	"github.com/tmknom/actdocs/internal/parse"
	"github.com/tmknom/actdocs/internal/read"
)

func Orchestrate(source string, formatter *conf.FormatterConfig, sort *conf.SortConfig) (string, error) {
	reader := &read.SourceReader{}
	yaml, err := reader.Read(source)
	if err != nil {
		return "", err
	}
	log.Printf("read: %s", source)

	factory := &parse.ParserFactory{Raw: yaml}
	parser, err := factory.Factory(formatter, sort)
	if err != nil {
		return "", err
	}
	log.Printf("selected parser: %T", parser)

	content, err := parser.ParseAST(yaml)
	if err != nil {
		return "", err
	}

	formatted := ""
	switch content.(type) {
	case *parse.ActionAST:
		formatter := action.NewActionFormatter(formatter)
		formatted = formatter.Format(content.(*parse.ActionAST))
	case *parse.WorkflowAST:
		formatter := format.NewWorkflowFormatter(formatter)
		formatted = formatter.Format(content.(*parse.WorkflowAST))
	default:
		return "", fmt.Errorf("unsupported AST type: %T", content)
	}
	return formatted, nil
}
