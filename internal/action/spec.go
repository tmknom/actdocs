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

	Omit bool `json:"-"`
}

func (s *Spec) ToJson() string {
	bytes, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(bytes)
}

func (s *Spec) ToMarkdown() string {
	var sb strings.Builder
	sb.WriteString(s.ToDescriptionMarkdown())
	sb.WriteString("\n\n")
	sb.WriteString(s.ToInputsMarkdown())
	sb.WriteString("\n\n")
	sb.WriteString(s.ToOutputsMarkdown())
	return strings.TrimSpace(sb.String())
}

func (s *Spec) ToDescriptionMarkdown() string {
	if s.Omit && !s.Description.IsValid() {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(DescriptionTitle)
	sb.WriteString("\n\n")
	sb.WriteString(strings.TrimSpace(s.Description.StringOrUpperNA()))
	return sb.String()
}

func (s *Spec) ToInputsMarkdown() string {
	if s.Omit && len(s.Inputs) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(InputsTitle)
	sb.WriteString("\n\n")
	if len(s.Inputs) != 0 {
		sb.WriteString(InputsColumnTitle)
		sb.WriteString("\n")
		sb.WriteString(InputsColumnSeparator)
		sb.WriteString("\n")
		for _, input := range s.Inputs {
			sb.WriteString(input.toMarkdown())
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString(util.UpperNAString)
	}
	return strings.TrimSpace(sb.String())
}

func (s *Spec) ToOutputsMarkdown() string {
	if s.Omit && len(s.Outputs) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(OutputsTitle)
	sb.WriteString("\n\n")
	if len(s.Outputs) != 0 {
		sb.WriteString(OutputsColumnTitle)
		sb.WriteString("\n")
		sb.WriteString(OutputsColumnSeparator)
		sb.WriteString("\n")
		for _, output := range s.Outputs {
			sb.WriteString(output.toMarkdown())
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString(util.UpperNAString)
	}
	return strings.TrimSpace(sb.String())
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
	DescriptionTitle = "## Description"

	InputsTitle           = "## Inputs"
	InputsColumnTitle     = "| Name | Description | Default | Required |"
	InputsColumnSeparator = "| :--- | :---------- | :------ | :------: |"

	OutputsTitle           = "## Outputs"
	OutputsColumnTitle     = "| Name | Description |"
	OutputsColumnSeparator = "| :--- | :---------- |"
)
