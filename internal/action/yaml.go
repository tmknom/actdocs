package action

type Yaml struct {
	Name        *string                `yaml:"name"`
	Description *string                `yaml:"description"`
	Inputs      map[string]*InputYaml  `yaml:"inputs"`
	Outputs     map[string]*OutputYaml `yaml:"outputs"`
	Runs        *RunsYaml              `yaml:"runs"`
}

func NewYaml() *Yaml {
	return &Yaml{
		Inputs:  map[string]*InputYaml{},
		Outputs: map[string]*OutputYaml{},
	}
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
