package actdocs

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAppRunWithGenerate(t *testing.T) {
	cases := []struct {
		args     []string
		expected string
	}{
		{
			args:     []string{"generate", "--sort", "testdata/valid-workflow.yml"},
			expected: expectedGenerateWithSortWorkflow,
		},
		{
			args:     []string{"generate", "--sort-by-name", "testdata/valid-workflow.yml"},
			expected: expectedGenerateWithSortByNameWorkflow,
		},
		{
			args:     []string{"generate", "--sort", "--format=json", "testdata/valid-workflow.yml"},
			expected: expectedGenerateWithSortFormatJsonWorkflow,
		},
		{
			args:     []string{"generate", "--sort", "testdata/valid-action.yml"},
			expected: expectedGenerateWithSortAction,
		},
		{
			args:     []string{"generate", "--sort-by-name", "testdata/valid-action.yml"},
			expected: expectedGenerateWithSortByNameAction,
		},
		{
			args:     []string{"generate", "--sort", "--format=json", "testdata/valid-action.yml"},
			expected: expectedGenerateWithSortFormatJsonAction,
		},
	}

	app := NewApp("test", "", "", "")
	for _, tc := range cases {
		outWriter := &bytes.Buffer{}
		inOut := NewIO(os.Stdin, outWriter, os.Stderr)
		err := app.Run(tc.args, inOut.InReader, inOut.OutWriter, inOut.ErrWriter)

		if err != nil {
			t.Fatalf("%s: unexpected error: %s", strings.Join(tc.args, " "), err)
		}

		if diff := cmp.Diff(outWriter.String(), tc.expected); diff != "" {
			t.Errorf("%q: unexpected out: \n%s", tc.args, diff)
		}
	}
}

const expectedGenerateWithSortWorkflow = `## Inputs

| Name | Description | Type | Default | Required |
| :--- | :---------- | :--- | :------ | :------: |
| full-string | The full string value. | ` + "`string`" + ` | ` + "``" + ` | yes |
| required-and-description | The required and description value. | n/a | n/a | yes |
| default-and-type |  | ` + "`string`" + ` | ` + "`foo`" + ` | no |
| empty |  | n/a | n/a | no |
| full-boolean | The full boolean value. | ` + "`boolean`" + ` | ` + "`true`" + ` | no |
| full-number | The full number value. | ` + "`number`" + ` | ` + "`5`" + ` | no |
`

const expectedGenerateWithSortByNameWorkflow = `## Inputs

| Name | Description | Type | Default | Required |
| :--- | :---------- | :--- | :------ | :------: |
| default-and-type |  | ` + "`string`" + ` | ` + "`foo`" + ` | no |
| empty |  | n/a | n/a | no |
| full-boolean | The full boolean value. | ` + "`boolean`" + ` | ` + "`true`" + ` | no |
| full-number | The full number value. | ` + "`number`" + ` | ` + "`5`" + ` | no |
| full-string | The full string value. | ` + "`string`" + ` | ` + "``" + ` | yes |
| required-and-description | The required and description value. | n/a | n/a | yes |
`

const expectedGenerateWithSortFormatJsonWorkflow = `{
  "inputs": [
    {
      "name": "full-string",
      "default": "",
      "description": "The full string value.",
      "required": "true",
      "type": "string"
    },
    {
      "name": "required-and-description",
      "default": null,
      "description": "The required and description value.",
      "required": "true",
      "type": null
    },
    {
      "name": "default-and-type",
      "default": "foo",
      "description": null,
      "required": null,
      "type": "string"
    },
    {
      "name": "empty",
      "default": null,
      "description": null,
      "required": null,
      "type": null
    },
    {
      "name": "full-boolean",
      "default": "true",
      "description": "The full boolean value.",
      "required": "false",
      "type": "boolean"
    },
    {
      "name": "full-number",
      "default": "5",
      "description": "The full number value.",
      "required": "false",
      "type": "number"
    }
  ]
}`

const expectedGenerateWithSortAction = `## Inputs

| Name | Description | Default | Required |
| :--- | :---------- | :------ | :------: |
| full-string | The full string value. | ` + "`Default value`" + ` | yes |
| description-only | The description without default and required. | n/a | no |
| empty |  | n/a | no |
| full-boolean | The full boolean value. | ` + "`true`" + ` | no |
| full-number | The full number value. | ` + "`5`" + ` | no |

## Outputs

| Name | Description |
| :--- | :---------- |
| only-value |  |
| with-description | The output value with description. |
`

const expectedGenerateWithSortByNameAction = `## Inputs

| Name | Description | Default | Required |
| :--- | :---------- | :------ | :------: |
| description-only | The description without default and required. | n/a | no |
| empty |  | n/a | no |
| full-boolean | The full boolean value. | ` + "`true`" + ` | no |
| full-number | The full number value. | ` + "`5`" + ` | no |
| full-string | The full string value. | ` + "`Default value`" + ` | yes |

## Outputs

| Name | Description |
| :--- | :---------- |
| only-value |  |
| with-description | The output value with description. |
`

const expectedGenerateWithSortFormatJsonAction = `{
  "description": "This is a test Custom Action for actdocs.",
  "inputs": [
    {
      "Name": "full-string",
      "Default": "Default value",
      "Description": "The full string value.",
      "Required": "true"
    },
    {
      "Name": "description-only",
      "Default": null,
      "Description": "The description without default and required.",
      "Required": null
    },
    {
      "Name": "empty",
      "Default": null,
      "Description": null,
      "Required": null
    },
    {
      "Name": "full-boolean",
      "Default": "true",
      "Description": "The full boolean value.",
      "Required": "false"
    },
    {
      "Name": "full-number",
      "Default": "5",
      "Description": "The full number value.",
      "Required": "false"
    }
  ],
  "outputs": [
    {
      "Name": "only-value",
      "Description": null
    },
    {
      "Name": "with-description",
      "Description": "The output value with description."
    }
  ]
}`

func TestAppRunWithInject(t *testing.T) {
	cases := []struct {
		args     []string
		expected string
	}{
		{
			args:     []string{"inject", "--sort", "--dry-run", "--file=testdata/output.md", "testdata/valid-workflow.yml"},
			expected: expectedInjectWithSortWorkflow,
		},
		{
			args:     []string{"inject", "--sort", "--dry-run", "--file=testdata/output.md", "testdata/valid-action.yml"},
			expected: expectedInjectWithSortAction,
		},
	}

	app := NewApp("test", "", "", "")
	for _, tc := range cases {
		outWriter := &bytes.Buffer{}
		inOut := NewIO(os.Stdin, outWriter, os.Stderr)
		err := app.Run(tc.args, inOut.InReader, inOut.OutWriter, inOut.ErrWriter)

		if err != nil {
			t.Fatalf("%s: unexpected error: %s", strings.Join(tc.args, " "), err)
		}

		if diff := cmp.Diff(outWriter.String(), tc.expected); diff != "" {
			t.Errorf("%q: unexpected out: \n%s", tc.args, diff)
		}
	}
}

const expectedInjectWithSortWorkflow = `# Output test

## Header

This is a header.

<!-- actdocs start -->
## Inputs

| Name | Description | Type | Default | Required |
| :--- | :---------- | :--- | :------ | :------: |
| full-string | The full string value. | ` + "`string`" + ` | ` + "``" + ` | yes |
| required-and-description | The required and description value. | n/a | n/a | yes |
| default-and-type |  | ` + "`string`" + ` | ` + "`foo`" + ` | no |
| empty |  | n/a | n/a | no |
| full-boolean | The full boolean value. | ` + "`boolean`" + ` | ` + "`true`" + ` | no |
| full-number | The full number value. | ` + "`number`" + ` | ` + "`5`" + ` | no |
<!-- actdocs end -->

## Footer

This is a footer.
`

const expectedInjectWithSortAction = `# Output test

## Header

This is a header.

<!-- actdocs start -->
## Inputs

| Name | Description | Default | Required |
| :--- | :---------- | :------ | :------: |
| full-string | The full string value. | ` + "`Default value`" + ` | yes |
| description-only | The description without default and required. | n/a | no |
| empty |  | n/a | no |
| full-boolean | The full boolean value. | ` + "`true`" + ` | no |
| full-number | The full number value. | ` + "`5`" + ` | no |

## Outputs

| Name | Description |
| :--- | :---------- |
| only-value |  |
| with-description | The output value with description. |
<!-- actdocs end -->

## Footer

This is a footer.
`
