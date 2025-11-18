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
	stmt := ast.VarDeclarationStatement{}
	p.currentStatement = &stmt

	isConstant := p.advance().Kind == lexer.Const
	errStr := "Inside variable declaration expected to find variable name"
	tok := p.expectError(lexer.Identifier, errStr)
	stmt.Tokens = append(stmt.Tokens, tok)
	variableName := tok.Value

	stmt.Tokens = append(stmt.Tokens, p.expect(lexer.Colon))
	explicitType, toks := parseType(p, DefaultBindingPower)
	stmt.ExplicitType = explicitType
	stmt.Tokens = append(stmt.Tokens, toks...)

	if p.currentToken().Kind != lexer.SemiColon {
		stmt.Tokens = append(stmt.Tokens, p.expect(lexer.Assignment))
		assigedValue, err = parseExpression(p, Assignment)
		if err != nil {
			return nil, err
		}

		stmt.Tokens = append(stmt.Tokens, assigedValue.TokenStream()...)
	} else if explicitType == nil {
		return nil, fmt.Errorf("Missing either right hand side in var declaration or exlicit type")
	}

	if isConstant && assigedValue == nil {
		return nil, fmt.Errorf("Cannot define constant wihtout a value")
	}

	stmt.Tokens = append(stmt.Tokens, p.expect(lexer.SemiColon))
	stmt.IsConstant = isConstant
	stmt.Value = assigedValue
	stmt.Name = variableName

	return stmt, nil
}
