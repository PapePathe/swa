package ast

// BlockStatement ...
type BlockStatement struct {
	Body []Statement // The body of the block statement
}

func (bs BlockStatement) statement() {}

// ExpressionStatement ...
type ExpressionStatement struct {
	Exp Expression // The expression
}

func (bs ExpressionStatement) statement() {}

// VarDeclarationStatement ...
type VarDeclarationStatement struct {
	Name       string     // The name of the variable
	IsConstant bool       // Wether or not the variable is a constant
	Value      Expression // The value assigned to the variable
}

func (bs VarDeclarationStatement) statement() {}
