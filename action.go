package actdocs

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v2"
)

type ActionCmd struct {
	*TemplateConfig
	// args is actual args parsed from flags.
	args []string
	// inReader is a reader defined by the user that replaces stdin
	inReader io.Reader
	// outWriter is a writer defined by the user that replaces stdout
	outWriter io.Writer
	// errWriter is a writer defined by the user that replaces stderr
	errWriter io.Writer
}

func NewActionCmd(config *TemplateConfig, args []string, inReader io.Reader, outWriter, errWriter io.Writer) *ActionCmd {
	return &ActionCmd{
		TemplateConfig: config,
		args:           args,
		inReader:       inReader,
		outWriter:      outWriter,
		errWriter:      errWriter,
	}
}

func (c *ActionCmd) Run() (err error) {
	filename := c.args[0]
	rawYaml, err := readYaml(filename)
	if err != nil {
		return err
	}

	action := NewAction(rawYaml)
	content, err := action.Generate()
	if err != nil {
		return err
	}

	template := NewTemplate(c.TemplateConfig)
	err = template.Render(content)
	if err != nil {
		return err
	}

	return nil
}

type Action struct {
	Name        *NullString
	Description *NullString
	Inputs      []*ActionInput
	Outputs     []*ActionOutput
	Runs        *ActionRuns
	rawYaml     rawYaml
}

func NewAction(rawYaml rawYaml) *Action {
	return &Action{
		Inputs:  []*ActionInput{},
		Outputs: []*ActionOutput{},
		rawYaml: rawYaml,
	}
}

func (a *Action) Generate() (string, error) {
	content := &ActionYamlContent{}
	err := yaml.Unmarshal(a.rawYaml, content)
	if err != nil {
		return "", err
	}

	a.Name = NewNullString(content.Name)
	a.Description = NewNullString(content.Description)
	a.Runs = NewActionRuns(content.Runs)

	for name, element := range content.inputs() {
		a.parseInput(name, element)
	}

	for name, element := range content.outputs() {
		a.parseOutput(name, element)
	}

	return a.String(), nil
}

func (a *Action) parseInput(name string, element *ActionYamlInput) {
	result := NewActionInput(name)
	if element != nil {
		result.Default = NewNullString(element.Default)
		result.Description = NewNullString(element.Description)
		result.Required = NewNullString(element.Required)
	}
	a.Inputs = append(a.Inputs, result)
}

func (a *Action) parseOutput(name string, element *ActionYamlOutput) {
	result := NewActionOutput(name)
	if element != nil {
		result.Description = NewNullString(element.Description)
	}
	a.Outputs = append(a.Outputs, result)
}

func (a *Action) hasInputs() bool {
	return len(a.Inputs) != 0
}

func (a *Action) hasOutputs() bool {
	return len(a.Outputs) != 0
}

func (a *Action) String() string {
	str := ""

	if a.hasInputs() {
		str += ActionTableHeader
		for _, input := range a.Inputs {
			str += input.String()
		}
		str += "\n"
	}

	if a.hasOutputs() {
		str += ActionOutputsTableHeader
		for _, output := range a.Outputs {
			str += output.String()
		}
		str += "\n"
	}
	return str
}

const ActionDescriptionHeader = `## Description

`

const ActionTableHeader = `## Inputs

| Name | Description | Default | Required |
| :--- | :---------- | :------ | :------: |
`

const ActionOutputsTableHeader = `## Outputs

| Name | Description |
| :--- | :---------- |
`

type ActionInput struct {
	Name        string
	Default     *NullString
	Description *NullString
	Required    *NullString
}

func NewActionInput(name string) *ActionInput {
	return &ActionInput{
		Name:        name,
		Default:     DefaultNullString,
		Description: DefaultNullString,
		Required:    DefaultNullString,
	}
}

func (i *ActionInput) String() string {
	str := TableSeparator
	str += fmt.Sprintf(" %s %s", i.Name, TableSeparator)
	str += fmt.Sprintf(" %s %s", i.Description.StringOrEmpty(), TableSeparator)
	str += fmt.Sprintf(" %s %s", i.Default.StringOrEmpty(), TableSeparator)
	str += fmt.Sprintf(" %s %s", i.Required.StringOrEmpty(), TableSeparator)
	str += "\n"
	return str
}

type ActionOutput struct {
	Name        string
	Description *NullString
}

func NewActionOutput(name string) *ActionOutput {
	return &ActionOutput{
		Name:        name,
		Description: DefaultNullString,
	}
}

func (o *ActionOutput) String() string {
	str := TableSeparator
	str += fmt.Sprintf(" %s %s", o.Name, TableSeparator)
	str += fmt.Sprintf(" %s %s", o.Description.StringOrEmpty(), TableSeparator)
	str += "\n"
	return str
}

type ActionRuns struct {
	Using string
	Steps []*interface{}
}

func NewActionRuns(runs *ActionYamlRuns) *ActionRuns {
	result := &ActionRuns{
		Using: "undefined",
		Steps: []*interface{}{},
	}

	if runs != nil {
		result.Using = runs.Using
		result.Steps = runs.Steps
	}
	return result
}

func (r *ActionRuns) String() string {
	str := ""
	str += fmt.Sprintf("Using: %s, ", r.Using)
	str += fmt.Sprintf("Steps: [")
	for _, step := range r.Steps {
		str += fmt.Sprintf("%#v, ", *step)
	}
	str += fmt.Sprintf("]")
	return str
}

type ActionYamlContent struct {
	Name        *string                      `yaml:"name"`
	Description *string                      `yaml:"description"`
	Inputs      map[string]*ActionYamlInput  `yaml:"inputs"`
	Outputs     map[string]*ActionYamlOutput `yaml:"outputs"`
	Runs        *ActionYamlRuns              `yaml:"runs"`
}

type ActionYamlInput struct {
	Default     *string `mapstructure:"default"`
	Description *string `mapstructure:"description"`
	Required    *string `mapstructure:"required"`
}

type ActionYamlOutput struct {
	Description *string `mapstructure:"description"`
}

type ActionYamlRuns struct {
	Using string         `yaml:"using"`
	Steps []*interface{} `yaml:"steps"`
}

func (c *ActionYamlContent) inputs() map[string]*ActionYamlInput {
	if c.Inputs == nil {
		return map[string]*ActionYamlInput{}
	}
	return c.Inputs
}

func (c *ActionYamlContent) outputs() map[string]*ActionYamlOutput {
	if c.Inputs == nil {
		return map[string]*ActionYamlOutput{}
	}
	return c.Outputs
}
