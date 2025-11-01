package parser

import (
	"fmt"
	"strconv"

	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParsePrintStatement(p *Parser) (ast.Statement, error) {
	values := []ast.Expression{}

	p.expect(lexer.Print)

	p.expect(lexer.OpenParen)

	for p.hasTokens() && p.currentToken().Kind != lexer.CloseParen {
		switch p.currentToken().Kind {
		case lexer.String:
			str := p.expect(lexer.String).Value
			values = append(values, ast.StringExpression{Value: str[1 : len(str)-1]})
		case lexer.Identifier:
			err, next := p.nextToken()
			if err != nil {
				values = append(values, ast.SymbolExpression{Value: p.expect(lexer.Identifier).Value})
			} else {
				switch next.Kind {
				case lexer.Dot:
					left := ast.SymbolExpression{Value: p.expect(lexer.Identifier).Value}
					memberCallExpr, err := ParseMemberCallExpression(p, left, Relational)
					if err != nil {
						return nil, err
					}
					values = append(values, memberCallExpr)
				case lexer.OpenBracket:
					left := ast.SymbolExpression{Value: p.expect(lexer.Identifier).Value}
					arrayAccessExpr, err := ParseArrayAccess(p, left, Relational)
					if err != nil {
						return nil, err
					}
					values = append(values, arrayAccessExpr)
				default:
					values = append(values, ast.SymbolExpression{Value: p.expect(lexer.Identifier).Value})
				}
			}
		case lexer.Number:
			value := p.expect(lexer.Number).Value

			number, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return ast.PrintStatetement{}, fmt.Errorf("Error parsing number expression <%s> in PRINT statement", err)
			}

			values = append(values, ast.NumberExpression{Value: number})
		default:
			return ast.PrintStatetement{}, fmt.Errorf("Token %s not supported in print statement", p.currentToken().Kind)
		}

		if p.currentToken().Kind == lexer.Comma {
			p.expect(lexer.Comma)
		}
	}

	p.expect(lexer.CloseParen)
	p.expect(lexer.SemiColon)

	return ast.PrintStatetement{Values: values}, nil
}
