package workflow

import "github.com/tmknom/actdocs/internal/parse"

func ConvertWorkflowSpec(ast *parse.WorkflowAST) *Spec {
	inputs := []*InputSpec{}
	for _, inputAst := range ast.Inputs {
		input := &InputSpec{
			Name:        inputAst.Name,
			Default:     inputAst.Default,
			Description: inputAst.Description,
			Required:    inputAst.Required,
			Type:        inputAst.Type,
		}
		inputs = append(inputs, input)
	}

	secrets := []*SecretSpec{}
	for _, secretAst := range ast.Secrets {
		secret := &SecretSpec{
			Name:        secretAst.Name,
			Description: secretAst.Description,
			Required:    secretAst.Required,
		}
		secrets = append(secrets, secret)
	}

	outputs := []*OutputSpec{}
	for _, outputAst := range ast.Outputs {
		output := &OutputSpec{
			Name:        outputAst.Name,
			Description: outputAst.Description,
		}
		outputs = append(outputs, output)
	}

	permissions := []*PermissionSpec{}
	for _, permissionAst := range ast.Permissions {
		permission := &PermissionSpec{
			Scope:  permissionAst.Scope,
			Access: permissionAst.Access,
		}
		permissions = append(permissions, permission)
	}

	return &Spec{Inputs: inputs, Secrets: secrets, Outputs: outputs, Permissions: permissions}
}
