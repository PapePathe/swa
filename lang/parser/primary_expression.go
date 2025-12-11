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
		expr := ast.NumberExpression{}
		p.currentExpression = &expr
		tok := p.advance()
		expr.Tokens = append(expr.Tokens, tok)

		number, err := strconv.ParseInt(tok.Value, 10, 64)
		if err != nil {
			e := err.(*strconv.NumError)
			return nil, fmt.Errorf("%s: %s while parsing number expression", e.Num, e.Err)
		}
		expr.Value = number

		return expr, nil
	case lexer.String:
		expr := ast.StringExpression{}
		p.currentExpression = &expr
		tok := p.advance()
		expr.Tokens = append(expr.Tokens, tok)
		expr.Value = tok.Value

		return expr, nil
	case lexer.Identifier:
		expr := ast.SymbolExpression{}
		p.currentExpression = &expr
		tok := p.advance()
		expr.Value = tok.Value
		expr.Tokens = append(expr.Tokens, tok)
		return expr, nil
	case lexer.Float:
		expr := ast.FloatExpression{}
		p.currentExpression = &expr
		tok := p.advance()
		expr.Tokens = append(expr.Tokens, tok)

		number, err := strconv.ParseFloat(tok.Value, 64)
		if err != nil {
			return nil, err
		}
		expr.Value = number

		return expr, nil
	default:
		return nil, fmt.Errorf("Cannot create PrimaryExpression from %s", p.currentToken().Kind)
	}
}
