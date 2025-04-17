package action

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tmknom/actdocs/internal/conf"
	"github.com/tmknom/actdocs/internal/util"
)

func TestSpec_ToJson(t *testing.T) {
	cases := []struct {
		name     string
		sut      *Spec
		expected string
	}{
		{
			name: "empty",
			sut: &Spec{
				Description: NewNullValue(),
				Inputs:      []*InputSpec{},
				Outputs:     []*OutputSpec{},
			},
			expected: emptyActionExpectedJson,
		},
		{
			name: "full",
			sut: &Spec{
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
		got := tc.sut.ToJson()
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

func TestSpec_ToMarkdown(t *testing.T) {
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
		got := tc.markdown.ToMarkdown(tc.config.Omit)
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

func TestSpec_ToDescriptionMarkdown(t *testing.T) {
	cases := []struct {
		name        string
		description *util.NullString
		omit        bool
		expected    string
	}{
		{
			name:        "omit",
			description: NewNullValue(),
			omit:        true,
			expected:    "",
		},
		{
			name:        "null value",
			description: NewNullValue(),
			omit:        false,
			expected:    "## Description\n\nN/A\n\n",
		},
		{
			name:        "valid value",
			description: NewNotNullValue("The valid."),
			omit:        false,
			expected:    "## Description\n\nThe valid.\n\n",
		},
	}

	for _, tc := range cases {
		spec := &Spec{Description: tc.description}
		got := spec.ToDescriptionMarkdown(tc.omit)

		if diff := cmp.Diff(got, tc.expected); diff != "" {
			t.Errorf("diff: %s", diff)
		}
	}
}

func TestSpec_toInputsMarkdown(t *testing.T) {
	cases := []struct {
		name     string
		inputs   []*InputSpec
		omit     bool
		expected string
	}{
		{
			name:     "empty",
			inputs:   []*InputSpec{},
			omit:     false,
			expected: "## Inputs\n\nN/A\n\n",
		},
		{
			name: "minimal",
			inputs: []*InputSpec{
				{Name: "minimal", Default: NewNullValue(), Description: NewNullValue(), Required: NewNullValue()},
			},
			omit:     false,
			expected: "## Inputs\n\n| Name | Description | Default | Required |\n| :--- | :---------- | :------ | :------: |\n| minimal |  | n/a | no |\n\n",
		},
		{
			name: "single",
			inputs: []*InputSpec{
				{Name: "single", Default: NewNotNullValue("5"), Description: NewNotNullValue("The number."), Required: NewNotNullValue("true")},
			},
			omit:     false,
			expected: "## Inputs\n\n| Name | Description | Default | Required |\n| :--- | :---------- | :------ | :------: |\n| single | The number. | `5` | yes |\n\n",
		},
		{
			name: "multiple",
			inputs: []*InputSpec{
				{Name: "multiple-1", Default: NewNotNullValue("The string"), Description: NewNotNullValue("1"), Required: NewNotNullValue("false")},
				{Name: "multiple-2", Default: NewNotNullValue("true"), Description: NewNotNullValue("2"), Required: NewNotNullValue("true")},
			},
			omit:     false,
			expected: "## Inputs\n\n| Name | Description | Default | Required |\n| :--- | :---------- | :------ | :------: |\n| multiple-1 | 1 | `The string` | no |\n| multiple-2 | 2 | `true` | yes |\n\n",
		},
	}

	for _, tc := range cases {
		spec := &Spec{Inputs: tc.inputs}
		got := spec.ToInputsMarkdown(tc.omit)

		if diff := cmp.Diff(got, tc.expected); diff != "" {
			t.Errorf("diff: %s", diff)
		}
	}
}

func TestSpec_toOutputsMarkdown(t *testing.T) {
	cases := []struct {
		name     string
		outputs  []*OutputSpec
		omit     bool
		expected string
	}{
		{
			name:     "empty",
			outputs:  []*OutputSpec{},
			omit:     false,
			expected: "## Outputs\n\nN/A\n\n",
		},
		{
			name: "minimal",
			outputs: []*OutputSpec{
				{Name: "minimal", Description: NewNullValue()},
			},
			omit:     false,
			expected: "## Outputs\n\n| Name | Description |\n| :--- | :---------- |\n| minimal |  |\n\n",
		},
		{
			name: "single",
			outputs: []*OutputSpec{
				{Name: "single", Description: NewNotNullValue("The test description.")},
			},
			omit:     false,
			expected: "## Outputs\n\n| Name | Description |\n| :--- | :---------- |\n| single | The test description. |\n\n",
		},
		{
			name: "multiple",
			outputs: []*OutputSpec{
				{Name: "multiple-1", Description: NewNotNullValue("1")},
				{Name: "multiple-2", Description: NewNotNullValue("2")},
			},
			omit:     false,
			expected: "## Outputs\n\n| Name | Description |\n| :--- | :---------- |\n| multiple-1 | 1 |\n| multiple-2 | 2 |\n\n",
		},
	}

	for _, tc := range cases {
		spec := &Spec{Outputs: tc.outputs}
		got := spec.ToOutputsMarkdown(tc.omit)

		if diff := cmp.Diff(got, tc.expected); diff != "" {
			t.Errorf("diff: %s", diff)
		}
	}
}

func TestInputSpec_toMarkdown(t *testing.T) {
	cases := []struct {
		name     string
		sut      *InputSpec
		expected string
	}{
		{
			name: "single line",
			sut: &InputSpec{
				Name:        "single-line",
				Default:     NewNotNullValue("Default value"),
				Description: NewNotNullValue("The test description."),
				Required:    NewNotNullValue("false"),
			},
			expected: "| single-line | The test description. | `Default value` | no |",
		},
		{
			name: "multi line",
			sut: &InputSpec{
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

func TestOutputSpec_toMarkdown(t *testing.T) {
	cases := []struct {
		name     string
		sut      *OutputSpec
		expected string
	}{
		{
			name: "single line",
			sut: &OutputSpec{
				Name:        "single-line",
				Description: NewNotNullValue("The test description."),
			},
			expected: "| single-line | The test description. |",
		},
		{
			name: "multi line",
			sut: &OutputSpec{
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
