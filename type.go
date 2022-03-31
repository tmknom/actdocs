package actdocs

import (
	"bytes"
	"io"
	"os"
)

const TableSeparator = "|"

type Generator interface {
	Generate() (string, error)
}

type rawYaml []byte

func (y rawYaml) IsReusableWorkflow() bool {
	return bytes.Contains(y, []byte("workflow_call:"))
}

func (y rawYaml) IsCustomActions() bool {
	return bytes.Contains(y, []byte("runs:"))
}

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
	return emptyString
}

func (s *NullString) QuoteStringOrNA() string {
	if s.Valid {
		return s.quoteString()
	}
	return naString
}

func (s *NullString) YesOrNo() string {
	if s.Valid && s.String == "true" {
		return yesString
	}
	return noString
}

func (s *NullString) quoteString() string {
	return "`" + s.String + "`"
}

const emptyString = ""
const naString = "n/a"
const yesString = "yes"
const noString = "no"
