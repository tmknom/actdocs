package format

type FormatterConfig struct {
	Format         string
	Omit           bool
	Sort           bool
	SortByName     bool
	SortByRequired bool
}

func DefaultFormatterConfig() *FormatterConfig {
	return &FormatterConfig{
		Format:         DefaultFormat,
		Omit:           DefaultOmit,
		Sort:           DefaultSort,
		SortByName:     DefaultSortByName,
		SortByRequired: DefaultSortByRequired,
	}
}

const (
	DefaultFormat         = "markdown"
	DefaultOmit           = false
	DefaultSort           = false
	DefaultSortByName     = false
	DefaultSortByRequired = false
)

func (c *FormatterConfig) IsJson() bool {
	return c.Format == "json"
}
