package parser

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

// ParseStatement ...
func ParseStatement(p *Parser) ast.Statement {
	statementFn, exists := statementLookup[p.currentToken().Kind]

	if exists {
		return statementFn(p)
	}

	expression := parseExpression(p, DefaultBindingPower)
	p.expect(lexer.SemiColon)

	return ast.ExpressionStatement{
		Exp: expression,
	}
}

// ParseVarDeclarationStatement ...
func ParseVarDeclarationStatement(p *Parser) ast.Statement {
	isConstant := p.advance().Kind == lexer.Const
	variableName := p.expectError(lexer.Identifier, "Inside variable declaration expected to find variable name").Value
	p.expect(lexer.Assignment)
	assigedValue := parseExpression(p, Assignment)
	p.expect(lexer.SemiColon)

	return ast.VarDeclarationStatement{
		IsConstant: isConstant,
		Value:      assigedValue,
		Name:       variableName,
	}
}
