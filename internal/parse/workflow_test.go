package parse

import (
	"strings"
	"testing"

	"github.com/tmknom/actdocs/internal/format"
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
		{"multiline parameter", multiLineWorkflowFixture, multiLineWorkflowExpected},
		{"invalid YAML", invalidWorkflowFixture, invalidWorkflowExpected},
	}

	for _, tc := range cases {
		parser := NewWorkflowParser(format.DefaultFormatterConfig())
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
## Inputs

| Name | Description | Type | Default | Required |
| :--- | :---------- | :--- | :------ | :------: |
| full-string | The full string value. | ` + "`string`" + ` | ` + "``" + ` | yes |
| full-boolean | The full boolean value. | ` + "`boolean`" + ` | ` + "`true`" + ` | no |
| empty |  | n/a | n/a | no |
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
## Inputs

| Name | Description | Type | Default | Required |
| :--- | :---------- | :--- | :------ | :------: |
| full-number | The full number value. | ` + "`number`" + ` | ` + "`5`" + ` | no |
`

const emptyWorkflowFixture = `
on:
  workflow_call:
    inputs:
      empty:
`

const emptyWorkflowExpected = `
## Inputs

| Name | Description | Type | Default | Required |
| :--- | :---------- | :--- | :------ | :------: |
| empty |  | n/a | n/a | no |
`

const multiLineWorkflowFixture = `
on:
  workflow_call:
    inputs:
      multiline-string:
        default: |
          {
            "key": "value"
          }
        required: false
        type: string
        description: |
          The Multiline string.
          Like this.
`

const multiLineWorkflowExpected = `
| Name | Description | Type | Default | Required |
| :--- | :---------- | :--- | :------ | :------: |
| multiline-string | <pre>The Multiline string.<br>Like this.</pre> | ` + "`string`" + ` | <pre>{<br>  "key": "value"<br>}</pre> | no |
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
