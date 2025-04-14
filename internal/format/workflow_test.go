package format

import (
	"strings"
	"testing"

	"github.com/tmknom/actdocs/internal/conf"
	"github.com/tmknom/actdocs/internal/parse"
)

func TestWorkflowFormatter_Format(t *testing.T) {
	cases := []struct {
		name     string
		ast      *parse.WorkflowAST
		expected string
	}{
		{
			name: "empty parameter",
			ast: &parse.WorkflowAST{
				Inputs: []*parse.WorkflowInput{
					{Name: "empty", Default: NewNullValue(), Description: NewNullValue(), Required: NewNullValue(), Type: NewNullValue()},
				},
				Secrets:     []*parse.WorkflowSecret{},
				Outputs:     []*parse.WorkflowOutput{},
				Permissions: []*parse.WorkflowPermission{},
			},
			expected: emptyWorkflowExpected,
		},
		{
			name: "full parameter",
			ast: &parse.WorkflowAST{
				Inputs: []*parse.WorkflowInput{
					{Name: "full-number", Default: NewNotNullValue("5"), Description: NewNotNullValue("The full number value."), Required: NewNotNullValue("false"), Type: NewNotNullValue("number")},
				},
				Secrets:     []*parse.WorkflowSecret{},
				Outputs:     []*parse.WorkflowOutput{},
				Permissions: []*parse.WorkflowPermission{},
			},
			expected: fullWorkflowExpected,
		},
		{
			name: "complex parameter",
			ast: &parse.WorkflowAST{
				Inputs: []*parse.WorkflowInput{
					{Name: "full-string", Default: NewNotNullValue(""), Description: NewNotNullValue("The full string value."), Required: NewNotNullValue("true"), Type: NewNotNullValue("string")},
					{Name: "full-boolean", Default: NewNotNullValue("true"), Description: NewNotNullValue("The full boolean value."), Required: NewNotNullValue("false"), Type: NewNotNullValue("boolean")},
					{Name: "empty", Default: NewNullValue(), Description: NewNullValue(), Required: NewNullValue(), Type: NewNullValue()},
				},
				Secrets:     []*parse.WorkflowSecret{},
				Outputs:     []*parse.WorkflowOutput{},
				Permissions: []*parse.WorkflowPermission{},
			},
			expected: complexWorkflowExpected,
		},
		{
			name: "multiline parameter",
			ast: &parse.WorkflowAST{
				Inputs: []*parse.WorkflowInput{
					{Name: "multiline-string", Default: NewNotNullValue("{\n  \"key\": \"value\"\n}"), Description: NewNotNullValue("The multiline string.\nLike this."), Required: NewNotNullValue("false"), Type: NewNotNullValue("string")},
				},
				Secrets:     []*parse.WorkflowSecret{},
				Outputs:     []*parse.WorkflowOutput{},
				Permissions: []*parse.WorkflowPermission{},
			},
			expected: multiLineWorkflowExpected,
		},
	}

	for _, tc := range cases {
		formatter := NewWorkflowFormatter(conf.DefaultFormatterConfig())
		got := formatter.Format(tc.ast)
		expected := strings.Split(tc.expected, "\n")
		for _, line := range expected {
			if !strings.Contains(got, line) {
				t.Fatalf("%s: not contained:\nexpected:\n%s\n\ngot:\n%s", tc.name, line, got)
			}
		}
	}

}

const emptyWorkflowExpected = `
## Inputs

| Name | Description | Type | Default | Required |
| :--- | :---------- | :--- | :------ | :------: |
| empty |  | n/a | n/a | no |
`

const fullWorkflowExpected = `
## Inputs

| Name | Description | Type | Default | Required |
| :--- | :---------- | :--- | :------ | :------: |
| full-number | The full number value. | ` + "`number`" + ` | ` + "`5`" + ` | no |
`

const complexWorkflowExpected = `
## Inputs

| Name | Description | Type | Default | Required |
| :--- | :---------- | :--- | :------ | :------: |
| full-string | The full string value. | ` + "`string`" + ` | ` + "``" + ` | yes |
| full-boolean | The full boolean value. | ` + "`boolean`" + ` | ` + "`true`" + ` | no |
| empty |  | n/a | n/a | no |
`

const multiLineWorkflowExpected = `
| Name | Description | Type | Default | Required |
| :--- | :---------- | :--- | :------ | :------: |
| multiline-string | <pre>The multiline string.<br>Like this.</pre> | ` + "`string`" + ` | <pre>{<br>  "key": "value"<br>}</pre> | no |
`
