package ast

import (
	"fmt"

	"tinygo.org/x/go-llvm"
)

type StructInitializationExpression struct {
	Name       string
	Properties map[string]Expression
}

var _ Expression = (*StructInitializationExpression)(nil)

func (si StructInitializationExpression) CompileLLVM(ctx *CompilerCtx) (error, *llvm.Value) {
	newtype, ok := ctx.StructSymbolTable[si.Name]

	if !ok {
		err := fmt.Errorf("StructInitializationExpression: Undefined struct named %s", si.Name)
		return err, nil
	}

	structInstance := ctx.Builder.CreateAlloca(newtype, "")

	index := 0
	for name, expr := range si.Properties {
		err, val := expr.CompileLLVM(ctx)
		if err != nil {
			return fmt.Errorf("StructInitializationExpression: %s", err), nil
		}

		switch expr.(type) {
		case StringExpression:
			field1Ptr := ctx.Builder.CreateStructGEP(newtype, structInstance, index, name)
			ctx.Builder.CreateStore(*val, field1Ptr)

		case NumberExpression:
			field1Ptr := ctx.Builder.CreateStructGEP(newtype, structInstance, index, name)
			ctx.Builder.CreateStore(*val, field1Ptr)
		}
		//		fmt.Printf("Initializing struct property %s at index %d with value %s\n", name, index, *val)

		index = index + 1
	}

	return nil, nil
}
