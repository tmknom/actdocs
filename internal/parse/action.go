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
	*ActionAST
	config *format.FormatterConfig
}

func NewActionParser(config *format.FormatterConfig) *ActionParser {
	return &ActionParser{
		ActionAST: &ActionAST{
			Inputs:  []*ActionInput{},
			Outputs: []*ActionOutput{},
		},
		config: config,
	}
}

type ActionAST struct {
	Name        *util.NullString
	Description *util.NullString
	Inputs      []*ActionInput
	Outputs     []*ActionOutput
	Runs        *ActionRuns
}

func (p *ActionParser) Parse(yamlBytes []byte) (string, error) {
	content := &ActionYaml{}
	err := yaml.Unmarshal(yamlBytes, content)
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

	formatter := NewActionFormatter(p.ActionAST, p.config)
	return formatter.Format(), nil
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

func (p *ActionParser) parseInput(name string, element *actionInputYaml) {
	result := NewActionInput(name)
	if element != nil {
		result.Default = util.NewNullString(element.Default)
		result.Description = util.NewNullString(element.Description)
		result.Required = util.NewNullString(element.Required)
	}
	p.Inputs = append(p.Inputs, result)
}

func (p *ActionParser) parseOutput(name string, element *actionOutputYaml) {
	result := NewActionOutput(name)
	if element != nil {
		result.Description = util.NewNullString(element.Description)
	}
	p.Outputs = append(p.Outputs, result)
}

type ActionFormatter struct {
	*ActionAST
	config *format.FormatterConfig
}

func NewActionFormatter(ast *ActionAST, config *format.FormatterConfig) *ActionFormatter {
	return &ActionFormatter{
		ActionAST: ast,
		config:    config,
	}
}

func (f *ActionFormatter) Format() string {
	if f.config.IsJson() {
		return f.toJson()
	}
	return f.toMarkdown()
}

func (f *ActionFormatter) toJson() string {
	action := &ActionJson{
		Description: f.Description,
		Inputs:      f.Inputs,
		Outputs:     f.Outputs,
	}

	bytes, err := json.MarshalIndent(action, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(bytes)
}

func (f *ActionFormatter) toMarkdown() string {
	var sb strings.Builder
	if f.hasDescription() || !f.config.Omit {
		sb.WriteString(f.toDescriptionMarkdown())
		sb.WriteString("\n\n")
	}

	if f.hasInputs() || !f.config.Omit {
		sb.WriteString(f.toInputsMarkdown())
		sb.WriteString("\n\n")
	}

	if f.hasOutputs() || !f.config.Omit {
		sb.WriteString(f.toOutputsMarkdown())
		sb.WriteString("\n\n")
	}
	return strings.TrimSpace(sb.String())
}

func (f *ActionFormatter) toDescriptionMarkdown() string {
	var sb strings.Builder
	sb.WriteString(ActionDescriptionTitle)
	sb.WriteString("\n\n")
	sb.WriteString(f.Description.StringOrUpperNA())
	return strings.TrimSpace(sb.String())
}

func (f *ActionFormatter) toInputsMarkdown() string {
	var sb strings.Builder
	sb.WriteString(ActionInputsTitle)
	sb.WriteString("\n\n")
	if f.hasInputs() {
		sb.WriteString(ActionInputsColumnTitle)
		sb.WriteString("\n")
		sb.WriteString(ActionInputsColumnSeparator)
		sb.WriteString("\n")
		for _, input := range f.Inputs {
			sb.WriteString(input.toMarkdown())
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString(util.UpperNAString)
	}
	return strings.TrimSpace(sb.String())
}

func (f *ActionFormatter) toOutputsMarkdown() string {
	var sb strings.Builder
	sb.WriteString(ActionOutputsTitle)
	sb.WriteString("\n\n")
	if f.hasOutputs() {
		sb.WriteString(ActionOutputsColumnTitle)
		sb.WriteString("\n")
		sb.WriteString(ActionOutputsColumnSeparator)
		sb.WriteString("\n")
		for _, output := range f.Outputs {
			sb.WriteString(output.toMarkdown())
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString(util.UpperNAString)
	}
	return strings.TrimSpace(sb.String())
}

func (f *ActionFormatter) hasDescription() bool {
	return f.Description.IsValid()
}

func (f *ActionFormatter) hasInputs() bool {
	return len(f.Inputs) != 0
}

func (f *ActionFormatter) hasOutputs() bool {
	return len(f.Outputs) != 0
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

func NewActionRuns(runs *actionRunsYaml) *ActionRuns {
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
