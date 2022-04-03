package actdocs

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

type Injector struct {
	*InjectorConfig
	*IO
	*YamlFile
}

func NewInjector(config *InjectorConfig, inOut *IO, yamlFile *YamlFile) *Injector {
	return &Injector{
		InjectorConfig: config,
		IO:             inOut,
		YamlFile:       yamlFile,
	}
}

type InjectorConfig struct {
	OutputFile string
	DryRun     bool
	*GlobalConfig
}

func NewInjectorConfig(globalConfig *GlobalConfig) *InjectorConfig {
	return &InjectorConfig{
		GlobalConfig: globalConfig,
	}
}

func (i *Injector) Run() error {
	content, err := i.FormatYaml(i.GlobalConfig)
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
