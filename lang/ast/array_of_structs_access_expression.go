package ast

import (
	"encoding/json"
	"fmt"

	"tinygo.org/x/go-llvm"
)

type ArrayOfStructsAccessExpression struct {
	Name     Expression
	Index    Expression
	Property Expression
}

func (expr ArrayOfStructsAccessExpression) findSymbolTableEntry(ctx *CompilerCtx) (error, *ArraySymbolTableEntry, *SymbolTableEntry, int) {
	varName, ok := expr.Name.(SymbolExpression)
	if !ok {
		key := "ArrayAccessExpression.NameNotASymbol"
		return ctx.Dialect.Error(key, varName.Value), nil, nil, 0
	}

	array, ok := ctx.SymbolTable[varName.Value]
	if !ok {
		key := "ArrayAccessExpression.NotFoundInSymbolTable"
		return ctx.Dialect.Error(key, varName.Value), nil, nil, 0
	}

	entry, ok := ctx.ArraysSymbolTable[varName.Value]
	if !ok {
		key := "ArrayAccessExpression.NotFoundInArraysSymbolTable"
		return ctx.Dialect.Error(key, varName.Value), nil, nil, 0
	}
	itemIndex, ok := expr.Index.(NumberExpression)
	if !ok {
		key := "ArrayAccessExpression.AccessedIndexIsNotANumber"
		return ctx.Dialect.Error(key, expr.Index), nil, nil, 0
	}
	if int(itemIndex.Value) > entry.ElementsCount-1 {
		key := "ArrayAccessExpression.IndexOutOfBounds"
		return ctx.Dialect.Error(key, itemIndex, varName.Value), nil, nil, 0
	}

	return nil, &entry, &array, int(itemIndex.Value)
}

func (expr ArrayOfStructsAccessExpression) CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult) {
	err, entry, array, itemIndex := expr.findSymbolTableEntry(ctx)
	if err != nil {
		return err, nil
	}

	itemPtr := ctx.Builder.CreateInBoundsGEP(
		entry.Type,
		array.Value,
		[]llvm.Value{
			llvm.ConstInt(llvm.GlobalContext().Int32Type(), uint64(0), false),
			llvm.ConstInt(llvm.GlobalContext().Int32Type(), uint64(itemIndex), false),
		},
		"",
	)

	propertyName, _ := expr.Property.(SymbolExpression)
	err, index := entry.UnderlyingTypeDef.Metadata.PropertyIndex(propertyName.Value)
	if err != nil {
		return fmt.Errorf("ArrayOfStructsAccessExpression: property %s not found", propertyName.Value), nil
	}

	structPtr := ctx.Builder.CreateStructGEP(
		entry.UnderlyingType,
		itemPtr,
		index,
		"",
	)

	load := ctx.Builder.CreateLoad(entry.UnderlyingTypeDef.PropertyTypes[index], structPtr, "")

	return nil, &CompilerResult{Value: &load}
}

func (cs ArrayOfStructsAccessExpression) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["Name"] = cs.Name
	m["Index"] = cs.Index
	m["Property"] = cs.Property

	res := make(map[string]any)
	res["ast.ArrayOfStructsAccessExpression"] = m

	return json.Marshal(res)
}
