package workflow

import (
	"github.com/tmknom/actdocs/internal/conf"
)

type Formatter struct {
	config *conf.FormatterConfig
	*Spec
}

func NewFormatter(config *conf.FormatterConfig) *Formatter {
	return &Formatter{
		config: config,
	}
}

func (f *Formatter) Format(spec *Spec) string {
	if f.config.IsJson() {
		return spec.ToJson()
	}
	return spec.ToMarkdown(f.config.Omit)
}
