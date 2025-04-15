package format

import (
	"encoding/json"
	"strings"

	"github.com/tmknom/actdocs/internal/conf"
	"github.com/tmknom/actdocs/internal/parse"
)

type ActionFormatter struct {
	config *conf.FormatterConfig
	*ActionSpec
}

func NewActionFormatter(config *conf.FormatterConfig) *ActionFormatter {
	return &ActionFormatter{
		config: config,
	}
}

func (f *ActionFormatter) Format(ast *parse.ActionAST) string {
	f.ActionSpec = ConvertActionSpec(ast)

	if f.config.IsJson() {
		return f.ToJson(f.ActionSpec)
	}
	return f.ToMarkdown(f.ActionSpec, f.config)
}

func (f *ActionFormatter) ToJson(actionSpec *ActionSpec) string {
	bytes, err := json.MarshalIndent(actionSpec, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(bytes)
}

func (f *ActionFormatter) ToMarkdown(actionSpec *ActionSpec, config *conf.FormatterConfig) string {
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
