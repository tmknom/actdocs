package format

import "github.com/tmknom/actdocs/internal/parse"

func ConvertActionSpec(ast *parse.ActionAST) *ActionSpec {
	inputs := []*ActionInputSpec{}
	for _, inputAst := range ast.Inputs {
		input := &ActionInputSpec{
			Name:        inputAst.Name,
			Default:     inputAst.Default,
			Description: inputAst.Description,
			Required:    inputAst.Required,
		}
		inputs = append(inputs, input)
	}

	outputs := []*ActionOutputSpec{}
	for _, outputAst := range ast.Outputs {
		output := &ActionOutputSpec{
			Name:        outputAst.Name,
			Description: outputAst.Description,
		}
		outputs = append(outputs, output)
	}

	return &ActionSpec{
		Description: ast.Description,
		Inputs:      inputs,
		Outputs:     outputs,
	}
}
