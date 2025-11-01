package parser

import (
	"fmt"
	"strconv"

	"swahili/lang/ast"
	"swahili/lang/lexer"
)

// ParsePrimaryExpression ...
func ParsePrimaryExpression(p *Parser) (ast.Expression, error) {
	switch p.currentToken().Kind {
	case lexer.Number:
		number, err := strconv.ParseFloat(p.advance().Value, 64)
		if err != nil {
			return nil, err
		}

		return ast.NumberExpression{
			Value: number,
		}, nil
	case lexer.String:
		value := p.advance().Value

		return ast.StringExpression{
			Value: value[1 : len(value)-1],
		}, nil
	case lexer.Identifier:
		return ast.SymbolExpression{
			Value: p.advance().Value,
		}, nil

	default:
		return nil, fmt.Errorf("Cannot create PrimaryExpression from %s", p.currentToken().Kind)
	}
}
