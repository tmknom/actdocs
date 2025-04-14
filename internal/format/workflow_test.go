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
					{"empty", NewNullValue(), NewNullValue(), NewNullValue(), NewNullValue()},
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
					{"full-number", NewNotNullValue("5"), NewNotNullValue("The full number value."), NewNotNullValue("false"), NewNotNullValue("number")},
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
					{"full-string", NewNotNullValue(""), NewNotNullValue("The full string value."), NewNotNullValue("true"), NewNotNullValue("string")},
					{"full-boolean", NewNotNullValue("true"), NewNotNullValue("The full boolean value."), NewNotNullValue("false"), NewNotNullValue("boolean")},
					{"empty", NewNullValue(), NewNullValue(), NewNullValue(), NewNullValue()},
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
					{"multiline-string", NewNotNullValue("{\n  \"key\": \"value\"\n}"), NewNotNullValue("The multiline string.\nLike this."), NewNotNullValue("false"), NewNotNullValue("string")},
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
