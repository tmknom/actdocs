package parse

import (
	"strings"
	"testing"

	"github.com/tmknom/actdocs/internal/config"
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
		action := NewAction(TestRawYaml(tc.fixture), config.DefaultGlobalConfig())
		got, err := action.Parse()
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

const complexActionExpected = `
## Description

This is a test Custom Action for actdocs.

## Inputs

| Name | Description | Default | Required |
| :--- | :---------- | :------ | :------: |
| full-string | The full string value. | ` + "`Default value`" + ` | yes |
| full-boolean | The full boolean value. | ` + "`true`" + ` | no |
| empty |  | n/a | no |

## Outputs

| Name | Description |
| :--- | :---------- |
| with-description | The Render value with description. |
| only-value |  |
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

const fullActionExpected = `
## Description

This is a test Custom Action for actdocs.

## Inputs

| Name | Description | Default | Required |
| :--- | :---------- | :------ | :------: |
| full-number | The full number value. | ` + "`5`" + ` | no |
## Outputs

| Name | Description |
| :--- | :---------- |
| with-description | The Render value with description. |
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

const emptyActionExpected = `
## Description

N/A

## Inputs

| Name | Description | Default | Required |
| :--- | :---------- | :------ | :------: |
| empty |  | n/a | no |
## Outputs

| Name | Description |
| :--- | :---------- |
| only-value |  |
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

const complexMultiLineActionExpected = `
## Description

This is a test Custom Action for actdocs.

## Inputs

| Name | Description | Default | Required |
| :--- | :---------- | :------ | :------: |
| multiline-string | <pre>The multiline string.<br>Like this.</pre> | <pre>{<br>  "key": "value"<br>}</pre> | yes |
| empty |  | n/a | no |

## Outputs

| Name | Description |
| :--- | :---------- |
| with-multiline-description | <pre>The Render value with multiline description.<br>Like this.</pre> |
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
