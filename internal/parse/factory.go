package parse

import (
	"fmt"
	"regexp"

	"github.com/tmknom/actdocs/internal/config"
)

type YamlParser interface {
	Parse() (string, error)
}

type ParserFactory struct {
	Raw []byte
}

func (f ParserFactory) Factory(globalConfig *config.GlobalConfig) (YamlParser, error) {
	if f.isReusableWorkflow() {
		return NewWorkflow(f.Raw, globalConfig), nil
	} else if f.isCustomActions() {
		return NewAction(f.Raw, globalConfig), nil
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
