package format

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tmknom/actdocs/internal/conf"
	"github.com/tmknom/actdocs/internal/parse"
	"github.com/tmknom/actdocs/internal/util"
)

type ActionFormatter struct {
	config *conf.FormatterConfig
	*ActionSpec
}

func NewActionFormatter(config *conf.FormatterConfig) *ActionFormatter {
	return &ActionFormatter{
		config: config,
	}
}

func (f *ActionFormatter) Format(ast *parse.ActionAST) string {
	f.ActionSpec = f.convertActionMarkdown(ast)

	if f.config.IsJson() {
		return f.ToJson(f.ActionSpec)
	}
	return f.ToMarkdown(f.ActionSpec, f.config)
}

func (f *ActionFormatter) ToJson(actionSpec *ActionSpec) string {
	bytes, err := json.MarshalIndent(actionSpec, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(bytes)
}

func (f *ActionFormatter) ToMarkdown(actionSpec *ActionSpec, config *conf.FormatterConfig) string {
	var sb strings.Builder
	if actionSpec.Description.IsValid() || !config.Omit {
		sb.WriteString(f.toDescriptionMarkdown(actionSpec.Description))
		sb.WriteString("\n\n")
	}

	if len(actionSpec.Inputs) != 0 || !config.Omit {
		sb.WriteString(f.toInputsMarkdown(actionSpec.Inputs))
		sb.WriteString("\n\n")
	}

	if len(actionSpec.Outputs) != 0 || !config.Omit {
		sb.WriteString(f.toOutputsMarkdown(actionSpec.Outputs))
		sb.WriteString("\n\n")
	}
	return strings.TrimSpace(sb.String())
}

func (f *ActionFormatter) convertActionMarkdown(ast *parse.ActionAST) *ActionSpec {
	inputs := []*ActionInputSpec{}
	for _, inputAst := range ast.Inputs {
		input := &ActionInputSpec{
			Name:        inputAst.Name,
			Default:     inputAst.Default,
			Description: inputAst.Description,
			Required:    inputAst.Required,
		}
		inputs = append(inputs, input)
	}

	outputs := []*ActionOutputSpec{}
	for _, outputAst := range ast.Outputs {
		output := &ActionOutputSpec{
			Name:        outputAst.Name,
			Description: outputAst.Description,
		}
		outputs = append(outputs, output)
	}

	return &ActionSpec{
		Description: ast.Description,
		Inputs:      inputs,
		Outputs:     outputs,
	}
}

func (f *ActionFormatter) toDescriptionMarkdown(description *util.NullString) string {
	var sb strings.Builder
	sb.WriteString(ActionDescriptionTitle)
	sb.WriteString("\n\n")
	sb.WriteString(description.StringOrUpperNA())
	return strings.TrimSpace(sb.String())
}

func (f *ActionFormatter) toInputsMarkdown(inputs []*ActionInputSpec) string {
	var sb strings.Builder
	sb.WriteString(ActionInputsTitle)
	sb.WriteString("\n\n")
	if len(inputs) != 0 {
		sb.WriteString(ActionInputsColumnTitle)
		sb.WriteString("\n")
		sb.WriteString(ActionInputsColumnSeparator)
		sb.WriteString("\n")
		for _, input := range inputs {
			sb.WriteString(input.toMarkdown())
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString(util.UpperNAString)
	}
	return strings.TrimSpace(sb.String())
}

func (f *ActionFormatter) toOutputsMarkdown(outputs []*ActionOutputSpec) string {
	var sb strings.Builder
	sb.WriteString(ActionOutputsTitle)
	sb.WriteString("\n\n")
	if len(outputs) != 0 {
		sb.WriteString(ActionOutputsColumnTitle)
		sb.WriteString("\n")
		sb.WriteString(ActionOutputsColumnSeparator)
		sb.WriteString("\n")
		for _, output := range outputs {
			sb.WriteString(output.toMarkdown())
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString(util.UpperNAString)
	}
	return strings.TrimSpace(sb.String())
}

const ActionDescriptionTitle = "## Description"

const ActionInputsTitle = "## Inputs"
const ActionInputsColumnTitle = "| Name | Description | Default | Required |"
const ActionInputsColumnSeparator = "| :--- | :---------- | :------ | :------: |"

const ActionOutputsTitle = "## Outputs"
const ActionOutputsColumnTitle = "| Name | Description |"
const ActionOutputsColumnSeparator = "| :--- | :---------- |"

type ActionSpec struct {
	Description *util.NullString    `json:"description"`
	Inputs      []*ActionInputSpec  `json:"inputs"`
	Outputs     []*ActionOutputSpec `json:"outputs"`
}

type ActionInputSpec struct {
	Name        string           `json:"name"`
	Default     *util.NullString `json:"default"`
	Description *util.NullString `json:"description"`
	Required    *util.NullString `json:"required"`
}

func (s *ActionInputSpec) toMarkdown() string {
	str := util.TableSeparator
	str += fmt.Sprintf(" %s %s", s.Name, util.TableSeparator)
	str += fmt.Sprintf(" %s %s", s.Description.StringOrEmpty(), util.TableSeparator)
	str += fmt.Sprintf(" %s %s", s.Default.QuoteStringOrLowerNA(), util.TableSeparator)
	str += fmt.Sprintf(" %s %s", s.Required.YesOrNo(), util.TableSeparator)
	return str
}

type ActionOutputSpec struct {
	Name        string           `json:"name"`
	Description *util.NullString `json:"description"`
}

func (s *ActionOutputSpec) toMarkdown() string {
	str := util.TableSeparator
	str += fmt.Sprintf(" %s %s", s.Name, util.TableSeparator)
	str += fmt.Sprintf(" %s %s", s.Description.StringOrEmpty(), util.TableSeparator)
	return str
}
