package parse

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"

	config2 "github.com/tmknom/actdocs/internal/format"
	"github.com/tmknom/actdocs/internal/util"
	"gopkg.in/yaml.v2"
)

type Workflow struct {
	Inputs      []*WorkflowInput
	Secrets     []*WorkflowSecret
	Outputs     []*WorkflowOutput
	Permissions []*WorkflowPermission
	config      *config2.GlobalConfig
	rawYaml     []byte
}

func NewWorkflow(rawYaml []byte, config *config2.GlobalConfig) *Workflow {
	return &Workflow{
		Inputs:      []*WorkflowInput{},
		Secrets:     []*WorkflowSecret{},
		Outputs:     []*WorkflowOutput{},
		Permissions: []*WorkflowPermission{},
		config:      config,
		rawYaml:     rawYaml,
	}
}

func (w *Workflow) Parse() (string, error) {
	log.Printf("config: %#v", w.config)

	content := &WorkflowYamlContent{}
	err := yaml.Unmarshal(w.rawYaml, content)
	if err != nil {
		return "", err
	}

	for name, value := range content.inputs() {
		input := w.parseInput(name, value)
		w.Inputs = append(w.Inputs, input)
	}

	for name, value := range content.outputs() {
		output := w.parseOutput(name, value)
		w.Outputs = append(w.Outputs, output)
	}

	for name, value := range content.secrets() {
		secret := w.parseSecret(name, value)
		w.Secrets = append(w.Secrets, secret)
	}

	for scope, access := range content.permissions() {
		permission := NewWorkflowPermission(scope.(string), access.(string))
		w.Permissions = append(w.Permissions, permission)
	}

	w.sort()
	return w.format(), nil
}

func (w *Workflow) sort() {
	switch {
	case w.config.Sort:
		w.sortInputs()
		w.sortSecrets()
		w.sortOutputsByName()
		w.sortPermissionsByScope()
	case w.config.SortByName:
		w.sortInputsByName()
		w.sortSecretsByName()
		w.sortOutputsByName()
		w.sortPermissionsByScope()
	case w.config.SortByRequired:
		w.sortInputsByRequired()
		w.sortSecretByRequired()
	}
}

func (w *Workflow) sortInputs() {
	log.Printf("sorted: inputs")

	//goland:noinspection GoPreferNilSlice
	required := []*WorkflowInput{}
	//goland:noinspection GoPreferNilSlice
	notRequired := []*WorkflowInput{}
	for _, input := range w.Inputs {
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
	w.Inputs = append(required, notRequired...)
}

func (w *Workflow) sortInputsByName() {
	log.Printf("sorted: inputs by name")
	item := w.Inputs
	sort.Slice(item, func(i, j int) bool {
		return item[i].Name < item[j].Name
	})
}

func (w *Workflow) sortInputsByRequired() {
	log.Printf("sorted: inputs by required")
	item := w.Inputs
	sort.Slice(item, func(i, j int) bool {
		return item[i].Required.IsTrue()
	})
}

func (w *Workflow) sortSecrets() {
	log.Printf("sorted: secrets")

	//goland:noinspection GoPreferNilSlice
	required := []*WorkflowSecret{}
	//goland:noinspection GoPreferNilSlice
	notRequired := []*WorkflowSecret{}
	for _, input := range w.Secrets {
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
	w.Secrets = append(required, notRequired...)
}

func (w *Workflow) sortSecretsByName() {
	log.Printf("sorted: secrets by name")
	item := w.Secrets
	sort.Slice(item, func(i, j int) bool {
		return item[i].Name < item[j].Name
	})
}

func (w *Workflow) sortSecretByRequired() {
	log.Printf("sorted: secrets by required")
	item := w.Secrets
	sort.Slice(item, func(i, j int) bool {
		return item[i].Required.IsTrue()
	})
}

func (w *Workflow) sortOutputsByName() {
	log.Printf("sorted: outputs by name")
	item := w.Outputs
	sort.Slice(item, func(i, j int) bool {
		return item[i].Name < item[j].Name
	})
}

func (w *Workflow) sortPermissionsByScope() {
	log.Printf("sorted: permission by scope")
	item := w.Permissions
	sort.Slice(item, func(i, j int) bool {
		return item[i].Scope < item[j].Scope
	})
}

func (w *Workflow) parseInput(name string, value *WorkflowYamlInput) *WorkflowInput {
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

func (w *Workflow) parseSecret(name string, value *WorkflowYamlSecret) *WorkflowSecret {
	result := NewWorkflowSecret(name)
	if value == nil {
		return result
	}

	result.Description = util.NewNullString(value.Description)
	result.Required = util.NewNullString(value.Required)

	return result
}

func (w *Workflow) parseOutput(name string, value *WorkflowYamlOutput) *WorkflowOutput {
	result := NewWorkflowOutput(name)
	if value == nil {
		return result
	}

	result.Description = util.NewNullString(value.Description)
	return result
}

func (w *Workflow) format() string {
	if w.config.IsJson() {
		return w.toJson()
	}
	return w.toMarkdown()
}

func (w *Workflow) toJson() string {
	bytes, err := json.MarshalIndent(&WorkflowJson{Inputs: w.Inputs, Secrets: w.Secrets, Outputs: w.Outputs, Permissions: w.Permissions}, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(bytes)
}

func (w *Workflow) toMarkdown() string {
	var sb strings.Builder
	if w.hasInputs() || !w.config.Omit {
		sb.WriteString(w.toInputsMarkdown())
		sb.WriteString("\n\n")
	}

	if w.hasSecrets() || !w.config.Omit {
		sb.WriteString(w.toSecretsMarkdown())
		sb.WriteString("\n\n")
	}

	if w.hasOutputs() || !w.config.Omit {
		sb.WriteString(w.toOutputsMarkdown())
		sb.WriteString("\n\n")
	}

	if w.hasPermissions() || !w.config.Omit {
		sb.WriteString(w.toPermissionsMarkdown())
		sb.WriteString("\n\n")
	}
	return strings.TrimSpace(sb.String())
}

func (w *Workflow) toInputsMarkdown() string {
	var sb strings.Builder
	sb.WriteString(WorkflowInputsTitle)
	sb.WriteString("\n\n")
	if w.hasInputs() {
		sb.WriteString(WorkflowInputsColumnTitle)
		sb.WriteString("\n")
		sb.WriteString(WorkflowInputsColumnSeparator)
		sb.WriteString("\n")
		for _, input := range w.Inputs {
			sb.WriteString(input.toMarkdown())
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString(util.UpperNAString)
	}
	return strings.TrimSpace(sb.String())
}

func (w *Workflow) toSecretsMarkdown() string {
	var sb strings.Builder
	sb.WriteString(WorkflowSecretsTitle)
	sb.WriteString("\n\n")
	if w.hasSecrets() {
		sb.WriteString(WorkflowSecretsColumnTitle)
		sb.WriteString("\n")
		sb.WriteString(WorkflowSecretsColumnSeparator)
		sb.WriteString("\n")
		for _, secret := range w.Secrets {
			sb.WriteString(secret.toMarkdown())
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString(util.UpperNAString)
	}
	return strings.TrimSpace(sb.String())
}

func (w *Workflow) toOutputsMarkdown() string {
	var sb strings.Builder
	sb.WriteString(WorkflowOutputsTitle)
	sb.WriteString("\n\n")
	if w.hasOutputs() {
		sb.WriteString(WorkflowOutputsColumnTitle)
		sb.WriteString("\n")
		sb.WriteString(WorkflowOutputsColumnSeparator)
		sb.WriteString("\n")
		for _, output := range w.Outputs {
			sb.WriteString(output.toMarkdown())
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString(util.UpperNAString)
	}
	return strings.TrimSpace(sb.String())
}

func (w *Workflow) toPermissionsMarkdown() string {
	var sb strings.Builder
	sb.WriteString(WorkflowPermissionsTitle)
	sb.WriteString("\n\n")
	if w.hasPermissions() {
		sb.WriteString(WorkflowPermissionsColumnTitle)
		sb.WriteString("\n")
		sb.WriteString(WorkflowPermissionsColumnSeparator)
		sb.WriteString("\n")
		for _, permission := range w.Permissions {
			sb.WriteString(permission.toMarkdown())
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString(util.UpperNAString)
	}
	return strings.TrimSpace(sb.String())
}

func (w *Workflow) hasInputs() bool {
	return len(w.Inputs) != 0
}

func (w *Workflow) hasSecrets() bool {
	return len(w.Secrets) != 0
}

func (w *Workflow) hasOutputs() bool {
	return len(w.Outputs) != 0
}

func (w *Workflow) hasPermissions() bool {
	return len(w.Permissions) != 0
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
