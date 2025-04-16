package workflow

type Yaml struct {
	On          *OnYaml     `yaml:"on"`
	Permissions interface{} `yaml:"permissions"`
}

type OnYaml struct {
	WorkflowCall *WorkflowCallYaml `yaml:"workflow_call"`
}

type WorkflowCallYaml struct {
	Inputs  map[string]*InputYaml  `yaml:"inputs"`
	Secrets map[string]*SecretYaml `yaml:"secrets"`
	Outputs map[string]*OutputYaml `yaml:"outputs"`
}

type InputYaml struct {
	Default     *string `mapstructure:"default"`
	Description *string `mapstructure:"description"`
	Required    *string `mapstructure:"required"`
	Type        *string `mapstructure:"type"`
}

type SecretYaml struct {
	Description *string `mapstructure:"description"`
	Required    *string `mapstructure:"required"`
}

type OutputYaml struct {
	Description *string `mapstructure:"description"`
}

func (y *Yaml) WorkflowInputs() map[string]*InputYaml {
	if y.On == nil || y.On.WorkflowCall == nil || y.On.WorkflowCall.Inputs == nil {
		return map[string]*InputYaml{}
	}
	return y.On.WorkflowCall.Inputs
}

func (y *Yaml) WorkflowSecrets() map[string]*SecretYaml {
	if y.On == nil || y.On.WorkflowCall == nil || y.On.WorkflowCall.Secrets == nil {
		return map[string]*SecretYaml{}
	}
	return y.On.WorkflowCall.Secrets
}

func (y *Yaml) WorkflowOutputs() map[string]*OutputYaml {
	if y.On == nil || y.On.WorkflowCall == nil || y.On.WorkflowCall.Outputs == nil {
		return map[string]*OutputYaml{}
	}
	return y.On.WorkflowCall.Outputs
}

func (y *Yaml) WorkflowPermissions() map[interface{}]interface{} {
	if y.Permissions == nil {
		return map[interface{}]interface{}{}
	}

	switch y.Permissions.(type) {
	case string:
		access := y.Permissions.(string)
		if access == ReadAllAccess || access == WriteAllAccess {
			return map[interface{}]interface{}{AllScope: access}
		}
	case map[interface{}]interface{}:
		return y.Permissions.(map[interface{}]interface{})
	}
	return map[interface{}]interface{}{}
}

const ReadAllAccess = "read-all"
const WriteAllAccess = "write-all"
const AllScope = "-"
