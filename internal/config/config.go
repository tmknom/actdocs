package config

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

type Ldflags struct {
	Name    string
	Version string
	Commit  string
	Date    string
}

func NewLdflags(name string, version string, commit string, date string) *Ldflags {
	return &Ldflags{
		Name:    name,
		Version: version,
		Commit:  commit,
		Date:    date,
	}
}
