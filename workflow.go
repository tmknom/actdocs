package actdocs

import (
	"fmt"
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
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var data YamlData
	if yaml.Unmarshal(bytes, &data) != nil {
		return err
	}

	for name, value := range data.On.WorkflowCall.Inputs {
		input := c.parseInput(name, &value)
		fmt.Fprint(command.OutOrStdout(), input.String())
	}

	return nil
}

func (c *WorkflowCmd) parseInput(name string, value *YamlInput) *WorkflowInput {
	input := &WorkflowInput{Name: name}

	if value.Default != nil {
		str := fmt.Sprint(value.Default)
		input.Default = &str
	}

	if value.Description != nil {
		str := fmt.Sprint(value.Description)
		input.Description = &str
	}

	if value.Required != nil {
		str := fmt.Sprint(value.Required)
		input.Required = &str
	}

	if value.Type != nil {
		str := fmt.Sprint(value.Type)
		input.Type = &str
	}

	return input
}

type WorkflowInput struct {
	Name        string
	Default     *string
	Description *string
	Required    *string
	Type        *string
}

func (i *WorkflowInput) String() string {
	str := fmt.Sprint("-----------\n")
	str += fmt.Sprintf("%s : {\n", i.Name)

	if i.Default != nil {
		str += fmt.Sprintf("   Default: %s\n", *i.Default)
	}
	if i.Description != nil {
		str += fmt.Sprintf("   Description: %s\n", *i.Description)
	}
	if i.Required != nil {
		str += fmt.Sprintf("   Required: %s\n", *i.Required)
	}
	if i.Type != nil {
		str += fmt.Sprintf("   Type: %s\n", *i.Type)
	}

	str += fmt.Sprint("}\n")
	return str
}

type YamlData struct {
	On YamlOn `yaml:"on"`
}

type YamlOn struct {
	WorkflowCall YamlWorkflowCall `yaml:"workflow_call"`
}

type YamlWorkflowCall struct {
	Inputs map[string]YamlInput `yaml:"inputs"`
}

type YamlInput struct {
	Default     interface{} `yaml:"default"`
	Description interface{} `yaml:"description"`
	Required    interface{} `yaml:"required"`
	Type        interface{} `yaml:"type"`
}
