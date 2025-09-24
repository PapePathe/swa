package ast

import (
	"fmt"

	"tinygo.org/x/go-llvm"
)

type StructInitializationExpression struct {
	Name       string
	Properties []string
	Values     []Expression
}

var _ Expression = (*StructInitializationExpression)(nil)

func (si StructInitializationExpression) CompileLLVM(ctx *CompilerCtx) (error, *llvm.Value) {
	newtype, ok := ctx.StructSymbolTable[si.Name]

	if !ok {
		err := fmt.Errorf("StructInitializationExpression: Undefined struct named %s", si.Name)

		return err, nil
	}

	structInstance := ctx.Builder.CreateAlloca(newtype.LLVMType, "")

	for _, name := range si.Properties {
		err, propIndex := newtype.Metadata.PropertyIndex(name)
		if err != nil {
			return fmt.Errorf("StructInitializationExpression: property %s not found", name), nil
		}

		expr := si.Values[propIndex]
		err, val := expr.CompileLLVM(ctx)

		if err != nil {
			return fmt.Errorf("StructInitializationExpression: %w", err), nil
		}

		switch expr.(type) {
		case StringExpression:
			glob := llvm.AddGlobal(*ctx.Module, val.Type(), "")
			glob.SetInitializer(*val)
			field1Ptr := ctx.Builder.CreateStructGEP(newtype.LLVMType, structInstance, propIndex, "")
			ctx.Builder.CreateStore(glob, field1Ptr)

		case NumberExpression:
			field1Ptr := ctx.Builder.CreateStructGEP(newtype.LLVMType, structInstance, propIndex, "")
			ctx.Builder.CreateStore(*val, field1Ptr)
		}
	}

	return nil, &structInstance
}
