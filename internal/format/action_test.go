package format

import (
	"strings"
	"testing"

	"github.com/tmknom/actdocs/internal/conf"
	"github.com/tmknom/actdocs/internal/parse"
	"github.com/tmknom/actdocs/internal/util"
)

func TestActionFormatter_Format(t *testing.T) {
	cases := []struct {
		name     string
		ast      *parse.ActionAST
		expected string
	}{
		{
			name: "empty parameter",
			ast: &parse.ActionAST{
				Name:        NewNullValue(),
				Description: NewNullValue(),
				Inputs: []*parse.ActionInput{
					{Name: "empty", Default: NewNullValue(), Description: NewNullValue(), Required: NewNullValue()},
				},
				Outputs: []*parse.ActionOutput{
					{Name: "only-value", Description: NewNullValue()},
				},
			},
			expected: emptyActionExpected,
		},
		{
			name: "full parameter",
			ast: &parse.ActionAST{
				Name:        NewNotNullValue("Test Fixture"),
				Description: NewNotNullValue("This is a test Custom Action for actdocs."),
				Inputs: []*parse.ActionInput{
					{Name: "full-number", Default: NewNotNullValue("5"), Description: NewNotNullValue("The full number value."), Required: NewNotNullValue("false")},
				},
				Outputs: []*parse.ActionOutput{
					{Name: "with-description", Description: NewNotNullValue("The Render value with description.")},
				},
			},
			expected: fullActionExpected,
		},
		{
			name: "complex parameter",
			ast: &parse.ActionAST{
				Name:        NewNotNullValue("Test Fixture"),
				Description: NewNotNullValue("This is a test Custom Action for actdocs."),
				Inputs: []*parse.ActionInput{
					{Name: "full-string", Default: NewNotNullValue("Default value"), Description: NewNotNullValue("The full string value."), Required: NewNotNullValue("true")},
					{Name: "full-boolean", Default: NewNotNullValue("true"), Description: NewNotNullValue("The full boolean value."), Required: NewNotNullValue("false")},
					{Name: "empty", Default: NewNullValue(), Description: NewNullValue(), Required: NewNullValue()},
				},
				Outputs: []*parse.ActionOutput{
					{Name: "with-description", Description: NewNotNullValue("The Render value with description.")},
					{Name: "only-value", Description: NewNullValue()},
				},
			},
			expected: complexActionExpected,
		},
		{
			name: "complex multiline parameter",
			ast: &parse.ActionAST{
				Name:        NewNotNullValue("Test Fixture"),
				Description: NewNotNullValue("This is a test Custom Action for actdocs."),
				Inputs: []*parse.ActionInput{
					{Name: "multiline-string", Default: NewNotNullValue("{\n  \"key\": \"value\"\n}"), Description: NewNotNullValue("The multiline string.\nLike this."), Required: NewNotNullValue("true")},
					{Name: "empty", Default: NewNullValue(), Description: NewNullValue(), Required: NewNullValue()},
				},
				Outputs: []*parse.ActionOutput{
					{Name: "with-multiline-description", Description: NewNotNullValue("The Render value with multiline description.\nLike this.")},
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
