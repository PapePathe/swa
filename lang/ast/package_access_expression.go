package ast

import (
	"encoding/json"
	"fmt"
	"strings"

	"swahili/lang/lexer"
)

// PackageAccessExpression ...
type PackageAccessExpression struct {
	Namespaces []Expression
	Arguments  []Expression
	Value      string
	Tokens     []lexer.Token
}

var _ Expression = (*PackageAccessExpression)(nil)

func (expr PackageAccessExpression) String() string {
	return expr.Value
}

func (expr PackageAccessExpression) Name() string {
	var names []string

	for _, ns := range expr.Namespaces {
		name, ok := ns.(SymbolExpression)
		if !ok {
			panic(fmt.Errorf("PackageAccessExpression %s is not a symbol expression", ns))
		}

		names = append(names, name.Value)

	}

	return strings.Join(names, ".")
}

func (expr PackageAccessExpression) CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult) {
	return nil, &CompilerResult{}
}

func (expr PackageAccessExpression) TokenStream() []lexer.Token {
	return expr.Tokens
}

func (expr PackageAccessExpression) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["Namespaces"] = expr.Value
	m["Arguments"] = expr.Value

	res := make(map[string]any)
	res["ast.PackageAccessExpression"] = m

	return json.Marshal(res)
}
