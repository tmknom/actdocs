package format

import (
	"fmt"
	"strings"

	"github.com/tmknom/actdocs/internal/util"
)

type ActionSpec struct {
	Description *util.NullString    `json:"description"`
	Inputs      []*ActionInputSpec  `json:"inputs"`
	Outputs     []*ActionOutputSpec `json:"outputs"`
}

func (s *ActionSpec) toDescriptionMarkdown() string {
	var sb strings.Builder
	sb.WriteString(ActionDescriptionTitle)
	sb.WriteString("\n\n")
	sb.WriteString(s.Description.StringOrUpperNA())
	return strings.TrimSpace(sb.String())
}

func (s *ActionSpec) toInputsMarkdown() string {
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
	return strings.TrimSpace(sb.String())
}

func (s *ActionSpec) toOutputsMarkdown() string {
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
	return strings.TrimSpace(sb.String())
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
