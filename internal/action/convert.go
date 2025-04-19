package action

func ConvertSpec(ast *AST, omit bool) *Spec {
	//goland:noinspection GoPreferNilSlice
	inputs := []*InputSpec{}
	for _, inputAst := range ast.Inputs {
		input := &InputSpec{
			Name:        inputAst.Name,
			Default:     inputAst.Default,
			Description: inputAst.Description,
			Required:    inputAst.Required,
		}
		inputs = append(inputs, input)
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

	return &Spec{
		Description: ast.Description,
		Inputs:      inputs,
		Outputs:     outputs,
		Omit:        omit,
	}
}
