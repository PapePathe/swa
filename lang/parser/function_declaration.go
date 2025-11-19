package parser

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParseFunctionDeclaration(p *Parser) (ast.Statement, error) {
	funDecl := ast.FuncDeclStatement{}
	p.currentStatement = &funDecl

	funDecl.Tokens = append(funDecl.Tokens, p.expect(lexer.Function))
	tok := p.expect(lexer.Identifier)
	funDecl.Name = tok.Value
	funDecl.Tokens = append(funDecl.Tokens, tok)
	funDecl.Tokens = append(funDecl.Tokens, p.expect(lexer.OpenParen))

	for p.hasTokens() && p.currentToken().Kind != lexer.CloseParen {
		tok := p.expect(lexer.Identifier)
		funDecl.Tokens = append(funDecl.Tokens, tok)
		funDecl.Tokens = append(funDecl.Tokens, p.expect(lexer.Colon))
		argType, toks := parseType(p, DefaultBindingPower)
		funDecl.Tokens = append(funDecl.Tokens, toks...)
		arg := ast.FuncArg{
			Name:    tok.Value,
			ArgType: argType,
		}
		funDecl.Args = append(funDecl.Args, arg)

		if p.currentToken().Kind == lexer.Comma {
			funDecl.Tokens = append(funDecl.Tokens, p.expect(lexer.Comma))
		}
	}

	funDecl.Tokens = append(funDecl.Tokens, p.expect(lexer.CloseParen))

	returnType, tokens := parseType(p, DefaultBindingPower)
	funDecl.ReturnType = returnType
	funDecl.Tokens = append(funDecl.Tokens, tokens...)

	body, err := ParseBlockStatement(p)
	if err != nil {
		return nil, err
	}

	funDecl.Tokens = append(funDecl.Tokens, body.TokenStream()...)
	funDecl.Body = body

	return funDecl, nil
}

func ParseFunctionCall(p *Parser, left ast.Expression, bp BindingPower) (ast.Expression, error) {
	expr := ast.FunctionCallExpression{}
	p.currentExpression = &expr
	expr.Tokens = append(expr.Tokens, left.TokenStream()...)
	expr.Tokens = append(expr.Tokens, p.expect(lexer.OpenParen))

	for p.hasTokens() && p.currentToken().Kind != lexer.CloseParen {
		arg, err := parseExpression(p, DefaultBindingPower)
		if err != nil {
			return nil, err
		}
		expr.Args = append(expr.Args, arg)
		expr.Tokens = append(expr.Tokens, arg.TokenStream()...)
		if p.currentToken().Kind == lexer.Comma {
			expr.Tokens = append(expr.Tokens, p.expect(lexer.Comma))
		}
	}

	expr.Tokens = append(expr.Tokens, p.expect(lexer.CloseParen))
	expr.Name = left

	return expr, nil
}
