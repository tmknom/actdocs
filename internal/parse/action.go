package parse

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"

	config2 "github.com/tmknom/actdocs/internal/format"
	"github.com/tmknom/actdocs/internal/util"
	"gopkg.in/yaml.v2"
)

type Action struct {
	Name        *util.NullString
	Description *util.NullString
	Inputs      []*ActionInput
	Outputs     []*ActionOutput
	Runs        *ActionRuns
	config      *config2.GlobalConfig
	rawYaml     []byte
}

func NewAction(rawYaml []byte, config *config2.GlobalConfig) *Action {
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
	log.Printf("unmarshal yaml: content = %#v\n", content)

	a.Name = util.NewNullString(content.Name)
	a.Description = util.NewNullString(content.Description)
	a.Runs = NewActionRuns(content.Runs)

	for name, element := range content.inputs() {
		a.parseInput(name, element)
	}

	for name, element := range content.outputs() {
		a.parseOutput(name, element)
	}

	a.sort()
	return a.format(), nil
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

	//goland:noinspection GoPreferNilSlice
	required := []*ActionInput{}
	//goland:noinspection GoPreferNilSlice
	notRequired := []*ActionInput{}
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
		result.Default = util.NewNullString(element.Default)
		result.Description = util.NewNullString(element.Description)
		result.Required = util.NewNullString(element.Required)
	}
	a.Inputs = append(a.Inputs, result)
}

func (a *Action) parseOutput(name string, element *ActionYamlOutput) {
	result := NewActionOutput(name)
	if element != nil {
		result.Description = util.NewNullString(element.Description)
	}
	a.Outputs = append(a.Outputs, result)
}

func (a *Action) format() string {
	if a.config.IsJson() {
		return a.toJson()
	}
	return a.toMarkdown()
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

func (a *Action) toMarkdown() string {
	var sb strings.Builder
	if a.hasDescription() || !a.config.Omit {
		sb.WriteString(a.toDescriptionMarkdown())
		sb.WriteString("\n\n")
	}

	if a.hasInputs() || !a.config.Omit {
		sb.WriteString(a.toInputsMarkdown())
		sb.WriteString("\n\n")
	}

	if a.hasOutputs() || !a.config.Omit {
		sb.WriteString(a.toOutputsMarkdown())
		sb.WriteString("\n\n")
	}
	return strings.TrimSpace(sb.String())
}

func (a *Action) toDescriptionMarkdown() string {
	var sb strings.Builder
	sb.WriteString(ActionDescriptionTitle)
	sb.WriteString("\n\n")
	sb.WriteString(a.Description.StringOrUpperNA())
	return strings.TrimSpace(sb.String())
}

func (a *Action) toInputsMarkdown() string {
	var sb strings.Builder
	sb.WriteString(ActionInputsTitle)
	sb.WriteString("\n\n")
	if a.hasInputs() {
		sb.WriteString(ActionInputsColumnTitle)
		sb.WriteString("\n")
		sb.WriteString(ActionInputsColumnSeparator)
		sb.WriteString("\n")
		for _, input := range a.Inputs {
			sb.WriteString(input.toMarkdown())
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString(util.UpperNAString)
	}
	return strings.TrimSpace(sb.String())
}

func (a *Action) toOutputsMarkdown() string {
	var sb strings.Builder
	sb.WriteString(ActionOutputsTitle)
	sb.WriteString("\n\n")
	if a.hasOutputs() {
		sb.WriteString(ActionOutputsColumnTitle)
		sb.WriteString("\n")
		sb.WriteString(ActionOutputsColumnSeparator)
		sb.WriteString("\n")
		for _, output := range a.Outputs {
			sb.WriteString(output.toMarkdown())
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString(util.UpperNAString)
	}
	return strings.TrimSpace(sb.String())
}

func (a *Action) hasDescription() bool {
	return a.Description.Valid
}

func (a *Action) hasInputs() bool {
	return len(a.Inputs) != 0
}

func (a *Action) hasOutputs() bool {
	return len(a.Outputs) != 0
}

const ActionDescriptionTitle = "## Description"

const ActionInputsTitle = "## Inputs"
const ActionInputsColumnTitle = "| Name | Description | Default | Required |"
const ActionInputsColumnSeparator = "| :--- | :---------- | :------ | :------: |"

const ActionOutputsTitle = "## Outputs"
const ActionOutputsColumnTitle = "| Name | Description |"
const ActionOutputsColumnSeparator = "| :--- | :---------- |"

type ActionJson struct {
	Description *util.NullString `json:"description"`
	Inputs      []*ActionInput   `json:"inputs"`
	Outputs     []*ActionOutput  `json:"outputs"`
}

type ActionInput struct {
	Name        string
	Default     *util.NullString
	Description *util.NullString
	Required    *util.NullString
}

func NewActionInput(name string) *ActionInput {
	return &ActionInput{
		Name:        name,
		Default:     util.DefaultNullString,
		Description: util.DefaultNullString,
		Required:    util.DefaultNullString,
	}
}

func (i *ActionInput) toMarkdown() string {
	str := util.TableSeparator
	str += fmt.Sprintf(" %s %s", i.Name, util.TableSeparator)
	str += fmt.Sprintf(" %s %s", i.Description.StringOrEmpty(), util.TableSeparator)
	str += fmt.Sprintf(" %s %s", i.Default.QuoteStringOrLowerNA(), util.TableSeparator)
	str += fmt.Sprintf(" %s %s", i.Required.YesOrNo(), util.TableSeparator)
	return str
}

type ActionOutput struct {
	Name        string
	Description *util.NullString
}

func NewActionOutput(name string) *ActionOutput {
	return &ActionOutput{
		Name:        name,
		Description: util.DefaultNullString,
	}
}

func (o *ActionOutput) toMarkdown() string {
	str := util.TableSeparator
	str += fmt.Sprintf(" %s %s", o.Name, util.TableSeparator)
	str += fmt.Sprintf(" %s %s", o.Description.StringOrEmpty(), util.TableSeparator)
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
	if c.Outputs == nil {
		return map[string]*ActionYamlOutput{}
	}
	return c.Outputs
}
