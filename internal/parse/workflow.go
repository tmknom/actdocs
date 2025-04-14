package parse

import (
	"github.com/tmknom/actdocs/internal/conf"
	"github.com/tmknom/actdocs/internal/util"
	"gopkg.in/yaml.v2"
	"log"
	"sort"
)

type WorkflowParser struct {
	*WorkflowAST
	config *conf.FormatterConfig
	*conf.SortConfig
}

func NewWorkflowParser(config *conf.FormatterConfig, sort *conf.SortConfig) *WorkflowParser {
	return &WorkflowParser{
		WorkflowAST: &WorkflowAST{
			Inputs:      []*WorkflowInput{},
			Secrets:     []*WorkflowSecret{},
			Outputs:     []*WorkflowOutput{},
			Permissions: []*WorkflowPermission{},
		},
		config:     config,
		SortConfig: sort,
	}
}

type WorkflowAST struct {
	Inputs      []*WorkflowInput
	Secrets     []*WorkflowSecret
	Outputs     []*WorkflowOutput
	Permissions []*WorkflowPermission
}

type WorkflowInput struct {
	Name        string
	Default     *util.NullString
	Description *util.NullString
	Required    *util.NullString
	Type        *util.NullString
}

func NewWorkflowInput(name string) *WorkflowInput {
	return &WorkflowInput{
		Name:        name,
		Default:     util.DefaultNullString,
		Description: util.DefaultNullString,
		Required:    util.DefaultNullString,
		Type:        util.DefaultNullString,
	}
}

type WorkflowSecret struct {
	Name        string
	Description *util.NullString
	Required    *util.NullString
}

func NewWorkflowSecret(name string) *WorkflowSecret {
	return &WorkflowSecret{
		Name:        name,
		Description: util.DefaultNullString,
		Required:    util.DefaultNullString,
	}
}

type WorkflowOutput struct {
	Name        string
	Description *util.NullString
}

func NewWorkflowOutput(name string) *WorkflowOutput {
	return &WorkflowOutput{
		Name:        name,
		Description: util.DefaultNullString,
	}
}

type WorkflowPermission struct {
	Scope  string
	Access string
}

func NewWorkflowPermission(scope string, access string) *WorkflowPermission {
	return &WorkflowPermission{
		Scope:  scope,
		Access: access,
	}
}

func (p *WorkflowParser) Parse(yamlBytes []byte) (string, error) {
	ast, err := p.ParseAST(yamlBytes)
	if err != nil {
		return "", err
	}

	formatter := NewWorkflowFormatter(p.config)
	return formatter.Format(ast), nil
}

func (p *WorkflowParser) ParseAST(yamlBytes []byte) (*WorkflowAST, error) {
	log.Printf("config: %#v", p.config)

	content := &WorkflowYaml{}
	err := yaml.Unmarshal(yamlBytes, content)
	if err != nil {
		return nil, err
	}

	for name, value := range content.inputs() {
		input := p.parseInput(name, value)
		p.Inputs = append(p.Inputs, input)
	}

	for name, value := range content.outputs() {
		output := p.parseOutput(name, value)
		p.Outputs = append(p.Outputs, output)
	}

	for name, value := range content.secrets() {
		secret := p.parseSecret(name, value)
		p.Secrets = append(p.Secrets, secret)
	}

	for scope, access := range content.permissions() {
		permission := NewWorkflowPermission(scope.(string), access.(string))
		p.Permissions = append(p.Permissions, permission)
	}

	p.sort()
	return p.WorkflowAST, nil
}

func (p *WorkflowParser) sort() {
	switch {
	case p.SortConfig.Sort:
		p.sortInputs()
		p.sortSecrets()
		p.sortOutputsByName()
		p.sortPermissionsByScope()
	case p.SortConfig.SortByName:
		p.sortInputsByName()
		p.sortSecretsByName()
		p.sortOutputsByName()
		p.sortPermissionsByScope()
	case p.SortConfig.SortByRequired:
		p.sortInputsByRequired()
		p.sortSecretByRequired()
	}
}

func (p *WorkflowParser) sortInputs() {
	log.Printf("sorted: inputs")

	//goland:noinspection GoPreferNilSlice
	required := []*WorkflowInput{}
	//goland:noinspection GoPreferNilSlice
	notRequired := []*WorkflowInput{}
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

func (p *WorkflowParser) sortInputsByName() {
	log.Printf("sorted: inputs by name")
	item := p.Inputs
	sort.Slice(item, func(i, j int) bool {
		return item[i].Name < item[j].Name
	})
}

func (p *WorkflowParser) sortInputsByRequired() {
	log.Printf("sorted: inputs by required")
	item := p.Inputs
	sort.Slice(item, func(i, j int) bool {
		return item[i].Required.IsTrue()
	})
}

func (p *WorkflowParser) sortSecrets() {
	log.Printf("sorted: secrets")

	//goland:noinspection GoPreferNilSlice
	required := []*WorkflowSecret{}
	//goland:noinspection GoPreferNilSlice
	notRequired := []*WorkflowSecret{}
	for _, input := range p.Secrets {
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
	p.Secrets = append(required, notRequired...)
}

func (p *WorkflowParser) sortSecretsByName() {
	log.Printf("sorted: secrets by name")
	item := p.Secrets
	sort.Slice(item, func(i, j int) bool {
		return item[i].Name < item[j].Name
	})
}

func (p *WorkflowParser) sortSecretByRequired() {
	log.Printf("sorted: secrets by required")
	item := p.Secrets
	sort.Slice(item, func(i, j int) bool {
		return item[i].Required.IsTrue()
	})
}

func (p *WorkflowParser) sortOutputsByName() {
	log.Printf("sorted: outputs by name")
	item := p.Outputs
	sort.Slice(item, func(i, j int) bool {
		return item[i].Name < item[j].Name
	})
}

func (p *WorkflowParser) sortPermissionsByScope() {
	log.Printf("sorted: permission by scope")
	item := p.Permissions
	sort.Slice(item, func(i, j int) bool {
		return item[i].Scope < item[j].Scope
	})
}

func (p *WorkflowParser) parseInput(name string, value *workflowInputYaml) *WorkflowInput {
	result := NewWorkflowInput(name)
	if value == nil {
		return result
	}

	result.Default = util.NewNullString(value.Default)
	result.Description = util.NewNullString(value.Description)
	result.Required = util.NewNullString(value.Required)
	result.Type = util.NewNullString(value.Type)

	return result
}

func (p *WorkflowParser) parseSecret(name string, value *workflowSecretYaml) *WorkflowSecret {
	result := NewWorkflowSecret(name)
	if value == nil {
		return result
	}

	result.Description = util.NewNullString(value.Description)
	result.Required = util.NewNullString(value.Required)

	return result
}

func (p *WorkflowParser) parseOutput(name string, value *workflowOutputYaml) *WorkflowOutput {
	result := NewWorkflowOutput(name)
	if value == nil {
		return result
	}

	result.Description = util.NewNullString(value.Description)
	return result
}
