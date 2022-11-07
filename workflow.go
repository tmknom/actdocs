package actdocs

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

type Workflow struct {
	Inputs  []*WorkflowInput
	Secrets []*WorkflowSecret
	config  *GlobalConfig
	rawYaml RawYaml
}

func NewWorkflow(rawYaml RawYaml, config *GlobalConfig) *Workflow {
	return &Workflow{
		Inputs:  []*WorkflowInput{},
		Secrets: []*WorkflowSecret{},
		config:  config,
		rawYaml: rawYaml,
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

	for name, value := range content.secrets() {
		secret := w.parseSecret(name, value)
		w.Secrets = append(w.Secrets, secret)
	}

	w.sort()
	return w.format(), nil
}

func (w *Workflow) sort() {
	switch {
	case w.config.Sort:
		w.sortInputs()
		w.sortSecrets()
	case w.config.SortByName:
		w.sortInputsByName()
		w.sortSecretsByName()
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

func (w *Workflow) parseInput(name string, value *WorkflowYamlInput) *WorkflowInput {
	result := NewWorkflowInput(name)
	if value == nil {
		return result
	}

	result.Default = NewNullString(value.Default)
	result.Description = NewNullString(value.Description)
	result.Required = NewNullString(value.Required)
	result.Type = NewNullString(value.Type)

	return result
}

func (w *Workflow) parseSecret(name string, value *WorkflowYamlSecret) *WorkflowSecret {
	result := NewWorkflowSecret(name)
	if value == nil {
		return result
	}

	result.Description = NewNullString(value.Description)
	result.Required = NewNullString(value.Required)

	return result
}

func (w *Workflow) format() string {
	if w.config.isJson() {
		return w.toJson()
	}
	return w.toMarkdown()
}

func (w *Workflow) toJson() string {
	bytes, err := json.MarshalIndent(&WorkflowJson{Inputs: w.Inputs, Secrets: w.Secrets}, "", "  ")
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
		sb.WriteString(UpperNAString)
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
		sb.WriteString(UpperNAString)
	}
	return strings.TrimSpace(sb.String())
}

func (w *Workflow) hasInputs() bool {
	return len(w.Inputs) != 0
}

func (w *Workflow) hasSecrets() bool {
	return len(w.Secrets) != 0
}

const WorkflowInputsTitle = "## Inputs"
const WorkflowInputsColumnTitle = "| Name | Description | Type | Default | Required |"
const WorkflowInputsColumnSeparator = "| :--- | :---------- | :--- | :------ | :------: |"

const WorkflowSecretsTitle = "## Secrets"
const WorkflowSecretsColumnTitle = "| Name | Description | Required |"
const WorkflowSecretsColumnSeparator = "| :--- | :---------- | :------: |"

type WorkflowJson struct {
	Inputs  []*WorkflowInput  `json:"inputs"`
	Secrets []*WorkflowSecret `json:"secrets"`
}

type WorkflowInput struct {
	Name        string      `json:"name"`
	Default     *NullString `json:"default"`
	Description *NullString `json:"description"`
	Required    *NullString `json:"required"`
	Type        *NullString `json:"type"`
}

func NewWorkflowInput(name string) *WorkflowInput {
	return &WorkflowInput{
		Name:        name,
		Default:     DefaultNullString,
		Description: DefaultNullString,
		Required:    DefaultNullString,
		Type:        DefaultNullString,
	}
}

func (i *WorkflowInput) toMarkdown() string {
	str := TableSeparator
	str += fmt.Sprintf(" %s %s", i.Name, TableSeparator)
	str += fmt.Sprintf(" %s %s", i.Description.StringOrEmpty(), TableSeparator)
	str += fmt.Sprintf(" %s %s", i.Type.QuoteStringOrLowerNA(), TableSeparator)
	str += fmt.Sprintf(" %s %s", i.Default.QuoteStringOrLowerNA(), TableSeparator)
	str += fmt.Sprintf(" %s %s", i.Required.YesOrNo(), TableSeparator)
	return str
}

type WorkflowSecret struct {
	Name        string      `json:"name"`
	Description *NullString `json:"description"`
	Required    *NullString `json:"required"`
}

func NewWorkflowSecret(name string) *WorkflowSecret {
	return &WorkflowSecret{
		Name:        name,
		Description: DefaultNullString,
		Required:    DefaultNullString,
	}
}

func (i *WorkflowSecret) toMarkdown() string {
	str := TableSeparator
	str += fmt.Sprintf(" %s %s", i.Name, TableSeparator)
	str += fmt.Sprintf(" %s %s", i.Description.StringOrEmpty(), TableSeparator)
	str += fmt.Sprintf(" %s %s", i.Required.YesOrNo(), TableSeparator)
	return str
}

type WorkflowYamlContent struct {
	On *WorkflowYamlOn `yaml:"on"`
}

type WorkflowYamlOn struct {
	WorkflowCall *WorkflowYamlWorkflowCall `yaml:"workflow_call"`
}

type WorkflowYamlWorkflowCall struct {
	Inputs  map[string]*WorkflowYamlInput  `yaml:"inputs"`
	Secrets map[string]*WorkflowYamlSecret `yaml:"secrets"`
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
