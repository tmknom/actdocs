package actdocs

import (
	"bufio"
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

	result := i.render(content, file)
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

	elements := []string{before, beginComment, strings.TrimSpace(content), endComment, after}
	return strings.Join(elements, "\n")
}

func (i *Injector) scanBefore(scanner *bufio.Scanner) string {
	result := ""
	for scanner.Scan() {
		str := scanner.Text()
		if str == beginComment {
			break
		}
		result += str + "\n"
	}
	return strings.TrimSpace(result) + "\n"
}

func (i *Injector) skipCurrentContent(scanner *bufio.Scanner) {
	for scanner.Scan() {
		if scanner.Text() == endComment {
			break
		}
	}
}

func (i *Injector) scanAfter(scanner *bufio.Scanner) string {
	result := ""
	for scanner.Scan() {
		result += scanner.Text() + "\n"
	}
	return result
}

const beginComment = "<!-- actdocs start -->"
const endComment = "<!-- actdocs end -->"
