package workflow

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
				Inputs: []*InputAST{
					{Name: "foo", Default: NewNotNullValue("Default"), Description: NewNotNullValue("The inputs."), Required: NewNotNullValue("false"), Type: NewNotNullValue("string")},
				},
				Secrets: []*SecretAST{
					{Name: "bar", Description: NewNotNullValue("The secrets."), Required: NewNotNullValue("false")},
				},
				Outputs: []*OutputAST{
					{Name: "baz", Description: NewNotNullValue("The outputs.")},
				},
				Permissions: []*PermissionAST{
					{Scope: "contents", Access: "write"},
				},
			},
			expected: basicWorkflowExpected,
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

const basicWorkflowExpected = `## Inputs

| Name | Description | Type | Default | Required |
| :--- | :---------- | :--- | :------ | :------: |
| foo | The inputs. | ` + "`string`" + ` | ` + "`Default`" + ` | no |

## Secrets

| Name | Description | Required |
| :--- | :---------- | :------: |
| bar | The secrets. | no |

## Outputs

| Name | Description |
| :--- | :---------- |
| baz | The outputs. |

## Permissions

| Scope | Access |
| :--- | :---- |
| contents | write |`
