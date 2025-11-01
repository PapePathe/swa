package parser

import (
	"fmt"
	"slices"

	"swahili/lang/ast"
	"swahili/lang/lexer"
)

// ParseStatement ...
func ParseStatement(p *Parser) (ast.Statement, error) {
	statementFn, exists := statementLookup[p.currentToken().Kind]

	if exists {
		return statementFn(p)
	}

	expression, err := parseExpression(p, DefaultBindingPower)
	if err != nil {
		return ast.ExpressionStatement{}, err
	}

	p.expect(lexer.SemiColon)

	return ast.ExpressionStatement{
		Exp: expression,
	}, nil
}

func ParseStructDeclarationStatement(p *Parser) (ast.Statement, error) {
	p.expect(lexer.Struct)
	structName := p.expect(lexer.Identifier).Value

	properties := []string{}
	types := []ast.Type{}

	p.expect(lexer.OpenCurly)

	for p.hasTokens() && p.currentToken().Kind != lexer.CloseCurly {
		var propertyName string

		if p.currentToken().Kind == lexer.Identifier {
			propertyName = p.expect(lexer.Identifier).Value
			p.expectError(lexer.Colon, "Expected to find colon following struct property name")
			propType := parseType(p, DefaultBindingPower)
			p.expect(lexer.Comma)

			if slices.Contains(properties, propertyName) {
				return ast.StructDeclarationStatement{}, fmt.Errorf("property %s has already been defined", propertyName)
			}

			properties = append(properties, propertyName)
			types = append(types, propType)

			continue
		}
	}

	p.expect(lexer.CloseCurly)

	return ast.StructDeclarationStatement{
		Name:       structName,
		Properties: properties,
		Types:      types,
	}, nil
}
