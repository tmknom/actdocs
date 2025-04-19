package workflow

import (
	"bufio"
	"io"
	"strings"
)

type Renderer struct {
	scanner *bufio.Scanner
	builder strings.Builder
	skip    bool
	Omit    bool
}

func NewRenderer(template io.Reader, omit bool) *Renderer {
	scanner := bufio.NewScanner(template)
	var builder strings.Builder
	return &Renderer{
		scanner: scanner,
		builder: builder,
		skip:    false,
		Omit:    omit,
	}
}

func (r *Renderer) Render(spec *Spec) string {
	for r.scanner.Scan() {
		text := r.scanner.Text()
		if !r.skip {
			r.appendTextWithNewline(text)
			r.tryStartContentInjection(spec, text)
		} else {
			r.tryEndContentInjection(text)
		}
	}
	return r.builder.String()
}

func (r *Renderer) tryStartContentInjection(spec *Spec, text string) {
	if r.isStartDirective(text) {
		r.skip = true
		content := r.generateMarkdown(spec, text)
		r.appendGeneratedMarkdown(content)
	}
}

func (r *Renderer) generateMarkdown(spec *Spec, text string) string {
	if text == BeginInputsDirective {
		return spec.ToInputsMarkdown()
	} else if text == BeginSecretsDirective {
		return spec.ToSecretsMarkdown()
	} else if text == BeginOutputsDirective {
		return spec.ToOutputsMarkdown()
	} else if text == BeginPermissionsDirective {
		return spec.ToPermissionsMarkdown()
	}
	return spec.ToMarkdown()
}

func (r *Renderer) appendGeneratedMarkdown(content string) {
	if content != "" {
		r.builder.WriteString("\n")
		r.appendTextWithNewline(content)
		r.builder.WriteString("\n")
	}
}

func (r *Renderer) tryEndContentInjection(text string) {
	if r.isEndDirective(text) {
		r.skip = false
		r.appendTextWithNewline(text)
	}
}

func (r *Renderer) isStartDirective(text string) bool {
	return text == BeginAllDirective || text == BeginInputsDirective || text == BeginSecretsDirective || text == BeginOutputsDirective || text == BeginPermissionsDirective
}

func (r *Renderer) isEndDirective(text string) bool {
	return text == EndAllDirective || text == EndInputsDirective || text == EndSecretsDirective || text == EndOutputsDirective || text == EndPermissionsDirective
}

func (r *Renderer) appendTextWithNewline(text string) {
	r.builder.WriteString(text)
	r.builder.WriteString("\n")
}

const (
	BeginAllDirective = "<!-- actdocs start -->"
	EndAllDirective   = "<!-- actdocs end -->"

	BeginInputsDirective = "<!-- actdocs inputs start -->"
	EndInputsDirective   = "<!-- actdocs inputs end -->"

	BeginSecretsDirective = "<!-- actdocs secrets start -->"
	EndSecretsDirective   = "<!-- actdocs secrets end -->"

	BeginOutputsDirective = "<!-- actdocs outputs start -->"
	EndOutputsDirective   = "<!-- actdocs outputs end -->"

	BeginPermissionsDirective = "<!-- actdocs permissions start -->"
	EndPermissionsDirective   = "<!-- actdocs permissions end -->"
)
