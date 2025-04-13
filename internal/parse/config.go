package parse

type SortConfig struct {
	Sort           bool
	SortByName     bool
	SortByRequired bool
}

func DefaultSortConfig() *SortConfig {
	return &SortConfig{
		Sort:           DefaultSort,
		SortByName:     DefaultSortByName,
		SortByRequired: DefaultSortByRequired,
	}
}

const (
	DefaultSort           = false
	DefaultSortByName     = false
	DefaultSortByRequired = false
)
