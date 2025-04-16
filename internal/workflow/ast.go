package workflow

import "github.com/tmknom/actdocs/internal/util"

type AST struct {
	Inputs      []*InputAST
	Secrets     []*SecretAST
	Outputs     []*OutputAST
	Permissions []*PermissionAST
}

type InputAST struct {
	Name        string
	Default     *util.NullString
	Description *util.NullString
	Required    *util.NullString
	Type        *util.NullString
}

func NewInputAST(name string) *InputAST {
	return &InputAST{
		Name:        name,
		Default:     util.DefaultNullString,
		Description: util.DefaultNullString,
		Required:    util.DefaultNullString,
		Type:        util.DefaultNullString,
	}
}

type SecretAST struct {
	Name        string
	Description *util.NullString
	Required    *util.NullString
}

func NewSecretAST(name string) *SecretAST {
	return &SecretAST{
		Name:        name,
		Description: util.DefaultNullString,
		Required:    util.DefaultNullString,
	}
}

type OutputAST struct {
	Name        string
	Description *util.NullString
}

func NewOutputAST(name string) *OutputAST {
	return &OutputAST{
		Name:        name,
		Description: util.DefaultNullString,
	}
}

type PermissionAST struct {
	Scope  string
	Access string
}

func NewPermissionAST(scope string, access string) *PermissionAST {
	return &PermissionAST{
		Scope:  scope,
		Access: access,
	}
}
