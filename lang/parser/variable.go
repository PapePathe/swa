package parser

import (
	"fmt"

	"swahili/lang/ast"
	"swahili/lang/lexer"
)

// ParseVarDeclarationStatement ...
func ParseVarDeclarationStatement(p *Parser) (ast.Statement, error) {
	var explicitType ast.Type
	var err error
	var assigedValue ast.Expression
	tokens := []lexer.Token{}

	isConstant := p.advance().Kind == lexer.Const
	errStr := "Inside variable declaration expected to find variable name"
	tok := p.expectError(lexer.Identifier, errStr)
	tokens = append(tokens, tok)
	variableName := tok.Value

	tokens = append(tokens, p.expect(lexer.Colon))
	explicitType, toks := parseType(p, DefaultBindingPower)
	tokens = append(tokens, toks...)

	if p.currentToken().Kind != lexer.SemiColon {
		tokens = append(tokens, p.expect(lexer.Assignment))

		assigedValue, err = parseExpression(p, Assignment)
		if err != nil {
			return nil, err
		}

		tokens = append(tokens, assigedValue.TokenStream()...)

	} else if explicitType == nil {
		return nil, fmt.Errorf("Missing either right hand side in var declaration or exlicit type")
	}

	if isConstant && assigedValue == nil {
		return nil, fmt.Errorf("Cannot define constant wihtout a value")
	}

	tokens = append(tokens, p.expect(lexer.SemiColon))

	return ast.VarDeclarationStatement{
		IsConstant:   isConstant,
		Value:        assigedValue,
		Name:         variableName,
		ExplicitType: explicitType,
		Tokens:       tokens,
	}, nil
}
