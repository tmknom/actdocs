package config

import (
	"io"
)

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

type IO struct {
	InReader  io.Reader
	OutWriter io.Writer
	ErrWriter io.Writer
}

func NewIO(inReader io.Reader, outWriter, errWriter io.Writer) *IO {
	return &IO{
		InReader:  inReader,
		OutWriter: outWriter,
		ErrWriter: errWriter,
	}
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
