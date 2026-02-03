package parser

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParseErrorExpression(p *Parser) (ast.Expression, error) {
	expr := ast.ErrorExpression{}
	p.currentExpression = &expr
	tok := p.expect(lexer.TypeError)
	expr.Tokens = append(expr.Tokens, tok)

	return &expr, nil
}
