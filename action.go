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
	rawYaml rawYaml
}

func NewAction(rawYaml rawYaml) *Action {
	return &Action{
		Inputs:  []*ActionInput{},
		rawYaml: rawYaml,
	}
}

func (a *Action) Generate() (string, error) {
	content := &ActionYamlContent{}
	err := yaml.Unmarshal(a.rawYaml, content)
	if err != nil {
		return "", err
	}

	for name, value := range content.inputs() {
		input := a.parseInput(name, value)
		a.appendInput(input)
	}

	return a.String(), nil
}

func (a *Action) parseInput(name string, value *ActionYamlInput) *ActionInput {
	input := NewActionInput(name)
	if value == nil {
		return input
	}

	input.Default = NewNullString(value.Default)
	input.Description = NewNullString(value.Description)
	input.Required = NewNullString(value.Required)

	return input
}

func (a *Action) appendInput(input *ActionInput) {
	a.Inputs = append(a.Inputs, input)
}

func (a *Action) String() string {
	str := ActionTableHeader
	for _, input := range a.Inputs {
		str += input.String()
	}
	return str
}

const ActionTableHeader = `
| Name | Description | Default | Required |
| :--- | :---------- | :------ | :------: |
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

type ActionYamlContent struct {
	Inputs map[string]*ActionYamlInput `yaml:"inputs"`
}

type ActionYamlInput struct {
	Default     interface{} `yaml:"default"`
	Description interface{} `yaml:"description"`
	Required    interface{} `yaml:"required"`
}

func (c *ActionYamlContent) inputs() map[string]*ActionYamlInput {
	if c.Inputs == nil {
		return map[string]*ActionYamlInput{}
	}
	return c.Inputs
}
