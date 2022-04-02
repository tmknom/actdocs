package actdocs

import (
	"bytes"
	"encoding/json"
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
	value string
	valid bool // valid is true if value is not NULL
}

func NewNullString(value *string) *NullString {
	var str string
	if value != nil {
		str = *value
	}
	return &NullString{
		value: str,
		valid: value != nil,
	}
}

var DefaultNullString = NewNullString(nil)

func (s *NullString) MarshalJSON() ([]byte, error) {
	if s.valid {
		return json.Marshal(s.value)
	}
	return json.Marshal(nil)
}

func (s *NullString) StringOrEmpty() string {
	if s.valid {
		return s.value
	}
	return emptyString
}

func (s *NullString) QuoteStringOrNA() string {
	if s.valid {
		return s.quoteString()
	}
	return naString
}

func (s *NullString) YesOrNo() string {
	if s.valid && s.value == "true" {
		return yesString
	}
	return noString
}

func (s *NullString) IsTrue() bool {
	return s.valid && s.value == "true"
}

func (s *NullString) quoteString() string {
	return "`" + s.value + "`"
}

const emptyString = ""
const naString = "n/a"
const yesString = "yes"
const noString = "no"
