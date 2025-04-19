package workflow

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tmknom/actdocs/internal/util"
)

type Spec struct {
	Inputs      []*InputSpec      `json:"inputs"`
	Secrets     []*SecretSpec     `json:"secrets"`
	Outputs     []*OutputSpec     `json:"outputs"`
	Permissions []*PermissionSpec `json:"permissions"`

	Omit bool `json:"-"`
}

func (s *Spec) ToJson() string {
	bytes, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(bytes)
}

func (s *Spec) ToMarkdown() string {
	var sb strings.Builder
	sb.WriteString(s.ToInputsMarkdown())
	sb.WriteString("\n\n")
	sb.WriteString(s.ToSecretsMarkdown())
	sb.WriteString("\n\n")
	sb.WriteString(s.ToOutputsMarkdown())
	sb.WriteString("\n\n")
	sb.WriteString(s.ToPermissionsMarkdown())
	return strings.TrimSpace(sb.String())
}

func (s *Spec) ToInputsMarkdown() string {
	if s.Omit && len(s.Inputs) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(InputsTitle)
	sb.WriteString("\n\n")
	if len(s.Inputs) != 0 {
		sb.WriteString(InputsColumnTitle)
		sb.WriteString("\n")
		sb.WriteString(InputsColumnSeparator)
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

func (s *Spec) ToSecretsMarkdown() string {
	if s.Omit && len(s.Secrets) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(SecretsTitle)
	sb.WriteString("\n\n")
	if len(s.Secrets) != 0 {
		sb.WriteString(SecretsColumnTitle)
		sb.WriteString("\n")
		sb.WriteString(SecretsColumnSeparator)
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

func (s *Spec) ToOutputsMarkdown() string {
	if s.Omit && len(s.Outputs) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(OutputsTitle)
	sb.WriteString("\n\n")
	if len(s.Outputs) != 0 {
		sb.WriteString(OutputsColumnTitle)
		sb.WriteString("\n")
		sb.WriteString(OutputsColumnSeparator)
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

func (s *Spec) ToPermissionsMarkdown() string {
	if s.Omit && len(s.Permissions) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(PermissionsTitle)
	sb.WriteString("\n\n")
	if len(s.Permissions) != 0 {
		sb.WriteString(PermissionsColumnTitle)
		sb.WriteString("\n")
		sb.WriteString(PermissionsColumnSeparator)
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

type InputSpec struct {
	Name        string           `json:"name"`
	Default     *util.NullString `json:"default"`
	Description *util.NullString `json:"description"`
	Required    *util.NullString `json:"required"`
	Type        *util.NullString `json:"type"`
}

func (s *InputSpec) toMarkdown() string {
	str := util.TableSeparator
	str += fmt.Sprintf(" %s %s", s.Name, util.TableSeparator)
	str += fmt.Sprintf(" %s %s", s.Description.StringOrEmpty(), util.TableSeparator)
	str += fmt.Sprintf(" %s %s", s.Type.QuoteStringOrLowerNA(), util.TableSeparator)
	str += fmt.Sprintf(" %s %s", s.Default.QuoteStringOrLowerNA(), util.TableSeparator)
	str += fmt.Sprintf(" %s %s", s.Required.YesOrNo(), util.TableSeparator)
	return str
}

type SecretSpec struct {
	Name        string           `json:"name"`
	Description *util.NullString `json:"description"`
	Required    *util.NullString `json:"required"`
}

func (s *SecretSpec) toMarkdown() string {
	str := util.TableSeparator
	str += fmt.Sprintf(" %s %s", s.Name, util.TableSeparator)
	str += fmt.Sprintf(" %s %s", s.Description.StringOrEmpty(), util.TableSeparator)
	str += fmt.Sprintf(" %s %s", s.Required.YesOrNo(), util.TableSeparator)
	return str
}

type OutputSpec struct {
	Name        string           `json:"name"`
	Description *util.NullString `json:"description"`
}

func (s *OutputSpec) toMarkdown() string {
	str := util.TableSeparator
	str += fmt.Sprintf(" %s %s", s.Name, util.TableSeparator)
	str += fmt.Sprintf(" %s %s", s.Description.StringOrEmpty(), util.TableSeparator)
	return str
}

type PermissionSpec struct {
	Scope  string `json:"scope"`
	Access string `json:"access"`
}

func (s *PermissionSpec) toMarkdown() string {
	str := util.TableSeparator
	str += fmt.Sprintf(" %s %s", s.Scope, util.TableSeparator)
	str += fmt.Sprintf(" %s %s", s.Access, util.TableSeparator)
	return str
}

const (
	InputsTitle           = "## Inputs"
	InputsColumnTitle     = "| Name | Description | Type | Default | Required |"
	InputsColumnSeparator = "| :--- | :---------- | :--- | :------ | :------: |"

	SecretsTitle           = "## Secrets"
	SecretsColumnTitle     = "| Name | Description | Required |"
	SecretsColumnSeparator = "| :--- | :---------- | :------: |"

	OutputsTitle           = "## Outputs"
	OutputsColumnTitle     = "| Name | Description |"
	OutputsColumnSeparator = "| :--- | :---------- |"

	PermissionsTitle           = "## Permissions"
	PermissionsColumnTitle     = "| Scope | Access |"
	PermissionsColumnSeparator = "| :--- | :---- |"
)
