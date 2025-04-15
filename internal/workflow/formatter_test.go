package workflow

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tmknom/actdocs/internal/conf"
)

func TestWorkflowFormatter_Format(t *testing.T) {
	cases := []struct {
		name     string
		ast      *WorkflowAST
		expected string
	}{
		{
			name: "basic",
			ast: &WorkflowAST{
				Inputs: []*WorkflowInput{
					{Name: "foo", Default: NewNotNullValue("Default"), Description: NewNotNullValue("The inputs."), Required: NewNotNullValue("false"), Type: NewNotNullValue("string")},
				},
				Secrets: []*WorkflowSecret{
					{Name: "bar", Description: NewNotNullValue("The secrets."), Required: NewNotNullValue("false")},
				},
				Outputs: []*WorkflowOutput{
					{Name: "baz", Description: NewNotNullValue("The outputs.")},
				},
				Permissions: []*WorkflowPermission{
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
		json     *Spec
		expected string
	}{
		{
			name: "empty",
			json: &Spec{
				Inputs:      []*InputSpec{},
				Secrets:     []*SecretSpec{},
				Outputs:     []*OutputSpec{},
				Permissions: []*PermissionSpec{},
			},
			expected: emptyWorkflowExpectedJson,
		},
		{
			name: "full",
			json: &Spec{
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
		formatter := NewWorkflowFormatter(conf.DefaultFormatterConfig())
		got := formatter.ToJson(tc.json)
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

func TestWorkflowFormatter_ToMarkdown(t *testing.T) {
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
		formatter := NewWorkflowFormatter(tc.config)
		formatter.Spec = tc.markdown
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
