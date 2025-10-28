package ast

import (
	"encoding/json"

	"tinygo.org/x/go-llvm"
)

type ArrayAccessExpression struct {
	Name  Expression
	Index Expression
}

var _ Expression = (*ArrayAccessExpression)(nil)

func (expr ArrayAccessExpression) findSymbolTableEntry(ctx *CompilerCtx) (error, *ArraySymbolTableEntry, *SymbolTableEntry, int) {
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

func (expr ArrayAccessExpression) CompileLLVMForPrint(ctx *CompilerCtx) (error, *llvm.Value) {
	err, entry, array, itemIndex := expr.findSymbolTableEntry(ctx)
	if err != nil {
		return err, nil
	}

	itemPtr := ctx.Builder.CreateInBoundsGEP(
		entry.UnderlyingType,
		array.Value,
		[]llvm.Value{llvm.ConstInt(llvm.GlobalContext().Int32Type(), uint64(itemIndex), false)},
		"",
	)

	if entry.UnderlyingType.TypeKind() == llvm.IntegerTypeKind {
		load := ctx.Builder.CreateLoad(entry.UnderlyingType, itemPtr, "")

		return nil, &load
	}

	return nil, &itemPtr
}

func (expr ArrayAccessExpression) CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult) {
	err, entry, array, itemIndex := expr.findSymbolTableEntry(ctx)
	if err != nil {
		return err, nil
	}

	itemPtr := ctx.Builder.CreateInBoundsGEP(
		entry.UnderlyingType,
		array.Value,
		[]llvm.Value{llvm.ConstInt(llvm.GlobalContext().Int32Type(), uint64(itemIndex), false)},
		"",
	)

	return nil, &CompilerResult{Value: &itemPtr}
}

func (cs ArrayAccessExpression) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["Name"] = cs.Name
	m["Index"] = cs.Index

	res := make(map[string]any)
	res["ast.ArrayAccessExpression"] = m

	return json.Marshal(res)
}
