package parse

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/tmknom/actdocs/internal/conf"
	"github.com/tmknom/actdocs/internal/util"
	"gopkg.in/yaml.v2"
)

type WorkflowParser struct {
	*WorkflowAST
	config *conf.FormatterConfig
	*conf.SortConfig
}

func NewWorkflowParser(config *conf.FormatterConfig, sort *conf.SortConfig) *WorkflowParser {
	return &WorkflowParser{
		WorkflowAST: &WorkflowAST{
			Inputs:      []*WorkflowInput{},
			Secrets:     []*WorkflowSecret{},
			Outputs:     []*WorkflowOutput{},
			Permissions: []*WorkflowPermission{},
		},
		config:     config,
		SortConfig: sort,
	}
}

type WorkflowAST struct {
	Inputs      []*WorkflowInput
	Secrets     []*WorkflowSecret
	Outputs     []*WorkflowOutput
	Permissions []*WorkflowPermission
}

func (p *WorkflowParser) Parse(yamlBytes []byte) (string, error) {
	log.Printf("config: %#v", p.config)

	content := &WorkflowYaml{}
	err := yaml.Unmarshal(yamlBytes, content)
	if err != nil {
		return "", err
	}

	for name, value := range content.inputs() {
		input := p.parseInput(name, value)
		p.Inputs = append(p.Inputs, input)
	}

	for name, value := range content.outputs() {
		output := p.parseOutput(name, value)
		p.Outputs = append(p.Outputs, output)
	}

	for name, value := range content.secrets() {
		secret := p.parseSecret(name, value)
		p.Secrets = append(p.Secrets, secret)
	}

	for scope, access := range content.permissions() {
		permission := NewWorkflowPermission(scope.(string), access.(string))
		p.Permissions = append(p.Permissions, permission)
	}

	p.sort()

	formatter := NewWorkflowFormatter(p.WorkflowAST, p.config)
	return formatter.Format(), nil
}

func (p *WorkflowParser) sort() {
	switch {
	case p.SortConfig.Sort:
		p.sortInputs()
		p.sortSecrets()
		p.sortOutputsByName()
		p.sortPermissionsByScope()
	case p.SortConfig.SortByName:
		p.sortInputsByName()
		p.sortSecretsByName()
		p.sortOutputsByName()
		p.sortPermissionsByScope()
	case p.SortConfig.SortByRequired:
		p.sortInputsByRequired()
		p.sortSecretByRequired()
	}
}

func (p *WorkflowParser) sortInputs() {
	log.Printf("sorted: inputs")

	//goland:noinspection GoPreferNilSlice
	required := []*WorkflowInput{}
	//goland:noinspection GoPreferNilSlice
	notRequired := []*WorkflowInput{}
	for _, input := range p.Inputs {
		if input.Required.IsTrue() {
			required = append(required, input)
		} else {
			notRequired = append(notRequired, input)
		}
	}

	sort.Slice(required, func(i, j int) bool {
		return required[i].Name < required[j].Name
	})
	sort.Slice(notRequired, func(i, j int) bool {
		return notRequired[i].Name < notRequired[j].Name
	})
	p.Inputs = append(required, notRequired...)
}

func (p *WorkflowParser) sortInputsByName() {
	log.Printf("sorted: inputs by name")
	item := p.Inputs
	sort.Slice(item, func(i, j int) bool {
		return item[i].Name < item[j].Name
	})
}

func (p *WorkflowParser) sortInputsByRequired() {
	log.Printf("sorted: inputs by required")
	item := p.Inputs
	sort.Slice(item, func(i, j int) bool {
		return item[i].Required.IsTrue()
	})
}

func (p *WorkflowParser) sortSecrets() {
	log.Printf("sorted: secrets")

	//goland:noinspection GoPreferNilSlice
	required := []*WorkflowSecret{}
	//goland:noinspection GoPreferNilSlice
	notRequired := []*WorkflowSecret{}
	for _, input := range p.Secrets {
		if input.Required.IsTrue() {
			required = append(required, input)
		} else {
			notRequired = append(notRequired, input)
		}
	}

	sort.Slice(required, func(i, j int) bool {
		return required[i].Name < required[j].Name
	})
	sort.Slice(notRequired, func(i, j int) bool {
		return notRequired[i].Name < notRequired[j].Name
	})
	p.Secrets = append(required, notRequired...)
}

func (p *WorkflowParser) sortSecretsByName() {
	log.Printf("sorted: secrets by name")
	item := p.Secrets
	sort.Slice(item, func(i, j int) bool {
		return item[i].Name < item[j].Name
	})
}

func (p *WorkflowParser) sortSecretByRequired() {
	log.Printf("sorted: secrets by required")
	item := p.Secrets
	sort.Slice(item, func(i, j int) bool {
		return item[i].Required.IsTrue()
	})
}

func (p *WorkflowParser) sortOutputsByName() {
	log.Printf("sorted: outputs by name")
	item := p.Outputs
	sort.Slice(item, func(i, j int) bool {
		return item[i].Name < item[j].Name
	})
}

func (p *WorkflowParser) sortPermissionsByScope() {
	log.Printf("sorted: permission by scope")
	item := p.Permissions
	sort.Slice(item, func(i, j int) bool {
		return item[i].Scope < item[j].Scope
	})
}

func (p *WorkflowParser) parseInput(name string, value *workflowInputYaml) *WorkflowInput {
	result := NewWorkflowInput(name)
	if value == nil {
		return result
	}

	result.Default = util.NewNullString(value.Default)
	result.Description = util.NewNullString(value.Description)
	result.Required = util.NewNullString(value.Required)
	result.Type = util.NewNullString(value.Type)

	return result
}

func (p *WorkflowParser) parseSecret(name string, value *workflowSecretYaml) *WorkflowSecret {
	result := NewWorkflowSecret(name)
	if value == nil {
		return result
	}

	result.Description = util.NewNullString(value.Description)
	result.Required = util.NewNullString(value.Required)

	return result
}

func (p *WorkflowParser) parseOutput(name string, value *workflowOutputYaml) *WorkflowOutput {
	result := NewWorkflowOutput(name)
	if value == nil {
		return result
	}

	result.Description = util.NewNullString(value.Description)
	return result
}

type WorkflowFormatter struct {
	*WorkflowAST
	config *conf.FormatterConfig
}

func NewWorkflowFormatter(ast *WorkflowAST, config *conf.FormatterConfig) *WorkflowFormatter {
	return &WorkflowFormatter{
		WorkflowAST: ast,
		config:      config,
	}
}

func (f *WorkflowFormatter) Format() string {
	if f.config.IsJson() {
		return f.toJson()
	}
	return f.toMarkdown()
}

func (f *WorkflowFormatter) toJson() string {
	bytes, err := json.MarshalIndent(&WorkflowJson{Inputs: f.Inputs, Secrets: f.Secrets, Outputs: f.Outputs, Permissions: f.Permissions}, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(bytes)
}

func (f *WorkflowFormatter) toMarkdown() string {
	var sb strings.Builder
	if f.hasInputs() || !f.config.Omit {
		sb.WriteString(f.toInputsMarkdown())
		sb.WriteString("\n\n")
	}

	if f.hasSecrets() || !f.config.Omit {
		sb.WriteString(f.toSecretsMarkdown())
		sb.WriteString("\n\n")
	}

	if f.hasOutputs() || !f.config.Omit {
		sb.WriteString(f.toOutputsMarkdown())
		sb.WriteString("\n\n")
	}

	if f.hasPermissions() || !f.config.Omit {
		sb.WriteString(f.toPermissionsMarkdown())
		sb.WriteString("\n\n")
	}
	return strings.TrimSpace(sb.String())
}

func (f *WorkflowFormatter) toInputsMarkdown() string {
	var sb strings.Builder
	sb.WriteString(WorkflowInputsTitle)
	sb.WriteString("\n\n")
	if f.hasInputs() {
		sb.WriteString(WorkflowInputsColumnTitle)
		sb.WriteString("\n")
		sb.WriteString(WorkflowInputsColumnSeparator)
		sb.WriteString("\n")
		for _, input := range f.Inputs {
			sb.WriteString(input.toMarkdown())
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString(util.UpperNAString)
	}
	return strings.TrimSpace(sb.String())
}

func (f *WorkflowFormatter) toSecretsMarkdown() string {
	var sb strings.Builder
	sb.WriteString(WorkflowSecretsTitle)
	sb.WriteString("\n\n")
	if f.hasSecrets() {
		sb.WriteString(WorkflowSecretsColumnTitle)
		sb.WriteString("\n")
		sb.WriteString(WorkflowSecretsColumnSeparator)
		sb.WriteString("\n")
		for _, secret := range f.Secrets {
			sb.WriteString(secret.toMarkdown())
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString(util.UpperNAString)
	}
	return strings.TrimSpace(sb.String())
}

func (f *WorkflowFormatter) toOutputsMarkdown() string {
	var sb strings.Builder
	sb.WriteString(WorkflowOutputsTitle)
	sb.WriteString("\n\n")
	if f.hasOutputs() {
		sb.WriteString(WorkflowOutputsColumnTitle)
		sb.WriteString("\n")
		sb.WriteString(WorkflowOutputsColumnSeparator)
		sb.WriteString("\n")
		for _, output := range f.Outputs {
			sb.WriteString(output.toMarkdown())
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString(util.UpperNAString)
	}
	return strings.TrimSpace(sb.String())
}

func (f *WorkflowFormatter) toPermissionsMarkdown() string {
	var sb strings.Builder
	sb.WriteString(WorkflowPermissionsTitle)
	sb.WriteString("\n\n")
	if f.hasPermissions() {
		sb.WriteString(WorkflowPermissionsColumnTitle)
		sb.WriteString("\n")
		sb.WriteString(WorkflowPermissionsColumnSeparator)
		sb.WriteString("\n")
		for _, permission := range f.Permissions {
			sb.WriteString(permission.toMarkdown())
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString(util.UpperNAString)
	}
	return strings.TrimSpace(sb.String())
}

func (f *WorkflowFormatter) hasInputs() bool {
	return len(f.Inputs) != 0
}

func (f *WorkflowFormatter) hasSecrets() bool {
	return len(f.Secrets) != 0
}

func (f *WorkflowFormatter) hasOutputs() bool {
	return len(f.Outputs) != 0
}

func (f *WorkflowFormatter) hasPermissions() bool {
	return len(f.Permissions) != 0
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
	Inputs      []*WorkflowInput      `json:"inputs"`
	Outputs     []*WorkflowOutput     `json:"outputs"`
	Secrets     []*WorkflowSecret     `json:"secrets"`
	Permissions []*WorkflowPermission `json:"permissions"`
}

type WorkflowInput struct {
	Name        string           `json:"name"`
	Default     *util.NullString `json:"default"`
	Description *util.NullString `json:"description"`
	Required    *util.NullString `json:"required"`
	Type        *util.NullString `json:"type"`
}

func NewWorkflowInput(name string) *WorkflowInput {
	return &WorkflowInput{
		Name:        name,
		Default:     util.DefaultNullString,
		Description: util.DefaultNullString,
		Required:    util.DefaultNullString,
		Type:        util.DefaultNullString,
	}
}

func (i *WorkflowInput) toMarkdown() string {
	str := util.TableSeparator
	str += fmt.Sprintf(" %s %s", i.Name, util.TableSeparator)
	str += fmt.Sprintf(" %s %s", i.Description.StringOrEmpty(), util.TableSeparator)
	str += fmt.Sprintf(" %s %s", i.Type.QuoteStringOrLowerNA(), util.TableSeparator)
	str += fmt.Sprintf(" %s %s", i.Default.QuoteStringOrLowerNA(), util.TableSeparator)
	str += fmt.Sprintf(" %s %s", i.Required.YesOrNo(), util.TableSeparator)
	return str
}

type WorkflowSecret struct {
	Name        string           `json:"name"`
	Description *util.NullString `json:"description"`
	Required    *util.NullString `json:"required"`
}

func NewWorkflowSecret(name string) *WorkflowSecret {
	return &WorkflowSecret{
		Name:        name,
		Description: util.DefaultNullString,
		Required:    util.DefaultNullString,
	}
}

func (i *WorkflowSecret) toMarkdown() string {
	str := util.TableSeparator
	str += fmt.Sprintf(" %s %s", i.Name, util.TableSeparator)
	str += fmt.Sprintf(" %s %s", i.Description.StringOrEmpty(), util.TableSeparator)
	str += fmt.Sprintf(" %s %s", i.Required.YesOrNo(), util.TableSeparator)
	return str
}

type WorkflowOutput struct {
	Name        string           `json:"name"`
	Description *util.NullString `json:"description"`
}

func NewWorkflowOutput(name string) *WorkflowOutput {
	return &WorkflowOutput{
		Name:        name,
		Description: util.DefaultNullString,
	}
}

func (i *WorkflowOutput) toMarkdown() string {
	str := util.TableSeparator
	str += fmt.Sprintf(" %s %s", i.Name, util.TableSeparator)
	str += fmt.Sprintf(" %s %s", i.Description.StringOrEmpty(), util.TableSeparator)
	return str
}

type WorkflowPermission struct {
	Scope  string `json:"scope"`
	Access string `json:"access"`
}

func NewWorkflowPermission(scope string, access string) *WorkflowPermission {
	return &WorkflowPermission{
		Scope:  scope,
		Access: access,
	}
}

func (i *WorkflowPermission) toMarkdown() string {
	str := util.TableSeparator
	str += fmt.Sprintf(" %s %s", i.Scope, util.TableSeparator)
	str += fmt.Sprintf(" %s %s", i.Access, util.TableSeparator)
	return str
}
