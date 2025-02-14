package parser

import (
	"fmt"
	"strconv"
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParsePrimaryExpression(p *Parser) ast.Expression {

	switch p.currentToken().Kind {
	case lexer.NUMBER:
		number, _ := strconv.ParseFloat(p.advance().Value, 64)

		return ast.NumberExpression{
			Value: number,
		}
	case lexer.STRING:
		return ast.StringExpression{
			Value: p.advance().Value,
		}
	case lexer.IDENTIFIER:
		return ast.SymbolExpression{
			Value: p.advance().Value,
		}

	default:
		panic(fmt.Sprintf("Cannot create PrimaryExpression from %s", p.currentToken().Kind))
	}
}

func ParseBinaryExpression(p *Parser, left ast.Expression, bp BindingPower) ast.Expression {
	operatorToken := p.advance()
	right := parseExpression(p, bp)
	return ast.BinaryExpression{
		Left:     left,
		Right:    right,
		Operator: operatorToken,
	}

}

func parseExpression(p *Parser, bp BindingPower) ast.Expression {
	tokenKind := p.currentToken().Kind
	nudFn, exists := nudLookup[tokenKind]

	if !exists {
		panic(fmt.Sprintf("nud handler expected for token %s\n", tokenKind))
	}

	left := nudFn(p)
	for bindingPowerLookup[p.currentToken().Kind] > bp {
		ledFn, exists := ledLookup[p.currentToken().Kind]

		if !exists {
			panic(fmt.Sprintf("led handler expected for token (%s: value(%s))\n", tokenKind, p.currentToken().Value))
		}

		left = ledFn(p, left, bp)
	}

	return left
}
