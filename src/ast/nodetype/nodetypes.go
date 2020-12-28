package nodetype

type NodeType string

const (
	Module               NodeType = "MODULE"
	FunctionDeclaration  NodeType = "FUNCTION_DECLARATION"
	ParameterDeclaration NodeType = "PARAMETER_DECLARATION"
	TypeExpression       NodeType = "TYPE_EXPRESSION"
	LocalDeclaration     NodeType = "LOCAL_DECLARATION"
	Statement            NodeType = "STATEMENT"
	Identifier           NodeType = "IDENTIFIER"
	Integer              NodeType = "INTEGER"
	Float                NodeType = "FLOAT"
	Boolean              NodeType = "BOOLEAN"
	Addition             NodeType = "ADDITION"
	Subtraction          NodeType = "SUBTRACTION"
	Multiplication       NodeType = "MULTIPLICATION"
	Division             NodeType = "DIVISION"
	IsEqualTo            NodeType = "IS_EQUAL_TO"
	NamedTypeExpression  NodeType = "NAMED_TYPE_EXPRESSION"
)
