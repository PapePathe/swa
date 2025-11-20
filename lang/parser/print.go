package parser

import (
	"fmt"
	"strconv"

	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParsePrintStatement(p *Parser) (ast.Statement, error) {
	stmt := ast.PrintStatetement{}

	stmt.Tokens = append(stmt.Tokens, p.expect(lexer.Print))
	stmt.Tokens = append(stmt.Tokens, p.expect(lexer.OpenParen))

	for p.hasTokens() && p.currentToken().Kind != lexer.CloseParen {
		switch p.currentToken().Kind {
		case lexer.String:
			tok := p.expect(lexer.String)
			stmt.Tokens = append(stmt.Tokens, tok)
			str := tok.Value
			stmt.Values = append(stmt.Values, ast.StringExpression{Value: str[1 : len(str)-1]})
		case lexer.Identifier:
			err, next := p.nextToken()
			if err != nil {
				tok := p.expect(lexer.Identifier)
				stmt.Tokens = append(stmt.Tokens, tok)
				stmt.Values = append(stmt.Values, ast.SymbolExpression{Value: tok.Value})
			} else {
				switch next.Kind {
				case lexer.Dot:
					tok := p.expect(lexer.Identifier)
					left := ast.SymbolExpression{Value: tok.Value}
					stmt.Tokens = append(stmt.Tokens, tok)
					memberCallExpr, err := ParseMemberCallExpression(p, left, Relational)
					if err != nil {
						return nil, err
					}
					stmt.Tokens = append(stmt.Tokens, memberCallExpr.TokenStream()...)
					stmt.Values = append(stmt.Values, memberCallExpr)
				case lexer.OpenBracket:
					tok := p.expect(lexer.Identifier)
					stmt.Tokens = append(stmt.Tokens, tok)
					left := ast.SymbolExpression{Value: tok.Value}
					arrayAccessExpr, err := ParseArrayAccess(p, left, Relational)
					if err != nil {
						return nil, err
					}

					stmt.Tokens = append(stmt.Tokens, arrayAccessExpr.TokenStream()...)
					stmt.Values = append(stmt.Values, arrayAccessExpr)
				default:
					tok := p.expect(lexer.Identifier)
					stmt.Tokens = append(stmt.Tokens, tok)
					stmt.Values = append(stmt.Values, ast.SymbolExpression{Value: tok.Value})
				}
			}
		case lexer.Float:
			tok := p.expect(lexer.Float)
			stmt.Tokens = append(stmt.Tokens, tok)
			value := tok.Value

			number, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, fmt.Errorf("Error parsing float expression <%s> in PRINT statement", err)
			}

			stmt.Values = append(stmt.Values, ast.FloatExpression{Value: number})
		case lexer.Number:
			tok := p.expect(lexer.Number)
			stmt.Tokens = append(stmt.Tokens, tok)
			value := tok.Value

			number, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("Error parsing number expression <%s> in PRINT statement", err)
			}

			stmt.Values = append(stmt.Values, ast.NumberExpression{Value: number})
		default:
			return nil, fmt.Errorf("ParsePrintStatement: Token %s not supported in print statement", p.currentToken().Kind)
		}

		if p.currentToken().Kind == lexer.Comma {
			stmt.Tokens = append(stmt.Tokens, p.expect(lexer.Comma))
		}
	}

	stmt.Tokens = append(stmt.Tokens, p.expect(lexer.CloseParen))
	stmt.Tokens = append(stmt.Tokens, p.expect(lexer.SemiColon))

	return stmt, nil
}
