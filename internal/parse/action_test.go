package parse

import (
	"strings"
	"testing"

	"github.com/tmknom/actdocs/internal/conf"
)

func TestActionParse(t *testing.T) {
	cases := []struct {
		name     string
		fixture  string
		expected string
	}{
		{"complex parameter", complexActionFixture, complexActionExpected},
		{"full parameter", fullActionFixture, fullActionExpected},
		{"empty parameter", emptyActionFixture, emptyActionExpected},
		{"complex multiline parameter", complexMultiLineActionFixture, complexMultiLineActionExpected},
		{"invalid YAML", invalidActionFixture, invalidActionExpected},
	}

	for _, tc := range cases {
		parser := NewActionParser(conf.DefaultFormatterConfig(), conf.DefaultSortConfig())
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

const emptyActionFixture = `
name:
description:

inputs:
  empty:

outputs:
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
