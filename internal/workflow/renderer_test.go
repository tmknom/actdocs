package workflow

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
				Inputs: []*InputSpec{
					{Name: "full-number", Default: NewNotNullValue("5"), Description: NewNotNullValue("The full number value."), Required: NewNotNullValue("false"), Type: NewNotNullValue("number")},
				},
				Secrets: []*SecretSpec{
					{Name: "full", Description: NewNotNullValue("The secret value."), Required: NewNotNullValue("true")},
				},
				Outputs: []*OutputSpec{
					{Name: "full", Description: NewNotNullValue("The output value.")},
				},
				Permissions: []*PermissionSpec{
					{Scope: "contents", Access: "write"},
				},
			},
			template: testBaseDir + "testdata/output.md",
			expected: fullRenderExpected,
		},
		{
			name: "sections",
			spec: &Spec{
				Inputs: []*InputSpec{
					{Name: "full-number", Default: NewNotNullValue("5"), Description: NewNotNullValue("The full number value."), Required: NewNotNullValue("false"), Type: NewNotNullValue("number")},
				},
				Secrets: []*SecretSpec{
					{Name: "full", Description: NewNotNullValue("The secret value."), Required: NewNotNullValue("true")},
				},
				Outputs: []*OutputSpec{
					{Name: "full", Description: NewNotNullValue("The output value.")},
				},
				Permissions: []*PermissionSpec{
					{Scope: "contents", Access: "write"},
				},
			},
			template: testBaseDir + "testdata/inject-workflow-sections.md",
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

## Inputs

| Name | Description | Type | Default | Required |
| :--- | :---------- | :--- | :------ | :------: |
| full-number | The full number value. | ` + "`number`" + ` | ` + "`5`" + ` | no |

## Secrets

| Name | Description | Required |
| :--- | :---------- | :------: |
| full | The secret value. | yes |

## Outputs

| Name | Description |
| :--- | :---------- |
| full | The output value. |

## Permissions

| Scope | Access |
| :--- | :---- |
| contents | write |

<!-- actdocs end -->

## Footer

This is a footer.
`

const sectionsRenderExpected = `# Output test

## Header

This is a header.

<!-- actdocs inputs start -->

## Inputs

| Name | Description | Type | Default | Required |
| :--- | :---------- | :--- | :------ | :------: |
| full-number | The full number value. | ` + "`number`" + ` | ` + "`5`" + ` | no |

<!-- actdocs inputs end -->

<!-- actdocs secrets start -->

## Secrets

| Name | Description | Required |
| :--- | :---------- | :------: |
| full | The secret value. | yes |

<!-- actdocs secrets end -->

<!-- actdocs outputs start -->

## Outputs

| Name | Description |
| :--- | :---------- |
| full | The output value. |

<!-- actdocs outputs end -->

<!-- actdocs permissions start -->

## Permissions

| Scope | Access |
| :--- | :---- |
| contents | write |

<!-- actdocs permissions end -->

## Footer

This is a footer.
`
