package workflow

func ConvertSpec(ast *AST, omit bool) *Spec {
	//goland:noinspection GoPreferNilSlice
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

	//goland:noinspection GoPreferNilSlice
	secrets := []*SecretSpec{}
	for _, secretAst := range ast.Secrets {
		secret := &SecretSpec{
			Name:        secretAst.Name,
			Description: secretAst.Description,
			Required:    secretAst.Required,
		}
		secrets = append(secrets, secret)
	}

	//goland:noinspection GoPreferNilSlice
	outputs := []*OutputSpec{}
	for _, outputAst := range ast.Outputs {
		output := &OutputSpec{
			Name:        outputAst.Name,
			Description: outputAst.Description,
		}
		outputs = append(outputs, output)
	}

	//goland:noinspection GoPreferNilSlice
	permissions := []*PermissionSpec{}
	for _, permissionAst := range ast.Permissions {
		permission := &PermissionSpec{
			Scope:  permissionAst.Scope,
			Access: permissionAst.Access,
		}
		permissions = append(permissions, permission)
	}

	return &Spec{
		Inputs:      inputs,
		Secrets:     secrets,
		Outputs:     outputs,
		Permissions: permissions,
		Omit:        omit,
	}
}
