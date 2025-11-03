package parser

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParseFunctionDeclaration(p *Parser) (ast.Statement, error) {
	tokens := []lexer.Token{}

	funDecl := ast.FuncDeclStatement{}
	args := []ast.FuncArg{}

	tokens = append(tokens, p.expect(lexer.Function))
	tok := p.expect(lexer.Identifier)
	tokens = append(tokens, tok)
	funDecl.Name = tok.Value
	tokens = append(tokens, p.expect(lexer.OpenParen))

	for p.hasTokens() && p.currentToken().Kind != lexer.CloseParen {
		tok := p.expect(lexer.Identifier)
		tokens = append(tokens, tok)
		tokens = append(tokens, p.expect(lexer.Colon))
		argType, toks := parseType(p, DefaultBindingPower)
		tokens = append(tokens, toks...)
		arg := ast.FuncArg{
			Name:    tok.Value,
			ArgType: argType,
		}
		args = append(args, arg)

		if p.currentToken().Kind == lexer.Comma {
			tokens = append(tokens, p.expect(lexer.Comma))
		}
	}

	funDecl.Args = args

	tokens = append(tokens, p.expect(lexer.CloseParen))

	returnType, tokens := parseType(p, DefaultBindingPower)
	tokens = append(tokens, tokens...)
	funDecl.ReturnType = returnType
	body, err := ParseBlockStatement(p)
	if err != nil {
		return nil, err
	}
	tokens = append(tokens, body.TokenStream()...)

	funDecl.Body = body
	funDecl.Tokens = tokens

	return funDecl, nil
}

func ParseFunctionCall(p *Parser, left ast.Expression, bp BindingPower) (ast.Expression, error) {
	tokens := []lexer.Token{}
	args := []ast.Expression{}
	tokens = append(tokens, left.TokenStream()...)
	tokens = append(tokens, p.expect(lexer.OpenParen))

	for p.hasTokens() && p.currentToken().Kind != lexer.CloseParen {
		arg, err := parseExpression(p, DefaultBindingPower)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
		tokens = append(tokens, arg.TokenStream()...)
		if p.currentToken().Kind == lexer.Comma {
			tokens = append(tokens, p.expect(lexer.Comma))
		}
	}

	tokens = append(tokens, p.expect(lexer.CloseParen))

	return ast.FunctionCallExpression{
		Name:   left,
		Args:   args,
		Tokens: tokens,
	}, nil
}
