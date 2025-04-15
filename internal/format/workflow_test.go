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
		json     *WorkflowSpec
		expected string
	}{
		{
			name: "empty",
			json: &WorkflowSpec{
				Inputs:      []*WorkflowInputSpec{},
				Secrets:     []*WorkflowSecretSpec{},
				Outputs:     []*WorkflowOutputSpec{},
				Permissions: []*WorkflowPermissionSpec{},
			},
			expected: emptyWorkflowExpectedJson,
		},
		{
			name: "full",
			json: &WorkflowSpec{
				Inputs: []*WorkflowInputSpec{
					{Name: "minimal", Default: NewNullValue(), Description: NewNullValue(), Required: NewNullValue(), Type: NewNullValue()},
					{Name: "full", Default: NewNotNullValue("true"), Description: NewNotNullValue("The input value."), Required: NewNotNullValue("true"), Type: NewNotNullValue("boolean")},
				},
				Secrets: []*WorkflowSecretSpec{
					{Name: "minimal", Description: NewNullValue(), Required: NewNullValue()},
					{Name: "full", Description: NewNotNullValue("The secret value."), Required: NewNotNullValue("true")},
				},
				Outputs: []*WorkflowOutputSpec{
					{Name: "minimal", Description: NewNullValue()},
					{Name: "full", Description: NewNotNullValue("The output value.")},
				},
				Permissions: []*WorkflowPermissionSpec{
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
		markdown *WorkflowSpec
		expected string
	}{
		{
			name:   "omit",
			config: &conf.FormatterConfig{Format: conf.DefaultFormat, Omit: true},
			markdown: &WorkflowSpec{
				Inputs:      []*WorkflowInputSpec{},
				Secrets:     []*WorkflowSecretSpec{},
				Outputs:     []*WorkflowOutputSpec{},
				Permissions: []*WorkflowPermissionSpec{},
			},
			expected: "",
		},
		{
			name:   "empty",
			config: conf.DefaultFormatterConfig(),
			markdown: &WorkflowSpec{
				Inputs:      []*WorkflowInputSpec{},
				Secrets:     []*WorkflowSecretSpec{},
				Outputs:     []*WorkflowOutputSpec{},
				Permissions: []*WorkflowPermissionSpec{},
			},
			expected: emptyWorkflowExpected,
		},
		{
			name:   "full",
			config: conf.DefaultFormatterConfig(),
			markdown: &WorkflowSpec{
				Inputs: []*WorkflowInputSpec{
					{Name: "single", Default: NewNotNullValue("5"), Description: NewNotNullValue("The number."), Required: NewNotNullValue("true"), Type: NewNotNullValue("number")},
				},
				Secrets: []*WorkflowSecretSpec{
					{Name: "single", Description: NewNotNullValue("The test description."), Required: NewNotNullValue("true")},
				},
				Outputs: []*WorkflowOutputSpec{
					{Name: "single", Description: NewNotNullValue("The test description.")},
				},
				Permissions: []*WorkflowPermissionSpec{
					{Scope: "contents", Access: "write"},
				},
			},
			expected: fullWorkflowExpected,
		},
	}

	for _, tc := range cases {
		formatter := NewWorkflowFormatter(tc.config)
		formatter.WorkflowSpec = tc.markdown
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
