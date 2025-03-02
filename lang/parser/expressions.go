package parser

import (
	"fmt"
	"swahili/lang/ast"
)

func parseExpression(p *Parser, bp BindingPower) ast.Expression {
	tokenKind := p.currentToken().Kind
	nudFn, exists := nudLookup[tokenKind]

	if !exists {
		panic(
			fmt.Sprintf(
				"nud handler expected for token %s and binding power %v \n",
				tokenKind,
				bp,
			),
		)
	}

	left := nudFn(p)

	for bindingPowerLookup[p.currentToken().Kind] > bp {
		ledFn, exists := ledLookup[p.currentToken().Kind]

		if !exists {
			panic(
				fmt.Sprintf(
					"led handler expected for token (%s: value(%s))\n",
					tokenKind,
					p.currentToken().Value,
				),
			)
		}

		left = ledFn(p, left, bindingPowerLookup[p.currentToken().Kind])
	}

	return left
}
