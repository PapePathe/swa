package ast

import (
	"fmt"

	"tinygo.org/x/go-llvm"
)

type MemberExpression struct {
	Object   Expression
	Property Expression
	Computed bool
}

var _ Expression = (*MemberExpression)(nil)

func (expr MemberExpression) CompileLLVM(ctx *CompilerCtx) (error, *llvm.Value) {
	obj, ok := expr.Object.(SymbolExpression)
	if !ok {
		return fmt.Errorf("struct object should be a symbol %v", obj), nil
	}

	varDef, ok := ctx.SymbolTable[obj.Value]
	if !ok {
		return fmt.Errorf("Variable %s of type Struct is not defined", obj.Value), nil
	}

	prop, ok := expr.Property.(SymbolExpression)
	if !ok {
		return fmt.Errorf("struct property should be a symbol %v", prop), nil
	}

	err, propIndex := varDef.Ref.Metadata.PropertyIndex(prop.Value)
	if err != nil {
		return fmt.Errorf("Struct %s does not have a field named %s", varDef.Ref.Metadata.Name, prop), nil
	}

	addr := ctx.Builder.CreateStructGEP(varDef.Ref.LLVMType, varDef.Value, propIndex, "")

	return nil, &addr
}
