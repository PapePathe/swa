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

func (expr ArrayAccessExpression) CompileLLVM(ctx *CompilerCtx) (error, *llvm.Value) {
	varName, ok := expr.Name.(SymbolExpression)
	if !ok {
		key := "ArrayAccessExpression.NameNotASymbol"
		return ctx.Dialect.Error(key, varName.Value), nil
	}

	array, ok := ctx.SymbolTable[varName.Value]
	if !ok {
		key := "ArrayAccessExpression.NotFoundInSymbolTable"
		return ctx.Dialect.Error(key, varName.Value), nil
	}

	itemIndex, ok := expr.Index.(NumberExpression)
	if !ok {
		key := "ArrayAccessExpression.AccessedIndexIsNotANumber"
		return ctx.Dialect.Error(key, expr.Index), nil
	}

	entry, ok := ctx.ArraysSymbolTable[varName.Value]
	if !ok {
		key := "ArrayAccessExpression.NotFoundInArraysSymbolTable"
		return ctx.Dialect.Error(key, varName.Value), nil
	}

	if int(itemIndex.Value) > entry.ElementsCount-1 {
		key := "ArrayAccessExpression.IndexOutOfBounds"
		return ctx.Dialect.Error(key, itemIndex, varName.Value), nil
	}

	itemPtr := ctx.Builder.CreateInBoundsGEP(
		ctx.Context.Int32Type(),
		array.Value,
		[]llvm.Value{
			llvm.ConstInt(llvm.GlobalContext().Int32Type(), uint64(itemIndex.Value), false),
		},
		"",
	)

	return nil, &itemPtr
}

func (cs ArrayAccessExpression) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["Name"] = cs.Name
	m["Index"] = cs.Index

	res := make(map[string]any)
	res["ast.ArrayAccessExpression"] = m

	return json.Marshal(res)
}
