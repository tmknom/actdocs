package parse

import (
	"github.com/google/go-cmp/cmp"
	"strings"
	"testing"

	"github.com/tmknom/actdocs/internal/conf"
)

func TestWorkflowParse(t *testing.T) {
	cases := []struct {
		name     string
		fixture  string
		expected string
	}{
		{"complex parameter", complexWorkflowFixture, complexWorkflowExpected},
		{"full parameter", fullWorkflowFixture, fullWorkflowExpected},
		{"empty parameter", emptyWorkflowFixture, emptyWorkflowExpected},
		{"invalid YAML", invalidWorkflowFixture, invalidWorkflowExpected},
	}

	for _, tc := range cases {
		parser := NewWorkflowParser(conf.DefaultFormatterConfig(), conf.DefaultSortConfig())
		got, err := parser.Parse(TestRawYaml(tc.fixture))
		if err != nil {
			t.Fatalf("%s: unexpected error: %s", tc.name, err)
		}

		expected := strings.Split(tc.expected, "\n")
		for _, line := range expected {
			if !strings.Contains(got, line) {
				t.Fatalf("%s: not contained:\nexpected:\n%s\n\ngot:\n%s", tc.name, line, got)
			}
		}
	}

}

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
		parser := NewWorkflowParser(conf.DefaultFormatterConfig(), conf.DefaultSortConfig())
		got, err := parser.ParseAST(TestRawYaml(tc.fixture))
		if err != nil {
			t.Fatalf("%s: unexpected error: %s", tc.name, err)
		}

		if diff := cmp.Diff(got, tc.expected); diff != "" {
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
