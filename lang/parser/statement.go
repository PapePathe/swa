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
		return nil, err
	}

	p.expect(lexer.SemiColon)

	return ast.ExpressionStatement{
		Exp: expression,
	}, nil
}

func ParseStructDeclarationStatement(p *Parser) (ast.Statement, error) {
	tokens := []lexer.Token{}

	tokens = append(tokens, p.expect(lexer.Struct))
	structName := p.expect(lexer.Identifier)
	tokens = append(tokens, structName)

	properties := []string{}
	types := []ast.Type{}

	tokens = append(tokens, p.expect(lexer.OpenCurly))

	for p.hasTokens() && p.currentToken().Kind != lexer.CloseCurly {
		var propertyName string

		if p.currentToken().Kind == lexer.Identifier {
			prop := p.expect(lexer.Identifier)
			tokens = append(tokens, prop)
			propertyName = prop.Value
			tok := p.expectError(lexer.Colon, "Expected to find colon following struct property name")
			tokens = append(tokens, tok)
			propType, toks := parseType(p, DefaultBindingPower)
			tokens = append(tokens, toks...)
			tokens = append(tokens, p.expect(lexer.Comma))

			if slices.Contains(properties, propertyName) {
				return nil, fmt.Errorf("property %s has already been defined", propertyName)
			}

			properties = append(properties, propertyName)
			types = append(types, propType)

			continue
		}
	}

	tokens = append(tokens, p.expect(lexer.CloseCurly))

	return ast.StructDeclarationStatement{
		Name:       structName.Value,
		Properties: properties,
		Types:      types,
		Tokens:     tokens,
	}, nil
}
