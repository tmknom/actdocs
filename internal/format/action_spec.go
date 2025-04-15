package format

import (
	"fmt"

	"github.com/tmknom/actdocs/internal/util"
)

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
