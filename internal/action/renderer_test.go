package action

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRenderer_Render(t *testing.T) {
	cases := []struct {
		name     string
		spec     *Spec
		template string
		expected string
	}{
		{
			name: "all",
			spec: &Spec{
				Description: NewNotNullValue("This is a test Custom Action for actdocs."),
				Inputs: []*InputSpec{
					{"full-number", NewNotNullValue("5"), NewNotNullValue("The full number value."), NewNotNullValue("false")},
				},
				Outputs: []*OutputSpec{
					{"with-description", NewNotNullValue("The Render value with description.")},
				},
			},
			template: testBaseDir + "testdata/output.md",
			expected: fullRenderExpected,
		},
		{
			name: "sections",
			spec: &Spec{
				Description: NewNotNullValue("This is a test Custom Action for actdocs."),
				Inputs: []*InputSpec{
					{"full-number", NewNotNullValue("5"), NewNotNullValue("The full number value."), NewNotNullValue("false")},
				},
				Outputs: []*OutputSpec{
					{"with-description", NewNotNullValue("The Render value with description.")},
				},
			},
			template: testBaseDir + "testdata/inject-sections.md",
			expected: sectionsRenderExpected,
		},
	}

	for _, tc := range cases {
		template, err := os.Open(tc.template)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		defer func(file *os.File) { err = file.Close() }(template)

		renderer := NewRenderer(template, false)
		got := renderer.Render(tc.spec)
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

const sectionsRenderExpected = `# Output test

## Header

This is a header.

<!-- actdocs description start -->

## Description

This is a test Custom Action for actdocs.

<!-- actdocs description end -->

<!-- actdocs inputs start -->

## Inputs

| Name | Description | Default | Required |
| :--- | :---------- | :------ | :------: |
| full-number | The full number value. | ` + "`5`" + ` | no |

<!-- actdocs inputs end -->

<!-- actdocs outputs start -->

## Outputs

| Name | Description |
| :--- | :---------- |
| with-description | The Render value with description. |

<!-- actdocs outputs end -->

## Footer

This is a footer.
`
