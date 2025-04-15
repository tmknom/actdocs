package format

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWorkflowInputSpec_toMarkdown(t *testing.T) {
	cases := []struct {
		name     string
		sut      *WorkflowInputSpec
		expected string
	}{
		{
			name: "single line",
			sut: &WorkflowInputSpec{
				Name:        "single-line",
				Default:     NewNotNullValue("Default value"),
				Description: NewNotNullValue("The test description."),
				Required:    NewNotNullValue("false"),
				Type:        NewNotNullValue("string"),
			},
			expected: "| single-line | The test description. | `string` | `Default value` | no |",
		},
		{
			name: "multi line",
			sut: &WorkflowInputSpec{
				Name:        "multi-line",
				Default:     NewNotNullValue("{\n  \"key\": \"value\"\n}"),
				Description: NewNotNullValue("one\ntwo\nthree"),
				Required:    NewNotNullValue("true"),
				Type:        NewNotNullValue("number"),
			},
			expected: "| multi-line | <pre>one<br>two<br>three</pre> | `number` | <pre>{<br>  \"key\": \"value\"<br>}</pre> | yes |",
		},
	}

	for _, tc := range cases {
		got := tc.sut.toMarkdown()

		if diff := cmp.Diff(got, tc.expected); diff != "" {
			t.Errorf("diff: %s", diff)
		}
	}
}

func TestWorkflowSecretSpec_toMarkdown(t *testing.T) {
	cases := []struct {
		name     string
		sut      *WorkflowSecretSpec
		expected string
	}{
		{
			name: "single line",
			sut: &WorkflowSecretSpec{
				Name:        "single-line",
				Description: NewNotNullValue("The test description."),
				Required:    NewNotNullValue("false"),
			},
			expected: "| single-line | The test description. | no |",
		},
		{
			name: "multi line",
			sut: &WorkflowSecretSpec{
				Name:        "multi-line",
				Description: NewNotNullValue("one\ntwo\nthree"),
				Required:    NewNotNullValue("true"),
			},
			expected: "| multi-line | <pre>one<br>two<br>three</pre> | yes |",
		},
	}

	for _, tc := range cases {
		got := tc.sut.toMarkdown()

		if diff := cmp.Diff(got, tc.expected); diff != "" {
			t.Errorf("diff: %s", diff)
		}
	}
}

func TestWorkflowOutputSpec_toMarkdown(t *testing.T) {
	cases := []struct {
		name     string
		sut      *WorkflowOutputSpec
		expected string
	}{
		{
			name: "single line",
			sut: &WorkflowOutputSpec{
				Name:        "single-line",
				Description: NewNotNullValue("The test description."),
			},
			expected: "| single-line | The test description. |",
		},
		{
			name: "multi line",
			sut: &WorkflowOutputSpec{
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

func TestWorkflowPermissionSpec_toMarkdown(t *testing.T) {
	cases := []struct {
		name     string
		sut      *WorkflowPermissionSpec
		expected string
	}{
		{
			name:     "valid",
			sut:      &WorkflowPermissionSpec{Scope: "contents", Access: "write"},
			expected: "| contents | write |",
		},
	}

	for _, tc := range cases {
		got := tc.sut.toMarkdown()

		if diff := cmp.Diff(got, tc.expected); diff != "" {
			t.Errorf("diff: %s", diff)
		}
	}
}
