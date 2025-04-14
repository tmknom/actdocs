package parse

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tmknom/actdocs/internal/conf"
	"github.com/tmknom/actdocs/internal/util"
)

type WorkflowFormatter struct {
	config *conf.FormatterConfig
	*WorkflowMarkdown
	*WorkflowJson
}

func NewWorkflowFormatter(config *conf.FormatterConfig) *WorkflowFormatter {
	return &WorkflowFormatter{
		config: config,
	}
}

func (f *WorkflowFormatter) Format(ast *WorkflowAST) string {
	f.WorkflowJson = f.convertWorkflowJson(ast)
	f.WorkflowMarkdown = f.convertWorkflowMarkdown(ast)
	if f.config.IsJson() {
		return f.ToJson(f.WorkflowJson)
	}
	return f.ToMarkdown(f.WorkflowMarkdown, f.config)
}

func (f *WorkflowFormatter) ToJson(workflowJson *WorkflowJson) string {
	bytes, err := json.MarshalIndent(workflowJson, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(bytes)
}

func (f *WorkflowFormatter) ToMarkdown(workflowMarkdown *WorkflowMarkdown, config *conf.FormatterConfig) string {
	var sb strings.Builder
	if f.hasInputs() || !config.Omit {
		sb.WriteString(f.toInputsMarkdown(workflowMarkdown.Inputs))
		sb.WriteString("\n\n")
	}

	if f.hasSecrets() || !config.Omit {
		sb.WriteString(f.toSecretsMarkdown(workflowMarkdown.Secrets))
		sb.WriteString("\n\n")
	}

	if f.hasOutputs() || !config.Omit {
		sb.WriteString(f.toOutputsMarkdown(workflowMarkdown.Outputs))
		sb.WriteString("\n\n")
	}

	if f.hasPermissions() || !config.Omit {
		sb.WriteString(f.toPermissionsMarkdown(workflowMarkdown.Permissions))
		sb.WriteString("\n\n")
	}
	return strings.TrimSpace(sb.String())
}

func (f *WorkflowFormatter) convertWorkflowJson(ast *WorkflowAST) *WorkflowJson {
	inputs := []*WorkflowInputJson{}
	for _, inputAst := range ast.Inputs {
		input := &WorkflowInputJson{
			Name:        inputAst.Name,
			Default:     inputAst.Default,
			Description: inputAst.Description,
			Required:    inputAst.Required,
			Type:        inputAst.Type,
		}
		inputs = append(inputs, input)
	}

	secrets := []*WorkflowSecretJson{}
	for _, secretAst := range ast.Secrets {
		secret := &WorkflowSecretJson{
			Name:        secretAst.Name,
			Description: secretAst.Description,
			Required:    secretAst.Required,
		}
		secrets = append(secrets, secret)
	}

	outputs := []*WorkflowOutputJson{}
	for _, outputAst := range ast.Outputs {
		output := &WorkflowOutputJson{
			Name:        outputAst.Name,
			Description: outputAst.Description,
		}
		outputs = append(outputs, output)
	}

	permissions := []*WorkflowPermissionJson{}
	for _, permissionAst := range ast.Permissions {
		permission := &WorkflowPermissionJson{
			Scope:  permissionAst.Scope,
			Access: permissionAst.Access,
		}
		permissions = append(permissions, permission)
	}

	return &WorkflowJson{Inputs: inputs, Secrets: secrets, Outputs: outputs, Permissions: permissions}
}

func (f *WorkflowFormatter) convertWorkflowMarkdown(ast *WorkflowAST) *WorkflowMarkdown {
	inputs := []*WorkflowInputMarkdown{}
	for _, inputAst := range ast.Inputs {
		input := &WorkflowInputMarkdown{
			Name:        inputAst.Name,
			Default:     inputAst.Default,
			Description: inputAst.Description,
			Required:    inputAst.Required,
			Type:        inputAst.Type,
		}
		inputs = append(inputs, input)
	}

	secrets := []*WorkflowSecretMarkdown{}
	for _, secretAst := range ast.Secrets {
		secret := &WorkflowSecretMarkdown{
			Name:        secretAst.Name,
			Description: secretAst.Description,
			Required:    secretAst.Required,
		}
		secrets = append(secrets, secret)
	}

	outputs := []*WorkflowOutputMarkdown{}
	for _, outputAst := range ast.Outputs {
		output := &WorkflowOutputMarkdown{
			Name:        outputAst.Name,
			Description: outputAst.Description,
		}
		outputs = append(outputs, output)
	}

	permissions := []*WorkflowPermissionMarkdown{}
	for _, permissionAst := range ast.Permissions {
		permission := &WorkflowPermissionMarkdown{
			Scope:  permissionAst.Scope,
			Access: permissionAst.Access,
		}
		permissions = append(permissions, permission)
	}

	return &WorkflowMarkdown{Inputs: inputs, Secrets: secrets, Outputs: outputs, Permissions: permissions}
}

func (f *WorkflowFormatter) toInputsMarkdown(inputs []*WorkflowInputMarkdown) string {
	var sb strings.Builder
	sb.WriteString(WorkflowInputsTitle)
	sb.WriteString("\n\n")
	if f.hasInputs() {
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

func (f *WorkflowFormatter) toSecretsMarkdown(secrets []*WorkflowSecretMarkdown) string {
	var sb strings.Builder
	sb.WriteString(WorkflowSecretsTitle)
	sb.WriteString("\n\n")
	if f.hasSecrets() {
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

func (f *WorkflowFormatter) toOutputsMarkdown(outputs []*WorkflowOutputMarkdown) string {
	var sb strings.Builder
	sb.WriteString(WorkflowOutputsTitle)
	sb.WriteString("\n\n")
	if f.hasOutputs() {
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

func (f *WorkflowFormatter) toPermissionsMarkdown(permissions []*WorkflowPermissionMarkdown) string {
	var sb strings.Builder
	sb.WriteString(WorkflowPermissionsTitle)
	sb.WriteString("\n\n")
	if f.hasPermissions() {
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

func (f *WorkflowFormatter) hasInputs() bool {
	return len(f.WorkflowMarkdown.Inputs) != 0
}

func (f *WorkflowFormatter) hasSecrets() bool {
	return len(f.WorkflowMarkdown.Secrets) != 0
}

func (f *WorkflowFormatter) hasOutputs() bool {
	return len(f.WorkflowMarkdown.Outputs) != 0
}

func (f *WorkflowFormatter) hasPermissions() bool {
	return len(f.WorkflowMarkdown.Permissions) != 0
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

type WorkflowJson struct {
	Inputs      []*WorkflowInputJson      `json:"inputs"`
	Outputs     []*WorkflowOutputJson     `json:"outputs"`
	Secrets     []*WorkflowSecretJson     `json:"secrets"`
	Permissions []*WorkflowPermissionJson `json:"permissions"`
}

type WorkflowInputJson struct {
	Name        string           `json:"name"`
	Default     *util.NullString `json:"default"`
	Description *util.NullString `json:"description"`
	Required    *util.NullString `json:"required"`
	Type        *util.NullString `json:"type"`
}

type WorkflowOutputJson struct {
	Name        string           `json:"name"`
	Description *util.NullString `json:"description"`
}

type WorkflowSecretJson struct {
	Name        string           `json:"name"`
	Description *util.NullString `json:"description"`
	Required    *util.NullString `json:"required"`
}

type WorkflowPermissionJson struct {
	Scope  string `json:"scope"`
	Access string `json:"access"`
}

type WorkflowMarkdown struct {
	Inputs      []*WorkflowInputMarkdown
	Secrets     []*WorkflowSecretMarkdown
	Outputs     []*WorkflowOutputMarkdown
	Permissions []*WorkflowPermissionMarkdown
}

type WorkflowInputMarkdown struct {
	Name        string
	Default     *util.NullString
	Description *util.NullString
	Required    *util.NullString
	Type        *util.NullString
}

func (i *WorkflowInputMarkdown) toMarkdown() string {
	str := util.TableSeparator
	str += fmt.Sprintf(" %s %s", i.Name, util.TableSeparator)
	str += fmt.Sprintf(" %s %s", i.Description.StringOrEmpty(), util.TableSeparator)
	str += fmt.Sprintf(" %s %s", i.Type.QuoteStringOrLowerNA(), util.TableSeparator)
	str += fmt.Sprintf(" %s %s", i.Default.QuoteStringOrLowerNA(), util.TableSeparator)
	str += fmt.Sprintf(" %s %s", i.Required.YesOrNo(), util.TableSeparator)
	return str
}

type WorkflowSecretMarkdown struct {
	Name        string
	Description *util.NullString
	Required    *util.NullString
}

func (i *WorkflowSecretMarkdown) toMarkdown() string {
	str := util.TableSeparator
	str += fmt.Sprintf(" %s %s", i.Name, util.TableSeparator)
	str += fmt.Sprintf(" %s %s", i.Description.StringOrEmpty(), util.TableSeparator)
	str += fmt.Sprintf(" %s %s", i.Required.YesOrNo(), util.TableSeparator)
	return str
}

type WorkflowOutputMarkdown struct {
	Name        string
	Description *util.NullString
}

func (i *WorkflowOutputMarkdown) toMarkdown() string {
	str := util.TableSeparator
	str += fmt.Sprintf(" %s %s", i.Name, util.TableSeparator)
	str += fmt.Sprintf(" %s %s", i.Description.StringOrEmpty(), util.TableSeparator)
	return str
}

type WorkflowPermissionMarkdown struct {
	Scope  string
	Access string
}

func (i *WorkflowPermissionMarkdown) toMarkdown() string {
	str := util.TableSeparator
	str += fmt.Sprintf(" %s %s", i.Scope, util.TableSeparator)
	str += fmt.Sprintf(" %s %s", i.Access, util.TableSeparator)
	return str
}
