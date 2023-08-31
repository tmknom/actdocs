package actdocs

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
)

type YamlFile string

func NewYamlFile(args []string) *YamlFile {
	filename := ""
	if len(args) > 0 {
		filename = args[0]
	}
	result := YamlFile(filename)
	return &result
}

func (f *YamlFile) FormatYaml(globalConfig *GlobalConfig) (string, error) {
	rawYaml, err := f.read()
	if err != nil {
		return "", err
	}
	log.Printf("read: %s", *f)

	parser, err := rawYaml.createYamlParser(globalConfig)
	if err != nil {
		return "", err
	}
	log.Printf("selected parser: %T", parser)

	content, err := parser.Parse()
	if err != nil {
		return "", err
	}
	return content, nil
}

func (f *YamlFile) read() (rawYaml RawYaml, err error) {
	file, err := os.Open(string(*f))
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) { err = file.Close() }(file)

	return io.ReadAll(file)
}

type RawYaml []byte

func (y RawYaml) createYamlParser(globalConfig *GlobalConfig) (YamlParser, error) {
	if y.isReusableWorkflow() {
		return NewWorkflow(y, globalConfig), nil
	} else if y.isCustomActions() {
		return NewAction(y, globalConfig), nil
	} else {
		return nil, fmt.Errorf("not found parser: invalid YAML file")
	}
}

func (y RawYaml) isReusableWorkflow() bool {
	r := regexp.MustCompile(`(?m)^[\s]*workflow_call:`)
	return r.Match(y)
}

func (y RawYaml) isCustomActions() bool {
	r := regexp.MustCompile(`(?m)^[\s]*runs:`)
	return r.Match(y)
}

type YamlParser interface {
	Parse() (string, error)
}
