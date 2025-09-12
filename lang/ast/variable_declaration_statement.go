package ast

import (
	"encoding/json"
	"fmt"

	"tinygo.org/x/go-llvm"
)

// VarDeclarationStatement ...
type VarDeclarationStatement struct {
	// The name of the variable
	Name string
	// Wether or not the variable is a constant
	IsConstant bool
	// The value assigned to the variable
	Value Expression
	// The explicit type of the variable
	ExplicitType Type
}

var _ Statement = (*VarDeclarationStatement)(nil)

func (vd VarDeclarationStatement) CompileLLVM(ctx *CompilerCtx) (error, *llvm.Value) {
	err, val := vd.Value.CompileLLVM(ctx)
	if err != nil {
		return err, nil
	}

	if vd.Name != "r1" {
		if val == nil {
			err := fmt.Errorf("VarDeclarationStatement: return value is nil <%s> <%s>", vd.Name, vd.Value)
			return err, nil
		}

		glob := llvm.AddGlobal(*ctx.Module, val.Type(), vd.Name)
		glob.SetInitializer(*val)

		ctx.SymbolTable[vd.Name] = glob

	}

	return nil, nil
}

func (cs VarDeclarationStatement) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["name"] = cs.Name
	m["is_constant"] = cs.IsConstant
	m["value"] = cs.Value
	m["explicit_type"] = cs.ExplicitType

	res := make(map[string]any)
	res["ast.VarDeclarationStatement"] = m

	return json.Marshal(res)
}
