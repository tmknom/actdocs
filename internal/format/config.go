package format

type FormatterConfig struct {
	Format string
	Omit   bool
}

func DefaultFormatterConfig() *FormatterConfig {
	return &FormatterConfig{
		Format: DefaultFormat,
		Omit:   DefaultOmit,
	}
}

const (
	DefaultFormat = "markdown"
	DefaultOmit   = false
)

func (c *FormatterConfig) IsJson() bool {
	return c.Format == "json"
}
