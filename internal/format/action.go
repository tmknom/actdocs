package format

import (
	"encoding/json"
	"strings"

	"github.com/tmknom/actdocs/internal/conf"
	"github.com/tmknom/actdocs/internal/parse"
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
	f.ActionSpec = ConvertActionSpec(ast)

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
		sb.WriteString(actionSpec.toDescriptionMarkdown())
		sb.WriteString("\n\n")
	}

	if len(actionSpec.Inputs) != 0 || !config.Omit {
		sb.WriteString(actionSpec.toInputsMarkdown())
		sb.WriteString("\n\n")
	}

	if len(actionSpec.Outputs) != 0 || !config.Omit {
		sb.WriteString(actionSpec.toOutputsMarkdown())
		sb.WriteString("\n\n")
	}
	return strings.TrimSpace(sb.String())
}

func ConvertActionSpec(ast *parse.ActionAST) *ActionSpec {
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

func (f *ActionFormatter) toDescriptionMarkdown(actionSpec *ActionSpec) string {
	return actionSpec.toDescriptionMarkdown()
}

func (f *ActionFormatter) toInputsMarkdown(actionSpec *ActionSpec) string {
	return actionSpec.toInputsMarkdown()
}

func (f *ActionFormatter) toOutputsMarkdown(actionSpec *ActionSpec) string {
	return actionSpec.toOutputsMarkdown()
}

const ActionDescriptionTitle = "## Description"

const ActionInputsTitle = "## Inputs"
const ActionInputsColumnTitle = "| Name | Description | Default | Required |"
const ActionInputsColumnSeparator = "| :--- | :---------- | :------ | :------: |"

const ActionOutputsTitle = "## Outputs"
const ActionOutputsColumnTitle = "| Name | Description |"
const ActionOutputsColumnSeparator = "| :--- | :---------- |"
