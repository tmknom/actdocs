package parse

import (
	"github.com/tmknom/actdocs/internal/util"
	"strings"
	"testing"

	"github.com/tmknom/actdocs/internal/conf"
)

func TestActionFormatter_Format(t *testing.T) {
	cases := []struct {
		name     string
		ast      *ActionAST
		expected string
	}{
		{
			name: "empty parameter",
			ast: &ActionAST{
				Name:        NewNullValue(),
				Description: NewNullValue(),
				Inputs: []*ActionInput{
					{"empty", NewNullValue(), NewNullValue(), NewNullValue()},
				},
				Outputs: []*ActionOutput{
					{"only-value", NewNullValue()},
				},
			},
			expected: emptyActionExpected,
		},
		{
			name: "full parameter",
			ast: &ActionAST{
				Name:        NewNotNullValue("Test Fixture"),
				Description: NewNotNullValue("This is a test Custom Action for actdocs."),
				Inputs: []*ActionInput{
					{"full-number", NewNotNullValue("5"), NewNotNullValue("The full number value."), NewNotNullValue("false")},
				},
				Outputs: []*ActionOutput{
					{"with-description", NewNotNullValue("The Render value with description.")},
				},
			},
			expected: fullActionExpected,
		},
		{
			name: "complex parameter",
			ast: &ActionAST{
				Name:        NewNotNullValue("Test Fixture"),
				Description: NewNotNullValue("This is a test Custom Action for actdocs."),
				Inputs: []*ActionInput{
					{"full-string", NewNotNullValue("Default value"), NewNotNullValue("The full string value."), NewNotNullValue("true")},
					{"full-boolean", NewNotNullValue("true"), NewNotNullValue("The full boolean value."), NewNotNullValue("false")},
					{"empty", NewNullValue(), NewNullValue(), NewNullValue()},
				},
				Outputs: []*ActionOutput{
					{"with-description", NewNotNullValue("The Render value with description.")},
					{"only-value", NewNullValue()},
				},
			},
			expected: complexActionExpected,
		},
		{
			name: "complex multiline parameter",
			ast: &ActionAST{
				Name:        NewNotNullValue("Test Fixture"),
				Description: NewNotNullValue("This is a test Custom Action for actdocs."),
				Inputs: []*ActionInput{
					{"multiline-string", NewNotNullValue("{\n  \"key\": \"value\"\n}"), NewNotNullValue("The multiline string.\nLike this."), NewNotNullValue("true")},
					{"empty", NewNullValue(), NewNullValue(), NewNullValue()},
				},
				Outputs: []*ActionOutput{
					{"with-multiline-description", NewNotNullValue("The Render value with multiline description.\nLike this.")},
				},
			},
			expected: complexMultiLineActionExpected,
		},
	}

	for _, tc := range cases {
		formatter := NewActionFormatter(conf.DefaultFormatterConfig())
		got := formatter.Format(tc.ast)
		expected := strings.Split(tc.expected, "\n")
		for _, line := range expected {
			if !strings.Contains(got, line) {
				t.Fatalf("%s: not contained:\nexpected:\n%s\n\ngot:\n%s", tc.name, line, got)
			}
		}
	}
}

func NewNullValue() *util.NullString {
	return util.NewNullString(nil)
}

func NewNotNullValue(value string) *util.NullString {
	return util.NewNullString(&value)
}

const emptyActionExpected = `
## Description

N/A

## Inputs

| Name | Description | Default | Required |
| :--- | :---------- | :------ | :------: |
| empty |  | n/a | no |
## Outputs

| Name | Description |
| :--- | :---------- |
| only-value |  |
`

const fullActionExpected = `
## Description

This is a test Custom Action for actdocs.

## Inputs

| Name | Description | Default | Required |
| :--- | :---------- | :------ | :------: |
| full-number | The full number value. | ` + "`5`" + ` | no |
## Outputs

| Name | Description |
| :--- | :---------- |
| with-description | The Render value with description. |
`

const complexActionExpected = `
## Description

This is a test Custom Action for actdocs.

## Inputs

| Name | Description | Default | Required |
| :--- | :---------- | :------ | :------: |
| full-string | The full string value. | ` + "`Default value`" + ` | yes |
| full-boolean | The full boolean value. | ` + "`true`" + ` | no |
| empty |  | n/a | no |

## Outputs

| Name | Description |
| :--- | :---------- |
| with-description | The Render value with description. |
| only-value |  |
`

const complexMultiLineActionExpected = `
## Description

This is a test Custom Action for actdocs.

## Inputs

| Name | Description | Default | Required |
| :--- | :---------- | :------ | :------: |
| multiline-string | <pre>The multiline string.<br>Like this.</pre> | <pre>{<br>  "key": "value"<br>}</pre> | yes |
| empty |  | n/a | no |

## Outputs

| Name | Description |
| :--- | :---------- |
| with-multiline-description | <pre>The Render value with multiline description.<br>Like this.</pre> |
`
