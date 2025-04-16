package action

type Yaml struct {
	Name        *string                `yaml:"name"`
	Description *string                `yaml:"description"`
	Inputs      map[string]*InputYaml  `yaml:"inputs"`
	Outputs     map[string]*OutputYaml `yaml:"outputs"`
	Runs        *RunsYaml              `yaml:"runs"`
}

type InputYaml struct {
	Default     *string `mapstructure:"default"`
	Description *string `mapstructure:"description"`
	Required    *string `mapstructure:"required"`
}

type OutputYaml struct {
	Description *string `mapstructure:"description"`
}

type RunsYaml struct {
	Using string         `yaml:"using"`
	Steps []*interface{} `yaml:"steps"`
}

func (y *Yaml) ActionInputs() map[string]*InputYaml {
	if y.Inputs == nil {
		return map[string]*InputYaml{}
	}
	return y.Inputs
}

func (y *Yaml) ActionOutputs() map[string]*OutputYaml {
	if y.Outputs == nil {
		return map[string]*OutputYaml{}
	}
	return y.Outputs
}
