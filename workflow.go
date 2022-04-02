package actdocs

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"

	"gopkg.in/yaml.v2"
)

type Workflow struct {
	Inputs  []*WorkflowInput
	config  *GeneratorConfig
	rawYaml rawYaml
}

func NewWorkflow(rawYaml rawYaml, config *GeneratorConfig) *Workflow {
	return &Workflow{
		Inputs:  []*WorkflowInput{},
		config:  config,
		rawYaml: rawYaml,
	}
}

func (w *Workflow) Generate() (string, error) {
	log.Printf("config: %#v", w.config)

	content := &WorkflowYamlContent{}
	err := yaml.Unmarshal(w.rawYaml, content)
	if err != nil {
		return "", err
	}

	for name, value := range content.inputs() {
		input := w.parseInput(name, value)
		w.appendInput(input)
	}

	w.sort()
	return w.String(), nil
}

func (w *Workflow) sort() {
	switch {
	case w.config.Sort:
		w.sortInputs()
	case w.config.SortByName:
		w.sortInputsByName()
	case w.config.SortByRequired:
		w.sortInputsByRequired()
	}
}

func (w *Workflow) sortInputs() {
	log.Printf("sorted: inputs")

	var required []*WorkflowInput
	var notRequired []*WorkflowInput
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

func (w *Workflow) parseInput(name string, value *WorkflowYamlInput) *WorkflowInput {
	input := NewWorkflowInput(name)
	if value == nil {
		return input
	}

	input.Default = NewNullString(value.Default)
	input.Description = NewNullString(value.Description)
	input.Required = NewNullString(value.Required)
	input.Type = NewNullString(value.Type)

	return input
}

func (w *Workflow) appendInput(input *WorkflowInput) {
	w.Inputs = append(w.Inputs, input)
}

func (w *Workflow) String() string {
	if w.config.isJson() {
		return w.toJson()
	}
	return w.toMarkdown()
}

func (w *Workflow) toMarkdown() string {
	str := ""

	if w.hasInputs() {
		str += WorkflowTableHeader
		for _, input := range w.Inputs {
			str += input.String()
		}
	}

	return str
}

func (w *Workflow) toJson() string {
	bytes, err := json.Marshal(&WorkflowJson{w.Inputs})
	if err != nil {
		return "{}"
	}
	return string(bytes)
}

func (w *Workflow) hasInputs() bool {
	return len(w.Inputs) != 0
}

const WorkflowTableHeader = `## Inputs

| Name | Description | Type | Default | Required |
| :--- | :---------- | :--- | :------ | :------: |
`

type WorkflowJson struct {
	Inputs []*WorkflowInput `json:"inputs"`
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

func (i *WorkflowInput) String() string {
	str := TableSeparator
	str += fmt.Sprintf(" %s %s", i.Name, TableSeparator)
	str += fmt.Sprintf(" %s %s", i.Description.StringOrEmpty(), TableSeparator)
	str += fmt.Sprintf(" %s %s", i.Type.QuoteStringOrNA(), TableSeparator)
	str += fmt.Sprintf(" %s %s", i.Default.QuoteStringOrNA(), TableSeparator)
	str += fmt.Sprintf(" %s %s", i.Required.YesOrNo(), TableSeparator)
	str += "\n"
	return str
}

type WorkflowYamlContent struct {
	On *WorkflowYamlOn `yaml:"on"`
}

type WorkflowYamlOn struct {
	WorkflowCall *WorkflowYamlWorkflowCall `yaml:"workflow_call"`
}

type WorkflowYamlWorkflowCall struct {
	Inputs map[string]*WorkflowYamlInput `yaml:"inputs"`
}

type WorkflowYamlInput struct {
	Default     *string `mapstructure:"default"`
	Description *string `mapstructure:"description"`
	Required    *string `mapstructure:"required"`
	Type        *string `mapstructure:"type"`
}

func (c *WorkflowYamlContent) inputs() map[string]*WorkflowYamlInput {
	if c.On == nil || c.On.WorkflowCall == nil || c.On.WorkflowCall.Inputs == nil {
		return map[string]*WorkflowYamlInput{}
	}
	return c.On.WorkflowCall.Inputs
}
