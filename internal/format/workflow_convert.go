package format

import "github.com/tmknom/actdocs/internal/parse"

func ConvertWorkflowSpec(ast *parse.WorkflowAST) *WorkflowSpec {
	inputs := []*WorkflowInputSpec{}
	for _, inputAst := range ast.Inputs {
		input := &WorkflowInputSpec{
			Name:        inputAst.Name,
			Default:     inputAst.Default,
			Description: inputAst.Description,
			Required:    inputAst.Required,
			Type:        inputAst.Type,
		}
		inputs = append(inputs, input)
	}

	secrets := []*WorkflowSecretSpec{}
	for _, secretAst := range ast.Secrets {
		secret := &WorkflowSecretSpec{
			Name:        secretAst.Name,
			Description: secretAst.Description,
			Required:    secretAst.Required,
		}
		secrets = append(secrets, secret)
	}

	outputs := []*WorkflowOutputSpec{}
	for _, outputAst := range ast.Outputs {
		output := &WorkflowOutputSpec{
			Name:        outputAst.Name,
			Description: outputAst.Description,
		}
		outputs = append(outputs, output)
	}

	permissions := []*WorkflowPermissionSpec{}
	for _, permissionAst := range ast.Permissions {
		permission := &WorkflowPermissionSpec{
			Scope:  permissionAst.Scope,
			Access: permissionAst.Access,
		}
		permissions = append(permissions, permission)
	}

	return &WorkflowSpec{Inputs: inputs, Secrets: secrets, Outputs: outputs, Permissions: permissions}
}
