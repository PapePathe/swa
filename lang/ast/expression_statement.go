package ast

import (
	"swahili/lang/lexer"
)

// ExpressionStatement ...
type ExpressionStatement struct {
	// The expression
	Exp    Expression
	Tokens []lexer.Token
}

var _ Statement = (*ExpressionStatement)(nil)

func (exp ExpressionStatement) CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult) {
	err, _ := exp.Exp.CompileLLVM(ctx)
	if err != nil {
		return err, nil
	}

	return nil, nil
}

func (es ExpressionStatement) Accept(g CodeGenerator) error {
	return g.VisitExpressionStatement(&es)
}

func (expr ExpressionStatement) TokenStream() []lexer.Token {
	return expr.Tokens
}
