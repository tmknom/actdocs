package actdocs

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type WorkflowCmd struct{}

func NewWorkflowCmd() *WorkflowCmd {
	return &WorkflowCmd{}
}

func (c *WorkflowCmd) Run(command *cobra.Command, args []string) error {
	filename := args[0]
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	command.SetIn(file)
	workflow := NewWorkflow(command.InOrStdin(), command.OutOrStdout(), command.ErrOrStderr())
	return workflow.Generate()
}

type Workflow struct {
	Inputs    []*WorkflowInput
	inStream  io.Reader
	outStream io.Writer
	errStream io.Writer
}

func NewWorkflow(inStream io.Reader, outStream, errStream io.Writer) *Workflow {
	return &Workflow{
		Inputs:    []*WorkflowInput{},
		inStream:  inStream,
		outStream: outStream,
		errStream: errStream,
	}
}

func (w *Workflow) Generate() error {
	content, err := w.readYaml()
	if err != nil {
		return err
	}

	for name, value := range content.inputs() {
		input := w.parseInput(name, value)
		w.appendInput(input)
	}
	w.String()

	return nil
}

func (w *Workflow) readYaml() (*YamlContent, error) {
	bytes, err := io.ReadAll(w.inStream)
	if err != nil {
		return nil, err
	}

	content := &YamlContent{}
	if yaml.Unmarshal(bytes, content) != nil {
		return nil, err
	}

	return content, nil
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

func (w *Workflow) String() {
	str := TableHeader
	for _, input := range w.Inputs {
		str += input.String()
	}
	fmt.Fprint(w.outStream, str)
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

const TableSeparator = "|"

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
