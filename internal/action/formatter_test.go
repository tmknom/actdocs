package action

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tmknom/actdocs/internal/conf"
)

func TestFormatter_Format(t *testing.T) {
	cases := []struct {
		name     string
		ast      *AST
		expected string
	}{
		{
			name: "basic",
			ast: &AST{
				Name:        NewNotNullValue("Test Fixture"),
				Description: NewNotNullValue("This is a test Custom Action for actdocs."),
				Inputs: []*InputAST{
					{Name: "full-string", Default: NewNotNullValue("Default value"), Description: NewNotNullValue("The full string value."), Required: NewNotNullValue("true")},
					{Name: "full-boolean", Default: NewNotNullValue("true"), Description: NewNotNullValue("The full boolean value."), Required: NewNotNullValue("false")},
				},
				Outputs: []*OutputAST{
					{Name: "with-description", Description: NewNotNullValue("The Render value with description.")},
					{Name: "no-description", Description: NewNullValue()},
				},
			},
			expected: formatExpected,
		},
	}

	for _, tc := range cases {
		formatter := NewFormatter(conf.DefaultFormatterConfig())
		spec := ConvertSpec(tc.ast)
		got := formatter.Format(spec)
		if diff := cmp.Diff(got, tc.expected); diff != "" {
			t.Errorf("diff: %s", diff)
		}
	}
}

const formatExpected = `## Description

This is a test Custom Action for actdocs.

## Inputs

| Name | Description | Default | Required |
| :--- | :---------- | :------ | :------: |
| full-string | The full string value. | ` + "`Default value`" + ` | yes |
| full-boolean | The full boolean value. | ` + "`true`" + ` | no |

## Outputs

| Name | Description |
| :--- | :---------- |
| with-description | The Render value with description. |
| no-description |  |`
