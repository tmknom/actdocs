package actdocs

import (
	"fmt"
)

type Generator struct {
	*GeneratorConfig
	*IO
	*YamlFile
}

func NewGenerator(config *GeneratorConfig, inOut *IO, yamlFile *YamlFile) *Generator {
	return &Generator{
		GeneratorConfig: config,
		IO:              inOut,
		YamlFile:        yamlFile,
	}
}

type GeneratorConfig struct {
	*GlobalConfig
}

func NewGeneratorConfig(globalConfig *GlobalConfig) *GeneratorConfig {
	return &GeneratorConfig{
		GlobalConfig: globalConfig,
	}
}

func (c *Generator) Run() error {
	content, err := c.FormatYaml(c.GlobalConfig)
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(c.OutWriter, content)
	return err
}
