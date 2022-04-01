package actdocs

import (
	"fmt"
	"io"
	"log"
)

type Config struct {
	*TemplateConfig
	*GeneratorConfig
}

func NewConfig(outWriter io.Writer) *Config {
	return &Config{
		TemplateConfig:  NewTemplateConfig(outWriter),
		GeneratorConfig: NewGeneratorConfig(),
	}
}

type GenerateCmd struct {
	*Config
	filename string
	// inReader is a reader defined by the user that replaces stdin
	inReader io.Reader
	// outWriter is a writer defined by the user that replaces stdout
	outWriter io.Writer
	// errWriter is a writer defined by the user that replaces stderr
	errWriter io.Writer
}

func NewGenerateCmd(config *Config, inReader io.Reader, outWriter, errWriter io.Writer) *GenerateCmd {
	return &GenerateCmd{
		Config:    config,
		inReader:  inReader,
		outWriter: outWriter,
		errWriter: errWriter,
	}
}

func (c *GenerateCmd) Run() error {
	log.Printf("read: %v", c.filename)
	rawYaml, err := readYaml(c.filename)
	if err != nil {
		return err
	}

	content, err := c.generate(rawYaml)
	if err != nil {
		return err
	}

	template := NewTemplate(c.TemplateConfig)
	return template.Render(content)
}

func (c *GenerateCmd) generate(rawYaml rawYaml) (string, error) {
	var generator Generator
	if rawYaml.IsReusableWorkflow() {
		generator = NewWorkflow(rawYaml, c.GeneratorConfig)
	} else if rawYaml.IsCustomActions() {
		generator = NewAction(rawYaml, c.GeneratorConfig)
	} else {
		return "", fmt.Errorf("invalid file: %s", c.filename)
	}
	log.Printf("selected generator: %T", generator)
	return generator.Generate()
}

type GeneratorConfig struct {
	Sort           bool
	SortByName     bool
	SortByRequired bool
}

func NewGeneratorConfig() *GeneratorConfig {
	return &GeneratorConfig{}
}
