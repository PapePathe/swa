package parser

import (
	"fmt"
	"slices"
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

// ParseStatement ...
func ParseStatement(p *Parser) (ast.Statement, error) {
	tokens := []lexer.Token{}
	statementFn, exists := statementLookup[p.currentToken().Kind]

	if exists {
		return statementFn(p)
	}

	expression, err := parseExpression(p, DefaultBindingPower)
	if err != nil {
		return nil, err
	}

	tokens = append(tokens, expression.TokenStream()...)
	tokens = append(tokens, p.expect(lexer.SemiColon))

	return ast.ExpressionStatement{
		Exp:    expression,
		Tokens: tokens,
	}, nil
}

func ParseStructDeclarationStatement(p *Parser) (ast.Statement, error) {
	stmt := ast.StructDeclarationStatement{}
	p.currentStatement = &stmt

	stmt.Tokens = append(stmt.Tokens, p.expect(lexer.Struct))
	structName := p.expect(lexer.Identifier)
	stmt.Tokens = append(stmt.Tokens, structName)

	properties := []string{}
	types := []ast.Type{}

	stmt.Tokens = append(stmt.Tokens, p.expect(lexer.OpenCurly))

	for p.hasTokens() && p.currentToken().Kind != lexer.CloseCurly {
		var propertyName string

		if p.currentToken().Kind == lexer.Identifier {
			prop := p.expect(lexer.Identifier)
			stmt.Tokens = append(stmt.Tokens, prop)
			propertyName = prop.Value
			tok := p.expect(lexer.Colon)
			stmt.Tokens = append(stmt.Tokens, tok)
			propType, toks := parseType(p, DefaultBindingPower)
			stmt.Tokens = append(stmt.Tokens, toks...)
			stmt.Tokens = append(stmt.Tokens, p.expect(lexer.Comma))

			if slices.Contains(properties, propertyName) {
				return nil, fmt.Errorf("property %s has already been defined", propertyName)
			}

			properties = append(properties, propertyName)
			types = append(types, propType)

			continue
		}
	}

	stmt.Tokens = append(stmt.Tokens, p.expect(lexer.CloseCurly))
	stmt.Name = structName.Value
	stmt.Properties = properties
	stmt.Types = types

	return stmt, nil
}
