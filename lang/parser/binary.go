package parser

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

// ParseBinaryExpression ...
func ParseBinaryExpression(p *Parser, left ast.Expression, bp BindingPower) (ast.Expression, error) {
	expr := ast.BinaryExpression{}
	p.currentExpression = &expr
	expr.Left = left
	expr.Tokens = append(expr.Tokens, left.TokenStream()...)
	operatorToken := p.advance()
	expr.Operator = operatorToken
	expr.Tokens = append(expr.Tokens, operatorToken)

	right, err := parseExpression(p, bp)
	if err != nil {
		return nil, err
	}

	expr.Right = right
	expr.Tokens = append(expr.Tokens, right.TokenStream()...)

	return expr, nil
}

// ParseBinaryExpression ...
func ParsePackageAccessExpression(p *Parser, left ast.Expression, bp BindingPower) (ast.Expression, error) {
	expr := ast.PackageAccessExpression{}
	p.currentExpression = &expr

	expr.Namespaces = append(expr.Namespaces, left)

	for p.hasTokens(); p.currentToken().Kind == lexer.DoubleColon; {
		p.expect(lexer.DoubleColon)
		id := p.expect(lexer.Identifier)
		expr.Namespaces = append(expr.Namespaces, ast.SymbolExpression{Value: id.Value})
	}

	if p.currentToken().Kind == lexer.OpenParen {
		return ParseFunctionCall(p, expr, DefaultBindingPower)
	}

	return expr, nil
}
