package actdocs

import (
	"encoding/json"
)

const TableSeparator = "|"

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

func (s *NullString) StringOrUpperNA() string {
	if s.valid {
		return s.value
	}
	return UpperNAString
}

func (s *NullString) QuoteStringOrLowerNA() string {
	if s.valid {
		return s.quoteString()
	}
	return LowerNAString
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
const yesString = "yes"
const noString = "no"

const LowerNAString = "n/a"
const UpperNAString = "N/A"
