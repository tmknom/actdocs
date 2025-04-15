package parse

import (
	"log"
	"sort"

	"github.com/tmknom/actdocs/internal/conf"
	"github.com/tmknom/actdocs/internal/util"
	"gopkg.in/yaml.v2"
)

type WorkflowParser struct {
	*WorkflowAST
	*conf.SortConfig
}

func NewWorkflowParser(sort *conf.SortConfig) *WorkflowParser {
	return &WorkflowParser{
		WorkflowAST: &WorkflowAST{
			Inputs:      []*WorkflowInput{},
			Secrets:     []*WorkflowSecret{},
			Outputs:     []*WorkflowOutput{},
			Permissions: []*WorkflowPermission{},
		},
		SortConfig: sort,
	}
}

type WorkflowAST struct {
	Inputs      []*WorkflowInput
	Secrets     []*WorkflowSecret
	Outputs     []*WorkflowOutput
	Permissions []*WorkflowPermission
}

func (a *WorkflowAST) AST() string {
	return "ActionAST"
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

func (p *WorkflowParser) ParseAST(yamlBytes []byte) (util.InterfaceAST, error) {
	content := &WorkflowYaml{}
	err := yaml.Unmarshal(yamlBytes, content)
	if err != nil {
		return nil, err
	}

	for name, value := range content.WorkflowInputs() {
		input := p.parseInput(name, value)
		p.Inputs = append(p.Inputs, input)
	}

	for name, value := range content.WorkflowOutputs() {
		output := p.parseOutput(name, value)
		p.Outputs = append(p.Outputs, output)
	}

	for name, value := range content.WorkflowSecrets() {
		secret := p.parseSecret(name, value)
		p.Secrets = append(p.Secrets, secret)
	}

	for scope, access := range content.WorkflowPermissions() {
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

func (p *WorkflowParser) parseInput(name string, value *WorkflowInputYaml) *WorkflowInput {
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

func (p *WorkflowParser) parseSecret(name string, value *WorkflowSecretYaml) *WorkflowSecret {
	result := NewWorkflowSecret(name)
	if value == nil {
		return result
	}

	result.Description = util.NewNullString(value.Description)
	result.Required = util.NewNullString(value.Required)

	return result
}

func (p *WorkflowParser) parseOutput(name string, value *WorkflowOutputYaml) *WorkflowOutput {
	result := NewWorkflowOutput(name)
	if value == nil {
		return result
	}

	result.Description = util.NewNullString(value.Description)
	return result
}

type WorkflowYaml struct {
	On          *WorkflowOnYaml `yaml:"on"`
	Permissions interface{}     `yaml:"permissions"`
}

type WorkflowOnYaml struct {
	WorkflowCall *WorkflowWorkflowCallYaml `yaml:"workflow_call"`
}

type WorkflowWorkflowCallYaml struct {
	Inputs  map[string]*WorkflowInputYaml  `yaml:"inputs"`
	Secrets map[string]*WorkflowSecretYaml `yaml:"secrets"`
	Outputs map[string]*WorkflowOutputYaml `yaml:"outputs"`
}

type WorkflowInputYaml struct {
	Default     *string `mapstructure:"default"`
	Description *string `mapstructure:"description"`
	Required    *string `mapstructure:"required"`
	Type        *string `mapstructure:"type"`
}

type WorkflowSecretYaml struct {
	Description *string `mapstructure:"description"`
	Required    *string `mapstructure:"required"`
}

type WorkflowOutputYaml struct {
	Description *string `mapstructure:"description"`
}

func (y *WorkflowYaml) WorkflowInputs() map[string]*WorkflowInputYaml {
	if y.On == nil || y.On.WorkflowCall == nil || y.On.WorkflowCall.Inputs == nil {
		return map[string]*WorkflowInputYaml{}
	}
	return y.On.WorkflowCall.Inputs
}

func (y *WorkflowYaml) WorkflowSecrets() map[string]*WorkflowSecretYaml {
	if y.On == nil || y.On.WorkflowCall == nil || y.On.WorkflowCall.Secrets == nil {
		return map[string]*WorkflowSecretYaml{}
	}
	return y.On.WorkflowCall.Secrets
}

func (y *WorkflowYaml) WorkflowOutputs() map[string]*WorkflowOutputYaml {
	if y.On == nil || y.On.WorkflowCall == nil || y.On.WorkflowCall.Outputs == nil {
		return map[string]*WorkflowOutputYaml{}
	}
	return y.On.WorkflowCall.Outputs
}

func (y *WorkflowYaml) WorkflowPermissions() map[interface{}]interface{} {
	if y.Permissions == nil {
		return map[interface{}]interface{}{}
	}

	switch y.Permissions.(type) {
	case string:
		access := y.Permissions.(string)
		if access == util.ReadAllAccess || access == util.WriteAllAccess {
			return map[interface{}]interface{}{util.AllScope: access}
		}
	case map[interface{}]interface{}:
		return y.Permissions.(map[interface{}]interface{})
	}
	return map[interface{}]interface{}{}
}
