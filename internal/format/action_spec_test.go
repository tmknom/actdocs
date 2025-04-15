package format

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

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
