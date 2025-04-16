package cli

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
)

type InjectRenderer struct{}

func (r *InjectRenderer) Render(docs string, dest *os.File) (string, error) {
	if docs != "" {
		return r.renderWithOverride(docs, dest), nil
	}
	return r.renderWithoutOverride(dest)
}

func (r *InjectRenderer) renderWithOverride(content string, reader io.Reader) string {
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

func (r *InjectRenderer) renderWithoutOverride(reader io.Reader) (string, error) {
	buf := bytes.Buffer{}
	_, err := buf.ReadFrom(reader)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (r *InjectRenderer) scanBefore(scanner *bufio.Scanner) string {
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

func (r *InjectRenderer) skipCurrentContent(scanner *bufio.Scanner) {
	for scanner.Scan() {
		if scanner.Text() == endComment {
			break
		}
	}
}

func (r *InjectRenderer) scanAfter(scanner *bufio.Scanner) string {
	var sb strings.Builder
	for scanner.Scan() {
		sb.WriteString(scanner.Text())
		sb.WriteString("\n")
	}
	return strings.TrimSpace(sb.String())
}

const beginComment = "<!-- actdocs start -->"
const endComment = "<!-- actdocs end -->"
