package actdocs

import (
	"fmt"
	"io"
	"log"
)

type GenerateCmd struct {
	*TemplateConfig
	filename string
	// inReader is a reader defined by the user that replaces stdin
	inReader io.Reader
	// outWriter is a writer defined by the user that replaces stdout
	outWriter io.Writer
	// errWriter is a writer defined by the user that replaces stderr
	errWriter io.Writer
}

func NewGenerateCmd(config *TemplateConfig, inReader io.Reader, outWriter, errWriter io.Writer) *GenerateCmd {
	return &GenerateCmd{
		TemplateConfig: config,
		inReader:       inReader,
		outWriter:      outWriter,
		errWriter:      errWriter,
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
		generator = NewWorkflow(rawYaml)
	} else if rawYaml.IsCustomActions() {
		generator = NewAction(rawYaml)
	} else {
		return "", fmt.Errorf("invalid file: %s", c.filename)
	}
	log.Printf("selected generator: %T", generator)
	return generator.Generate()
}
