package cli

import (
	"fmt"
	"log"
	"regexp"

	"github.com/tmknom/actdocs/internal/action"
	"github.com/tmknom/actdocs/internal/conf"
	"github.com/tmknom/actdocs/internal/read"
	"github.com/tmknom/actdocs/internal/util"
	"github.com/tmknom/actdocs/internal/workflow"
)

func Orchestrate(source string, formatter *conf.FormatterConfig, sort *conf.SortConfig) (string, error) {
	reader := &read.SourceReader{}
	yaml, err := reader.Read(source)
	if err != nil {
		return "", err
	}
	log.Printf("read: %s", source)

	factory := &ParserFactory{Raw: yaml}
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
	case *action.ActionAST:
		formatter := action.NewActionFormatter(formatter)
		formatted = formatter.Format(content.(*action.ActionAST))
	case *workflow.WorkflowAST:
		formatter := workflow.NewWorkflowFormatter(formatter)
		formatted = formatter.Format(content.(*workflow.WorkflowAST))
	default:
		return "", fmt.Errorf("unsupported AST type: %T", content)
	}
	return formatted, nil
}

type ParserFactory struct {
	Raw []byte
}

func (f ParserFactory) Factory(config *conf.FormatterConfig, sort *conf.SortConfig) (util.YamlParser, error) {
	if f.isReusableWorkflow() {
		return workflow.NewWorkflowParser(sort), nil
	} else if f.isCustomActions() {
		return action.NewActionParser(sort), nil
	} else {
		return nil, fmt.Errorf("not found parser: invalid YAML file")
	}
}

func (f ParserFactory) isReusableWorkflow() bool {
	r := regexp.MustCompile(`(?m)^[\s]*workflow_call:`)
	return r.Match(f.Raw)
}

func (f ParserFactory) isCustomActions() bool {
	r := regexp.MustCompile(`(?m)^[\s]*runs:`)
	return r.Match(f.Raw)
}
