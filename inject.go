package actdocs

import (
	"fmt"
	"io"
	"log"
)

type InjectCmd struct {
	*Config
	filename string
	// inReader is a reader defined by the user that replaces stdin
	inReader io.Reader
	// outWriter is a writer defined by the user that replaces stdout
	outWriter io.Writer
	// errWriter is a writer defined by the user that replaces stderr
	errWriter io.Writer
}

func NewInjectCmd(config *Config, inReader io.Reader, outWriter, errWriter io.Writer) *InjectCmd {
	return &InjectCmd{
		Config:    config,
		inReader:  inReader,
		outWriter: outWriter,
		errWriter: errWriter,
	}
}

func (c *InjectCmd) Run() error {
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

func (c *InjectCmd) generate(rawYaml rawYaml) (string, error) {
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
