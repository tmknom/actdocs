package action

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
			fixture: emptyActionFixture,
			expected: &AST{
				Name:        NewNullValue(),
				Description: NewNullValue(),
				Inputs: []*InputAST{
					{"empty", NewNullValue(), NewNullValue(), NewNullValue()},
				},
				Outputs: []*OutputAST{
					{"only-value", NewNullValue()},
				},
				Runs: &RunsAST{Using: "undefined", Steps: []*any{}},
			},
		},
		{
			name:    "full parameter",
			fixture: fullActionFixture,
			expected: &AST{
				Name:        NewNotNullValue("Test Fixture"),
				Description: NewNotNullValue("This is a test Custom Action for actdocs."),
				Inputs: []*InputAST{
					{"full-number", NewNotNullValue("5"), NewNotNullValue("The full number value."), NewNotNullValue("false")},
				},
				Outputs: []*OutputAST{
					{"with-description", NewNotNullValue("The Render value with description.")},
				},
				Runs: &RunsAST{Using: "undefined", Steps: []*any{}},
			},
		},
		{
			name:    "complex parameter",
			fixture: complexActionFixture,
			expected: &AST{
				Name:        NewNotNullValue("Test Fixture"),
				Description: NewNotNullValue("This is a test Custom Action for actdocs."),
				Inputs: []*InputAST{
					{"full-string", NewNotNullValue("Default value"), NewNotNullValue("The full string value."), NewNotNullValue("true")},
					{"full-boolean", NewNotNullValue("true"), NewNotNullValue("The full boolean value."), NewNotNullValue("false")},
					{"empty", NewNullValue(), NewNullValue(), NewNullValue()},
				},
				Outputs: []*OutputAST{
					{"with-description", NewNotNullValue("The Render value with description.")},
					{"only-value", NewNullValue()},
				},
				Runs: &RunsAST{Using: "undefined", Steps: []*any{}},
			},
		},
		{
			name:    "invalid YAML",
			fixture: invalidActionFixture,
			expected: &AST{
				Name:        NewNotNullValue("Test"),
				Description: NewNullValue(),
				Inputs:      []*InputAST{},
				Outputs:     []*OutputAST{},
				Runs:        &RunsAST{Using: "undefined", Steps: []*any{}},
			},
		},
	}

	for _, tc := range cases {
		parser := NewParser(conf.DefaultSortConfig())
		got, err := parser.Parse(TestRawYaml(tc.fixture))
		if err != nil {
			t.Fatalf("%s: unexpected error: %s", tc.name, err)
		}

		sortInput := func(a, b *InputAST) bool { return a.Name < b.Name }
		sortOutput := func(a, b *OutputAST) bool { return a.Name < b.Name }
		if diff := cmp.Diff(got, tc.expected, cmpopts.SortSlices(sortInput), cmpopts.SortSlices(sortOutput)); diff != "" {
			t.Errorf("%s: diff: %s", tc.name, diff)
		}
	}
}

const emptyActionFixture = `
name:
description:

inputs:
  empty:

outputs:
  only-value:
    value: "The Render value without description."
`

const fullActionFixture = `
name: Test Fixture
description: This is a test Custom Action for actdocs.

inputs:
  full-number:
    default: 5
    required: false
    description: "The full number value."

outputs:
  with-description:
    description: "The Render value with description."
    value: ${{ inputs.description-only }}
`

const complexActionFixture = `
name: Test Fixture
description: This is a test Custom Action for actdocs.

inputs:
  full-string:
    default: "Default value"
    required: true
    description: "The full string value."
  full-boolean:
    default: true
    required: false
    description: "The full boolean value."
  empty:

outputs:
  with-description:
    description: "The Render value with description."
    value: ${{ inputs.description-only }}
  only-value:
    value: "The Render value without description."
`

const invalidActionFixture = `
name: Test
on:
  workflow_call:
    inputs:
      full-number:
        type: number
        description: "The full number value."
`
