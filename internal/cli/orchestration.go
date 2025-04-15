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
		return action.Orchestrate(yaml, formatter, sort)
	} else if regexp.MustCompile(WorkflowRegex).Match(yaml) {
		return workflow.Orchestrate(yaml, formatter, sort)
	}
	return "", fmt.Errorf("not found parser: invalid YAML file")
}

const (
	ActionRegex   = `(?m)^[\s]*runs:`
	WorkflowRegex = `(?m)^[\s]*workflow_call:`
)
