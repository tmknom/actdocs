package workflow

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

func NewWorkflowFormatter(config *conf.FormatterConfig) *Formatter {
	return &Formatter{
		config: config,
	}
}

func (f *Formatter) Format(ast *parse.WorkflowAST) string {
	f.Spec = ConvertWorkflowSpec(ast)
	if f.config.IsJson() {
		return f.ToJson(f.Spec)
	}
	return f.ToMarkdown(f.Spec, f.config)
}

func (f *Formatter) ToJson(workflowSpec *Spec) string {
	bytes, err := json.MarshalIndent(workflowSpec, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(bytes)
}

func (f *Formatter) ToMarkdown(workflowSpec *Spec, config *conf.FormatterConfig) string {
	var sb strings.Builder
	if len(workflowSpec.Inputs) != 0 || !config.Omit {
		sb.WriteString(workflowSpec.toInputsMarkdown())
		sb.WriteString("\n\n")
	}

	if len(workflowSpec.Secrets) != 0 || !config.Omit {
		sb.WriteString(workflowSpec.toSecretsMarkdown())
		sb.WriteString("\n\n")
	}

	if len(workflowSpec.Outputs) != 0 || !config.Omit {
		sb.WriteString(workflowSpec.toOutputsMarkdown())
		sb.WriteString("\n\n")
	}

	if len(workflowSpec.Permissions) != 0 || !config.Omit {
		sb.WriteString(workflowSpec.toPermissionsMarkdown())
		sb.WriteString("\n\n")
	}
	return strings.TrimSpace(sb.String())
}
