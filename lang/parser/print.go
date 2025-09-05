package parser

import (
	"fmt"
	"strconv"
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParseFunctionDeclaration(p *Parser) ast.Statement {
	funDecl := ast.FuncDeclStatement{}
	args := []ast.FuncArg{}

	p.expect(lexer.Function)
	funDecl.Name = p.expect(lexer.Identifier).Value
	p.expect(lexer.OpenParen)

	for p.hasTokens() && p.currentToken().Kind != lexer.CloseParen {
		name := p.expect(lexer.Identifier).Value
		p.expect(lexer.Colon)
		argType := p.expect(lexer.Identifier).Value
		arg := ast.FuncArg{
			Name:    name,
			ArgType: argType,
		}
		args = append(args, arg)

		if p.currentToken().Kind == lexer.Comma {
			p.expect(lexer.Comma)
		}
	}

	funDecl.Args = args

	p.expect(lexer.CloseParen)
	funDecl.ReturnType = p.expect(lexer.Identifier).Value
	funDecl.Body = ParseBlockStatement(p)

	return funDecl
}

func ParseReturnStatement(p *Parser) ast.Statement {
	p.expect(lexer.Return)

	rs := ast.ReturnStatement{}
	rs.Value = parseExpression(p, DefaultBindingPower)

	p.expect(lexer.SemiColon)

	return rs
}

func ParseMainStatement(p *Parser) ast.Statement {
	ms := ast.MainStatement{}

	p.expect(lexer.Main)
	p.expect(lexer.OpenParen)
	p.expect(lexer.CloseParen)
	p.expect(lexer.Identifier)

	ms.Body = ParseBlockStatement(p)

	return ms
}

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
			values = append(values, ast.SymbolExpression{Value: p.expect(lexer.Identifier).Value})
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
