package action

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRenderer_scan(t *testing.T) {
	cases := []struct {
		name     string
		spec     *Spec
		expected string
	}{
		{
			name: "full parameter",
			spec: &Spec{
				Description: NewNotNullValue("This is a test Custom Action for actdocs."),
				Inputs: []*InputSpec{
					{"full-number", NewNotNullValue("5"), NewNotNullValue("The full number value."), NewNotNullValue("false")},
				},
				Outputs: []*OutputSpec{
					{"with-description", NewNotNullValue("The Render value with description.")},
				},
			},
			expected: fullRenderExpected,
		},
	}

	template, err := os.Open(testBaseDir + "testdata/output.md")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	defer func(file *os.File) { err = file.Close() }(template)

	for _, tc := range cases {
		renderer := NewRenderer(template, false)
		got := renderer.scan(tc.spec)
		if diff := cmp.Diff(got, tc.expected); diff != "" {
			t.Errorf("%s: diff: %s", tc.name, diff)
		}
	}
}

const testBaseDir = "../../"

const fullRenderExpected = `# Output test

## Header

This is a header.

<!-- actdocs start -->

## Description

This is a test Custom Action for actdocs.

## Inputs

| Name | Description | Default | Required |
| :--- | :---------- | :------ | :------: |
| full-number | The full number value. | ` + "`5`" + ` | no |

## Outputs

| Name | Description |
| :--- | :---------- |
| with-description | The Render value with description. |

<!-- actdocs end -->

## Footer

This is a footer.
`
