package action

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

type Renderer struct {
	BeginComment string
	EndComment   string
}

func NewAllInjectRenderer() *Renderer {
	return &Renderer{
		BeginComment: beginAllComment,
		EndComment:   endAllComment,
	}
}

func (r *Renderer) Render(content string, reader io.Reader) (string, error) {
	if content != "" {
		return r.renderWithOverride(content, reader), nil
	}
	return r.renderWithoutOverride(reader)
}

func (r *Renderer) renderWithOverride(content string, reader io.Reader) string {
	scanner := bufio.NewScanner(reader)

	before := r.scanBefore(scanner)
	r.skipCurrentContent(scanner)
	after := r.scanAfter(scanner)

	var sb strings.Builder
	sb.WriteString(before)
	sb.WriteString("\n\n")
	sb.WriteString(r.BeginComment)
	sb.WriteString("\n\n")
	sb.WriteString(strings.TrimSpace(content))
	sb.WriteString("\n\n")
	sb.WriteString(r.EndComment)
	sb.WriteString("\n\n")
	sb.WriteString(after)
	sb.WriteString("\n")
	return sb.String()
}

func (r *Renderer) renderWithoutOverride(reader io.Reader) (string, error) {
	buf := bytes.Buffer{}
	_, err := buf.ReadFrom(reader)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (r *Renderer) scanBefore(scanner *bufio.Scanner) string {
	var sb strings.Builder
	for scanner.Scan() {
		str := scanner.Text()
		if str == r.BeginComment {
			break
		}
		sb.WriteString(str)
		sb.WriteString("\n")
	}
	return strings.TrimSpace(sb.String())
}

func (r *Renderer) skipCurrentContent(scanner *bufio.Scanner) {
	for scanner.Scan() {
		if scanner.Text() == r.EndComment {
			break
		}
	}
}

func (r *Renderer) scanAfter(scanner *bufio.Scanner) string {
	var sb strings.Builder
	for scanner.Scan() {
		sb.WriteString(scanner.Text())
		sb.WriteString("\n")
	}
	return strings.TrimSpace(sb.String())
}

func (r *Renderer) isEndComment(text string) bool {
	return text == endAllComment || text == endDescriptionComment || text == endInputsComment || text == endOutputsComment
}

const (
	beginAllComment = "<!-- actdocs start -->"
	endAllComment   = "<!-- actdocs end -->"

	beginDescriptionComment = "<!-- actdocs description start -->"
	endDescriptionComment   = "<!-- actdocs description end -->"

	beginInputsComment = "<!-- actdocs inputs start -->"
	endInputsComment   = "<!-- actdocs inputs end -->"

	beginOutputsComment = "<!-- actdocs outputs start -->"
	endOutputsComment   = "<!-- actdocs outputs end -->"
)
