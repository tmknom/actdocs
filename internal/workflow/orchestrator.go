package workflow

import (
	"io"

	"github.com/tmknom/actdocs/internal/conf"
)

func Inject(yaml []byte, template io.Reader, formatter *conf.FormatterConfig, sortConfig *conf.SortConfig) (string, error) {
	ast, err := NewParser(sortConfig).Parse(yaml)
	if err != nil {
		return "", err
	}

	spec := ConvertSpec(ast, formatter.Omit)
	return NewRenderer(template, formatter.Omit).Render(spec), nil
}

func Generate(yaml []byte, formatter *conf.FormatterConfig, sortConfig *conf.SortConfig) (string, error) {
	ast, err := NewParser(sortConfig).Parse(yaml)
	if err != nil {
		return "", err
	}

	spec := ConvertSpec(ast, formatter.Omit)
	if formatter.IsJson() {
		return spec.ToJson(), nil
	}
	return spec.ToMarkdown(), nil
}
