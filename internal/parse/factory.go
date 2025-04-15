package parse

import (
	"fmt"
	"regexp"

	"github.com/tmknom/actdocs/internal/util"

	"github.com/tmknom/actdocs/internal/conf"
)

type ParserFactory struct {
	Raw []byte
}

func (f ParserFactory) Factory(config *conf.FormatterConfig, sort *conf.SortConfig) (util.YamlParser, error) {
	if f.isReusableWorkflow() {
		return NewWorkflowParser(sort), nil
	} else if f.isCustomActions() {
		return NewActionParser(sort), nil
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
