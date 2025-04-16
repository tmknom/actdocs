package action

import "github.com/tmknom/actdocs/internal/util"

func NewNullValue() *util.NullString {
	return util.NewNullString(nil)
}

func NewNotNullValue(value string) *util.NullString {
	return util.NewNullString(&value)
}

type TestRawYaml []byte
