package workflow

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

//type Spec struct {
//	Inputs      []*InputSpec      `json:"inputs"`
//	Secrets     []*SecretSpec     `json:"secrets"`
//	Outputs     []*OutputSpec     `json:"outputs"`
//	Permissions []*PermissionSpec `json:"permissions"`
//}

func TestRenderer_scan(t *testing.T) {
	cases := []struct {
		name     string
		spec     *Spec
		expected string
	}{
		{
			name: "full parameter",
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
