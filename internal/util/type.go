package util

import (
	"encoding/json"
	"fmt"
	"strings"
)

const TableSeparator = "|"

// NullString represents a string that may be null.
type NullString struct {
	Value string
	Valid bool // Valid is true if Value is not NULL
}

func NewNullString(value *string) *NullString {
	var str string
	if value != nil {
		str = *value
	}
	return &NullString{
		Value: str,
		Valid: value != nil,
	}
}

var DefaultNullString = NewNullString(nil)

func (s *NullString) MarshalJSON() ([]byte, error) {
	if s.Valid {
		return json.Marshal(s.Value)
	}
	return json.Marshal(nil)
}

func (s *NullString) StringOrEmpty() string {
	if s.Valid {
		if strings.Contains(s.Value, "\n") {
			return s.sanitizeString()
		}
		return s.Value
	}
	return emptyString
}

func (s *NullString) StringOrUpperNA() string {
	if s.Valid {
		return s.Value
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
	if s.Valid && s.Value == "true" {
		return yesString
	}
	return noString
}

func (s *NullString) IsTrue() bool {
	return s.Valid && s.Value == "true"
}

func (s *NullString) IsValid() bool {
	return s.Valid
}

func (s *NullString) quoteString() string {
	if strings.Contains(s.Value, "\n") {
		return s.sanitizeString()
	}
	return "`" + s.Value + "`"
}

func (s *NullString) sanitizeString() string {
	var str string
	str = strings.TrimSuffix(s.Value, "\n")
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
