package ast

import (
	"encoding/json"
	"fmt"
	"swahili/lang/lexer"

	"tinygo.org/x/go-llvm"
)

type StructInitializationExpression struct {
	Name       string
	Properties []string
	Values     []Expression
	Tokens     []lexer.Token
}

var _ Expression = (*StructInitializationExpression)(nil)

type StructItemValue struct {
	Position int
	Value    *llvm.Value
}

func (si StructInitializationExpression) CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult) {
	err, newtype := ctx.FindStructSymbol(si.Name)
	if err != nil {
		err := fmt.Errorf("StructInitializationExpression: Undefined struct named %s", si.Name)

		return err, nil
	}

	structInstance := ctx.Builder.CreateAlloca(newtype.LLVMType, fmt.Sprintf("%s.instance", si.Name))

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

		field1Ptr := ctx.Builder.CreateStructGEP(newtype.LLVMType, structInstance, propIndex, "")

		switch expr.(type) {
		case StringExpression:
			glob := llvm.AddGlobal(*ctx.Module, val.Value.Type(), "")
			glob.SetInitializer(*val.Value)
			ctx.Builder.CreateStore(glob, field1Ptr)
		case NumberExpression, FloatExpression:
			ctx.Builder.CreateStore(*val.Value, field1Ptr)
		case SymbolExpression:
			if val.SymbolTableEntry.Address != nil {
				if val.Value.Type().TypeKind() == llvm.PointerTypeKind {
					ctx.Builder.CreateStore(*val.SymbolTableEntry.Address, field1Ptr)

					break
				}

				ctx.Builder.CreateStore(*val.Value, field1Ptr)

				break
			}

			if val.SymbolTableEntry.Ref != nil {
				load := ctx.Builder.CreateLoad(val.SymbolTableEntry.Ref.LLVMType, *val.Value, "")
				ctx.Builder.CreateStore(load, field1Ptr)

				break
			}

			ctx.Builder.CreateStore(*val.Value, field1Ptr)
		case StructInitializationExpression:
			nestedStructType := newtype.PropertyTypes[propIndex]
			loadedNestedStruct := ctx.Builder.CreateLoad(nestedStructType, *val.Value, "")
			ctx.Builder.CreateStore(loadedNestedStruct, field1Ptr)
		case ArrayInitializationExpression:
			astType := newtype.Metadata.Types[propIndex]
			if _, isPointer := astType.(PointerType); isPointer {
				ctx.Builder.CreateStore(*val.Value, field1Ptr)
			} else {
				arrayType := newtype.PropertyTypes[propIndex]
				loadedArray := ctx.Builder.CreateLoad(arrayType, *val.Value, "")
				ctx.Builder.CreateStore(loadedArray, field1Ptr)
			}
		default:
			return fmt.Errorf("StructInitializationExpression: expression %v not implemented", expr), nil
		}
	}

	return nil, &CompilerResult{Value: &structInstance}
}

func (si StructInitializationExpression) InitValues(ctx *CompilerCtx) (error, []StructItemValue) {
	err, newtype := ctx.FindStructSymbol(si.Name)
	if err != nil {
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
		case NumberExpression, FloatExpression:
			fieldValues = append(fieldValues, StructItemValue{Position: propIndex, Value: val.Value})
		case StructInitializationExpression:
			fieldValues = append(fieldValues, StructItemValue{Position: propIndex, Value: val.Value})
		default:
			return fmt.Errorf("StructInitializationExpression expression: %v is not a known field type", expr), nil
		}
	}
	return nil, fieldValues
}

func (expr StructInitializationExpression) TokenStream() []lexer.Token {
	return expr.Tokens
}

func (expr StructInitializationExpression) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["Name"] = expr.Name
	m["Properties"] = expr.Properties
	m["Values"] = expr.Values

	res := make(map[string]any)
	res["ast.StructInitializationExpression"] = m

	return json.Marshal(res)
}
