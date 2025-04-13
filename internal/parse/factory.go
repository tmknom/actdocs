package parse

import (
	"fmt"
	"regexp"

	"github.com/tmknom/actdocs/internal/conf"
)

type YamlParser interface {
	Parse(yamlBytes []byte) (string, error)
}

type ParserFactory struct {
	Raw []byte
}

func (f ParserFactory) Factory(config *conf.FormatterConfig, sort *conf.SortConfig) (YamlParser, error) {
	if f.isReusableWorkflow() {
		return NewWorkflowParser(config, sort), nil
	} else if f.isCustomActions() {
		return NewActionParser(config, sort), nil
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
