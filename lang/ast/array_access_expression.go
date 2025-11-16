package ast

import (
	"encoding/json"
	"swahili/lang/lexer"

	"tinygo.org/x/go-llvm"
)

type ArrayAccessExpression struct {
	Name   Expression
	Index  Expression
	Tokens []lexer.Token
}

var _ Expression = (*ArrayAccessExpression)(nil)

func (expr ArrayAccessExpression) findSymbolTableEntry(ctx *CompilerCtx) (error, *ArraySymbolTableEntry, *SymbolTableEntry, []llvm.Value) {
	varName, ok := expr.Name.(SymbolExpression)
	if !ok {
		key := "ArrayAccessExpression.NameNotASymbol"
		return ctx.Dialect.Error(key, varName.Value), nil, nil, nil
	}

	err, array := ctx.FindSymbol(varName.Value)
	if err != nil {
		key := "ArrayAccessExpression.NotFoundInSymbolTable"
		return ctx.Dialect.Error(key, varName.Value), nil, nil, nil
	}

	err, entry := ctx.FindArraySymbol(varName.Value)
	if err != nil {
		key := "ArrayAccessExpression.NotFoundInArraysSymbolTable"
		return ctx.Dialect.Error(key, varName.Value), nil, nil, nil
	}

	var indices []llvm.Value
	switch expr.Index.(type) {
	case NumberExpression:
		idx, _ := expr.Index.(NumberExpression)

		if int(idx.Value) > entry.ElementsCount-1 {
			key := "ArrayAccessExpression.IndexOutOfBounds"
			return ctx.Dialect.Error(key, int(idx.Value), varName.Value), nil, nil, nil
		}
		indices = []llvm.Value{llvm.ConstInt(llvm.GlobalContext().Int32Type(), uint64(idx.Value), false)}
	case IntegerExpression:
		idx, _ := expr.Index.(IntegerExpression)

		if int(idx.Value) > entry.ElementsCount-1 {
			key := "ArrayAccessExpression.IndexOutOfBounds"
			return ctx.Dialect.Error(key, int(idx.Value), varName.Value), nil, nil, nil
		}
		indices = []llvm.Value{llvm.ConstInt(llvm.GlobalContext().Int32Type(), uint64(idx.Value), false)}
	case SymbolExpression:
		err, res := expr.Index.CompileLLVM(ctx)
		if err != nil {
			return err, nil, nil, nil
		}
		indices = []llvm.Value{*res.Value}
	default:
		key := "ArrayAccessExpression.AccessedIndexIsNotANumber"
		return ctx.Dialect.Error(key, expr.Index), nil, nil, nil
	}

	return nil, entry, array, indices
}

func (expr ArrayAccessExpression) CompileLLVMForPrint(ctx *CompilerCtx) (error, *llvm.Value) {
	err, entry, array, itemIndex := expr.findSymbolTableEntry(ctx)
	if err != nil {
		return err, nil
	}

	itemPtr := ctx.Builder.CreateInBoundsGEP(
		entry.UnderlyingType,
		array.Value,
		itemIndex,
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
		itemIndex,
		"",
	)

	// Load the value from the pointer for use in expressions
	loadedValue := ctx.Builder.CreateLoad(entry.UnderlyingType, itemPtr, "")

	return nil, &CompilerResult{Value: &loadedValue}
}

func (expr ArrayAccessExpression) TokenStream() []lexer.Token {
	return expr.Tokens
}

func (cs ArrayAccessExpression) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["Name"] = cs.Name
	m["Index"] = cs.Index

	res := make(map[string]any)
	res["ast.ArrayAccessExpression"] = m

	return json.Marshal(res)
}
