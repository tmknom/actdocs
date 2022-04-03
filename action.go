package actdocs

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

type Action struct {
	Name        *NullString
	Description *NullString
	Inputs      []*ActionInput
	Outputs     []*ActionOutput
	Runs        *ActionRuns
	config      *GlobalConfig
	rawYaml     RawYaml
}

func NewAction(rawYaml RawYaml, config *GlobalConfig) *Action {
	return &Action{
		Inputs:  []*ActionInput{},
		Outputs: []*ActionOutput{},
		config:  config,
		rawYaml: rawYaml,
	}
}

func (a *Action) Parse() (string, error) {
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

	a.sort()
	return a.String(), nil
}

func (a *Action) sort() {
	switch {
	case a.config.Sort:
		a.sortInputs()
		a.sortOutputsByName()
	case a.config.SortByName:
		a.sortInputsByName()
		a.sortOutputsByName()
	case a.config.SortByRequired:
		a.sortInputsByRequired()
	}
}

func (a *Action) sortInputs() {
	log.Printf("sorted: inputs")

	var required []*ActionInput
	var notRequired []*ActionInput
	for _, input := range a.Inputs {
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
	a.Inputs = append(required, notRequired...)
}

func (a *Action) sortInputsByName() {
	log.Printf("sorted: inputs by name")
	item := a.Inputs
	sort.Slice(item, func(i, j int) bool {
		return item[i].Name < item[j].Name
	})
}

func (a *Action) sortInputsByRequired() {
	log.Printf("sorted: inputs by required")
	item := a.Inputs
	sort.Slice(item, func(i, j int) bool {
		return item[i].Required.IsTrue()
	})
}

func (a *Action) sortOutputsByName() {
	log.Printf("sorted: outputs by name")
	item := a.Outputs
	sort.Slice(item, func(i, j int) bool {
		return item[i].Name < item[j].Name
	})
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

func (a *Action) String() string {
	if a.config.isJson() {
		return a.toJson()
	}
	return a.toMarkdown()
}

func (a *Action) toMarkdown() string {
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
	return strings.TrimSpace(str) + "\n"
}

func (a *Action) toJson() string {
	action := &ActionJson{
		Description: a.Description,
		Inputs:      a.Inputs,
		Outputs:     a.Outputs,
	}

	bytes, err := json.MarshalIndent(action, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(bytes)
}

func (a *Action) hasInputs() bool {
	return len(a.Inputs) != 0
}

func (a *Action) hasOutputs() bool {
	return len(a.Outputs) != 0
}

const ActionTableHeader = `## Inputs

| Name | Description | Default | Required |
| :--- | :---------- | :------ | :------: |
`

const ActionOutputsTableHeader = `## Outputs

| Name | Description |
| :--- | :---------- |
`

type ActionJson struct {
	Description *NullString     `json:"description"`
	Inputs      []*ActionInput  `json:"inputs"`
	Outputs     []*ActionOutput `json:"outputs"`
}

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
	str += fmt.Sprintf(" %s %s", i.Default.QuoteStringOrNA(), TableSeparator)
	str += fmt.Sprintf(" %s %s", i.Required.YesOrNo(), TableSeparator)
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
