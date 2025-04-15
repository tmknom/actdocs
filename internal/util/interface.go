package util

type YamlParser interface {
	ParseAST(yamlBytes []byte) (InterfaceAST, error)
}

type InterfaceAST interface {
	AST() string
}

const ReadAllAccess = "read-all"
const WriteAllAccess = "write-all"
const AllScope = "-"
