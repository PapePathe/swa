package parser

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParseGroupingExpression(p *Parser) ast.Expression {
	p.advance() // move past start of grouping expression

	expr := parseExpression(p, DefaultBindingPower)
	p.expect(lexer.CloseParen) // move past end of grouping expression

	return expr
}
