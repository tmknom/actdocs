package cli

import (
	"fmt"
	"log"

	"github.com/tmknom/actdocs/internal/config"
	"github.com/tmknom/actdocs/internal/parse"
	"github.com/tmknom/actdocs/internal/read"
)

type Generator struct {
	*GeneratorConfig
	*config.IO
	YamlFile string
}

func NewGenerator(config *GeneratorConfig, inOut *config.IO, yamlFile string) *Generator {
	return &Generator{
		GeneratorConfig: config,
		IO:              inOut,
		YamlFile:        yamlFile,
	}
}

type GeneratorConfig struct {
	*config.GlobalConfig
}

func NewGeneratorConfig(globalConfig *config.GlobalConfig) *GeneratorConfig {
	return &GeneratorConfig{
		GlobalConfig: globalConfig,
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
	parser, err := factory.Factory(c.GlobalConfig)
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
