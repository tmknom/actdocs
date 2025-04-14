package parse

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"sort"

	"github.com/tmknom/actdocs/internal/conf"
	"github.com/tmknom/actdocs/internal/util"
)

type ActionParser struct {
	*ActionAST
	*conf.SortConfig
}

func NewActionParser(sort *conf.SortConfig) *ActionParser {
	return &ActionParser{
		ActionAST: &ActionAST{
			Inputs:  []*ActionInput{},
			Outputs: []*ActionOutput{},
		},
		SortConfig: sort,
	}
}

type ActionAST struct {
	Name        *util.NullString
	Description *util.NullString
	Inputs      []*ActionInput
	Outputs     []*ActionOutput
	Runs        *ActionRuns
}

func (a *ActionAST) AST() string {
	return "ActionAST"
}

func (p *ActionParser) ParseAST(yamlBytes []byte) (InterfaceAST, error) {
	actionYaml := &ActionYaml{}
	err := yaml.Unmarshal(yamlBytes, actionYaml)
	if err != nil {
		return nil, err
	}
	log.Printf("unmarshal yaml: actionYaml = %#v\n", actionYaml)

	p.Name = util.NewNullString(actionYaml.Name)
	p.Description = util.NewNullString(actionYaml.Description)
	p.Runs = NewActionRuns(actionYaml.Runs)

	for name, element := range actionYaml.inputs() {
		p.parseInput(name, element)
	}

	for name, element := range actionYaml.outputs() {
		p.parseOutput(name, element)
	}

	p.sort()
	return p.ActionAST, nil
}

func (p *ActionParser) sort() {
	switch {
	case p.SortConfig.Sort:
		p.sortInputs()
		p.sortOutputsByName()
	case p.SortConfig.SortByName:
		p.sortInputsByName()
		p.sortOutputsByName()
	case p.SortConfig.SortByRequired:
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
