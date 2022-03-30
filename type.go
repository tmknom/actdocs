package actdocs

import (
	"io"
	"os"
)

const TableSeparator = "|"

type rawYaml []byte

func readYaml(filename string) (rawYaml rawYaml, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) { err = file.Close() }(file)

	return io.ReadAll(file)
}

// NullString represents a string that may be null.
type NullString struct {
	String string
	Valid  bool // Valid is true if String is not NULL
}

func NewNullString(value *string) *NullString {
	var str string
	if value != nil {
		str = *value
	}
	return &NullString{
		String: str,
		Valid:  value != nil,
	}
}

var DefaultNullString = NewNullString(nil)

func (s *NullString) StringOrEmpty() string {
	if s.Valid {
		return s.String
	}
	return ""
}
