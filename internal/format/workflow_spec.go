package format

import (
	"fmt"
	"strings"

	"github.com/tmknom/actdocs/internal/util"
)

type WorkflowSpec struct {
	Inputs      []*WorkflowInputSpec      `json:"inputs"`
	Secrets     []*WorkflowSecretSpec     `json:"secrets"`
	Outputs     []*WorkflowOutputSpec     `json:"outputs"`
	Permissions []*WorkflowPermissionSpec `json:"permissions"`
}

func (s *WorkflowSpec) toInputsMarkdown() string {
	var sb strings.Builder
	sb.WriteString(WorkflowInputsTitle)
	sb.WriteString("\n\n")
	if len(s.Inputs) != 0 {
		sb.WriteString(WorkflowInputsColumnTitle)
		sb.WriteString("\n")
		sb.WriteString(WorkflowInputsColumnSeparator)
		sb.WriteString("\n")
		for _, input := range s.Inputs {
			sb.WriteString(input.toMarkdown())
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString(util.UpperNAString)
	}
	return strings.TrimSpace(sb.String())
}

func (s *WorkflowSpec) toSecretsMarkdown() string {
	var sb strings.Builder
	sb.WriteString(WorkflowSecretsTitle)
	sb.WriteString("\n\n")
	if len(s.Secrets) != 0 {
		sb.WriteString(WorkflowSecretsColumnTitle)
		sb.WriteString("\n")
		sb.WriteString(WorkflowSecretsColumnSeparator)
		sb.WriteString("\n")
		for _, secret := range s.Secrets {
			sb.WriteString(secret.toMarkdown())
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString(util.UpperNAString)
	}
	return strings.TrimSpace(sb.String())
}

func (s *WorkflowSpec) toOutputsMarkdown() string {
	var sb strings.Builder
	sb.WriteString(WorkflowOutputsTitle)
	sb.WriteString("\n\n")
	if len(s.Outputs) != 0 {
		sb.WriteString(WorkflowOutputsColumnTitle)
		sb.WriteString("\n")
		sb.WriteString(WorkflowOutputsColumnSeparator)
		sb.WriteString("\n")
		for _, output := range s.Outputs {
			sb.WriteString(output.toMarkdown())
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString(util.UpperNAString)
	}
	return strings.TrimSpace(sb.String())
}

func (s *WorkflowSpec) toPermissionsMarkdown() string {
	var sb strings.Builder
	sb.WriteString(WorkflowPermissionsTitle)
	sb.WriteString("\n\n")
	if len(s.Permissions) != 0 {
		sb.WriteString(WorkflowPermissionsColumnTitle)
		sb.WriteString("\n")
		sb.WriteString(WorkflowPermissionsColumnSeparator)
		sb.WriteString("\n")
		for _, permission := range s.Permissions {
			sb.WriteString(permission.toMarkdown())
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString(util.UpperNAString)
	}
	return strings.TrimSpace(sb.String())
}

type WorkflowInputSpec struct {
	Name        string           `json:"name"`
	Default     *util.NullString `json:"default"`
	Description *util.NullString `json:"description"`
	Required    *util.NullString `json:"required"`
	Type        *util.NullString `json:"type"`
}

func (s *WorkflowInputSpec) toMarkdown() string {
	str := util.TableSeparator
	str += fmt.Sprintf(" %s %s", s.Name, util.TableSeparator)
	str += fmt.Sprintf(" %s %s", s.Description.StringOrEmpty(), util.TableSeparator)
	str += fmt.Sprintf(" %s %s", s.Type.QuoteStringOrLowerNA(), util.TableSeparator)
	str += fmt.Sprintf(" %s %s", s.Default.QuoteStringOrLowerNA(), util.TableSeparator)
	str += fmt.Sprintf(" %s %s", s.Required.YesOrNo(), util.TableSeparator)
	return str
}

type WorkflowSecretSpec struct {
	Name        string           `json:"name"`
	Description *util.NullString `json:"description"`
	Required    *util.NullString `json:"required"`
}

func (s *WorkflowSecretSpec) toMarkdown() string {
	str := util.TableSeparator
	str += fmt.Sprintf(" %s %s", s.Name, util.TableSeparator)
	str += fmt.Sprintf(" %s %s", s.Description.StringOrEmpty(), util.TableSeparator)
	str += fmt.Sprintf(" %s %s", s.Required.YesOrNo(), util.TableSeparator)
	return str
}

type WorkflowOutputSpec struct {
	Name        string           `json:"name"`
	Description *util.NullString `json:"description"`
}

func (s *WorkflowOutputSpec) toMarkdown() string {
	str := util.TableSeparator
	str += fmt.Sprintf(" %s %s", s.Name, util.TableSeparator)
	str += fmt.Sprintf(" %s %s", s.Description.StringOrEmpty(), util.TableSeparator)
	return str
}

type WorkflowPermissionSpec struct {
	Scope  string `json:"scope"`
	Access string `json:"access"`
}

func (s *WorkflowPermissionSpec) toMarkdown() string {
	str := util.TableSeparator
	str += fmt.Sprintf(" %s %s", s.Scope, util.TableSeparator)
	str += fmt.Sprintf(" %s %s", s.Access, util.TableSeparator)
	return str
}

const (
	WorkflowInputsTitle           = "## Inputs"
	WorkflowInputsColumnTitle     = "| Name | Description | Type | Default | Required |"
	WorkflowInputsColumnSeparator = "| :--- | :---------- | :--- | :------ | :------: |"

	WorkflowSecretsTitle           = "## Secrets"
	WorkflowSecretsColumnTitle     = "| Name | Description | Required |"
	WorkflowSecretsColumnSeparator = "| :--- | :---------- | :------: |"

	WorkflowOutputsTitle           = "## Outputs"
	WorkflowOutputsColumnTitle     = "| Name | Description |"
	WorkflowOutputsColumnSeparator = "| :--- | :---------- |"

	WorkflowPermissionsTitle           = "## Permissions"
	WorkflowPermissionsColumnTitle     = "| Scope | Access |"
	WorkflowPermissionsColumnSeparator = "| :--- | :---- |"
)
