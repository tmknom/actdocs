package parse

type ActionYaml struct {
	Name        *string                      `yaml:"name"`
	Description *string                      `yaml:"description"`
	Inputs      map[string]*actionInputYaml  `yaml:"inputs"`
	Outputs     map[string]*actionOutputYaml `yaml:"outputs"`
	Runs        *actionRunsYaml              `yaml:"runs"`
}

type actionInputYaml struct {
	Default     *string `mapstructure:"default"`
	Description *string `mapstructure:"description"`
	Required    *string `mapstructure:"required"`
}

type actionOutputYaml struct {
	Description *string `mapstructure:"description"`
}

type actionRunsYaml struct {
	Using string         `yaml:"using"`
	Steps []*interface{} `yaml:"steps"`
}

func (y *ActionYaml) inputs() map[string]*actionInputYaml {
	if y.Inputs == nil {
		return map[string]*actionInputYaml{}
	}
	return y.Inputs
}

func (y *ActionYaml) outputs() map[string]*actionOutputYaml {
	if y.Outputs == nil {
		return map[string]*actionOutputYaml{}
	}
	return y.Outputs
}

type WorkflowYaml struct {
	On          *workflowOnYaml `yaml:"on"`
	Permissions interface{}     `yaml:"permissions"`
}

type workflowOnYaml struct {
	WorkflowCall *workflowWorkflowCallYaml `yaml:"workflow_call"`
}

type workflowWorkflowCallYaml struct {
	Inputs  map[string]*workflowInputYaml  `yaml:"inputs"`
	Secrets map[string]*workflowSecretYaml `yaml:"secrets"`
	Outputs map[string]*workflowOutputYaml `yaml:"outputs"`
}

type workflowInputYaml struct {
	Default     *string `mapstructure:"default"`
	Description *string `mapstructure:"description"`
	Required    *string `mapstructure:"required"`
	Type        *string `mapstructure:"type"`
}

type workflowSecretYaml struct {
	Description *string `mapstructure:"description"`
	Required    *string `mapstructure:"required"`
}

type workflowOutputYaml struct {
	Description *string `mapstructure:"description"`
}

func (y *WorkflowYaml) inputs() map[string]*workflowInputYaml {
	if y.On == nil || y.On.WorkflowCall == nil || y.On.WorkflowCall.Inputs == nil {
		return map[string]*workflowInputYaml{}
	}
	return y.On.WorkflowCall.Inputs
}

func (y *WorkflowYaml) secrets() map[string]*workflowSecretYaml {
	if y.On == nil || y.On.WorkflowCall == nil || y.On.WorkflowCall.Secrets == nil {
		return map[string]*workflowSecretYaml{}
	}
	return y.On.WorkflowCall.Secrets
}

func (y *WorkflowYaml) outputs() map[string]*workflowOutputYaml {
	if y.On == nil || y.On.WorkflowCall == nil || y.On.WorkflowCall.Outputs == nil {
		return map[string]*workflowOutputYaml{}
	}
	return y.On.WorkflowCall.Outputs
}

func (y *WorkflowYaml) permissions() map[interface{}]interface{} {
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
