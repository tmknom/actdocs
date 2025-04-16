package workflow

import "github.com/tmknom/actdocs/internal/conf"

func Orchestrate(yaml []byte, sortConfig *conf.SortConfig) (*Spec, error) {
	ast, err := NewParser(sortConfig).Parse(yaml)
	if err != nil {
		return nil, err
	}
	return ConvertSpec(ast), nil
}

func Generate(yaml []byte, formatter *conf.FormatterConfig, sortConfig *conf.SortConfig) (string, error) {
	ast, err := NewParser(sortConfig).Parse(yaml)
	if err != nil {
		return "nil", err
	}

	spec := ConvertSpec(ast)
	return NewFormatter(formatter).Format(spec), nil
}
