package action

type ActionYaml struct {
	Name        *string                      `yaml:"name"`
	Description *string                      `yaml:"description"`
	Inputs      map[string]*ActionInputYaml  `yaml:"inputs"`
	Outputs     map[string]*ActionOutputYaml `yaml:"outputs"`
	Runs        *ActionRunsYaml              `yaml:"runs"`
}

type ActionInputYaml struct {
	Default     *string `mapstructure:"default"`
	Description *string `mapstructure:"description"`
	Required    *string `mapstructure:"required"`
}

type ActionOutputYaml struct {
	Description *string `mapstructure:"description"`
}

type ActionRunsYaml struct {
	Using string         `yaml:"using"`
	Steps []*interface{} `yaml:"steps"`
}

func (y *ActionYaml) ActionInputs() map[string]*ActionInputYaml {
	if y.Inputs == nil {
		return map[string]*ActionInputYaml{}
	}
	return y.Inputs
}

func (y *ActionYaml) ActionOutputs() map[string]*ActionOutputYaml {
	if y.Outputs == nil {
		return map[string]*ActionOutputYaml{}
	}
	return y.Outputs
}
