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
	return r.scan(spec)
}

func (r *Renderer) scan(spec *Spec) string {
	for r.scanner.Scan() {
		text := r.scanner.Text()
		if !r.skip {
			r.addTextWithLinefeed(text)
			r.checkBeginComment(spec, text)
		} else {
			r.checkEndComment(text)
		}
	}
	return r.builder.String()
}

func (r *Renderer) checkBeginComment(spec *Spec, text string) {
	if r.isBeginComment(text) {
		r.skip = true
		content := r.generateContent(spec, text)
		r.injectContent(content)
	}
}

func (r *Renderer) generateContent(spec *Spec, text string) string {
	if text == beginInputsComment {
		return spec.ToInputsMarkdown(r.Omit)
	} else if text == beginSecretsComment {
		return spec.ToSecretsMarkdown(r.Omit)
	} else if text == beginOutputsComment {
		return spec.ToOutputsMarkdown(r.Omit)
	} else if text == beginPermissionsComment {
		return spec.ToPermissionsMarkdown(r.Omit)
	}
	return spec.ToMarkdown(r.Omit)
}

func (r *Renderer) injectContent(content string) {
	if content != "" {
		r.builder.WriteString("\n")
		r.addTextWithLinefeed(content)
		r.builder.WriteString("\n")
	}
}

func (r *Renderer) checkEndComment(text string) {
	if r.isEndComment(text) {
		r.addTextWithLinefeed(text)
		r.skip = false
	}
}

func (r *Renderer) isBeginComment(text string) bool {
	return text == beginAllComment || text == beginInputsComment || text == beginSecretsComment || text == beginOutputsComment || text == beginPermissionsComment
}

func (r *Renderer) isEndComment(text string) bool {
	return text == endAllComment || text == endInputsComment || text == endSecretsComment || text == endOutputsComment || text == endPermissionsComment
}

func (r *Renderer) addTextWithLinefeed(text string) {
	r.builder.WriteString(text)
	r.builder.WriteString("\n")
}

const (
	beginAllComment = "<!-- actdocs start -->"
	endAllComment   = "<!-- actdocs end -->"

	beginInputsComment = "<!-- actdocs inputs start -->"
	endInputsComment   = "<!-- actdocs inputs end -->"

	beginSecretsComment = "<!-- actdocs secrets start -->"
	endSecretsComment   = "<!-- actdocs secrets end -->"

	beginOutputsComment = "<!-- actdocs outputs start -->"
	endOutputsComment   = "<!-- actdocs outputs end -->"

	beginPermissionsComment = "<!-- actdocs permissions start -->"
	endPermissionsComment   = "<!-- actdocs permissions end -->"
)
