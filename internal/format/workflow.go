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
		sb.WriteString(f.toInputsMarkdown(workflowSpec))
		sb.WriteString("\n\n")
	}

	if len(workflowSpec.Secrets) != 0 || !config.Omit {
		sb.WriteString(f.toSecretsMarkdown(workflowSpec))
		sb.WriteString("\n\n")
	}

	if len(workflowSpec.Outputs) != 0 || !config.Omit {
		sb.WriteString(f.toOutputsMarkdown(workflowSpec))
		sb.WriteString("\n\n")
	}

	if len(workflowSpec.Permissions) != 0 || !config.Omit {
		sb.WriteString(f.toPermissionsMarkdown(workflowSpec))
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

func (f *WorkflowFormatter) toInputsMarkdown(workflowSpec *WorkflowSpec) string {
	return workflowSpec.toInputsMarkdown()
}

func (f *WorkflowFormatter) toSecretsMarkdown(workflowSpec *WorkflowSpec) string {
	return workflowSpec.toSecretsMarkdown()
}

func (f *WorkflowFormatter) toOutputsMarkdown(workflowSpec *WorkflowSpec) string {
	return workflowSpec.toOutputsMarkdown()
}

func (f *WorkflowFormatter) toPermissionsMarkdown(workflowSpec *WorkflowSpec) string {
	return workflowSpec.toPermissionsMarkdown()
}

const WorkflowInputsTitle = "## Inputs"
const WorkflowInputsColumnTitle = "| Name | Description | Type | Default | Required |"
const WorkflowInputsColumnSeparator = "| :--- | :---------- | :--- | :------ | :------: |"

const WorkflowSecretsTitle = "## Secrets"
const WorkflowSecretsColumnTitle = "| Name | Description | Required |"
const WorkflowSecretsColumnSeparator = "| :--- | :---------- | :------: |"

const WorkflowOutputsTitle = "## Outputs"
const WorkflowOutputsColumnTitle = "| Name | Description |"
const WorkflowOutputsColumnSeparator = "| :--- | :---------- |"

const WorkflowPermissionsTitle = "## Permissions"
const WorkflowPermissionsColumnTitle = "| Scope | Access |"
const WorkflowPermissionsColumnSeparator = "| :--- | :---- |"
