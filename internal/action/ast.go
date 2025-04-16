package action

import (
	"fmt"

	"github.com/tmknom/actdocs/internal/util"
)

type AST struct {
	Name        *util.NullString
	Description *util.NullString
	Inputs      []*InputAST
	Outputs     []*OutputAST
	Runs        *RunsAST
}

type InputAST struct {
	Name        string
	Default     *util.NullString
	Description *util.NullString
	Required    *util.NullString
}

func NewInputAST(name string) *InputAST {
	return &InputAST{
		Name:        name,
		Default:     util.DefaultNullString,
		Description: util.DefaultNullString,
		Required:    util.DefaultNullString,
	}
}

type OutputAST struct {
	Name        string
	Description *util.NullString
}

func NewOutputAST(name string) *OutputAST {
	return &OutputAST{
		Name:        name,
		Description: util.DefaultNullString,
	}
}

type RunsAST struct {
	Using string
	Steps []*interface{}
}

func NewRunsAST(runs *RunsYaml) *RunsAST {
	result := &RunsAST{
		Using: "undefined",
		Steps: []*interface{}{},
	}

	if runs != nil {
		result.Using = runs.Using
		result.Steps = runs.Steps
	}
	return result
}

func (r *RunsAST) String() string {
	str := ""
	str += fmt.Sprintf("Using: %s, ", r.Using)
	str += fmt.Sprintf("Steps: [")
	for _, step := range r.Steps {
		str += fmt.Sprintf("%#v, ", *step)
	}
	str += fmt.Sprintf("]")
	return str
}
