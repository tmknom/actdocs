package workflow

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
