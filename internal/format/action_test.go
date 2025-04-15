package format

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tmknom/actdocs/internal/conf"
	"github.com/tmknom/actdocs/internal/parse"
	"github.com/tmknom/actdocs/internal/util"
)

func TestActionFormatter_Format(t *testing.T) {
	cases := []struct {
		name     string
		ast      *parse.ActionAST
		expected string
	}{
		{
			name: "basic",
			ast: &parse.ActionAST{
				Name:        NewNotNullValue("Test Fixture"),
				Description: NewNotNullValue("This is a test Custom Action for actdocs."),
				Inputs: []*parse.ActionInput{
					{Name: "full-string", Default: NewNotNullValue("Default value"), Description: NewNotNullValue("The full string value."), Required: NewNotNullValue("true")},
					{Name: "full-boolean", Default: NewNotNullValue("true"), Description: NewNotNullValue("The full boolean value."), Required: NewNotNullValue("false")},
				},
				Outputs: []*parse.ActionOutput{
					{Name: "with-description", Description: NewNotNullValue("The Render value with description.")},
					{Name: "no-description", Description: NewNullValue()},
				},
			},
			expected: formatExpected,
		},
	}

	for _, tc := range cases {
		formatter := NewActionFormatter(conf.DefaultFormatterConfig())
		got := formatter.Format(tc.ast)
		if diff := cmp.Diff(got, tc.expected); diff != "" {
			t.Errorf("diff: %s", diff)
		}
	}
}

const formatExpected = `## Description

This is a test Custom Action for actdocs.

## Inputs

| Name | Description | Default | Required |
| :--- | :---------- | :------ | :------: |
| full-string | The full string value. | ` + "`Default value`" + ` | yes |
| full-boolean | The full boolean value. | ` + "`true`" + ` | no |

## Outputs

| Name | Description |
| :--- | :---------- |
| with-description | The Render value with description. |
| no-description |  |`

func TestActionFormatter_ToJson(t *testing.T) {
	cases := []struct {
		name     string
		json     *ActionSpec
		expected string
	}{
		{
			name: "empty",
			json: &ActionSpec{
				Description: NewNullValue(),
				Inputs:      []*ActionInputSpec{},
				Outputs:     []*ActionOutputSpec{},
			},
			expected: emptyActionExpectedJson,
		},
		{
			name: "full",
			json: &ActionSpec{
				Description: NewNotNullValue("This is a test Custom Action for actdocs."),
				Inputs: []*ActionInputSpec{
					{Name: "minimal", Default: NewNullValue(), Description: NewNullValue(), Required: NewNullValue()},
					{Name: "full", Default: NewNotNullValue("The string"), Description: NewNotNullValue("The input value."), Required: NewNotNullValue("true")},
				},
				Outputs: []*ActionOutputSpec{
					{Name: "minimal", Description: NewNullValue()},
					{Name: "full", Description: NewNotNullValue("The output value.")},
				},
			},
			expected: fullActionExpectedJson,
		},
	}

	for _, tc := range cases {
		formatter := NewActionFormatter(conf.DefaultFormatterConfig())
		got := formatter.ToJson(tc.json)
		if diff := cmp.Diff(got, tc.expected); diff != "" {
			t.Errorf("diff: %s", diff)
		}
	}
}

const emptyActionExpectedJson = `{
  "description": null,
  "inputs": [],
  "outputs": []
}`

const fullActionExpectedJson = `{
  "description": "This is a test Custom Action for actdocs.",
  "inputs": [
    {
      "name": "minimal",
      "default": null,
      "description": null,
      "required": null
    },
    {
      "name": "full",
      "default": "The string",
      "description": "The input value.",
      "required": "true"
    }
  ],
  "outputs": [
    {
      "name": "minimal",
      "description": null
    },
    {
      "name": "full",
      "description": "The output value."
    }
  ]
}`

func TestActionFormatter_ToMarkdown(t *testing.T) {
	cases := []struct {
		name     string
		config   *conf.FormatterConfig
		markdown *ActionSpec
		expected string
	}{
		{
			name:   "omit",
			config: &conf.FormatterConfig{Format: conf.DefaultFormat, Omit: true},
			markdown: &ActionSpec{
				Description: NewNullValue(),
				Inputs:      []*ActionInputSpec{},
				Outputs:     []*ActionOutputSpec{},
			},
			expected: "",
		},
		{
			name:   "empty",
			config: conf.DefaultFormatterConfig(),
			markdown: &ActionSpec{
				Description: NewNullValue(),
				Inputs:      []*ActionInputSpec{},
				Outputs:     []*ActionOutputSpec{},
			},
			expected: emptyActionExpected,
		},
		{
			name:   "full",
			config: conf.DefaultFormatterConfig(),
			markdown: &ActionSpec{
				Description: NewNotNullValue("This is a test Custom Action for actdocs."),
				Inputs: []*ActionInputSpec{
					{Name: "full-number", Default: NewNotNullValue("5"), Description: NewNotNullValue("The full number value."), Required: NewNotNullValue("false")},
				},
				Outputs: []*ActionOutputSpec{
					{Name: "with-description", Description: NewNotNullValue("The Render value with description.")},
				},
			},
			expected: fullActionExpected,
		},
	}

	for _, tc := range cases {
		formatter := NewActionFormatter(tc.config)
		formatter.ActionSpec = tc.markdown
		got := formatter.ToMarkdown(tc.markdown, tc.config)
		if diff := cmp.Diff(got, tc.expected); diff != "" {
			t.Errorf("diff: %s", diff)
		}
	}
}

const emptyActionExpected = `## Description

N/A

## Inputs

N/A

## Outputs

N/A`

const fullActionExpected = `## Description

This is a test Custom Action for actdocs.

## Inputs

| Name | Description | Default | Required |
| :--- | :---------- | :------ | :------: |
| full-number | The full number value. | ` + "`5`" + ` | no |

## Outputs

| Name | Description |
| :--- | :---------- |
| with-description | The Render value with description. |`

func TestActionFormatter_toDescriptionMarkdown(t *testing.T) {
	cases := []struct {
		name        string
		description *util.NullString
		expected    string
	}{
		{
			name:        "null value",
			description: NewNullValue(),
			expected:    "## Description\n\nN/A",
		},
		{
			name:        "valid value",
			description: NewNotNullValue("The valid."),
			expected:    "## Description\n\nThe valid.",
		},
	}

	for _, tc := range cases {
		formatter := NewActionFormatter(conf.DefaultFormatterConfig())
		got := formatter.toDescriptionMarkdown(tc.description)

		if diff := cmp.Diff(got, tc.expected); diff != "" {
			t.Errorf("diff: %s", diff)
		}
	}
}

func TestActionFormatter_toInputsMarkdown(t *testing.T) {
	cases := []struct {
		name     string
		inputs   []*ActionInputSpec
		expected string
	}{
		{
			name:     "empty",
			inputs:   []*ActionInputSpec{},
			expected: "## Inputs\n\nN/A",
		},
		{
			name: "minimal",
			inputs: []*ActionInputSpec{
				{Name: "minimal", Default: NewNullValue(), Description: NewNullValue(), Required: NewNullValue()},
			},
			expected: "## Inputs\n\n| Name | Description | Default | Required |\n| :--- | :---------- | :------ | :------: |\n| minimal |  | n/a | no |",
		},
		{
			name: "single",
			inputs: []*ActionInputSpec{
				{Name: "single", Default: NewNotNullValue("5"), Description: NewNotNullValue("The number."), Required: NewNotNullValue("true")},
			},
			expected: "## Inputs\n\n| Name | Description | Default | Required |\n| :--- | :---------- | :------ | :------: |\n| single | The number. | `5` | yes |",
		},
		{
			name: "multiple",
			inputs: []*ActionInputSpec{
				{Name: "multiple-1", Default: NewNotNullValue("The string"), Description: NewNotNullValue("1"), Required: NewNotNullValue("false")},
				{Name: "multiple-2", Default: NewNotNullValue("true"), Description: NewNotNullValue("2"), Required: NewNotNullValue("true")},
			},
			expected: "## Inputs\n\n| Name | Description | Default | Required |\n| :--- | :---------- | :------ | :------: |\n| multiple-1 | 1 | `The string` | no |\n| multiple-2 | 2 | `true` | yes |",
		},
	}

	for _, tc := range cases {
		formatter := NewActionFormatter(conf.DefaultFormatterConfig())
		got := formatter.toInputsMarkdown(tc.inputs)

		if diff := cmp.Diff(got, tc.expected); diff != "" {
			t.Errorf("diff: %s", diff)
		}
	}
}

func TestActionFormatter_toOutputsMarkdown(t *testing.T) {
	cases := []struct {
		name     string
		outputs  []*ActionOutputSpec
		expected string
	}{
		{
			name:     "empty",
			outputs:  []*ActionOutputSpec{},
			expected: "## Outputs\n\nN/A",
		},
		{
			name: "minimal",
			outputs: []*ActionOutputSpec{
				{Name: "minimal", Description: NewNullValue()},
			},
			expected: "## Outputs\n\n| Name | Description |\n| :--- | :---------- |\n| minimal |  |",
		},
		{
			name: "single",
			outputs: []*ActionOutputSpec{
				{Name: "single", Description: NewNotNullValue("The test description.")},
			},
			expected: "## Outputs\n\n| Name | Description |\n| :--- | :---------- |\n| single | The test description. |",
		},
		{
			name: "multiple",
			outputs: []*ActionOutputSpec{
				{Name: "multiple-1", Description: NewNotNullValue("1")},
				{Name: "multiple-2", Description: NewNotNullValue("2")},
			},
			expected: "## Outputs\n\n| Name | Description |\n| :--- | :---------- |\n| multiple-1 | 1 |\n| multiple-2 | 2 |",
		},
	}

	for _, tc := range cases {
		formatter := NewActionFormatter(conf.DefaultFormatterConfig())
		got := formatter.toOutputsMarkdown(tc.outputs)

		if diff := cmp.Diff(got, tc.expected); diff != "" {
			t.Errorf("diff: %s", diff)
		}
	}
}

func NewNullValue() *util.NullString {
	return util.NewNullString(nil)
}

func NewNotNullValue(value string) *util.NullString {
	return util.NewNullString(&value)
}
