package actdocs

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

type Workflow struct {
	Inputs  []*WorkflowInput
	rawYaml rawYaml
}

func NewWorkflow(rawYaml rawYaml) *Workflow {
	return &Workflow{
		Inputs:  []*WorkflowInput{},
		rawYaml: rawYaml,
	}
}

func (w *Workflow) Generate() (string, error) {
	content := &WorkflowYamlContent{}
	err := yaml.Unmarshal(w.rawYaml, content)
	if err != nil {
		return "", err
	}

	for name, value := range content.inputs() {
		input := w.parseInput(name, value)
		w.appendInput(input)
	}

	return w.String(), nil
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
	str := ""

	if w.hasInputs() {
		str += WorkflowTableHeader
		for _, input := range w.Inputs {
			str += input.String()
		}
	}

	return str
}

func (w *Workflow) hasInputs() bool {
	return len(w.Inputs) != 0
}

const WorkflowTableHeader = `## Inputs

| Name | Description | Type | Default | Required |
| :--- | :---------- | :--- | :------ | :------: |
`

type WorkflowInput struct {
	Name        string
	Default     *NullString
	Description *NullString
	Required    *NullString
	Type        *NullString
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
