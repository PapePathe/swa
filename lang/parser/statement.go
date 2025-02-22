package parser

import (
	"fmt"

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
	var explicitType ast.Type

	var assigedValue ast.Expression

	isConstant := p.advance().Kind == lexer.Const
	errStr := "Inside variable declaration expected to find variable name"
	variableName := p.expectError(lexer.Identifier, errStr).Value

	if p.currentToken().Kind == lexer.Colon {
		p.advance()
		explicitType = parseType(p, DefaultBindingPower)
	}

	if p.currentToken().Kind != lexer.SemiColon {
		p.expect(lexer.Assignment)
		assigedValue = parseExpression(p, Assignment)
	} else if explicitType == nil {
		panic("Missing either right hand side in var declaration or exlicit type")
	}

	if isConstant && assigedValue == nil {
		panic("Cannot define constant wihtout a value")
	}

	p.expect(lexer.SemiColon)

	return ast.VarDeclarationStatement{
		IsConstant:   isConstant,
		Value:        assigedValue,
		Name:         variableName,
		ExplicitType: explicitType,
	}
}

func ParseStructDeclarationStatement(p *Parser) ast.Statement {
	p.expect(lexer.Struct)
	structName := p.expect(lexer.Identifier).Value

	propertes := map[string]ast.StructProperty{}

	p.expect(lexer.OpenCurly)

	for p.hasTokens() && p.currentToken().Kind != lexer.CloseCurly {
		var propertyName string

		if p.currentToken().Kind == lexer.Identifier {
			propertyName = p.expect(lexer.Identifier).Value
			p.expectError(lexer.Colon, "Expected to find colon following struct property name")
			propType := parseType(p, DefaultBindingPower)
			p.expect(lexer.Comma)

			if _, exists := propertes[propertyName]; exists {
				panic(fmt.Sprintf("property %s has already been defined", propertyName))
			}

			propertes[propertyName] = ast.StructProperty{
				PropType: propType,
			}

			continue
		}
	}

	p.expect(lexer.CloseCurly)

	return ast.StructDeclarationStatement{
		Name:       structName,
		Properties: propertes,
	}
}
