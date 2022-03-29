package actdocs

import (
	"strings"
	"testing"
)

func TestActionGenerate(t *testing.T) {
	cases := []struct {
		name     string
		fixture  string
		expected string
	}{
		{"complex parameter", complexParameterFixture, complexParameterExpected},
		{"full parameter", fullParameterFixture, fullParameterExpected},
		{"empty parameter", emptyParameterFixture, emptyParameterExpected},
		{"invalid action", invalidActionFixture, invalidParameterExpected},
	}

	for _, tc := range cases {
		action := NewAction(rawYaml(tc.fixture))
		got, err := action.Generate()
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

const complexParameterFixture = `
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
    description: "The output value with description."
    value: ${{ inputs.description-only }}
  only-value:
    value: "The output value without description."
`

const complexParameterExpected = `
## Inputs

| Name | Description | Default | Required |
| :--- | :---------- | :------ | :------: |
| full-string | The full string value. | Default value | true |
| full-boolean | The full boolean value. | true | false |
| empty |  |  |  |

## Outputs

| Name | Description |
| :--- | :---------- |
| with-description | The output value with description. |
| only-value |  |
`

const fullParameterFixture = `
inputs:
  full-number:
    default: 5
    required: false
    description: "The full number value."

outputs:
  with-description:
    description: "The output value with description."
    value: ${{ inputs.description-only }}
`

const fullParameterExpected = `
## Inputs

| Name | Description | Default | Required |
| :--- | :---------- | :------ | :------: |
| full-number | The full number value. | 5 | false |
## Outputs

| Name | Description |
| :--- | :---------- |
| with-description | The output value with description. |
`

const emptyParameterFixture = `
inputs:
  empty:

outputs:
  only-value:
    value: "The output value without description."
`

const emptyParameterExpected = `
## Inputs

| Name | Description | Default | Required |
| :--- | :---------- | :------ | :------: |
| empty |  |  |  |
## Outputs

| Name | Description |
| :--- | :---------- |
| only-value |  |
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

const invalidParameterExpected = ``
