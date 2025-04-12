package cli

import (
	"fmt"
	"log"

	"github.com/tmknom/actdocs/internal/format"
	"github.com/tmknom/actdocs/internal/parse"
	"github.com/tmknom/actdocs/internal/read"
)

type Generator struct {
	*GeneratorConfig
	*IO
	YamlFile string
}

func NewGenerator(config *GeneratorConfig, inOut *IO, yamlFile string) *Generator {
	return &Generator{
		GeneratorConfig: config,
		IO:              inOut,
		YamlFile:        yamlFile,
	}
}

type GeneratorConfig struct {
	*format.FormatterConfig
}

func NewGeneratorConfig(config *format.FormatterConfig) *GeneratorConfig {
	return &GeneratorConfig{
		FormatterConfig: config,
	}
}

func (c *Generator) Run() error {
	reader := &read.YamlReader{Filename: c.YamlFile}
	yaml, err := reader.Read()
	if err != nil {
		return err
	}
	log.Printf("read: %s", c.YamlFile)

	factory := &parse.ParserFactory{Raw: yaml}
	parser, err := factory.Factory(c.FormatterConfig)
	if err != nil {
		return err
	}
	log.Printf("selected parser: %T", parser)

	content, err := parser.Parse()
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(c.OutWriter, content)
	return err
}
