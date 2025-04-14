package parse

import (
	"strings"
	"testing"

	"github.com/tmknom/actdocs/internal/conf"
)

func TestWorkflowFormatter_Format(t *testing.T) {
	cases := []struct {
		name     string
		ast      *WorkflowAST
		expected string
	}{
		{
			name: "empty parameter",
			ast: &WorkflowAST{
				Inputs: []*WorkflowInput{
					{"empty", NewNullValue(), NewNullValue(), NewNullValue(), NewNullValue()},
				},
				Secrets:     []*WorkflowSecret{},
				Outputs:     []*WorkflowOutput{},
				Permissions: []*WorkflowPermission{},
			},
			expected: emptyWorkflowExpected,
		},
		{
			name: "full parameter",
			ast: &WorkflowAST{
				Inputs: []*WorkflowInput{
					{"full-number", NewNotNullValue("5"), NewNotNullValue("The full number value."), NewNotNullValue("false"), NewNotNullValue("number")},
				},
				Secrets:     []*WorkflowSecret{},
				Outputs:     []*WorkflowOutput{},
				Permissions: []*WorkflowPermission{},
			},
			expected: fullWorkflowExpected,
		},
		{
			name: "complex parameter",
			ast: &WorkflowAST{
				Inputs: []*WorkflowInput{
					{"full-string", NewNotNullValue(""), NewNotNullValue("The full string value."), NewNotNullValue("true"), NewNotNullValue("string")},
					{"full-boolean", NewNotNullValue("true"), NewNotNullValue("The full boolean value."), NewNotNullValue("false"), NewNotNullValue("boolean")},
					{"empty", NewNullValue(), NewNullValue(), NewNullValue(), NewNullValue()},
				},
				Secrets:     []*WorkflowSecret{},
				Outputs:     []*WorkflowOutput{},
				Permissions: []*WorkflowPermission{},
			},
			expected: complexWorkflowExpected,
		},
		{
			name: "multiline parameter",
			ast: &WorkflowAST{
				Inputs: []*WorkflowInput{
					{"multiline-string", NewNotNullValue("{\n  \"key\": \"value\"\n}"), NewNotNullValue("The multiline string.\nLike this."), NewNotNullValue("false"), NewNotNullValue("string")},
				},
				Secrets:     []*WorkflowSecret{},
				Outputs:     []*WorkflowOutput{},
				Permissions: []*WorkflowPermission{},
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
