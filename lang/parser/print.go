package parser

import (
	"fmt"
	"strconv"
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParsePrintStatement(p *Parser) ast.Statement {
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
					values = append(values, ParseMemberCallExpression(p, left, Relational))
				case lexer.OpenBracket:
					left := ast.SymbolExpression{Value: p.expect(lexer.Identifier).Value}
					values = append(values, ParseArrayAccess(p, left, Relational))
				default:
					values = append(values, ast.SymbolExpression{Value: p.expect(lexer.Identifier).Value})
				}
			}
		case lexer.Number:
			value := p.expect(lexer.Number).Value

			number, err := strconv.ParseFloat(value, 64)
			if err != nil { // change this to return error when we feel stable
				panic(fmt.Sprintf("Error parsing number expression <%s> in PRINT statement", err))
			}

			values = append(values, ast.NumberExpression{Value: number})
		default: // change this to return error when we feel stable
			panic(fmt.Sprintf("Token %s not supported in print statement", p.currentToken().Kind))
		}

		if p.currentToken().Kind == lexer.Comma {
			p.expect(lexer.Comma)
		}
	}

	p.expect(lexer.CloseParen)
	p.expect(lexer.SemiColon)

	return ast.PrintStatetement{Values: values}
}
