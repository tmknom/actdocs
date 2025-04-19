package workflow

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tmknom/actdocs/internal/conf"
)

func TestFormatter_Format(t *testing.T) {
	cases := []struct {
		name     string
		fixture  string
		expected string
	}{
		{
			name:     "basic",
			fixture:  complexWorkflowFixture,
			expected: basicWorkflowExpected,
		},
	}

	sortConfig := &conf.SortConfig{Sort: true}
	for _, tc := range cases {
		got, err := Generate(TestRawYaml(tc.fixture), conf.DefaultFormatterConfig(), sortConfig)
		if err != nil {
			t.Fatalf("%s: unexpected error: %s", tc.name, err)
		}

		if diff := cmp.Diff(got, tc.expected); diff != "" {
			t.Errorf("diff: %s", diff)
		}
	}
}

const basicWorkflowExpected = `## Inputs

| Name | Description | Type | Default | Required |
| :--- | :---------- | :--- | :------ | :------: |
| full-string | The full string value. | ` + "`string`" + ` | ` + "``" + ` | yes |
| empty |  | n/a | n/a | no |
| full-boolean | The full boolean value. | ` + "`boolean`" + ` | ` + "`true`" + ` | no |

## Secrets

N/A

## Outputs

N/A

## Permissions

N/A`
