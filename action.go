package actdocs

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v2"
)

type ActionCmd struct {
	// args is actual args parsed from flags.
	args []string
	// inReader is a reader defined by the user that replaces stdin
	inReader io.Reader
	// outWriter is a writer defined by the user that replaces stdout
	outWriter io.Writer
	// errWriter is a writer defined by the user that replaces stderr
	errWriter io.Writer
}

func NewActionCmd(args []string, inReader io.Reader, outWriter, errWriter io.Writer) *ActionCmd {
	return &ActionCmd{
		args:      args,
		inReader:  inReader,
		outWriter: outWriter,
		errWriter: errWriter,
	}
}

func (c *ActionCmd) Run() (err error) {
	filename := c.args[0]
	rawYaml, err := readYaml(filename)
	if err != nil {
		return err
	}

	action := NewAction(rawYaml)
	result, err := action.Generate()
	if err != nil {
		return err
	}
	fmt.Fprint(c.outWriter, result)

	return nil
}

type Action struct {
	Inputs  []*ActionInput
	Outputs []*ActionOutput
	rawYaml rawYaml
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
	}

	if a.hasOutputs() {
		str += ActionOutputsTableHeader
		for _, output := range a.Outputs {
			str += output.String()
		}
	}
	return str
}

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

type ActionYamlContent struct {
	Inputs  map[string]*ActionYamlInput  `yaml:"inputs"`
	Outputs map[string]*ActionYamlOutput `yaml:"outputs"`
}

type ActionYamlInput struct {
	Default     interface{} `yaml:"default"`
	Description interface{} `yaml:"description"`
	Required    interface{} `yaml:"required"`
}

type ActionYamlOutput struct {
	Description interface{} `yaml:"description"`
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
