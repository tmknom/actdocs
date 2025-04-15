package format

import (
	"encoding/json"
	"strings"

	"github.com/tmknom/actdocs/internal/conf"
	"github.com/tmknom/actdocs/internal/parse"
)

type WorkflowFormatter struct {
	config *conf.FormatterConfig
	*WorkflowSpec
}

func NewWorkflowFormatter(config *conf.FormatterConfig) *WorkflowFormatter {
	return &WorkflowFormatter{
		config: config,
	}
}

func (f *WorkflowFormatter) Format(ast *parse.WorkflowAST) string {
	f.WorkflowSpec = ConvertWorkflowSpec(ast)
	if f.config.IsJson() {
		return f.ToJson(f.WorkflowSpec)
	}
	return f.ToMarkdown(f.WorkflowSpec, f.config)
}

func (f *WorkflowFormatter) ToJson(workflowSpec *WorkflowSpec) string {
	bytes, err := json.MarshalIndent(workflowSpec, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(bytes)
}

func (f *WorkflowFormatter) ToMarkdown(workflowSpec *WorkflowSpec, config *conf.FormatterConfig) string {
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

func ConvertWorkflowSpec(ast *parse.WorkflowAST) *WorkflowSpec {
	inputs := []*WorkflowInputSpec{}
	for _, inputAst := range ast.Inputs {
		input := &WorkflowInputSpec{
			Name:        inputAst.Name,
			Default:     inputAst.Default,
			Description: inputAst.Description,
			Required:    inputAst.Required,
			Type:        inputAst.Type,
		}
		inputs = append(inputs, input)
	}

	secrets := []*WorkflowSecretSpec{}
	for _, secretAst := range ast.Secrets {
		secret := &WorkflowSecretSpec{
			Name:        secretAst.Name,
			Description: secretAst.Description,
			Required:    secretAst.Required,
		}
		secrets = append(secrets, secret)
	}

	outputs := []*WorkflowOutputSpec{}
	for _, outputAst := range ast.Outputs {
		output := &WorkflowOutputSpec{
			Name:        outputAst.Name,
			Description: outputAst.Description,
		}
		outputs = append(outputs, output)
	}

	permissions := []*WorkflowPermissionSpec{}
	for _, permissionAst := range ast.Permissions {
		permission := &WorkflowPermissionSpec{
			Scope:  permissionAst.Scope,
			Access: permissionAst.Access,
		}
		permissions = append(permissions, permission)
	}

	return &WorkflowSpec{Inputs: inputs, Secrets: secrets, Outputs: outputs, Permissions: permissions}
}
