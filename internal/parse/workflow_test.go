package parse

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/tmknom/actdocs/internal/conf"
)

func TestWorkflowParser_ParseAST(t *testing.T) {
	cases := []struct {
		name     string
		fixture  string
		expected *WorkflowAST
	}{
		{
			name:    "empty parameter",
			fixture: emptyWorkflowFixture,
			expected: &WorkflowAST{
				Inputs: []*WorkflowInput{
					{"empty", NewNullValue(), NewNullValue(), NewNullValue(), NewNullValue()},
				},
				Secrets:     []*WorkflowSecret{},
				Outputs:     []*WorkflowOutput{},
				Permissions: []*WorkflowPermission{},
			},
		},
		{
			name:    "full parameter",
			fixture: fullWorkflowFixture,
			expected: &WorkflowAST{
				Inputs: []*WorkflowInput{
					{"full-number", NewNotNullValue("5"), NewNotNullValue("The full number value."), NewNotNullValue("false"), NewNotNullValue("number")},
				},
				Secrets:     []*WorkflowSecret{},
				Outputs:     []*WorkflowOutput{},
				Permissions: []*WorkflowPermission{},
			},
		},
		{
			name:    "complex parameter",
			fixture: complexWorkflowFixture,
			expected: &WorkflowAST{
				Inputs: []*WorkflowInput{
					{"full-string", NewNotNullValue(""), NewNotNullValue("The full string value."), NewNotNullValue("true"), NewNotNullValue("string")},
					{"full-boolean", NewNotNullValue("true"), NewNotNullValue("The full boolean value."), NewNotNullValue("false"), NewNotNullValue("boolean")},
					{"empty", NewNullValue(), NewNullValue(), NewNullValue(), NewNullValue()},
				},
				Secrets:     []*WorkflowSecret{},
				Outputs:     []*WorkflowOutput{},
				Permissions: []*WorkflowPermission{},
			},
		},
		{
			name:    "invalid YAML",
			fixture: invalidWorkflowFixture,
			expected: &WorkflowAST{
				Inputs:      []*WorkflowInput{},
				Secrets:     []*WorkflowSecret{},
				Outputs:     []*WorkflowOutput{},
				Permissions: []*WorkflowPermission{},
			},
		},
	}

	for _, tc := range cases {
		parser := NewWorkflowParser(conf.DefaultSortConfig())
		got, err := parser.ParseAST(TestRawYaml(tc.fixture))
		if err != nil {
			t.Fatalf("%s: unexpected error: %s", tc.name, err)
		}

		sort := func(a, b *WorkflowInput) bool { return a.Name < b.Name }
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

const invalidWorkflowExpected = ""

type TestRawYaml []byte
