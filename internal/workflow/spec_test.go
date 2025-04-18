package workflow

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tmknom/actdocs/internal/conf"
)

func TestSpec_ToJson(t *testing.T) {
	cases := []struct {
		name     string
		sut      *Spec
		expected string
	}{
		{
			name: "empty",
			sut: &Spec{
				Inputs:      []*InputSpec{},
				Secrets:     []*SecretSpec{},
				Outputs:     []*OutputSpec{},
				Permissions: []*PermissionSpec{},
			},
			expected: emptyWorkflowExpectedJson,
		},
		{
			name: "full",
			sut: &Spec{
				Inputs: []*InputSpec{
					{Name: "minimal", Default: NewNullValue(), Description: NewNullValue(), Required: NewNullValue(), Type: NewNullValue()},
					{Name: "full", Default: NewNotNullValue("true"), Description: NewNotNullValue("The input value."), Required: NewNotNullValue("true"), Type: NewNotNullValue("boolean")},
				},
				Secrets: []*SecretSpec{
					{Name: "minimal", Description: NewNullValue(), Required: NewNullValue()},
					{Name: "full", Description: NewNotNullValue("The secret value."), Required: NewNotNullValue("true")},
				},
				Outputs: []*OutputSpec{
					{Name: "minimal", Description: NewNullValue()},
					{Name: "full", Description: NewNotNullValue("The output value.")},
				},
				Permissions: []*PermissionSpec{
					{Scope: "contents", Access: "write"},
					{Scope: "pull-requests", Access: "read"},
				},
			},
			expected: fullWorkflowExpectedJson,
		},
	}

	for _, tc := range cases {
		got := tc.sut.ToJson()
		if diff := cmp.Diff(got, tc.expected); diff != "" {
			t.Errorf("diff: %s", diff)
		}
	}
}

const emptyWorkflowExpectedJson = `{
  "inputs": [],
  "secrets": [],
  "outputs": [],
  "permissions": []
}`

const fullWorkflowExpectedJson = `{
  "inputs": [
    {
      "name": "minimal",
      "default": null,
      "description": null,
      "required": null,
      "type": null
    },
    {
      "name": "full",
      "default": "true",
      "description": "The input value.",
      "required": "true",
      "type": "boolean"
    }
  ],
  "secrets": [
    {
      "name": "minimal",
      "description": null,
      "required": null
    },
    {
      "name": "full",
      "description": "The secret value.",
      "required": "true"
    }
  ],
  "outputs": [
    {
      "name": "minimal",
      "description": null
    },
    {
      "name": "full",
      "description": "The output value."
    }
  ],
  "permissions": [
    {
      "scope": "contents",
      "access": "write"
    },
    {
      "scope": "pull-requests",
      "access": "read"
    }
  ]
}`

func TestSpec_ToMarkdown(t *testing.T) {
	cases := []struct {
		name     string
		config   *conf.FormatterConfig
		markdown *Spec
		expected string
	}{
		{
			name:   "omit",
			config: &conf.FormatterConfig{Format: conf.DefaultFormat, Omit: true},
			markdown: &Spec{
				Inputs:      []*InputSpec{},
				Secrets:     []*SecretSpec{},
				Outputs:     []*OutputSpec{},
				Permissions: []*PermissionSpec{},
			},
			expected: "",
		},
		{
			name:   "empty",
			config: conf.DefaultFormatterConfig(),
			markdown: &Spec{
				Inputs:      []*InputSpec{},
				Secrets:     []*SecretSpec{},
				Outputs:     []*OutputSpec{},
				Permissions: []*PermissionSpec{},
			},
			expected: emptyWorkflowExpected,
		},
		{
			name:   "full",
			config: conf.DefaultFormatterConfig(),
			markdown: &Spec{
				Inputs: []*InputSpec{
					{Name: "single", Default: NewNotNullValue("5"), Description: NewNotNullValue("The number."), Required: NewNotNullValue("true"), Type: NewNotNullValue("number")},
				},
				Secrets: []*SecretSpec{
					{Name: "single", Description: NewNotNullValue("The test description."), Required: NewNotNullValue("true")},
				},
				Outputs: []*OutputSpec{
					{Name: "single", Description: NewNotNullValue("The test description.")},
				},
				Permissions: []*PermissionSpec{
					{Scope: "contents", Access: "write"},
				},
			},
			expected: fullWorkflowExpected,
		},
	}

	for _, tc := range cases {
		got := tc.markdown.ToMarkdown(tc.config.Omit)
		if diff := cmp.Diff(got, tc.expected); diff != "" {
			t.Errorf("diff: %s", diff)
		}
	}
}

const emptyWorkflowExpected = `## Inputs

N/A

## Secrets

N/A

## Outputs

N/A

## Permissions

N/A`

const fullWorkflowExpected = `## Inputs

| Name | Description | Type | Default | Required |
| :--- | :---------- | :--- | :------ | :------: |
| single | The number. | ` + "`number`" + ` | ` + "`5`" + ` | yes |

## Secrets

| Name | Description | Required |
| :--- | :---------- | :------: |
| single | The test description. | yes |

## Outputs

| Name | Description |
| :--- | :---------- |
| single | The test description. |

## Permissions

| Scope | Access |
| :--- | :---- |
| contents | write |`

func TestSpec_ToInputsMarkdown(t *testing.T) {
	cases := []struct {
		name     string
		inputs   []*InputSpec
		omit     bool
		expected string
	}{
		{
			name:     "empty",
			inputs:   []*InputSpec{},
			omit:     false,
			expected: "## Inputs\n\nN/A\n\n",
		},
		{
			name: "minimal",
			inputs: []*InputSpec{
				{Name: "minimal", Default: NewNullValue(), Description: NewNullValue(), Required: NewNullValue(), Type: NewNullValue()},
			},
			omit:     false,
			expected: "## Inputs\n\n| Name | Description | Type | Default | Required |\n| :--- | :---------- | :--- | :------ | :------: |\n| minimal |  | n/a | n/a | no |\n\n",
		},
		{
			name: "single",
			inputs: []*InputSpec{
				{Name: "single", Default: NewNotNullValue("5"), Description: NewNotNullValue("The number."), Required: NewNotNullValue("true"), Type: NewNotNullValue("number")},
			},
			omit:     false,
			expected: "## Inputs\n\n| Name | Description | Type | Default | Required |\n| :--- | :---------- | :--- | :------ | :------: |\n| single | The number. | `number` | `5` | yes |\n\n",
		},
		{
			name: "multiple",
			inputs: []*InputSpec{
				{Name: "multiple-1", Default: NewNotNullValue("The string"), Description: NewNotNullValue("1"), Required: NewNotNullValue("false"), Type: NewNotNullValue("string")},
				{Name: "multiple-2", Default: NewNotNullValue("true"), Description: NewNotNullValue("2"), Required: NewNotNullValue("true"), Type: NewNotNullValue("boolean")},
			},
			omit:     false,
			expected: "## Inputs\n\n| Name | Description | Type | Default | Required |\n| :--- | :---------- | :--- | :------ | :------: |\n| multiple-1 | 1 | `string` | `The string` | no |\n| multiple-2 | 2 | `boolean` | `true` | yes |\n\n",
		},
	}

	for _, tc := range cases {
		spec := &Spec{Inputs: tc.inputs}
		got := spec.ToInputsMarkdown(tc.omit)

		if diff := cmp.Diff(got, tc.expected); diff != "" {
			t.Errorf("diff: %s", diff)
		}
	}
}

func TestSpec_ToSecretsMarkdown(t *testing.T) {
	cases := []struct {
		name     string
		secrets  []*SecretSpec
		omit     bool
		expected string
	}{
		{
			name:     "empty",
			secrets:  []*SecretSpec{},
			omit:     false,
			expected: "## Secrets\n\nN/A\n\n",
		},
		{
			name: "minimal",
			secrets: []*SecretSpec{
				{Name: "minimal", Description: NewNullValue(), Required: NewNullValue()},
			},
			omit:     false,
			expected: "## Secrets\n\n| Name | Description | Required |\n| :--- | :---------- | :------: |\n| minimal |  | no |\n\n",
		},
		{
			name: "single",
			secrets: []*SecretSpec{
				{Name: "single", Description: NewNotNullValue("The test description."), Required: NewNotNullValue("true")},
			},
			omit:     false,
			expected: "## Secrets\n\n| Name | Description | Required |\n| :--- | :---------- | :------: |\n| single | The test description. | yes |\n\n",
		},
		{
			name: "multiple",
			secrets: []*SecretSpec{
				{Name: "multiple-1", Description: NewNotNullValue("1"), Required: NewNotNullValue("false")},
				{Name: "multiple-2", Description: NewNotNullValue("2"), Required: NewNotNullValue("true")},
			},
			omit:     false,
			expected: "## Secrets\n\n| Name | Description | Required |\n| :--- | :---------- | :------: |\n| multiple-1 | 1 | no |\n| multiple-2 | 2 | yes |\n\n",
		},
	}

	for _, tc := range cases {
		spec := &Spec{Secrets: tc.secrets}
		got := spec.ToSecretsMarkdown(tc.omit)

		if diff := cmp.Diff(got, tc.expected); diff != "" {
			t.Errorf("diff: %s", diff)
		}
	}
}

func TestSpec_ToOutputsMarkdown(t *testing.T) {
	cases := []struct {
		name     string
		outputs  []*OutputSpec
		omit     bool
		expected string
	}{
		{
			name:     "empty",
			outputs:  []*OutputSpec{},
			omit:     false,
			expected: "## Outputs\n\nN/A\n\n",
		},
		{
			name: "minimal",
			outputs: []*OutputSpec{
				{Name: "minimal", Description: NewNullValue()},
			},
			omit:     false,
			expected: "## Outputs\n\n| Name | Description |\n| :--- | :---------- |\n| minimal |  |\n\n",
		},
		{
			name: "single",
			outputs: []*OutputSpec{
				{Name: "single", Description: NewNotNullValue("The test description.")},
			},
			omit:     false,
			expected: "## Outputs\n\n| Name | Description |\n| :--- | :---------- |\n| single | The test description. |\n\n",
		},
		{
			name: "multiple",
			outputs: []*OutputSpec{
				{Name: "multiple-1", Description: NewNotNullValue("1")},
				{Name: "multiple-2", Description: NewNotNullValue("2")},
			},
			omit:     false,
			expected: "## Outputs\n\n| Name | Description |\n| :--- | :---------- |\n| multiple-1 | 1 |\n| multiple-2 | 2 |\n\n",
		},
	}

	for _, tc := range cases {
		spec := &Spec{Outputs: tc.outputs}
		got := spec.ToOutputsMarkdown(tc.omit)

		if diff := cmp.Diff(got, tc.expected); diff != "" {
			t.Errorf("diff: %s", diff)
		}
	}
}

func TestSpec_ToPermissionsMarkdown(t *testing.T) {
	cases := []struct {
		name        string
		permissions []*PermissionSpec
		omit        bool
		expected    string
	}{
		{
			name:        "empty",
			permissions: []*PermissionSpec{},
			omit:        false,
			expected:    "## Permissions\n\nN/A\n\n",
		},
		{
			name: "single",
			permissions: []*PermissionSpec{
				{Scope: "contents", Access: "write"},
			},
			omit:     false,
			expected: "## Permissions\n\n| Scope | Access |\n| :--- | :---- |\n| contents | write |\n\n",
		},
		{
			name: "multiple",
			permissions: []*PermissionSpec{
				{Scope: "contents", Access: "write"},
				{Scope: "pull-requests", Access: "read"},
			},
			omit:     false,
			expected: "## Permissions\n\n| Scope | Access |\n| :--- | :---- |\n| contents | write |\n| pull-requests | read |\n\n",
		},
	}

	for _, tc := range cases {
		spec := &Spec{Permissions: tc.permissions}
		got := spec.ToPermissionsMarkdown(tc.omit)

		if diff := cmp.Diff(got, tc.expected); diff != "" {
			t.Errorf("diff: %s", diff)
		}
	}
}

func TestInputSpec_toMarkdown(t *testing.T) {
	cases := []struct {
		name     string
		sut      *InputSpec
		expected string
	}{
		{
			name: "single line",
			sut: &InputSpec{
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
			sut: &InputSpec{
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

func TestSecretSpec_toMarkdown(t *testing.T) {
	cases := []struct {
		name     string
		sut      *SecretSpec
		expected string
	}{
		{
			name: "single line",
			sut: &SecretSpec{
				Name:        "single-line",
				Description: NewNotNullValue("The test description."),
				Required:    NewNotNullValue("false"),
			},
			expected: "| single-line | The test description. | no |",
		},
		{
			name: "multi line",
			sut: &SecretSpec{
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

func TestOutputSpec_toMarkdown(t *testing.T) {
	cases := []struct {
		name     string
		sut      *OutputSpec
		expected string
	}{
		{
			name: "single line",
			sut: &OutputSpec{
				Name:        "single-line",
				Description: NewNotNullValue("The test description."),
			},
			expected: "| single-line | The test description. |",
		},
		{
			name: "multi line",
			sut: &OutputSpec{
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

func TestPermissionSpec_toMarkdown(t *testing.T) {
	cases := []struct {
		name     string
		sut      *PermissionSpec
		expected string
	}{
		{
			name:     "valid",
			sut:      &PermissionSpec{Scope: "contents", Access: "write"},
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
