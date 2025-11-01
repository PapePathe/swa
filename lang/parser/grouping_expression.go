package parser

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParseGroupingExpression(p *Parser) (ast.Expression, error) {
	tokens := []lexer.Token{}
	tokens = append(tokens, p.advance())

	expr, err := parseExpression(p, DefaultBindingPower)
	if err != nil {
		return nil, err
	}
	tokens = append(tokens, p.expect(lexer.CloseParen))

	return expr, nil
}
