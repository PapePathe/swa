package ast

import (
	"encoding/json"
	"fmt"

	"tinygo.org/x/go-llvm"
)

type StructInitializationExpression struct {
	Name       string
	Properties []string
	Values     []Expression
}

var _ Expression = (*StructInitializationExpression)(nil)

type StructItemValue struct {
	Position int
	Value    *llvm.Value
}

func (si StructInitializationExpression) InitValues(ctx *CompilerCtx) (error, []StructItemValue) {
	newtype, ok := ctx.StructSymbolTable[si.Name]

	if !ok {
		err := fmt.Errorf("StructInitializationExpression: Undefined struct named %s", si.Name)

		return err, nil
	}

	fieldValues := []StructItemValue{}
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
			glob := llvm.AddGlobal(*ctx.Module, val.Value.Type(), "")
			glob.SetInitializer(*val.Value)
			fieldValues = append(fieldValues, StructItemValue{Position: propIndex, Value: &glob})

		case NumberExpression:
			fieldValues = append(fieldValues, StructItemValue{Position: propIndex, Value: val.Value})
		}
	}
	return nil, fieldValues
}

func (si StructInitializationExpression) CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult) {
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
			glob := llvm.AddGlobal(*ctx.Module, val.Value.Type(), "")
			glob.SetInitializer(*val.Value)
			field1Ptr := ctx.Builder.CreateStructGEP(newtype.LLVMType, structInstance, propIndex, "")
			ctx.Builder.CreateStore(glob, field1Ptr)

		case NumberExpression:
			field1Ptr := ctx.Builder.CreateStructGEP(newtype.LLVMType, structInstance, propIndex, "")
			ctx.Builder.CreateStore(*val.Value, field1Ptr)
		}
	}

	return nil, &CompilerResult{Value: &structInstance}
}

func (cs StructInitializationExpression) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["Name"] = cs.Name
	m["Properties"] = cs.Properties
	m["Values"] = cs.Values

	res := make(map[string]any)
	res["ast.StructInitializationExpression"] = m

	return json.Marshal(res)
}

