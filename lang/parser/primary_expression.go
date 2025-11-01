package parser

import (
	"fmt"
	"strconv"

	"swahili/lang/ast"
	"swahili/lang/lexer"
)

// ParsePrimaryExpression ...
func ParsePrimaryExpression(p *Parser) (ast.Expression, error) {
	tokens := []lexer.Token{}

	switch p.currentToken().Kind {
	case lexer.Number:
		tok := p.advance()
		tokens = append(tokens, tok)
		number, err := strconv.ParseFloat(tok.Value, 64)
		if err != nil {
			return nil, err
		}

		return ast.NumberExpression{
			Value: number,
		}, nil
	case lexer.String:
		tok := p.advance()
		tokens = append(tokens, tok)
		value := tok.Value

		return ast.StringExpression{
			Value: value[1 : len(value)-1],
		}, nil
	case lexer.Identifier:
		tok := p.advance()
		tokens = append(tokens, tok)
		return ast.SymbolExpression{
			Value: tok.Value,
		}, nil

	default:
		return nil, fmt.Errorf("Cannot create PrimaryExpression from %s", p.currentToken().Kind)
	}
}
