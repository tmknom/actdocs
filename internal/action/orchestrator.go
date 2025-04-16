package action

import (
	"io"

	"github.com/tmknom/actdocs/internal/conf"
	"github.com/tmknom/actdocs/internal/render"
)

func Inject(yaml []byte, reader io.Reader, formatter *conf.FormatterConfig, sortConfig *conf.SortConfig) (string, error) {
	renderer := render.NewAllInjectRenderer()
	ast, err := NewParser(sortConfig).Parse(yaml)
	if err != nil {
		return "", err
	}

	spec := ConvertSpec(ast)
	formatted := NewFormatter(formatter).Format(spec)
	return renderer.Render(formatted, reader)
}

func Generate(yaml []byte, formatter *conf.FormatterConfig, sortConfig *conf.SortConfig) (string, error) {
	ast, err := NewParser(sortConfig).Parse(yaml)
	if err != nil {
		return "nil", err
	}

	spec := ConvertSpec(ast)
	return NewFormatter(formatter).Format(spec), nil
}
