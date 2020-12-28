package childgroup

type ChildGroup string

const (
	Name                  ChildGroup = "NAME"
	Type                  ChildGroup = "TYPE"
	FunctionDeclarations  ChildGroup = "FUNCTION_DECLARATIONS"
	ParameterDeclarations ChildGroup = "PARAMETER_DECLARATIONS"
	LocalDeclarations     ChildGroup = "LOCAL_DECLARATIONS"
	Statements            ChildGroup = "STATEMENTS"
	Initializer           ChildGroup = "INITIALIZER"
	Arguments             ChildGroup = "ARGUMENTS"
	Rhs                   ChildGroup = "RHS"
	Lhs                   ChildGroup = "LHS"
)
