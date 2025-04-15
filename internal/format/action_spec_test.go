package format

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tmknom/actdocs/internal/util"
)

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
		spec := &ActionSpec{Description: tc.description}
		got := spec.toDescriptionMarkdown()

		if diff := cmp.Diff(got, tc.expected); diff != "" {
			t.Errorf("diff: %s", diff)
		}
	}
}

func TestActionSpec_toInputsMarkdown(t *testing.T) {
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
		spec := &ActionSpec{Inputs: tc.inputs}
		got := spec.toInputsMarkdown()

		if diff := cmp.Diff(got, tc.expected); diff != "" {
			t.Errorf("diff: %s", diff)
		}
	}
}

func TestActionSpec_toOutputsMarkdown(t *testing.T) {
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
		spec := &ActionSpec{Outputs: tc.outputs}
		got := spec.toOutputsMarkdown()

		if diff := cmp.Diff(got, tc.expected); diff != "" {
			t.Errorf("diff: %s", diff)
		}
	}
}

func TestActionInputSpec_toMarkdown(t *testing.T) {
	cases := []struct {
		name     string
		sut      *ActionInputSpec
		expected string
	}{
		{
			name: "single line",
			sut: &ActionInputSpec{
				Name:        "single-line",
				Default:     NewNotNullValue("Default value"),
				Description: NewNotNullValue("The test description."),
				Required:    NewNotNullValue("false"),
			},
			expected: "| single-line | The test description. | `Default value` | no |",
		},
		{
			name: "multi line",
			sut: &ActionInputSpec{
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

func TestActionOutputSpec_toMarkdown(t *testing.T) {
	cases := []struct {
		name     string
		sut      *ActionOutputSpec
		expected string
	}{
		{
			name: "single line",
			sut: &ActionOutputSpec{
				Name:        "single-line",
				Description: NewNotNullValue("The test description."),
			},
			expected: "| single-line | The test description. |",
		},
		{
			name: "multi line",
			sut: &ActionOutputSpec{
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
