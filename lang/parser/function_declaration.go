package parser

import (
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
