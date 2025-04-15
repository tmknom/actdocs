package format

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWorkflowSpec_toInputsMarkdown(t *testing.T) {
	cases := []struct {
		name     string
		inputs   []*WorkflowInputSpec
		expected string
	}{
		{
			name:     "empty",
			inputs:   []*WorkflowInputSpec{},
			expected: "## Inputs\n\nN/A",
		},
		{
			name: "minimal",
			inputs: []*WorkflowInputSpec{
				{Name: "minimal", Default: NewNullValue(), Description: NewNullValue(), Required: NewNullValue(), Type: NewNullValue()},
			},
			expected: "## Inputs\n\n| Name | Description | Type | Default | Required |\n| :--- | :---------- | :--- | :------ | :------: |\n| minimal |  | n/a | n/a | no |",
		},
		{
			name: "single",
			inputs: []*WorkflowInputSpec{
				{Name: "single", Default: NewNotNullValue("5"), Description: NewNotNullValue("The number."), Required: NewNotNullValue("true"), Type: NewNotNullValue("number")},
			},
			expected: "## Inputs\n\n| Name | Description | Type | Default | Required |\n| :--- | :---------- | :--- | :------ | :------: |\n| single | The number. | `number` | `5` | yes |",
		},
		{
			name: "multiple",
			inputs: []*WorkflowInputSpec{
				{Name: "multiple-1", Default: NewNotNullValue("The string"), Description: NewNotNullValue("1"), Required: NewNotNullValue("false"), Type: NewNotNullValue("string")},
				{Name: "multiple-2", Default: NewNotNullValue("true"), Description: NewNotNullValue("2"), Required: NewNotNullValue("true"), Type: NewNotNullValue("boolean")},
			},
			expected: "## Inputs\n\n| Name | Description | Type | Default | Required |\n| :--- | :---------- | :--- | :------ | :------: |\n| multiple-1 | 1 | `string` | `The string` | no |\n| multiple-2 | 2 | `boolean` | `true` | yes |",
		},
	}

	for _, tc := range cases {
		spec := &WorkflowSpec{Inputs: tc.inputs}
		got := spec.toInputsMarkdown()

		if diff := cmp.Diff(got, tc.expected); diff != "" {
			t.Errorf("diff: %s", diff)
		}
	}
}

func TestWorkflowSpec_toSecretsMarkdown(t *testing.T) {
	cases := []struct {
		name     string
		secrets  []*WorkflowSecretSpec
		expected string
	}{
		{
			name:     "empty",
			secrets:  []*WorkflowSecretSpec{},
			expected: "## Secrets\n\nN/A",
		},
		{
			name: "minimal",
			secrets: []*WorkflowSecretSpec{
				{Name: "minimal", Description: NewNullValue(), Required: NewNullValue()},
			},
			expected: "## Secrets\n\n| Name | Description | Required |\n| :--- | :---------- | :------: |\n| minimal |  | no |",
		},
		{
			name: "single",
			secrets: []*WorkflowSecretSpec{
				{Name: "single", Description: NewNotNullValue("The test description."), Required: NewNotNullValue("true")},
			},
			expected: "## Secrets\n\n| Name | Description | Required |\n| :--- | :---------- | :------: |\n| single | The test description. | yes |",
		},
		{
			name: "multiple",
			secrets: []*WorkflowSecretSpec{
				{Name: "multiple-1", Description: NewNotNullValue("1"), Required: NewNotNullValue("false")},
				{Name: "multiple-2", Description: NewNotNullValue("2"), Required: NewNotNullValue("true")},
			},
			expected: "## Secrets\n\n| Name | Description | Required |\n| :--- | :---------- | :------: |\n| multiple-1 | 1 | no |\n| multiple-2 | 2 | yes |",
		},
	}

	for _, tc := range cases {
		spec := &WorkflowSpec{Secrets: tc.secrets}
		got := spec.toSecretsMarkdown()

		if diff := cmp.Diff(got, tc.expected); diff != "" {
			t.Errorf("diff: %s", diff)
		}
	}
}

func TestWorkflowSpec_toOutputsMarkdown(t *testing.T) {
	cases := []struct {
		name     string
		outputs  []*WorkflowOutputSpec
		expected string
	}{
		{
			name:     "empty",
			outputs:  []*WorkflowOutputSpec{},
			expected: "## Outputs\n\nN/A",
		},
		{
			name: "minimal",
			outputs: []*WorkflowOutputSpec{
				{Name: "minimal", Description: NewNullValue()},
			},
			expected: "## Outputs\n\n| Name | Description |\n| :--- | :---------- |\n| minimal |  |",
		},
		{
			name: "single",
			outputs: []*WorkflowOutputSpec{
				{Name: "single", Description: NewNotNullValue("The test description.")},
			},
			expected: "## Outputs\n\n| Name | Description |\n| :--- | :---------- |\n| single | The test description. |",
		},
		{
			name: "multiple",
			outputs: []*WorkflowOutputSpec{
				{Name: "multiple-1", Description: NewNotNullValue("1")},
				{Name: "multiple-2", Description: NewNotNullValue("2")},
			},
			expected: "## Outputs\n\n| Name | Description |\n| :--- | :---------- |\n| multiple-1 | 1 |\n| multiple-2 | 2 |",
		},
	}

	for _, tc := range cases {
		spec := &WorkflowSpec{Outputs: tc.outputs}
		got := spec.toOutputsMarkdown()

		if diff := cmp.Diff(got, tc.expected); diff != "" {
			t.Errorf("diff: %s", diff)
		}
	}
}

func TestWorkflowSpec_toPermissionsMarkdown(t *testing.T) {
	cases := []struct {
		name        string
		permissions []*WorkflowPermissionSpec
		expected    string
	}{
		{
			name:        "empty",
			permissions: []*WorkflowPermissionSpec{},
			expected:    "## Permissions\n\nN/A",
		},
		{
			name: "single",
			permissions: []*WorkflowPermissionSpec{
				{Scope: "contents", Access: "write"},
			},
			expected: "## Permissions\n\n| Scope | Access |\n| :--- | :---- |\n| contents | write |",
		},
		{
			name: "multiple",
			permissions: []*WorkflowPermissionSpec{
				{Scope: "contents", Access: "write"},
				{Scope: "pull-requests", Access: "read"},
			},
			expected: "## Permissions\n\n| Scope | Access |\n| :--- | :---- |\n| contents | write |\n| pull-requests | read |",
		},
	}

	for _, tc := range cases {
		spec := &WorkflowSpec{Permissions: tc.permissions}
		got := spec.toPermissionsMarkdown()

		if diff := cmp.Diff(got, tc.expected); diff != "" {
			t.Errorf("diff: %s", diff)
		}
	}
}

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
