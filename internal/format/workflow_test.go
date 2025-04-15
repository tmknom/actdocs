package format

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tmknom/actdocs/internal/conf"
	"github.com/tmknom/actdocs/internal/parse"
)

func TestWorkflowFormatter_Format(t *testing.T) {
	cases := []struct {
		name     string
		ast      *parse.WorkflowAST
		expected string
	}{
		{
			name: "basic",
			ast: &parse.WorkflowAST{
				Inputs: []*parse.WorkflowInput{
					{Name: "foo", Default: NewNotNullValue("Default"), Description: NewNotNullValue("The inputs."), Required: NewNotNullValue("false"), Type: NewNotNullValue("string")},
				},
				Secrets: []*parse.WorkflowSecret{
					{Name: "bar", Description: NewNotNullValue("The secrets."), Required: NewNotNullValue("false")},
				},
				Outputs: []*parse.WorkflowOutput{
					{Name: "baz", Description: NewNotNullValue("The outputs.")},
				},
				Permissions: []*parse.WorkflowPermission{
					{Scope: "contents", Access: "write"},
				},
			},
			expected: basicWorkflowExpected,
		},
	}

	for _, tc := range cases {
		formatter := NewWorkflowFormatter(conf.DefaultFormatterConfig())
		got := formatter.Format(tc.ast)
		if diff := cmp.Diff(got, tc.expected); diff != "" {
			t.Errorf("diff: %s", diff)
		}
	}
}

const basicWorkflowExpected = `## Inputs

| Name | Description | Type | Default | Required |
| :--- | :---------- | :--- | :------ | :------: |
| foo | The inputs. | ` + "`string`" + ` | ` + "`Default`" + ` | no |

## Secrets

| Name | Description | Required |
| :--- | :---------- | :------: |
| bar | The secrets. | no |

## Outputs

| Name | Description |
| :--- | :---------- |
| baz | The outputs. |

## Permissions

| Scope | Access |
| :--- | :---- |
| contents | write |`

func TestWorkflowFormatter_ToJson(t *testing.T) {
	cases := []struct {
		name     string
		json     *WorkflowJson
		expected string
	}{
		{
			name: "empty",
			json: &WorkflowJson{
				Inputs:      []*WorkflowInputJson{},
				Secrets:     []*WorkflowSecretJson{},
				Outputs:     []*WorkflowOutputJson{},
				Permissions: []*WorkflowPermissionJson{},
			},
			expected: emptyWorkflowExpectedJson,
		},
		{
			name: "full",
			json: &WorkflowJson{
				Inputs: []*WorkflowInputJson{
					{Name: "minimal", Default: NewNullValue(), Description: NewNullValue(), Required: NewNullValue(), Type: NewNullValue()},
					{Name: "full", Default: NewNotNullValue("true"), Description: NewNotNullValue("The input value."), Required: NewNotNullValue("true"), Type: NewNotNullValue("boolean")},
				},
				Secrets: []*WorkflowSecretJson{
					{Name: "minimal", Description: NewNullValue(), Required: NewNullValue()},
					{Name: "full", Description: NewNotNullValue("The secret value."), Required: NewNotNullValue("true")},
				},
				Outputs: []*WorkflowOutputJson{
					{Name: "minimal", Description: NewNullValue()},
					{Name: "full", Description: NewNotNullValue("The output value.")},
				},
				Permissions: []*WorkflowPermissionJson{
					{Scope: "contents", Access: "write"},
					{Scope: "pull-requests", Access: "read"},
				},
			},
			expected: fullWorkflowExpectedJson,
		},
	}

	for _, tc := range cases {
		formatter := NewWorkflowFormatter(conf.DefaultFormatterConfig())
		got := formatter.ToJson(tc.json)
		if diff := cmp.Diff(got, tc.expected); diff != "" {
			t.Errorf("diff: %s", diff)
		}
	}
}

const emptyWorkflowExpectedJson = `{
  "inputs": [],
  "outputs": [],
  "secrets": [],
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

func TestWorkflowFormatter_ToMarkdown(t *testing.T) {
	cases := []struct {
		name     string
		config   *conf.FormatterConfig
		markdown *WorkflowMarkdown
		expected string
	}{
		{
			name:   "omit",
			config: &conf.FormatterConfig{Format: conf.DefaultFormat, Omit: true},
			markdown: &WorkflowMarkdown{
				Inputs:      []*WorkflowInputMarkdown{},
				Secrets:     []*WorkflowSecretMarkdown{},
				Outputs:     []*WorkflowOutputMarkdown{},
				Permissions: []*WorkflowPermissionMarkdown{},
			},
			expected: "",
		},
		{
			name:   "empty",
			config: conf.DefaultFormatterConfig(),
			markdown: &WorkflowMarkdown{
				Inputs:      []*WorkflowInputMarkdown{},
				Secrets:     []*WorkflowSecretMarkdown{},
				Outputs:     []*WorkflowOutputMarkdown{},
				Permissions: []*WorkflowPermissionMarkdown{},
			},
			expected: emptyWorkflowExpected,
		},
		{
			name:   "full",
			config: conf.DefaultFormatterConfig(),
			markdown: &WorkflowMarkdown{
				Inputs: []*WorkflowInputMarkdown{
					{Name: "single", Default: NewNotNullValue("5"), Description: NewNotNullValue("The number."), Required: NewNotNullValue("true"), Type: NewNotNullValue("number")},
				},
				Secrets: []*WorkflowSecretMarkdown{
					{Name: "single", Description: NewNotNullValue("The test description."), Required: NewNotNullValue("true")},
				},
				Outputs: []*WorkflowOutputMarkdown{
					{Name: "single", Description: NewNotNullValue("The test description.")},
				},
				Permissions: []*WorkflowPermissionMarkdown{
					{Scope: "contents", Access: "write"},
				},
			},
			expected: fullWorkflowExpected,
		},
	}

	for _, tc := range cases {
		formatter := NewWorkflowFormatter(tc.config)
		formatter.WorkflowMarkdown = tc.markdown
		got := formatter.ToMarkdown(tc.markdown, tc.config)
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

func TestWorkflowFormatter_toInputsMarkdown(t *testing.T) {
	cases := []struct {
		name     string
		inputs   []*WorkflowInputMarkdown
		expected string
	}{
		{
			name:     "empty",
			inputs:   []*WorkflowInputMarkdown{},
			expected: "## Inputs\n\nN/A",
		},
		{
			name: "minimal",
			inputs: []*WorkflowInputMarkdown{
				{Name: "minimal", Default: NewNullValue(), Description: NewNullValue(), Required: NewNullValue(), Type: NewNullValue()},
			},
			expected: "## Inputs\n\n| Name | Description | Type | Default | Required |\n| :--- | :---------- | :--- | :------ | :------: |\n| minimal |  | n/a | n/a | no |",
		},
		{
			name: "single",
			inputs: []*WorkflowInputMarkdown{
				{Name: "single", Default: NewNotNullValue("5"), Description: NewNotNullValue("The number."), Required: NewNotNullValue("true"), Type: NewNotNullValue("number")},
			},
			expected: "## Inputs\n\n| Name | Description | Type | Default | Required |\n| :--- | :---------- | :--- | :------ | :------: |\n| single | The number. | `number` | `5` | yes |",
		},
		{
			name: "multiple",
			inputs: []*WorkflowInputMarkdown{
				{Name: "multiple-1", Default: NewNotNullValue("The string"), Description: NewNotNullValue("1"), Required: NewNotNullValue("false"), Type: NewNotNullValue("string")},
				{Name: "multiple-2", Default: NewNotNullValue("true"), Description: NewNotNullValue("2"), Required: NewNotNullValue("true"), Type: NewNotNullValue("boolean")},
			},
			expected: "## Inputs\n\n| Name | Description | Type | Default | Required |\n| :--- | :---------- | :--- | :------ | :------: |\n| multiple-1 | 1 | `string` | `The string` | no |\n| multiple-2 | 2 | `boolean` | `true` | yes |",
		},
	}

	for _, tc := range cases {
		formatter := NewWorkflowFormatter(conf.DefaultFormatterConfig())
		got := formatter.toInputsMarkdown(tc.inputs)

		if diff := cmp.Diff(got, tc.expected); diff != "" {
			t.Errorf("diff: %s", diff)
		}
	}
}

func TestWorkflowFormatter_toSecretsMarkdown(t *testing.T) {
	cases := []struct {
		name     string
		outputs  []*WorkflowSecretMarkdown
		expected string
	}{
		{
			name:     "empty",
			outputs:  []*WorkflowSecretMarkdown{},
			expected: "## Secrets\n\nN/A",
		},
		{
			name: "minimal",
			outputs: []*WorkflowSecretMarkdown{
				{Name: "minimal", Description: NewNullValue(), Required: NewNullValue()},
			},
			expected: "## Secrets\n\n| Name | Description | Required |\n| :--- | :---------- | :------: |\n| minimal |  | no |",
		},
		{
			name: "single",
			outputs: []*WorkflowSecretMarkdown{
				{Name: "single", Description: NewNotNullValue("The test description."), Required: NewNotNullValue("true")},
			},
			expected: "## Secrets\n\n| Name | Description | Required |\n| :--- | :---------- | :------: |\n| single | The test description. | yes |",
		},
		{
			name: "multiple",
			outputs: []*WorkflowSecretMarkdown{
				{Name: "multiple-1", Description: NewNotNullValue("1"), Required: NewNotNullValue("false")},
				{Name: "multiple-2", Description: NewNotNullValue("2"), Required: NewNotNullValue("true")},
			},
			expected: "## Secrets\n\n| Name | Description | Required |\n| :--- | :---------- | :------: |\n| multiple-1 | 1 | no |\n| multiple-2 | 2 | yes |",
		},
	}

	for _, tc := range cases {
		formatter := NewWorkflowFormatter(conf.DefaultFormatterConfig())
		got := formatter.toSecretsMarkdown(tc.outputs)

		if diff := cmp.Diff(got, tc.expected); diff != "" {
			t.Errorf("diff: %s", diff)
		}
	}
}

func TestWorkflowFormatter_toOutputsMarkdown(t *testing.T) {
	cases := []struct {
		name     string
		outputs  []*WorkflowOutputMarkdown
		expected string
	}{
		{
			name:     "empty",
			outputs:  []*WorkflowOutputMarkdown{},
			expected: "## Outputs\n\nN/A",
		},
		{
			name: "minimal",
			outputs: []*WorkflowOutputMarkdown{
				{Name: "minimal", Description: NewNullValue()},
			},
			expected: "## Outputs\n\n| Name | Description |\n| :--- | :---------- |\n| minimal |  |",
		},
		{
			name: "single",
			outputs: []*WorkflowOutputMarkdown{
				{Name: "single", Description: NewNotNullValue("The test description.")},
			},
			expected: "## Outputs\n\n| Name | Description |\n| :--- | :---------- |\n| single | The test description. |",
		},
		{
			name: "multiple",
			outputs: []*WorkflowOutputMarkdown{
				{Name: "multiple-1", Description: NewNotNullValue("1")},
				{Name: "multiple-2", Description: NewNotNullValue("2")},
			},
			expected: "## Outputs\n\n| Name | Description |\n| :--- | :---------- |\n| multiple-1 | 1 |\n| multiple-2 | 2 |",
		},
	}

	for _, tc := range cases {
		formatter := NewWorkflowFormatter(conf.DefaultFormatterConfig())
		got := formatter.toOutputsMarkdown(tc.outputs)

		if diff := cmp.Diff(got, tc.expected); diff != "" {
			t.Errorf("diff: %s", diff)
		}
	}
}

func TestWorkflowFormatter_toPermissionsMarkdown(t *testing.T) {
	cases := []struct {
		name     string
		outputs  []*WorkflowPermissionMarkdown
		expected string
	}{
		{
			name:     "empty",
			outputs:  []*WorkflowPermissionMarkdown{},
			expected: "## Permissions\n\nN/A",
		},
		{
			name: "single",
			outputs: []*WorkflowPermissionMarkdown{
				{Scope: "contents", Access: "write"},
			},
			expected: "## Permissions\n\n| Scope | Access |\n| :--- | :---- |\n| contents | write |",
		},
		{
			name: "multiple",
			outputs: []*WorkflowPermissionMarkdown{
				{Scope: "contents", Access: "write"},
				{Scope: "pull-requests", Access: "read"},
			},
			expected: "## Permissions\n\n| Scope | Access |\n| :--- | :---- |\n| contents | write |\n| pull-requests | read |",
		},
	}

	for _, tc := range cases {
		formatter := NewWorkflowFormatter(conf.DefaultFormatterConfig())
		got := formatter.toPermissionsMarkdown(tc.outputs)

		if diff := cmp.Diff(got, tc.expected); diff != "" {
			t.Errorf("diff: %s", diff)
		}
	}
}

func TestWorkflowInputMarkdown_toMarkdown(t *testing.T) {
	cases := []struct {
		name     string
		sut      *WorkflowInputMarkdown
		expected string
	}{
		{
			name: "single line",
			sut: &WorkflowInputMarkdown{
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
			sut: &WorkflowInputMarkdown{
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

func TestWorkflowSecretMarkdown_toMarkdown(t *testing.T) {
	cases := []struct {
		name     string
		sut      *WorkflowSecretMarkdown
		expected string
	}{
		{
			name: "single line",
			sut: &WorkflowSecretMarkdown{
				Name:        "single-line",
				Description: NewNotNullValue("The test description."),
				Required:    NewNotNullValue("false"),
			},
			expected: "| single-line | The test description. | no |",
		},
		{
			name: "multi line",
			sut: &WorkflowSecretMarkdown{
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

func TestWorkflowOutputMarkdown_toMarkdown(t *testing.T) {
	cases := []struct {
		name     string
		sut      *WorkflowOutputMarkdown
		expected string
	}{
		{
			name: "single line",
			sut: &WorkflowOutputMarkdown{
				Name:        "single-line",
				Description: NewNotNullValue("The test description."),
			},
			expected: "| single-line | The test description. |",
		},
		{
			name: "multi line",
			sut: &WorkflowOutputMarkdown{
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

func TestWorkflowPermissionMarkdown_toMarkdown(t *testing.T) {
	cases := []struct {
		name     string
		sut      *WorkflowPermissionMarkdown
		expected string
	}{
		{
			name:     "valid",
			sut:      &WorkflowPermissionMarkdown{Scope: "contents", Access: "write"},
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
