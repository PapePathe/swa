package parser

import (
	"fmt"

	"swahili/lang/ast"
	"swahili/lang/lexer"
)

// ParseVarDeclarationStatement ...
func ParseVarDeclarationStatement(p *Parser) (ast.Statement, error) {
	var explicitType ast.Type

	var assigedValue ast.Expression

	isConstant := p.advance().Kind == lexer.Const
	errStr := "Inside variable declaration expected to find variable name"
	variableName := p.expectError(lexer.Identifier, errStr).Value

	p.expect(lexer.Colon)
	explicitType = parseType(p, DefaultBindingPower)

	if p.currentToken().Kind != lexer.SemiColon {
		p.expect(lexer.Assignment)
		assigedValue, _ = parseExpression(p, Assignment)
	} else if explicitType == nil {
		return nil, fmt.Errorf("Missing either right hand side in var declaration or exlicit type")
	}

	if isConstant && assigedValue == nil {
		return nil, fmt.Errorf("Cannot define constant wihtout a value")
	}

	p.expect(lexer.SemiColon)

	return ast.VarDeclarationStatement{
		IsConstant:   isConstant,
		Value:        assigedValue,
		Name:         variableName,
		ExplicitType: explicitType,
	}, nil
}
