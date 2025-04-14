package parse

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tmknom/actdocs/internal/conf"
)

func TestActionParser_ParseAST(t *testing.T) {
	cases := []struct {
		name     string
		fixture  string
		expected *ActionAST
	}{
		{
			name:    "empty parameter",
			fixture: emptyActionFixture,
			expected: &ActionAST{
				Name:        NewNullValue(),
				Description: NewNullValue(),
				Inputs: []*ActionInput{
					{"empty", NewNullValue(), NewNullValue(), NewNullValue()},
				},
				Outputs: []*ActionOutput{
					{"only-value", NewNullValue()},
				},
				Runs: &ActionRuns{Using: "undefined", Steps: []*any{}},
			},
		},
		{
			name:    "full parameter",
			fixture: fullActionFixture,
			expected: &ActionAST{
				Name:        NewNotNullValue("Test Fixture"),
				Description: NewNotNullValue("This is a test Custom Action for actdocs."),
				Inputs: []*ActionInput{
					{"full-number", NewNotNullValue("5"), NewNotNullValue("The full number value."), NewNotNullValue("false")},
				},
				Outputs: []*ActionOutput{
					{"with-description", NewNotNullValue("The Render value with description.")},
				},
				Runs: &ActionRuns{Using: "undefined", Steps: []*any{}},
			},
		},
		{
			name:    "complex parameter",
			fixture: complexActionFixture,
			expected: &ActionAST{
				Name:        NewNotNullValue("Test Fixture"),
				Description: NewNotNullValue("This is a test Custom Action for actdocs."),
				Inputs: []*ActionInput{
					{"full-string", NewNotNullValue("Default value"), NewNotNullValue("The full string value."), NewNotNullValue("true")},
					{"full-boolean", NewNotNullValue("true"), NewNotNullValue("The full boolean value."), NewNotNullValue("false")},
					{"empty", NewNullValue(), NewNullValue(), NewNullValue()},
				},
				Outputs: []*ActionOutput{
					{"with-description", NewNotNullValue("The Render value with description.")},
					{"only-value", NewNullValue()},
				},
				Runs: &ActionRuns{Using: "undefined", Steps: []*any{}},
			},
		},
		{
			name:    "invalid YAML",
			fixture: invalidActionFixture,
			expected: &ActionAST{
				Name:        NewNotNullValue("Test"),
				Description: NewNullValue(),
				Inputs:      []*ActionInput{},
				Outputs:     []*ActionOutput{},
				Runs:        &ActionRuns{Using: "undefined", Steps: []*any{}},
			},
		},
	}

	for _, tc := range cases {
		parser := NewActionParser(conf.DefaultSortConfig())
		got, err := parser.ParseAST(TestRawYaml(tc.fixture))
		if err != nil {
			t.Fatalf("%s: unexpected error: %s", tc.name, err)
		}

		if diff := cmp.Diff(got, tc.expected); diff != "" {
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

const complexMultiLineActionFixture = `
name: Test Fixture
description: This is a test Custom Action for actdocs.

inputs:
  multiline-string:
    default: |
      {
        "key": "value"
      }
    required: true
    description: |
      The multiline string.
      Like this.
  empty:

outputs:
  with-multiline-description:
    description: |
      The Render value with multiline description.
      Like this.
    value: ${{ inputs.description-only }}
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

const invalidActionExpected = ""
