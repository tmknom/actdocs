package cli

import (
	"fmt"
	"log"
	"regexp"

	"github.com/tmknom/actdocs/internal/action"
	"github.com/tmknom/actdocs/internal/conf"
	"github.com/tmknom/actdocs/internal/read"
	"github.com/tmknom/actdocs/internal/workflow"
)

func Orchestrate(source string, formatter *conf.FormatterConfig, sort *conf.SortConfig) (string, error) {
	reader := &read.SourceReader{}
	yaml, err := reader.Read(source)
	if err != nil {
		return "", err
	}
	log.Printf("read: %s", source)

	if regexp.MustCompile(ActionRegex).Match(yaml) {
		spec, err := action.Orchestrate(yaml, sort)
		if err != nil {
			return "", err
		}
		formatted := action.NewFormatter(formatter).Format(spec)
		return formatted, nil
	} else if regexp.MustCompile(WorkflowRegex).Match(yaml) {
		spec, err := workflow.Orchestrate(yaml, sort)
		if err != nil {
			return "", err
		}
		formatted := workflow.NewFormatter(formatter).Format(spec)
		return formatted, nil
	}
	return "", fmt.Errorf("not found parser: invalid YAML file")
}

const (
	ActionRegex   = `(?m)^[\s]*runs:`
	WorkflowRegex = `(?m)^[\s]*workflow_call:`
)
