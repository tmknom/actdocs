package action

import (
	"io"

	"github.com/tmknom/actdocs/internal/conf"
)

func Inject(yaml []byte, template io.Reader, formatter *conf.FormatterConfig, sortConfig *conf.SortConfig) (string, error) {
	ast, err := NewParser(sortConfig).Parse(yaml)
	if err != nil {
		return "", err
	}

	spec := ConvertSpec(ast)
	return NewRenderer(template, formatter.Omit).Render(spec), nil
}

func Generate(yaml []byte, formatter *conf.FormatterConfig, sortConfig *conf.SortConfig) (string, error) {
	ast, err := NewParser(sortConfig).Parse(yaml)
	if err != nil {
		return "nil", err
	}

	spec := ConvertSpec(ast)
	return NewFormatter(formatter).Format(spec), nil
}
