package ast

// BlockStatement ...
type BlockStatement struct {
	// The body of the block statement
	Body []Statement `json:"ast.Body"`
}

func (bs BlockStatement) statement() {}

// ExpressionStatement ...
type ExpressionStatement struct {
	// The expression
	Exp Expression `json:"ast.ExpressionStatement"`
}

func (bs ExpressionStatement) statement() {}

// VarDeclarationStatement ...
type VarDeclarationStatement struct {
	// The name of the variable
	Name string
	// Wether or not the variable is a constant
	IsConstant bool
	// The value assigned to the variable
	Value Expression
	// The explicit type of the variable
	ExplicitType Type
}

func (bs VarDeclarationStatement) statement() {}

type StructProperty struct {
	PropType Type
}
type StructDeclarationStatement struct {
	Name       string
	Properties map[string]StructProperty `json:"ast.StructProperty"`
}

func (s StructDeclarationStatement) statement() {}
