package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/tmknom/actdocs/internal/config"
	"github.com/tmknom/actdocs/internal/parse"
	"github.com/tmknom/actdocs/internal/read"
)

type Injector struct {
	*InjectorConfig
	*IO
	YamlFile string
}

func NewInjector(config *InjectorConfig, inOut *IO, yamlFile string) *Injector {
	return &Injector{
		InjectorConfig: config,
		IO:             inOut,
		YamlFile:       yamlFile,
	}
}

type InjectorConfig struct {
	OutputFile string
	DryRun     bool
	*config.GlobalConfig
}

func NewInjectorConfig(globalConfig *config.GlobalConfig) *InjectorConfig {
	return &InjectorConfig{
		GlobalConfig: globalConfig,
	}
}

func (i *Injector) Run() error {
	reader := &read.YamlReader{Filename: i.YamlFile}
	yaml, err := reader.Read()
	if err != nil {
		return err
	}
	log.Printf("read: %s", i.YamlFile)

	factory := &parse.ParserFactory{Raw: yaml}
	parser, err := factory.Factory(i.GlobalConfig)
	if err != nil {
		return err
	}
	log.Printf("selected parser: %T", parser)

	content, err := parser.Parse()
	if err != nil {
		return err
	}

	file, err := os.Open(i.OutputFile)
	if err != nil {
		return err
	}
	defer func(file *os.File) { err = file.Close() }(file)

	var result string
	if content != "" {
		result = i.render(content, file)
	} else {
		result, err = i.renderWithoutOverride(file)
		if err != nil {
			return err
		}
	}

	if i.DryRun {
		_, err = fmt.Fprintf(i.OutWriter, result)
		return err
	}
	return os.WriteFile(i.OutputFile, []byte(result), 0644)
}

func (i *Injector) render(content string, reader io.Reader) string {
	scanner := bufio.NewScanner(reader)

	before := i.scanBefore(scanner)
	i.skipCurrentContent(scanner)
	after := i.scanAfter(scanner)

	var sb strings.Builder
	sb.WriteString(before)
	sb.WriteString("\n\n")
	sb.WriteString(beginComment)
	sb.WriteString("\n\n")
	sb.WriteString(strings.TrimSpace(content))
	sb.WriteString("\n\n")
	sb.WriteString(endComment)
	sb.WriteString("\n\n")
	sb.WriteString(after)
	sb.WriteString("\n")
	return sb.String()
}

func (i *Injector) scanBefore(scanner *bufio.Scanner) string {
	var sb strings.Builder
	for scanner.Scan() {
		str := scanner.Text()
		if str == beginComment {
			break
		}
		sb.WriteString(str)
		sb.WriteString("\n")
	}
	return strings.TrimSpace(sb.String())
}

func (i *Injector) skipCurrentContent(scanner *bufio.Scanner) {
	for scanner.Scan() {
		if scanner.Text() == endComment {
			break
		}
	}
}

func (i *Injector) scanAfter(scanner *bufio.Scanner) string {
	var sb strings.Builder
	for scanner.Scan() {
		sb.WriteString(scanner.Text())
		sb.WriteString("\n")
	}
	return strings.TrimSpace(sb.String())
}

func (i *Injector) renderWithoutOverride(reader io.Reader) (string, error) {
	buf := bytes.Buffer{}
	_, err := buf.ReadFrom(reader)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

const beginComment = "<!-- actdocs start -->"
const endComment = "<!-- actdocs end -->"
