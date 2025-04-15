package format

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tmknom/actdocs/internal/conf"
	"github.com/tmknom/actdocs/internal/parse"
	"github.com/tmknom/actdocs/internal/util"
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
	f.WorkflowSpec = f.convertWorkflowMarkdown(ast)
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
		sb.WriteString(f.toInputsMarkdown(workflowSpec.Inputs))
		sb.WriteString("\n\n")
	}

	if len(workflowSpec.Secrets) != 0 || !config.Omit {
		sb.WriteString(f.toSecretsMarkdown(workflowSpec.Secrets))
		sb.WriteString("\n\n")
	}

	if len(workflowSpec.Outputs) != 0 || !config.Omit {
		sb.WriteString(f.toOutputsMarkdown(workflowSpec.Outputs))
		sb.WriteString("\n\n")
	}

	if len(workflowSpec.Permissions) != 0 || !config.Omit {
		sb.WriteString(f.toPermissionsMarkdown(workflowSpec.Permissions))
		sb.WriteString("\n\n")
	}
	return strings.TrimSpace(sb.String())
}

func (f *WorkflowFormatter) convertWorkflowMarkdown(ast *parse.WorkflowAST) *WorkflowSpec {
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

func (f *WorkflowFormatter) toInputsMarkdown(inputs []*WorkflowInputSpec) string {
	var sb strings.Builder
	sb.WriteString(WorkflowInputsTitle)
	sb.WriteString("\n\n")
	if len(inputs) != 0 {
		sb.WriteString(WorkflowInputsColumnTitle)
		sb.WriteString("\n")
		sb.WriteString(WorkflowInputsColumnSeparator)
		sb.WriteString("\n")
		for _, input := range inputs {
			sb.WriteString(input.toMarkdown())
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString(util.UpperNAString)
	}
	return strings.TrimSpace(sb.String())
}

func (f *WorkflowFormatter) toSecretsMarkdown(secrets []*WorkflowSecretSpec) string {
	var sb strings.Builder
	sb.WriteString(WorkflowSecretsTitle)
	sb.WriteString("\n\n")
	if len(secrets) != 0 {
		sb.WriteString(WorkflowSecretsColumnTitle)
		sb.WriteString("\n")
		sb.WriteString(WorkflowSecretsColumnSeparator)
		sb.WriteString("\n")
		for _, secret := range secrets {
			sb.WriteString(secret.toMarkdown())
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString(util.UpperNAString)
	}
	return strings.TrimSpace(sb.String())
}

func (f *WorkflowFormatter) toOutputsMarkdown(outputs []*WorkflowOutputSpec) string {
	var sb strings.Builder
	sb.WriteString(WorkflowOutputsTitle)
	sb.WriteString("\n\n")
	if len(outputs) != 0 {
		sb.WriteString(WorkflowOutputsColumnTitle)
		sb.WriteString("\n")
		sb.WriteString(WorkflowOutputsColumnSeparator)
		sb.WriteString("\n")
		for _, output := range outputs {
			sb.WriteString(output.toMarkdown())
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString(util.UpperNAString)
	}
	return strings.TrimSpace(sb.String())
}

func (f *WorkflowFormatter) toPermissionsMarkdown(permissions []*WorkflowPermissionSpec) string {
	var sb strings.Builder
	sb.WriteString(WorkflowPermissionsTitle)
	sb.WriteString("\n\n")
	if len(permissions) != 0 {
		sb.WriteString(WorkflowPermissionsColumnTitle)
		sb.WriteString("\n")
		sb.WriteString(WorkflowPermissionsColumnSeparator)
		sb.WriteString("\n")
		for _, permission := range permissions {
			sb.WriteString(permission.toMarkdown())
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString(util.UpperNAString)
	}
	return strings.TrimSpace(sb.String())
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

type WorkflowSpec struct {
	Inputs      []*WorkflowInputSpec      `json:"inputs"`
	Secrets     []*WorkflowSecretSpec     `json:"secrets"`
	Outputs     []*WorkflowOutputSpec     `json:"outputs"`
	Permissions []*WorkflowPermissionSpec `json:"permissions"`
}

type WorkflowInputSpec struct {
	Name        string           `json:"name"`
	Default     *util.NullString `json:"default"`
	Description *util.NullString `json:"description"`
	Required    *util.NullString `json:"required"`
	Type        *util.NullString `json:"type"`
}

func (s *WorkflowInputSpec) toMarkdown() string {
	str := util.TableSeparator
	str += fmt.Sprintf(" %s %s", s.Name, util.TableSeparator)
	str += fmt.Sprintf(" %s %s", s.Description.StringOrEmpty(), util.TableSeparator)
	str += fmt.Sprintf(" %s %s", s.Type.QuoteStringOrLowerNA(), util.TableSeparator)
	str += fmt.Sprintf(" %s %s", s.Default.QuoteStringOrLowerNA(), util.TableSeparator)
	str += fmt.Sprintf(" %s %s", s.Required.YesOrNo(), util.TableSeparator)
	return str
}

type WorkflowSecretSpec struct {
	Name        string           `json:"name"`
	Description *util.NullString `json:"description"`
	Required    *util.NullString `json:"required"`
}

func (s *WorkflowSecretSpec) toMarkdown() string {
	str := util.TableSeparator
	str += fmt.Sprintf(" %s %s", s.Name, util.TableSeparator)
	str += fmt.Sprintf(" %s %s", s.Description.StringOrEmpty(), util.TableSeparator)
	str += fmt.Sprintf(" %s %s", s.Required.YesOrNo(), util.TableSeparator)
	return str
}

type WorkflowOutputSpec struct {
	Name        string           `json:"name"`
	Description *util.NullString `json:"description"`
}

func (s *WorkflowOutputSpec) toMarkdown() string {
	str := util.TableSeparator
	str += fmt.Sprintf(" %s %s", s.Name, util.TableSeparator)
	str += fmt.Sprintf(" %s %s", s.Description.StringOrEmpty(), util.TableSeparator)
	return str
}

type WorkflowPermissionSpec struct {
	Scope  string `json:"scope"`
	Access string `json:"access"`
}

func (s *WorkflowPermissionSpec) toMarkdown() string {
	str := util.TableSeparator
	str += fmt.Sprintf(" %s %s", s.Scope, util.TableSeparator)
	str += fmt.Sprintf(" %s %s", s.Access, util.TableSeparator)
	return str
}
