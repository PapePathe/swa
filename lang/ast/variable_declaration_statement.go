package ast

import (
	"encoding/json"
	"fmt"
	"swahili/lang/lexer"

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
	Tokens       []lexer.Token
}

var _ Statement = (*VarDeclarationStatement)(nil)

func (vd VarDeclarationStatement) CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult) {
	err, val := vd.Value.CompileLLVM(ctx)
	if err != nil {
		return err, nil
	}

	if val == nil {
		err := fmt.Errorf("VarDeclarationStatement: return value is nil <%s> <%s>", vd.Name, vd.Value)

		return err, nil
	}

	switch vd.Value.(type) {
	case StructInitializationExpression:
		explicitType, ok := vd.ExplicitType.(SymbolType)
		if !ok {
			return fmt.Errorf("explicit type is not a symbol %v", vd.ExplicitType), nil
		}

		err, typeDef := ctx.FindStructSymbol(explicitType.Name)
		if err != nil {
			return fmt.Errorf("Could not find typedef for %s in structs symbol table", explicitType.Name), nil
		}

		ctx.AddSymbol(vd.Name, &SymbolTableEntry{Value: *val.Value, Ref: typeDef})
	case StringExpression:
		glob := llvm.AddGlobal(*ctx.Module, val.Value.Type(), "")
		glob.SetInitializer(*val.Value)
		ctx.AddSymbol(vd.Name, &SymbolTableEntry{Value: glob, Address: &glob})
	case NumberExpression:
		alloc := ctx.Builder.CreateAlloca(val.Value.Type(), "")
		ctx.Builder.CreateStore(*val.Value, alloc)
		ctx.AddSymbol(vd.Name, &SymbolTableEntry{Value: *val.Value, Address: &alloc})
	case ArrayInitializationExpression:
		ctx.AddSymbol(vd.Name, &SymbolTableEntry{Value: *val.Value, Address: val.Value})
		ctx.AddArraySymbol(vd.Name, val.ArraySymbolTableEntry)
	default:
		panic(fmt.Sprintf("VarDeclarationStatement: Unhandled expression type (%v)", vd.Value))
	}

	return nil, nil
}

func (expr VarDeclarationStatement) TokenStream() []lexer.Token {
	return expr.Tokens
}

func (cs VarDeclarationStatement) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["Name"] = cs.Name
	m["IsConstant"] = cs.IsConstant
	m["Value"] = cs.Value
	m["ExplicitType"] = cs.ExplicitType

	res := make(map[string]any)
	res["ast.VarDeclarationStatement"] = m

	return json.Marshal(res)
}
