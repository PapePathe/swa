package parser

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParseStatement(p *Parser) ast.Statement {
	statementFn, exists := statementLookup[p.currentToken().Kind]

	if exists {
		return statementFn(p)
	}

	expression := parseExpression(p, DefaultBindingPower)
	p.expect(lexer.SEMI_COLON)

	return ast.ExpressionStatement{
		Exp: expression,
	}
}

func ParseVarDeclarationStatement(p *Parser) ast.Statement {
	isConstant := p.advance().Kind == lexer.CONST
	variableName := p.expectError(lexer.IDENTIFIER, "Inside variable declaration expected to find variable name").Value
	p.expect(lexer.ASSIGNMENT)
	assigedValue := parseExpression(p, Assignment)
	p.expect(lexer.SEMI_COLON)

	return ast.VarDeclarationStatement{
		IsConstant: isConstant,
		Value:      assigedValue,
		Name:       variableName,
	}
}
