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

func (expr ArrayOfStructsAccessExpression) findSymbolTableEntry(ctx *CompilerCtx) (error, *ArraySymbolTableEntry, *SymbolTableEntry, []llvm.Value) {
	varName, ok := expr.Name.(SymbolExpression)
	if !ok {
		key := "ArrayAccessExpression.NameNotASymbol"

		return ctx.Dialect.Error(key, expr.Name), nil, nil, nil
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

		if int(idx.Value) < 0 {
			key := "ArrayAccessExpression.AccessedIndexIsNotANumber"

			return ctx.Dialect.Error(key, expr.Index), nil, nil, nil
		}

		if int(idx.Value) > entry.ElementsCount-1 {
			key := "ArrayAccessExpression.IndexOutOfBounds"

			return ctx.Dialect.Error(key, int(idx.Value), varName.Value), nil, nil, nil
		}

		indices = []llvm.Value{
			llvm.ConstInt(llvm.GlobalContext().Int32Type(), uint64(0), false),
			llvm.ConstInt(llvm.GlobalContext().Int32Type(), uint64(idx.Value), false),
		}
	case SymbolExpression, BinaryExpression:
		err, res := expr.Index.CompileLLVM(ctx)
		if err != nil {
			return err, nil, nil, nil
		}

		indices = []llvm.Value{
			llvm.ConstInt(llvm.GlobalContext().Int32Type(), uint64(0), false),
			*res.Value,
		}
	default:
		key := "ArrayAccessExpression.AccessedIndexIsNotANumber"

		return ctx.Dialect.Error(key, expr.Index), nil, nil, nil
	}

	return nil, entry, array, indices
}

func (expr ArrayOfStructsAccessExpression) CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult) {
	err, entry, array, itemIndex := expr.findSymbolTableEntry(ctx)
	if err != nil {
		return err, nil
	}

	itemPtr := ctx.Builder.CreateInBoundsGEP(
		entry.Type,
		array.Value,
		itemIndex,
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

	var ref *StructSymbolTableEntry

	propType := entry.UnderlyingTypeDef.Metadata.Types[index]
	if symbolType, ok := propType.(SymbolType); ok {
		_, ref = ctx.FindStructSymbol(symbolType.Name)
	}

	return nil, &CompilerResult{
		Value: &load,
		SymbolTableEntry: &SymbolTableEntry{
			Address: &structPtr,
			Ref:     ref,
		},
	}
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
