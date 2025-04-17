package action

import (
	"bufio"
	"io"
	"strings"
)

type Renderer struct {
	scanner *bufio.Scanner
	builder strings.Builder
	Omit    bool
}

func NewRenderer(template io.Reader, omit bool) *Renderer {
	scanner := bufio.NewScanner(template)
	var builder strings.Builder
	return &Renderer{
		scanner: scanner,
		builder: builder,
		Omit:    omit,
	}
}

func (r *Renderer) Render(spec *Spec) string {
	//return r.scan(spec)
	return r.scanAndRender(spec)
}

func (r *Renderer) scan(spec *Spec) string {
	r.scanBeforeBeginComment(beginAllComment)

	content := spec.ToMarkdown(r.Omit)
	if content != "" {
		r.builder.WriteString("\n")
	}
	r.builder.WriteString(strings.TrimSpace(content))
	if content != "" {
		r.builder.WriteString("\n\n")
	}

	for r.scanner.Scan() {
		if r.scanner.Text() == endAllComment {
			break
		}
	}

	r.builder.WriteString(endAllComment)
	r.builder.WriteString("\n")
	for r.scanner.Scan() {
		r.builder.WriteString(r.scanner.Text())
		r.builder.WriteString("\n")
	}

	return r.builder.String()
}

func (r *Renderer) scanBeforeBeginComment(beginComment string) {
	for r.scanner.Scan() {
		str := r.scanner.Text()
		if str == beginComment {
			r.builder.WriteString(beginComment)
			r.builder.WriteString("\n")
			break
		}
		r.builder.WriteString(str)
		r.builder.WriteString("\n")
	}
}

func (r *Renderer) scanAndRender(spec *Spec) string {
	before := r.scanBefore()
	r.skipCurrentContent()
	after := r.scanAfter()

	content := spec.ToMarkdown(r.Omit)

	var sb strings.Builder
	sb.WriteString(before)
	sb.WriteString("\n\n")
	sb.WriteString(beginAllComment)
	sb.WriteString("\n")
	if content != "" {
		sb.WriteString("\n")
	}
	sb.WriteString(strings.TrimSpace(content))
	if content != "" {
		sb.WriteString("\n\n")
	}
	sb.WriteString(endAllComment)
	sb.WriteString("\n\n")
	sb.WriteString(after)
	sb.WriteString("\n")
	return sb.String()
}

func (r *Renderer) scanBefore() string {
	var sb strings.Builder
	for r.scanner.Scan() {
		str := r.scanner.Text()
		if str == beginAllComment {
			break
		}
		sb.WriteString(str)
		sb.WriteString("\n")
	}
	return strings.TrimSpace(sb.String())
}

func (r *Renderer) skipCurrentContent() {
	for r.scanner.Scan() {
		if r.scanner.Text() == endAllComment {
			break
		}
	}
}

func (r *Renderer) scanAfter() string {
	var sb strings.Builder
	for r.scanner.Scan() {
		sb.WriteString(r.scanner.Text())
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
