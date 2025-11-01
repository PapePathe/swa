package parser

import (
	"fmt"
	"strconv"

	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParsePrintStatement(p *Parser) (ast.Statement, error) {
	values := []ast.Expression{}
	tokens := []lexer.Token{}

	tokens = append(tokens, p.expect(lexer.Print))
	tokens = append(tokens, p.expect(lexer.OpenParen))

	for p.hasTokens() && p.currentToken().Kind != lexer.CloseParen {
		switch p.currentToken().Kind {
		case lexer.String:
			tok := p.expect(lexer.String)
			tokens = append(tokens, tok)
			str := tok.Value
			values = append(values, ast.StringExpression{Value: str[1 : len(str)-1]})
		case lexer.Identifier:
			err, next := p.nextToken()
			if err != nil {
				tok := p.expect(lexer.Identifier)
				tokens = append(tokens, tok)
				values = append(values, ast.SymbolExpression{Value: tok.Value})
			} else {
				switch next.Kind {
				case lexer.Dot:
					tok := p.expect(lexer.Identifier)
					tokens = append(tokens, tok)
					left := ast.SymbolExpression{Value: tok.Value}
					memberCallExpr, err := ParseMemberCallExpression(p, left, Relational)
					if err != nil {
						return nil, err
					}
					values = append(values, memberCallExpr)
				case lexer.OpenBracket:
					tok := p.expect(lexer.Identifier)
					tokens = append(tokens, tok)
					left := ast.SymbolExpression{Value: tok.Value}
					arrayAccessExpr, err := ParseArrayAccess(p, left, Relational)
					if err != nil {
						return nil, err
					}
					values = append(values, arrayAccessExpr)
				default:
					tok := p.expect(lexer.Identifier)
					tokens = append(tokens, tok)
					values = append(values, ast.SymbolExpression{Value: tok.Value})
				}
			}
		case lexer.Number:
			tok := p.expect(lexer.Number)
			tokens = append(tokens, tok)
			value := tok.Value

			number, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, fmt.Errorf("Error parsing number expression <%s> in PRINT statement", err)
			}

			values = append(values, ast.NumberExpression{Value: number})
		default:
			return nil, fmt.Errorf("Token %s not supported in print statement", p.currentToken().Kind)
		}

		if p.currentToken().Kind == lexer.Comma {
			tokens = append(tokens, p.expect(lexer.Comma))
		}
	}

	tokens = append(tokens, p.expect(lexer.CloseParen))
	tokens = append(tokens, p.expect(lexer.SemiColon))

	return ast.PrintStatetement{Values: values, Tokens: tokens}, nil
}
