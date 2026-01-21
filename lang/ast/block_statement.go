package ast

import (
	"swahili/lang/lexer"
)

// BlockStatement ...
type BlockStatement struct {
	// The body of the block statement
	Body   []Statement
	Tokens []lexer.Token
}

var _ Statement = (*BlockStatement)(nil)

func (bs BlockStatement) CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult) {
	newCtx := NewCompilerContext(
		ctx.Context,
		ctx.Builder,
		ctx.Module,
		ctx.Dialect,
		ctx,
	)
	for _, stmt := range bs.Body {
		err, _ := stmt.CompileLLVM(newCtx)
		if err != nil {
			return err, nil
		}
	}

	return nil, nil
}

func (bs BlockStatement) Accept(g CodeGenerator) error {
	return g.VisitBlockStatement(&bs)
}

func (expr BlockStatement) TokenStream() []lexer.Token {
	return expr.Tokens
}
