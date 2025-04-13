package parse

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/tmknom/actdocs/internal/format"
	"github.com/tmknom/actdocs/internal/util"
	"gopkg.in/yaml.v2"
)

type WorkflowParser struct {
	Inputs      []*WorkflowInput
	Secrets     []*WorkflowSecret
	Outputs     []*WorkflowOutput
	Permissions []*WorkflowPermission
	config      *format.FormatterConfig
	rawYaml     []byte
}

func NewWorkflowParser(rawYaml []byte, config *format.FormatterConfig) *WorkflowParser {
	return &WorkflowParser{
		Inputs:      []*WorkflowInput{},
		Secrets:     []*WorkflowSecret{},
		Outputs:     []*WorkflowOutput{},
		Permissions: []*WorkflowPermission{},
		config:      config,
		rawYaml:     rawYaml,
	}
}

func (p *WorkflowParser) Parse() (string, error) {
	log.Printf("config: %#v", p.config)

	content := &WorkflowYamlContent{}
	err := yaml.Unmarshal(p.rawYaml, content)
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
	return p.format(), nil
}

func (p *WorkflowParser) sort() {
	switch {
	case p.config.Sort:
		p.sortInputs()
		p.sortSecrets()
		p.sortOutputsByName()
		p.sortPermissionsByScope()
	case p.config.SortByName:
		p.sortInputsByName()
		p.sortSecretsByName()
		p.sortOutputsByName()
		p.sortPermissionsByScope()
	case p.config.SortByRequired:
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

func (p *WorkflowParser) parseInput(name string, value *WorkflowYamlInput) *WorkflowInput {
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

func (p *WorkflowParser) parseSecret(name string, value *WorkflowYamlSecret) *WorkflowSecret {
	result := NewWorkflowSecret(name)
	if value == nil {
		return result
	}

	result.Description = util.NewNullString(value.Description)
	result.Required = util.NewNullString(value.Required)

	return result
}

func (p *WorkflowParser) parseOutput(name string, value *WorkflowYamlOutput) *WorkflowOutput {
	result := NewWorkflowOutput(name)
	if value == nil {
		return result
	}

	result.Description = util.NewNullString(value.Description)
	return result
}

func (p *WorkflowParser) format() string {
	if p.config.IsJson() {
		return p.toJson()
	}
	return p.toMarkdown()
}

func (p *WorkflowParser) toJson() string {
	bytes, err := json.MarshalIndent(&WorkflowJson{Inputs: p.Inputs, Secrets: p.Secrets, Outputs: p.Outputs, Permissions: p.Permissions}, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(bytes)
}

func (p *WorkflowParser) toMarkdown() string {
	var sb strings.Builder
	if p.hasInputs() || !p.config.Omit {
		sb.WriteString(p.toInputsMarkdown())
		sb.WriteString("\n\n")
	}

	if p.hasSecrets() || !p.config.Omit {
		sb.WriteString(p.toSecretsMarkdown())
		sb.WriteString("\n\n")
	}

	if p.hasOutputs() || !p.config.Omit {
		sb.WriteString(p.toOutputsMarkdown())
		sb.WriteString("\n\n")
	}

	if p.hasPermissions() || !p.config.Omit {
		sb.WriteString(p.toPermissionsMarkdown())
		sb.WriteString("\n\n")
	}
	return strings.TrimSpace(sb.String())
}

func (p *WorkflowParser) toInputsMarkdown() string {
	var sb strings.Builder
	sb.WriteString(WorkflowInputsTitle)
	sb.WriteString("\n\n")
	if p.hasInputs() {
		sb.WriteString(WorkflowInputsColumnTitle)
		sb.WriteString("\n")
		sb.WriteString(WorkflowInputsColumnSeparator)
		sb.WriteString("\n")
		for _, input := range p.Inputs {
			sb.WriteString(input.toMarkdown())
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString(util.UpperNAString)
	}
	return strings.TrimSpace(sb.String())
}

func (p *WorkflowParser) toSecretsMarkdown() string {
	var sb strings.Builder
	sb.WriteString(WorkflowSecretsTitle)
	sb.WriteString("\n\n")
	if p.hasSecrets() {
		sb.WriteString(WorkflowSecretsColumnTitle)
		sb.WriteString("\n")
		sb.WriteString(WorkflowSecretsColumnSeparator)
		sb.WriteString("\n")
		for _, secret := range p.Secrets {
			sb.WriteString(secret.toMarkdown())
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString(util.UpperNAString)
	}
	return strings.TrimSpace(sb.String())
}

func (p *WorkflowParser) toOutputsMarkdown() string {
	var sb strings.Builder
	sb.WriteString(WorkflowOutputsTitle)
	sb.WriteString("\n\n")
	if p.hasOutputs() {
		sb.WriteString(WorkflowOutputsColumnTitle)
		sb.WriteString("\n")
		sb.WriteString(WorkflowOutputsColumnSeparator)
		sb.WriteString("\n")
		for _, output := range p.Outputs {
			sb.WriteString(output.toMarkdown())
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString(util.UpperNAString)
	}
	return strings.TrimSpace(sb.String())
}

func (p *WorkflowParser) toPermissionsMarkdown() string {
	var sb strings.Builder
	sb.WriteString(WorkflowPermissionsTitle)
	sb.WriteString("\n\n")
	if p.hasPermissions() {
		sb.WriteString(WorkflowPermissionsColumnTitle)
		sb.WriteString("\n")
		sb.WriteString(WorkflowPermissionsColumnSeparator)
		sb.WriteString("\n")
		for _, permission := range p.Permissions {
			sb.WriteString(permission.toMarkdown())
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString(util.UpperNAString)
	}
	return strings.TrimSpace(sb.String())
}

func (p *WorkflowParser) hasInputs() bool {
	return len(p.Inputs) != 0
}

func (p *WorkflowParser) hasSecrets() bool {
	return len(p.Secrets) != 0
}

func (p *WorkflowParser) hasOutputs() bool {
	return len(p.Outputs) != 0
}

func (p *WorkflowParser) hasPermissions() bool {
	return len(p.Permissions) != 0
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

type WorkflowYamlContent struct {
	On          *WorkflowYamlOn `yaml:"on"`
	Permissions interface{}     `yaml:"permissions"`
}

type WorkflowYamlOn struct {
	WorkflowCall *WorkflowYamlWorkflowCall `yaml:"workflow_call"`
}

type WorkflowYamlWorkflowCall struct {
	Inputs  map[string]*WorkflowYamlInput  `yaml:"inputs"`
	Secrets map[string]*WorkflowYamlSecret `yaml:"secrets"`
	Outputs map[string]*WorkflowYamlOutput `yaml:"outputs"`
}

type WorkflowYamlInput struct {
	Default     *string `mapstructure:"default"`
	Description *string `mapstructure:"description"`
	Required    *string `mapstructure:"required"`
	Type        *string `mapstructure:"type"`
}

type WorkflowYamlSecret struct {
	Description *string `mapstructure:"description"`
	Required    *string `mapstructure:"required"`
}

type WorkflowYamlOutput struct {
	Description *string `mapstructure:"description"`
}

func (c *WorkflowYamlContent) inputs() map[string]*WorkflowYamlInput {
	if c.On == nil || c.On.WorkflowCall == nil || c.On.WorkflowCall.Inputs == nil {
		return map[string]*WorkflowYamlInput{}
	}
	return c.On.WorkflowCall.Inputs
}

func (c *WorkflowYamlContent) secrets() map[string]*WorkflowYamlSecret {
	if c.On == nil || c.On.WorkflowCall == nil || c.On.WorkflowCall.Secrets == nil {
		return map[string]*WorkflowYamlSecret{}
	}
	return c.On.WorkflowCall.Secrets
}

func (c *WorkflowYamlContent) outputs() map[string]*WorkflowYamlOutput {
	if c.On == nil || c.On.WorkflowCall == nil || c.On.WorkflowCall.Outputs == nil {
		return map[string]*WorkflowYamlOutput{}
	}
	return c.On.WorkflowCall.Outputs
}

func (c *WorkflowYamlContent) permissions() map[interface{}]interface{} {
	if c.Permissions == nil {
		return map[interface{}]interface{}{}
	}

	switch c.Permissions.(type) {
	case string:
		access := c.Permissions.(string)
		if access == ReadAllAccess || access == WriteAllAccess {
			return map[interface{}]interface{}{AllScope: access}
		}
	case map[interface{}]interface{}:
		return c.Permissions.(map[interface{}]interface{})
	}
	return map[interface{}]interface{}{}
}

const ReadAllAccess = "read-all"
const WriteAllAccess = "write-all"
const AllScope = "-"
