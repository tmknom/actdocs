package actdocs

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

type rawYaml []byte

type WorkflowCmd struct {
	// args is actual args parsed from flags.
	args []string
	// inReader is a reader defined by the user that replaces stdin
	inReader io.Reader
	// outWriter is a writer defined by the user that replaces stdout
	outWriter io.Writer
	// errWriter is a writer defined by the user that replaces stderr
	errWriter io.Writer
}

func NewWorkflowCmd(args []string, inReader io.Reader, outWriter, errWriter io.Writer) *WorkflowCmd {
	return &WorkflowCmd{
		args:      args,
		inReader:  inReader,
		outWriter: outWriter,
		errWriter: errWriter,
	}
}

func (c *WorkflowCmd) Run() (err error) {
	filename := c.args[0]
	rawYaml, err := readYaml(filename)
	if err != nil {
		return err
	}

	workflow := NewWorkflow(rawYaml)
	result, err := workflow.Generate()
	if err != nil {
		return err
	}
	fmt.Fprint(c.outWriter, result)

	return nil
}

func readYaml(filename string) (rawYaml rawYaml, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) { err = file.Close() }(file)

	return io.ReadAll(file)
}

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
	content := &YamlContent{}
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

func (w *Workflow) parseInput(name string, value *YamlInput) *WorkflowInput {
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
	str := TableHeader
	for _, input := range w.Inputs {
		str += input.String()
	}
	return str
}

const TableHeader = `
| Name | Description | Default | Type  | Required |
| :--- | :---------- | :------ | :---: | :------: |
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
	str += fmt.Sprintf(" %s %s", i.Default.StringOrEmpty(), TableSeparator)
	str += fmt.Sprintf(" %s %s", i.Type.StringOrEmpty(), TableSeparator)
	str += fmt.Sprintf(" %s %s", i.Required.StringOrEmpty(), TableSeparator)
	str += "\n"
	return str
}

const TableSeparator = "|"

// NullString represents a string that may be null.
type NullString struct {
	String string
	Valid  bool // Valid is true if String is not NULL
}

func NewNullString(value interface{}) *NullString {
	return &NullString{
		String: fmt.Sprint(value),
		Valid:  value != nil,
	}
}

var DefaultNullString = NewNullString(nil)

func (s *NullString) StringOrEmpty() string {
	if s.Valid {
		return s.String
	}
	return ""
}

type YamlContent struct {
	On *YamlOn `yaml:"on"`
}

type YamlOn struct {
	WorkflowCall *YamlWorkflowCall `yaml:"workflow_call"`
}

type YamlWorkflowCall struct {
	Inputs map[string]*YamlInput `yaml:"inputs"`
}

type YamlInput struct {
	Default     interface{} `yaml:"default"`
	Description interface{} `yaml:"description"`
	Required    interface{} `yaml:"required"`
	Type        interface{} `yaml:"type"`
}

func (c *YamlContent) inputs() map[string]*YamlInput {
	if c.On == nil || c.On.WorkflowCall == nil || c.On.WorkflowCall.Inputs == nil {
		return map[string]*YamlInput{}
	}
	return c.On.WorkflowCall.Inputs
}
