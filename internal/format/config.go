package format

type GlobalConfig struct {
	Format         string
	Omit           bool
	Sort           bool
	SortByName     bool
	SortByRequired bool
}

func DefaultGlobalConfig() *GlobalConfig {
	return &GlobalConfig{
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

func (c *GlobalConfig) IsJson() bool {
	return c.Format == "json"
}
