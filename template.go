package actdocs

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Template struct {
	*TemplateConfig
}

func NewTemplate(config *TemplateConfig) *Template {
	return &Template{
		TemplateConfig: config,
	}
}

type TemplateConfig struct {
	OutputFile string
	outWriter  io.Writer
}

func NewTemplateConfig(outWriter io.Writer) *TemplateConfig {
	return &TemplateConfig{
		OutputFile: "",
		outWriter:  outWriter,
	}
}

const beginComment = "<!-- actdocs start -->"
const endComment = "<!-- actdocs end -->"

func (t *Template) Render(content string) (err error) {
	if t.OutputFile == "" {
		fmt.Fprint(t.outWriter, content)
		return nil
	}

	file, err := os.Open(t.OutputFile)
	if err != nil {
		return err
	}
	defer func(file *os.File) { err = file.Close() }(file)

	result := t.renderContent(content, file)
	return t.write(result)
}

func (t *Template) scanBefore(scanner *bufio.Scanner) string {
	result := ""
	for scanner.Scan() {
		str := scanner.Text()
		if str == beginComment {
			break
		}
		result += str + "\n"
	}
	return result
}

func (t *Template) skipCurrentContent(scanner *bufio.Scanner) {
	for scanner.Scan() {
		if scanner.Text() == endComment {
			break
		}
	}
}

func (t *Template) scanAfter(scanner *bufio.Scanner) string {
	result := ""
	for scanner.Scan() {
		result += scanner.Text() + "\n"
	}
	return result
}

func (t *Template) renderContent(content string, reader io.Reader) string {
	scanner := bufio.NewScanner(reader)

	before := t.scanBefore(scanner)
	t.skipCurrentContent(scanner)
	after := t.scanAfter(scanner)

	elements := []string{before, beginComment, strings.TrimSpace(content), endComment, after}
	return strings.Join(elements, "\n")
}

func (t *Template) write(result string) error {
	log.Printf("generated:\n%s", result)
	err := os.WriteFile(t.OutputFile, []byte(result), 0644)
	if err != nil {
		return err
	}
	return nil
}
