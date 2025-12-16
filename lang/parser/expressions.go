package parser

import (
	"fmt"

	"swahili/lang/ast"
)

func parseExpression(p *Parser, bp BindingPower) (ast.Expression, error) {
	tokenKind := p.currentToken().Kind
	nudFn, exists := nudLookup[tokenKind]

	if !exists {
		return nil, fmt.Errorf(
			"nud handler expected for token %s and binding power %v at line %d",
			tokenKind,
			bp,
			p.currentToken().Line,
		)
	}

	left, err := nudFn(p)
	if err != nil {
		return nil, err
	}

	for bindingPowerLookup[p.currentToken().Kind] > bp {
		ledFn, exists := ledLookup[p.currentToken().Kind]

		if !exists {
			return nil, fmt.Errorf(
				"led handler expected for token (%s: value(%s)) at line %d",
				tokenKind,
				p.currentToken().Value,
				p.currentToken().Line,
			)
		}

		left, err = ledFn(p, left, bindingPowerLookup[p.currentToken().Kind])
		if err != nil {
			return nil, err
		}
	}

	return left, nil
}
