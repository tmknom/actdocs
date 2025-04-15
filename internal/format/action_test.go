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
		json     *ActionMarkdown
		expected string
	}{
		{
			name: "empty",
			json: &ActionMarkdown{
				Description: NewNullValue(),
				Inputs:      []*ActionInputMarkdown{},
				Outputs:     []*ActionOutputMarkdown{},
			},
			expected: emptyActionExpectedJson,
		},
		{
			name: "full",
			json: &ActionMarkdown{
				Description: NewNotNullValue("This is a test Custom Action for actdocs."),
				Inputs: []*ActionInputMarkdown{
					{Name: "minimal", Default: NewNullValue(), Description: NewNullValue(), Required: NewNullValue()},
					{Name: "full", Default: NewNotNullValue("The string"), Description: NewNotNullValue("The input value."), Required: NewNotNullValue("true")},
				},
				Outputs: []*ActionOutputMarkdown{
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
		markdown *ActionMarkdown
		expected string
	}{
		{
			name:   "omit",
			config: &conf.FormatterConfig{Format: conf.DefaultFormat, Omit: true},
			markdown: &ActionMarkdown{
				Description: NewNullValue(),
				Inputs:      []*ActionInputMarkdown{},
				Outputs:     []*ActionOutputMarkdown{},
			},
			expected: "",
		},
		{
			name:   "empty",
			config: conf.DefaultFormatterConfig(),
			markdown: &ActionMarkdown{
				Description: NewNullValue(),
				Inputs:      []*ActionInputMarkdown{},
				Outputs:     []*ActionOutputMarkdown{},
			},
			expected: emptyActionExpected,
		},
		{
			name:   "full",
			config: conf.DefaultFormatterConfig(),
			markdown: &ActionMarkdown{
				Description: NewNotNullValue("This is a test Custom Action for actdocs."),
				Inputs: []*ActionInputMarkdown{
					{Name: "full-number", Default: NewNotNullValue("5"), Description: NewNotNullValue("The full number value."), Required: NewNotNullValue("false")},
				},
				Outputs: []*ActionOutputMarkdown{
					{Name: "with-description", Description: NewNotNullValue("The Render value with description.")},
				},
			},
			expected: fullActionExpected,
		},
	}

	for _, tc := range cases {
		formatter := NewActionFormatter(tc.config)
		formatter.ActionMarkdown = tc.markdown
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
		inputs   []*ActionInputMarkdown
		expected string
	}{
		{
			name:     "empty",
			inputs:   []*ActionInputMarkdown{},
			expected: "## Inputs\n\nN/A",
		},
		{
			name: "minimal",
			inputs: []*ActionInputMarkdown{
				{Name: "minimal", Default: NewNullValue(), Description: NewNullValue(), Required: NewNullValue()},
			},
			expected: "## Inputs\n\n| Name | Description | Default | Required |\n| :--- | :---------- | :------ | :------: |\n| minimal |  | n/a | no |",
		},
		{
			name: "single",
			inputs: []*ActionInputMarkdown{
				{Name: "single", Default: NewNotNullValue("5"), Description: NewNotNullValue("The number."), Required: NewNotNullValue("true")},
			},
			expected: "## Inputs\n\n| Name | Description | Default | Required |\n| :--- | :---------- | :------ | :------: |\n| single | The number. | `5` | yes |",
		},
		{
			name: "multiple",
			inputs: []*ActionInputMarkdown{
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
		outputs  []*ActionOutputMarkdown
		expected string
	}{
		{
			name:     "empty",
			outputs:  []*ActionOutputMarkdown{},
			expected: "## Outputs\n\nN/A",
		},
		{
			name: "minimal",
			outputs: []*ActionOutputMarkdown{
				{Name: "minimal", Description: NewNullValue()},
			},
			expected: "## Outputs\n\n| Name | Description |\n| :--- | :---------- |\n| minimal |  |",
		},
		{
			name: "single",
			outputs: []*ActionOutputMarkdown{
				{Name: "single", Description: NewNotNullValue("The test description.")},
			},
			expected: "## Outputs\n\n| Name | Description |\n| :--- | :---------- |\n| single | The test description. |",
		},
		{
			name: "multiple",
			outputs: []*ActionOutputMarkdown{
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

func TestActionInputMarkdown_toMarkdown(t *testing.T) {
	cases := []struct {
		name     string
		sut      *ActionInputMarkdown
		expected string
	}{
		{
			name: "single line",
			sut: &ActionInputMarkdown{
				Name:        "single-line",
				Default:     NewNotNullValue("Default value"),
				Description: NewNotNullValue("The test description."),
				Required:    NewNotNullValue("false"),
			},
			expected: "| single-line | The test description. | `Default value` | no |",
		},
		{
			name: "multi line",
			sut: &ActionInputMarkdown{
				Name:        "multi-line",
				Default:     NewNotNullValue("{\n  \"key\": \"value\"\n}"),
				Description: NewNotNullValue("one\ntwo\nthree"),
				Required:    NewNotNullValue("true"),
			},
			expected: "| multi-line | <pre>one<br>two<br>three</pre> | <pre>{<br>  \"key\": \"value\"<br>}</pre> | yes |",
		},
	}

	for _, tc := range cases {
		got := tc.sut.toMarkdown()

		if diff := cmp.Diff(got, tc.expected); diff != "" {
			t.Errorf("diff: %s", diff)
		}
	}
}

func TestActionOutputMarkdown_toMarkdown(t *testing.T) {
	cases := []struct {
		name     string
		sut      *ActionOutputMarkdown
		expected string
	}{
		{
			name: "single line",
			sut: &ActionOutputMarkdown{
				Name:        "single-line",
				Description: NewNotNullValue("The test description."),
			},
			expected: "| single-line | The test description. |",
		},
		{
			name: "multi line",
			sut: &ActionOutputMarkdown{
				Name:        "multi-line",
				Description: NewNotNullValue("one\ntwo\nthree"),
			},
			expected: "| multi-line | <pre>one<br>two<br>three</pre> |",
		},
	}

	for _, tc := range cases {
		got := tc.sut.toMarkdown()

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
