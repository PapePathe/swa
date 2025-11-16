package ast

import (
	"encoding/json"
	"fmt"
	"swahili/lang/lexer"

	"tinygo.org/x/go-llvm"
)

type ArrayOfStructsAccessExpression struct {
	Name     Expression
	Index    Expression
	Property Expression
	Tokens   []lexer.Token
}

func (expr ArrayOfStructsAccessExpression) findSymbolTableEntry(ctx *CompilerCtx) (error, *ArraySymbolTableEntry, *SymbolTableEntry, int) {
	varName, ok := expr.Name.(SymbolExpression)
	if !ok {
		key := "ArrayAccessExpression.NameNotASymbol"
		return ctx.Dialect.Error(key, varName.Value), nil, nil, 0
	}

	err, array := ctx.FindSymbol(varName.Value)
	if err != nil {
		key := "ArrayAccessExpression.NotFoundInSymbolTable"
		return ctx.Dialect.Error(key, varName.Value), nil, nil, 0
	}

	err, entry := ctx.FindArraySymbol(varName.Value)
	if err != nil {
		key := "ArrayAccessExpression.NotFoundInArraysSymbolTable"
		return ctx.Dialect.Error(key, varName.Value), nil, nil, 0
	}
	var indexValue int
	switch idx := expr.Index.(type) {
	case NumberExpression:
		indexValue = int(idx.Value)
	case IntegerExpression:
		indexValue = int(idx.Value)
	default:
		key := "ArrayAccessExpression.AccessedIndexIsNotANumber"
		return ctx.Dialect.Error(key, expr.Index), nil, nil, 0
	}

	if indexValue > entry.ElementsCount-1 {
		key := "ArrayAccessExpression.IndexOutOfBounds"
		return ctx.Dialect.Error(key, indexValue, varName.Value), nil, nil, 0
	}

	return nil, entry, array, indexValue
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

func (expr ArrayOfStructsAccessExpression) TokenStream() []lexer.Token {
	return expr.Tokens
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
