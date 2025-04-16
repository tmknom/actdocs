package action

import "github.com/tmknom/actdocs/internal/conf"

func Orchestrate(yaml []byte, sortConfig *conf.SortConfig) (*Spec, error) {
	ast, err := NewParser(sortConfig).Parse(yaml)
	if err != nil {
		return nil, err
	}
	return ConvertSpec(ast), nil
}
