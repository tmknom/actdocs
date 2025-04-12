package parse

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/tmknom/actdocs/internal/format"
	"github.com/tmknom/actdocs/internal/util"
	"gopkg.in/yaml.v2"
)

type ActionParser struct {
	Name        *util.NullString
	Description *util.NullString
	Inputs      []*ActionInput
	Outputs     []*ActionOutput
	Runs        *ActionRuns
	config      *format.FormatterConfig
	rawYaml     []byte
}

func NewActionParser(rawYaml []byte, config *format.FormatterConfig) *ActionParser {
	return &ActionParser{
		Inputs:  []*ActionInput{},
		Outputs: []*ActionOutput{},
		config:  config,
		rawYaml: rawYaml,
	}
}

func (p *ActionParser) Parse() (string, error) {
	content := &ActionYamlContent{}
	err := yaml.Unmarshal(p.rawYaml, content)
	if err != nil {
		return "", err
	}
	log.Printf("unmarshal yaml: content = %#v\n", content)

	p.Name = util.NewNullString(content.Name)
	p.Description = util.NewNullString(content.Description)
	p.Runs = NewActionRuns(content.Runs)

	for name, element := range content.inputs() {
		p.parseInput(name, element)
	}

	for name, element := range content.outputs() {
		p.parseOutput(name, element)
	}

	p.sort()
	return p.format(), nil
}

func (p *ActionParser) sort() {
	switch {
	case p.config.Sort:
		p.sortInputs()
		p.sortOutputsByName()
	case p.config.SortByName:
		p.sortInputsByName()
		p.sortOutputsByName()
	case p.config.SortByRequired:
		p.sortInputsByRequired()
	}
}

func (p *ActionParser) sortInputs() {
	log.Printf("sorted: inputs")

	//goland:noinspection GoPreferNilSlice
	required := []*ActionInput{}
	//goland:noinspection GoPreferNilSlice
	notRequired := []*ActionInput{}
	for _, input := range p.Inputs {
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
	p.Inputs = append(required, notRequired...)
}

func (p *ActionParser) sortInputsByName() {
	log.Printf("sorted: inputs by name")
	item := p.Inputs
	sort.Slice(item, func(i, j int) bool {
		return item[i].Name < item[j].Name
	})
}

func (p *ActionParser) sortInputsByRequired() {
	log.Printf("sorted: inputs by required")
	item := p.Inputs
	sort.Slice(item, func(i, j int) bool {
		return item[i].Required.IsTrue()
	})
}

func (p *ActionParser) sortOutputsByName() {
	log.Printf("sorted: outputs by name")
	item := p.Outputs
	sort.Slice(item, func(i, j int) bool {
		return item[i].Name < item[j].Name
	})
}

func (p *ActionParser) parseInput(name string, element *ActionYamlInput) {
	result := NewActionInput(name)
	if element != nil {
		result.Default = util.NewNullString(element.Default)
		result.Description = util.NewNullString(element.Description)
		result.Required = util.NewNullString(element.Required)
	}
	p.Inputs = append(p.Inputs, result)
}

func (p *ActionParser) parseOutput(name string, element *ActionYamlOutput) {
	result := NewActionOutput(name)
	if element != nil {
		result.Description = util.NewNullString(element.Description)
	}
	p.Outputs = append(p.Outputs, result)
}

func (p *ActionParser) format() string {
	if p.config.IsJson() {
		return p.toJson()
	}
	return p.toMarkdown()
}

func (p *ActionParser) toJson() string {
	action := &ActionJson{
		Description: p.Description,
		Inputs:      p.Inputs,
		Outputs:     p.Outputs,
	}

	bytes, err := json.MarshalIndent(action, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(bytes)
}

func (p *ActionParser) toMarkdown() string {
	var sb strings.Builder
	if p.hasDescription() || !p.config.Omit {
		sb.WriteString(p.toDescriptionMarkdown())
		sb.WriteString("\n\n")
	}

	if p.hasInputs() || !p.config.Omit {
		sb.WriteString(p.toInputsMarkdown())
		sb.WriteString("\n\n")
	}

	if p.hasOutputs() || !p.config.Omit {
		sb.WriteString(p.toOutputsMarkdown())
		sb.WriteString("\n\n")
	}
	return strings.TrimSpace(sb.String())
}

func (p *ActionParser) toDescriptionMarkdown() string {
	var sb strings.Builder
	sb.WriteString(ActionDescriptionTitle)
	sb.WriteString("\n\n")
	sb.WriteString(p.Description.StringOrUpperNA())
	return strings.TrimSpace(sb.String())
}

func (p *ActionParser) toInputsMarkdown() string {
	var sb strings.Builder
	sb.WriteString(ActionInputsTitle)
	sb.WriteString("\n\n")
	if p.hasInputs() {
		sb.WriteString(ActionInputsColumnTitle)
		sb.WriteString("\n")
		sb.WriteString(ActionInputsColumnSeparator)
		sb.WriteString("\n")
		for _, input := range p.Inputs {
			sb.WriteString(input.toMarkdown())
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString(util.UpperNAString)
	}
	return strings.TrimSpace(sb.String())
}

func (p *ActionParser) toOutputsMarkdown() string {
	var sb strings.Builder
	sb.WriteString(ActionOutputsTitle)
	sb.WriteString("\n\n")
	if p.hasOutputs() {
		sb.WriteString(ActionOutputsColumnTitle)
		sb.WriteString("\n")
		sb.WriteString(ActionOutputsColumnSeparator)
		sb.WriteString("\n")
		for _, output := range p.Outputs {
			sb.WriteString(output.toMarkdown())
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString(util.UpperNAString)
	}
	return strings.TrimSpace(sb.String())
}

func (p *ActionParser) hasDescription() bool {
	return p.Description.IsValid()
}

func (p *ActionParser) hasInputs() bool {
	return len(p.Inputs) != 0
}

func (p *ActionParser) hasOutputs() bool {
	return len(p.Outputs) != 0
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
