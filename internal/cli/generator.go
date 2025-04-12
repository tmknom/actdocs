package cli

import (
	"fmt"
	"log"

	"github.com/tmknom/actdocs/internal/format"
	"github.com/tmknom/actdocs/internal/parse"
	"github.com/tmknom/actdocs/internal/read"
)

type GenerateRunner struct {
	*GeneratorConfig
	*IO
	YamlFile string
}

func NewGenerateRunner(config *GeneratorConfig, inOut *IO, yamlFile string) *GenerateRunner {
	return &GenerateRunner{
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

func (r *GenerateRunner) Run() error {
	reader := &read.YamlReader{Filename: r.YamlFile}
	yaml, err := reader.Read()
	if err != nil {
		return err
	}
	log.Printf("read: %s", r.YamlFile)

	factory := &parse.ParserFactory{Raw: yaml}
	parser, err := factory.Factory(r.FormatterConfig)
	if err != nil {
		return err
	}
	log.Printf("selected parser: %T", parser)

	content, err := parser.Parse()
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(r.OutWriter, content)
	return err
}
