package format

import (
	"fmt"

	"github.com/tmknom/actdocs/internal/util"
)

type WorkflowSpec struct {
	Inputs      []*WorkflowInputSpec      `json:"inputs"`
	Secrets     []*WorkflowSecretSpec     `json:"secrets"`
	Outputs     []*WorkflowOutputSpec     `json:"outputs"`
	Permissions []*WorkflowPermissionSpec `json:"permissions"`
}

type WorkflowInputSpec struct {
	Name        string           `json:"name"`
	Default     *util.NullString `json:"default"`
	Description *util.NullString `json:"description"`
	Required    *util.NullString `json:"required"`
	Type        *util.NullString `json:"type"`
}

func (s *WorkflowInputSpec) toMarkdown() string {
	str := util.TableSeparator
	str += fmt.Sprintf(" %s %s", s.Name, util.TableSeparator)
	str += fmt.Sprintf(" %s %s", s.Description.StringOrEmpty(), util.TableSeparator)
	str += fmt.Sprintf(" %s %s", s.Type.QuoteStringOrLowerNA(), util.TableSeparator)
	str += fmt.Sprintf(" %s %s", s.Default.QuoteStringOrLowerNA(), util.TableSeparator)
	str += fmt.Sprintf(" %s %s", s.Required.YesOrNo(), util.TableSeparator)
	return str
}

type WorkflowSecretSpec struct {
	Name        string           `json:"name"`
	Description *util.NullString `json:"description"`
	Required    *util.NullString `json:"required"`
}

func (s *WorkflowSecretSpec) toMarkdown() string {
	str := util.TableSeparator
	str += fmt.Sprintf(" %s %s", s.Name, util.TableSeparator)
	str += fmt.Sprintf(" %s %s", s.Description.StringOrEmpty(), util.TableSeparator)
	str += fmt.Sprintf(" %s %s", s.Required.YesOrNo(), util.TableSeparator)
	return str
}

type WorkflowOutputSpec struct {
	Name        string           `json:"name"`
	Description *util.NullString `json:"description"`
}

func (s *WorkflowOutputSpec) toMarkdown() string {
	str := util.TableSeparator
	str += fmt.Sprintf(" %s %s", s.Name, util.TableSeparator)
	str += fmt.Sprintf(" %s %s", s.Description.StringOrEmpty(), util.TableSeparator)
	return str
}

type WorkflowPermissionSpec struct {
	Scope  string `json:"scope"`
	Access string `json:"access"`
}

func (s *WorkflowPermissionSpec) toMarkdown() string {
	str := util.TableSeparator
	str += fmt.Sprintf(" %s %s", s.Scope, util.TableSeparator)
	str += fmt.Sprintf(" %s %s", s.Access, util.TableSeparator)
	return str
}
