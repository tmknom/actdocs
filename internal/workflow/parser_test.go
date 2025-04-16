package workflow

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/tmknom/actdocs/internal/conf"
)

func TestParser_Parse(t *testing.T) {
	cases := []struct {
		name     string
		fixture  string
		expected *AST
	}{
		{
			name:    "empty parameter",
			fixture: emptyWorkflowFixture,
			expected: &AST{
				Inputs: []*InputAST{
					{"empty", NewNullValue(), NewNullValue(), NewNullValue(), NewNullValue()},
				},
				Secrets:     []*SecretAST{},
				Outputs:     []*OutputAST{},
				Permissions: []*PermissionAST{},
			},
		},
		{
			name:    "full parameter",
			fixture: fullWorkflowFixture,
			expected: &AST{
				Inputs: []*InputAST{
					{"full-number", NewNotNullValue("5"), NewNotNullValue("The full number value."), NewNotNullValue("false"), NewNotNullValue("number")},
				},
				Secrets:     []*SecretAST{},
				Outputs:     []*OutputAST{},
				Permissions: []*PermissionAST{},
			},
		},
		{
			name:    "complex parameter",
			fixture: complexWorkflowFixture,
			expected: &AST{
				Inputs: []*InputAST{
					{"full-string", NewNotNullValue(""), NewNotNullValue("The full string value."), NewNotNullValue("true"), NewNotNullValue("string")},
					{"full-boolean", NewNotNullValue("true"), NewNotNullValue("The full boolean value."), NewNotNullValue("false"), NewNotNullValue("boolean")},
					{"empty", NewNullValue(), NewNullValue(), NewNullValue(), NewNullValue()},
				},
				Secrets:     []*SecretAST{},
				Outputs:     []*OutputAST{},
				Permissions: []*PermissionAST{},
			},
		},
		{
			name:    "invalid YAML",
			fixture: invalidWorkflowFixture,
			expected: &AST{
				Inputs:      []*InputAST{},
				Secrets:     []*SecretAST{},
				Outputs:     []*OutputAST{},
				Permissions: []*PermissionAST{},
			},
		},
	}

	for _, tc := range cases {
		parser := NewParser(conf.DefaultSortConfig())
		got, err := parser.Parse(TestRawYaml(tc.fixture))
		if err != nil {
			t.Fatalf("%s: unexpected error: %s", tc.name, err)
		}

		sort := func(a, b *InputAST) bool { return a.Name < b.Name }
		if diff := cmp.Diff(got, tc.expected, cmpopts.SortSlices(sort)); diff != "" {
			t.Errorf("%s: diff: %s", tc.name, diff)
		}
	}
}

const emptyWorkflowFixture = `
on:
  workflow_call:
    inputs:
      empty:
`

const fullWorkflowFixture = `
on:
  workflow_call:
    inputs:
      full-number:
        default: 5
        required: false
        type: number
        description: "The full number value."
`

const complexWorkflowFixture = `
on:
  workflow_call:
    inputs:
      full-string:
        default: ""
        required: true
        type: string
        description: "The full string value."
      full-boolean:
        default: true
        required: false
        type: boolean
        description: "The full boolean value."
      empty:
`

const invalidWorkflowFixture = `
name: Test
inputs:
  full-number:
    default: 5
    required: false
    description: "The full number value."
`
