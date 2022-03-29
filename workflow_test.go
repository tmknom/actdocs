package actdocs

import (
	"strings"
	"testing"
)

func TestWorkflowGenerate(t *testing.T) {
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
		workflow := NewWorkflow(rawYaml(tc.fixture))
		got, err := workflow.Generate()
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
const complexWorkflowExpected = `
| Name | Description | Default | Type  | Required |
| :--- | :---------- | :------ | :---: | :------: |
| full-string | The full string value. |  | string | true |
| full-boolean | The full boolean value. | true | boolean | false |
| empty |  |  |  |  |
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

const fullWorkflowExpected = `
| Name | Description | Default | Type  | Required |
| :--- | :---------- | :------ | :---: | :------: |
| full-number | The full number value. | 5 | number | false |
`

const emptyWorkflowFixture = `
on:
  workflow_call:
    inputs:
      empty:
`

const emptyWorkflowExpected = `
| Name | Description | Default | Type  | Required |
| :--- | :---------- | :------ | :---: | :------: |
`

const invalidWorkflowFixture = `
name: Test
inputs:
  full-number:
    default: 5
    required: false
    description: "The full number value."
`

const invalidWorkflowExpected = `
| Name | Description | Default | Type  | Required |
| :--- | :---------- | :------ | :---: | :------: |
`
