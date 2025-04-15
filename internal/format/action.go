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
	*ActionMarkdown
}

func NewActionFormatter(config *conf.FormatterConfig) *ActionFormatter {
	return &ActionFormatter{
		config: config,
	}
}

func (f *ActionFormatter) Format(ast *parse.ActionAST) string {
	f.ActionMarkdown = f.convertActionMarkdown(ast)

	if f.config.IsJson() {
		return f.ToJson(f.ActionMarkdown)
	}
	return f.ToMarkdown(f.ActionMarkdown, f.config)
}

func (f *ActionFormatter) ToJson(actionJson *ActionMarkdown) string {
	bytes, err := json.MarshalIndent(actionJson, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(bytes)
}

func (f *ActionFormatter) ToMarkdown(actionMarkdown *ActionMarkdown, config *conf.FormatterConfig) string {
	var sb strings.Builder
	if actionMarkdown.Description.IsValid() || !config.Omit {
		sb.WriteString(f.toDescriptionMarkdown(actionMarkdown.Description))
		sb.WriteString("\n\n")
	}

	if len(actionMarkdown.Inputs) != 0 || !config.Omit {
		sb.WriteString(f.toInputsMarkdown(actionMarkdown.Inputs))
		sb.WriteString("\n\n")
	}

	if len(actionMarkdown.Outputs) != 0 || !config.Omit {
		sb.WriteString(f.toOutputsMarkdown(actionMarkdown.Outputs))
		sb.WriteString("\n\n")
	}
	return strings.TrimSpace(sb.String())
}

func (f *ActionFormatter) convertActionMarkdown(ast *parse.ActionAST) *ActionMarkdown {
	inputs := []*ActionInputMarkdown{}
	for _, inputAst := range ast.Inputs {
		input := &ActionInputMarkdown{
			Name:        inputAst.Name,
			Default:     inputAst.Default,
			Description: inputAst.Description,
			Required:    inputAst.Required,
		}
		inputs = append(inputs, input)
	}

	outputs := []*ActionOutputMarkdown{}
	for _, outputAst := range ast.Outputs {
		output := &ActionOutputMarkdown{
			Name:        outputAst.Name,
			Description: outputAst.Description,
		}
		outputs = append(outputs, output)
	}

	return &ActionMarkdown{
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

func (f *ActionFormatter) toInputsMarkdown(inputs []*ActionInputMarkdown) string {
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

func (f *ActionFormatter) toOutputsMarkdown(outputs []*ActionOutputMarkdown) string {
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

type ActionMarkdown struct {
	Description *util.NullString        `json:"description"`
	Inputs      []*ActionInputMarkdown  `json:"inputs"`
	Outputs     []*ActionOutputMarkdown `json:"outputs"`
}

type ActionInputMarkdown struct {
	Name        string           `json:"name"`
	Default     *util.NullString `json:"default"`
	Description *util.NullString `json:"description"`
	Required    *util.NullString `json:"required"`
}

func (i *ActionInputMarkdown) toMarkdown() string {
	str := util.TableSeparator
	str += fmt.Sprintf(" %s %s", i.Name, util.TableSeparator)
	str += fmt.Sprintf(" %s %s", i.Description.StringOrEmpty(), util.TableSeparator)
	str += fmt.Sprintf(" %s %s", i.Default.QuoteStringOrLowerNA(), util.TableSeparator)
	str += fmt.Sprintf(" %s %s", i.Required.YesOrNo(), util.TableSeparator)
	return str
}

type ActionOutputMarkdown struct {
	Name        string           `json:"name"`
	Description *util.NullString `json:"description"`
}

func (o *ActionOutputMarkdown) toMarkdown() string {
	str := util.TableSeparator
	str += fmt.Sprintf(" %s %s", o.Name, util.TableSeparator)
	str += fmt.Sprintf(" %s %s", o.Description.StringOrEmpty(), util.TableSeparator)
	return str
}
