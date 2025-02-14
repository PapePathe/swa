package ast

type BlockStatement struct {
	Body []Statement
}

func (bs BlockStatement) statement() {}

type ExpressionStatement struct {
	Exp Expression
}

func (bs ExpressionStatement) statement() {}

type VarDeclarationStatement struct {
	Name       string
	IsConstant bool
	Value      Expression
}

func (bs VarDeclarationStatement) statement() {}
