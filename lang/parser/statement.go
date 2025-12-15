package parser

import (
	"fmt"
	"slices"
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

// ParseStatement ...
func ParseStatement(p *Parser) (ast.Statement, error) {
	stmt := ast.ExpressionStatement{}
	p.currentStatement = &stmt
	statementFn, exists := statementLookup[p.currentToken().Kind]

	if exists {
		val, err := statementFn(p)
		if err != nil {
			return nil, err
		}

		return val, nil
	}

	expression, err := parseExpression(p, DefaultBindingPower)
	if err != nil {
		return nil, err
	}

	stmt.Exp = expression
	stmt.Tokens = append(stmt.Tokens, expression.TokenStream()...)
	stmt.Tokens = append(stmt.Tokens, p.expect(lexer.SemiColon))

	return stmt, nil
}

func ParseStructDeclarationStatement(p *Parser) (ast.Statement, error) {
	stmt := ast.StructDeclarationStatement{}
	p.currentStatement = &stmt

	stmt.Tokens = append(stmt.Tokens, p.expect(lexer.Struct))
	structName := p.expect(lexer.Identifier)
	stmt.Name = structName.Value
	stmt.Tokens = append(stmt.Tokens, structName)
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

			if slices.Contains(stmt.Properties, propertyName) {
				return nil, fmt.Errorf("property %s has already been defined", propertyName)
			}

			stmt.Properties = append(stmt.Properties, propertyName)
			stmt.Types = append(stmt.Types, propType)

			continue
		}
	}

	stmt.Tokens = append(stmt.Tokens, p.expect(lexer.CloseCurly))

	return stmt, nil
}
