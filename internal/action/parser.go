package action

import (
	"log"
	"sort"

	"github.com/tmknom/actdocs/internal/conf"
	"github.com/tmknom/actdocs/internal/util"
	"gopkg.in/yaml.v2"
)

type Parser struct {
	*AST
	*conf.SortConfig
}

func NewParser(sort *conf.SortConfig) *Parser {
	return &Parser{
		AST: &AST{
			Inputs:  []*InputAST{},
			Outputs: []*OutputAST{},
		},
		SortConfig: sort,
	}
}

func (p *Parser) Parse(yamlBytes []byte) (*AST, error) {
	actionYaml := NewYaml()
	err := yaml.Unmarshal(yamlBytes, actionYaml)
	if err != nil {
		return nil, err
	}
	log.Printf("unmarshal yaml: actionYaml = %#v\n", actionYaml)

	p.Name = util.NewNullString(actionYaml.Name)
	p.Description = util.NewNullString(actionYaml.Description)
	p.Runs = NewRunsAST(actionYaml.Runs)

	for name, element := range actionYaml.Inputs {
		p.parseInput(name, element)
	}

	for name, element := range actionYaml.Outputs {
		p.parseOutput(name, element)
	}

	p.sort()
	return p.AST, nil
}

func (p *Parser) sort() {
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

func (p *Parser) sortInputs() {
	log.Printf("sorted: inputs")

	//goland:noinspection GoPreferNilSlice
	required := []*InputAST{}
	//goland:noinspection GoPreferNilSlice
	notRequired := []*InputAST{}
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

func (p *Parser) sortInputsByName() {
	log.Printf("sorted: inputs by name")
	item := p.Inputs
	sort.Slice(item, func(i, j int) bool {
		return item[i].Name < item[j].Name
	})
}

func (p *Parser) sortInputsByRequired() {
	log.Printf("sorted: inputs by required")
	item := p.Inputs
	sort.Slice(item, func(i, j int) bool {
		return item[i].Required.IsTrue()
	})
}

func (p *Parser) sortOutputsByName() {
	log.Printf("sorted: outputs by name")
	item := p.Outputs
	sort.Slice(item, func(i, j int) bool {
		return item[i].Name < item[j].Name
	})
}

func (p *Parser) parseInput(name string, element *InputYaml) {
	result := NewInputAST(name)
	if element != nil {
		result.Default = util.NewNullString(element.Default)
		result.Description = util.NewNullString(element.Description)
		result.Required = util.NewNullString(element.Required)
	}
	p.Inputs = append(p.Inputs, result)
}

func (p *Parser) parseOutput(name string, element *OutputYaml) {
	result := NewOutputAST(name)
	if element != nil {
		result.Description = util.NewNullString(element.Description)
	}
	p.Outputs = append(p.Outputs, result)
}
