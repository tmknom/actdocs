package action

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tmknom/actdocs/internal/conf"
	"github.com/tmknom/actdocs/internal/parse"
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
		json     *Spec
		expected string
	}{
		{
			name: "empty",
			json: &Spec{
				Description: NewNullValue(),
				Inputs:      []*InputSpec{},
				Outputs:     []*OutputSpec{},
			},
			expected: emptyActionExpectedJson,
		},
		{
			name: "full",
			json: &Spec{
				Description: NewNotNullValue("This is a test Custom Action for actdocs."),
				Inputs: []*InputSpec{
					{Name: "minimal", Default: NewNullValue(), Description: NewNullValue(), Required: NewNullValue()},
					{Name: "full", Default: NewNotNullValue("The string"), Description: NewNotNullValue("The input value."), Required: NewNotNullValue("true")},
				},
				Outputs: []*OutputSpec{
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
		markdown *Spec
		expected string
	}{
		{
			name:   "omit",
			config: &conf.FormatterConfig{Format: conf.DefaultFormat, Omit: true},
			markdown: &Spec{
				Description: NewNullValue(),
				Inputs:      []*InputSpec{},
				Outputs:     []*OutputSpec{},
			},
			expected: "",
		},
		{
			name:   "empty",
			config: conf.DefaultFormatterConfig(),
			markdown: &Spec{
				Description: NewNullValue(),
				Inputs:      []*InputSpec{},
				Outputs:     []*OutputSpec{},
			},
			expected: emptyActionExpected,
		},
		{
			name:   "full",
			config: conf.DefaultFormatterConfig(),
			markdown: &Spec{
				Description: NewNotNullValue("This is a test Custom Action for actdocs."),
				Inputs: []*InputSpec{
					{Name: "full-number", Default: NewNotNullValue("5"), Description: NewNotNullValue("The full number value."), Required: NewNotNullValue("false")},
				},
				Outputs: []*OutputSpec{
					{Name: "with-description", Description: NewNotNullValue("The Render value with description.")},
				},
			},
			expected: fullActionExpected,
		},
	}

	for _, tc := range cases {
		formatter := NewActionFormatter(tc.config)
		formatter.Spec = tc.markdown
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
