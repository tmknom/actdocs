package action

import "github.com/tmknom/actdocs/internal/conf"

func Orchestrate(yaml []byte, formatterConfig *conf.FormatterConfig, sortConfig *conf.SortConfig) (string, error) {
	content, err := NewActionParser(sortConfig).ParseAST(yaml)
	if err != nil {
		return "", err
	}
	formatted := NewActionFormatter(formatterConfig).Format(content)
	return formatted, nil
}
