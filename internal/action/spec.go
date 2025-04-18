package action

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tmknom/actdocs/internal/util"
)

type Spec struct {
	Description *util.NullString `json:"description"`
	Inputs      []*InputSpec     `json:"inputs"`
	Outputs     []*OutputSpec    `json:"outputs"`
}

func (s *Spec) ToJson() string {
	bytes, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(bytes)
}

func (s *Spec) ToMarkdown(omit bool) string {
	var sb strings.Builder
	sb.WriteString(s.ToDescriptionMarkdown(omit))
	sb.WriteString(s.ToInputsMarkdown(omit))
	sb.WriteString(s.ToOutputsMarkdown(omit))
	return strings.TrimSpace(sb.String())
}

func (s *Spec) ToDescriptionMarkdown(omit bool) string {
	if omit && !s.Description.IsValid() {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(ActionDescriptionTitle)
	sb.WriteString("\n\n")
	sb.WriteString(strings.TrimSpace(s.Description.StringOrUpperNA()))
	sb.WriteString("\n\n")
	return sb.String()
}

func (s *Spec) ToInputsMarkdown(omit bool) string {
	if omit && len(s.Inputs) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(ActionInputsTitle)
	sb.WriteString("\n\n")
	if len(s.Inputs) != 0 {
		sb.WriteString(ActionInputsColumnTitle)
		sb.WriteString("\n")
		sb.WriteString(ActionInputsColumnSeparator)
		sb.WriteString("\n")
		for _, input := range s.Inputs {
			sb.WriteString(input.toMarkdown())
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString(util.UpperNAString)
	}
	return strings.TrimSpace(sb.String()) + "\n\n"
}

func (s *Spec) ToOutputsMarkdown(omit bool) string {
	if omit && len(s.Outputs) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(ActionOutputsTitle)
	sb.WriteString("\n\n")
	if len(s.Outputs) != 0 {
		sb.WriteString(ActionOutputsColumnTitle)
		sb.WriteString("\n")
		sb.WriteString(ActionOutputsColumnSeparator)
		sb.WriteString("\n")
		for _, output := range s.Outputs {
			sb.WriteString(output.toMarkdown())
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString(util.UpperNAString)
	}
	return strings.TrimSpace(sb.String()) + "\n\n"
}

type InputSpec struct {
	Name        string           `json:"name"`
	Default     *util.NullString `json:"default"`
	Description *util.NullString `json:"description"`
	Required    *util.NullString `json:"required"`
}

func (s *InputSpec) toMarkdown() string {
	str := util.TableSeparator
	str += fmt.Sprintf(" %s %s", s.Name, util.TableSeparator)
	str += fmt.Sprintf(" %s %s", s.Description.StringOrEmpty(), util.TableSeparator)
	str += fmt.Sprintf(" %s %s", s.Default.QuoteStringOrLowerNA(), util.TableSeparator)
	str += fmt.Sprintf(" %s %s", s.Required.YesOrNo(), util.TableSeparator)
	return str
}

type OutputSpec struct {
	Name        string           `json:"name"`
	Description *util.NullString `json:"description"`
}

func (s *OutputSpec) toMarkdown() string {
	str := util.TableSeparator
	str += fmt.Sprintf(" %s %s", s.Name, util.TableSeparator)
	str += fmt.Sprintf(" %s %s", s.Description.StringOrEmpty(), util.TableSeparator)
	return str
}

const (
	ActionDescriptionTitle = "## Description"

	ActionInputsTitle           = "## Inputs"
	ActionInputsColumnTitle     = "| Name | Description | Default | Required |"
	ActionInputsColumnSeparator = "| :--- | :---------- | :------ | :------: |"

	ActionOutputsTitle           = "## Outputs"
	ActionOutputsColumnTitle     = "| Name | Description |"
	ActionOutputsColumnSeparator = "| :--- | :---------- |"
)
