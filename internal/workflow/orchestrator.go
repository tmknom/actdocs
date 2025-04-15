package workflow

import "github.com/tmknom/actdocs/internal/conf"

func Orchestrate(yaml []byte, formatterConfig *conf.FormatterConfig, sortConfig *conf.SortConfig) (string, error) {
	content, err := NewWorkflowParser(sortConfig).ParseAST(yaml)
	if err != nil {
		return "", err
	}
	formatted := NewWorkflowFormatter(formatterConfig).Format(content)
	return formatted, nil
}
