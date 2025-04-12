package util

import (
	"encoding/json"
	"fmt"
	"strings"
)

const TableSeparator = "|"

// NullString represents a string that may be null.
type NullString struct {
	value string
	Valid bool // Valid is true if value is not NULL
}

func NewNullString(value *string) *NullString {
	var str string
	if value != nil {
		str = *value
	}
	return &NullString{
		value: str,
		Valid: value != nil,
	}
}

var DefaultNullString = NewNullString(nil)

func (s *NullString) MarshalJSON() ([]byte, error) {
	if s.Valid {
		return json.Marshal(s.value)
	}
	return json.Marshal(nil)
}

func (s *NullString) StringOrEmpty() string {
	if s.Valid {
		if strings.Contains(s.value, "\n") {
			return s.sanitizeString()
		}
		return s.value
	}
	return emptyString
}

func (s *NullString) StringOrUpperNA() string {
	if s.Valid {
		return s.value
	}
	return UpperNAString
}

func (s *NullString) QuoteStringOrLowerNA() string {
	if s.Valid {
		return s.quoteString()
	}
	return LowerNAString
}

func (s *NullString) YesOrNo() string {
	if s.Valid && s.value == "true" {
		return yesString
	}
	return noString
}

func (s *NullString) IsTrue() bool {
	return s.Valid && s.value == "true"
}

func (s *NullString) quoteString() string {
	if strings.Contains(s.value, "\n") {
		return s.sanitizeString()
	}
	return "`" + s.value + "`"
}

func (s *NullString) sanitizeString() string {
	var str string
	str = strings.TrimSuffix(s.value, "\n")
	str = strings.ReplaceAll(str, "\n", lineBreak)
	str = strings.ReplaceAll(str, "\r", "")
	return fmt.Sprintf("%s%s%s", codeStart, str, codeEnd)
}

const emptyString = ""
const yesString = "yes"
const noString = "no"

const LowerNAString = "n/a"
const UpperNAString = "N/A"

const codeStart = "<pre>"
const codeEnd = "</pre>"
const lineBreak = "<br>"
