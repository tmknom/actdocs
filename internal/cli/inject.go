package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tmknom/actdocs/internal/conf"
)

func NewInjectCommand(formatter *conf.FormatterConfig, sort *conf.SortConfig, io *IO) *cobra.Command {
	option := &InjectOption{IO: io}
	command := &cobra.Command{
		Use:   "inject",
		Short: "Inject generated documentation to existing file",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.SetPrefix(fmt.Sprintf("[%s] [%s] ", AppName, cmd.Name()))
			log.Printf("start: command = %s, option = %#v", cmd.Name(), option)
			if len(args) > 0 {
				runner := NewInjectRunner(args[0], formatter, sort, option)
				return runner.Run()
			}
			return cmd.Usage()
		},
	}

	command.PersistentFlags().StringVarP(&option.OutputFile, "file", "f", "", "file path to insert output into (default \"\")")
	command.PersistentFlags().BoolVar(&option.DryRun, "dry-run", false, "dry run")
	return command
}

type InjectRunner struct {
	source string
	*conf.FormatterConfig
	*conf.SortConfig
	*InjectOption
}

func NewInjectRunner(source string, formatter *conf.FormatterConfig, sort *conf.SortConfig, option *InjectOption) *InjectRunner {
	return &InjectRunner{
		source:          source,
		FormatterConfig: formatter,
		SortConfig:      sort,
		InjectOption:    option,
	}
}

type InjectOption struct {
	OutputFile string
	DryRun     bool
	*IO
}

func (r *InjectRunner) Run() error {
	formatted, err := Orchestrate(r.source, r.FormatterConfig, r.SortConfig)
	if err != nil {
		return err
	}

	file, err := os.Open(r.OutputFile)
	if err != nil {
		return err
	}
	defer func(file *os.File) { err = file.Close() }(file)

	var result string
	if formatted != "" {
		result = r.render(formatted, file)
	} else {
		result, err = r.renderWithoutOverride(file)
		if err != nil {
			return err
		}
	}

	if r.DryRun {
		_, err = fmt.Fprintf(r.OutWriter, result)
		return err
	}
	return os.WriteFile(r.OutputFile, []byte(result), 0644)
}

func (r *InjectRunner) render(content string, reader io.Reader) string {
	scanner := bufio.NewScanner(reader)

	before := r.scanBefore(scanner)
	r.skipCurrentContent(scanner)
	after := r.scanAfter(scanner)

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

func (r *InjectRunner) scanBefore(scanner *bufio.Scanner) string {
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

func (r *InjectRunner) skipCurrentContent(scanner *bufio.Scanner) {
	for scanner.Scan() {
		if scanner.Text() == endComment {
			break
		}
	}
}

func (r *InjectRunner) scanAfter(scanner *bufio.Scanner) string {
	var sb strings.Builder
	for scanner.Scan() {
		sb.WriteString(scanner.Text())
		sb.WriteString("\n")
	}
	return strings.TrimSpace(sb.String())
}

func (r *InjectRunner) renderWithoutOverride(reader io.Reader) (string, error) {
	buf := bytes.Buffer{}
	_, err := buf.ReadFrom(reader)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

const beginComment = "<!-- actdocs start -->"
const endComment = "<!-- actdocs end -->"
