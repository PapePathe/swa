package parser

import (
	"fmt"
	"strconv"
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

// ParsePrimaryExpression ...
func ParsePrimaryExpression(p *Parser) ast.Expression {
	switch p.currentToken().Kind {
	case lexer.Number:
		number, _ := strconv.ParseFloat(p.advance().Value, 64)

		return ast.NumberExpression{
			Value: number,
		}
	case lexer.String:
		value := p.advance().Value

		return ast.StringExpression{
			Value: value[1 : len(value)-1],
		}
	case lexer.Identifier:
		return ast.SymbolExpression{
			Value: p.advance().Value,
		}

	default:
		panic(fmt.Sprintf("Cannot create PrimaryExpression from %s", p.currentToken().Kind))
	}
}
