package cli

import (
	"fmt"
	"log"

	"github.com/tmknom/actdocs/internal/format"
	"github.com/tmknom/actdocs/internal/parse"
	"github.com/tmknom/actdocs/internal/read"
)

type GenerateRunner struct {
	source string
	*format.FormatterConfig
	*IO
}

func NewGenerateRunner(source string, config *format.FormatterConfig, inOut *IO) *GenerateRunner {
	return &GenerateRunner{
		source:          source,
		FormatterConfig: config,
		IO:              inOut,
	}
}

func (r *GenerateRunner) Run() error {
	reader := &read.YamlReader{Filename: r.source}
	yaml, err := reader.Read()
	if err != nil {
		return err
	}
	log.Printf("read: %s", r.source)

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
