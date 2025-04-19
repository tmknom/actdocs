package action

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tmknom/actdocs/internal/conf"
)

func TestGenerate(t *testing.T) {
	cases := []struct {
		name     string
		fixture  string
		expected string
	}{
		{
			name:     "basic",
			fixture:  complexActionFixture,
			expected: formatExpected,
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

const formatExpected = `## Description

This is a test Custom Action for actdocs.

## Inputs

| Name | Description | Default | Required |
| :--- | :---------- | :------ | :------: |
| full-string | The full string value. | ` + "`Default value`" + ` | yes |
| empty |  | n/a | no |
| full-boolean | The full boolean value. | ` + "`true`" + ` | no |

## Outputs

| Name | Description |
| :--- | :---------- |
| only-value |  |
| with-description | The Render value with description. |`
