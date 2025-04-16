package workflow

import "github.com/tmknom/actdocs/internal/conf"

func Orchestrate(yaml []byte, formatterConfig *conf.FormatterConfig, sortConfig *conf.SortConfig) (string, error) {
	ast, err := NewParser(sortConfig).Parse(yaml)
	if err != nil {
		return "", err
	}

	spec := ConvertSpec(ast)
	formatted := NewFormatter(formatterConfig).Format(spec)
	return formatted, nil
}
