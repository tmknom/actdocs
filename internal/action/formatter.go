package action

import (
	"encoding/json"
	"strings"

	"github.com/tmknom/actdocs/internal/conf"
	"github.com/tmknom/actdocs/internal/parse"
)

type Formatter struct {
	config *conf.FormatterConfig
	*Spec
}

func NewActionFormatter(config *conf.FormatterConfig) *Formatter {
	return &Formatter{
		config: config,
	}
}

func (f *Formatter) Format(ast *parse.ActionAST) string {
	f.Spec = ConvertActionSpec(ast)

	if f.config.IsJson() {
		return f.ToJson(f.Spec)
	}
	return f.ToMarkdown(f.Spec, f.config)
}

func (f *Formatter) ToJson(actionSpec *Spec) string {
	bytes, err := json.MarshalIndent(actionSpec, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(bytes)
}

func (f *Formatter) ToMarkdown(actionSpec *Spec, config *conf.FormatterConfig) string {
	var sb strings.Builder
	if actionSpec.Description.IsValid() || !config.Omit {
		sb.WriteString(actionSpec.toDescriptionMarkdown())
		sb.WriteString("\n\n")
	}

	if len(actionSpec.Inputs) != 0 || !config.Omit {
		sb.WriteString(actionSpec.toInputsMarkdown())
		sb.WriteString("\n\n")
	}

	if len(actionSpec.Outputs) != 0 || !config.Omit {
		sb.WriteString(actionSpec.toOutputsMarkdown())
		sb.WriteString("\n\n")
	}
	return strings.TrimSpace(sb.String())
}
