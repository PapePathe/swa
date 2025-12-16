package ast

import (
	"encoding/json"
	"fmt"
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
	var name string

	var array *SymbolTableEntry

	var entry *ArraySymbolTableEntry

	var err error

	switch expr.Name.(type) {
	case SymbolExpression:
		varName, ok := expr.Name.(SymbolExpression)
		if !ok {
			key := "ArrayAccessExpression.NameNotASymbol"

			return ctx.Dialect.Error(key, varName.Value), nil, nil, nil
		}

		err, array = ctx.FindSymbol(varName.Value)
		if err != nil {
			key := "ArrayAccessExpression.NotFoundInSymbolTable"

			return ctx.Dialect.Error(key, varName.Value), nil, nil, nil
		}

		err, entry = ctx.FindArraySymbol(varName.Value)
		if err != nil {
			key := "ArrayAccessExpression.NotFoundInArraySymbolTable"

			return ctx.Dialect.Error(key, varName.Value), nil, nil, nil
		}

		name = varName.Value
	case MemberExpression:
		err, val := expr.Name.CompileLLVM(ctx)
		if err != nil {
			return err, nil, nil, nil
		}

		var elementType llvm.Type

		var arrayType llvm.Type

		var elementsCount int = -1

		var isPointerType bool = false

		if val.SymbolTableEntry == nil && val.SymbolTableEntry.Ref == nil {
			return fmt.Errorf("ArrayAccessExpression Missing SymbolTableEntry"), nil, nil, nil
		}

		propExpr, ok := expr.Name.(MemberExpression)
		if !ok {
			format := "ArrayAccessExpression expected expr name to be a MemberExpression"

			return fmt.Errorf(format), nil, nil, nil
		}

		propSym, ok := propExpr.Property.(SymbolExpression)
		if !ok {
			format := "ArrayAccessExpression expected expr property to be a SymbolExpression"

			return fmt.Errorf(format), nil, nil, nil
		}

		if val.SymbolTableEntry.Ref == nil {
			format := "ArrayAccessExpression property %s is not an array"

			return fmt.Errorf(format, propSym.Value), nil, nil, nil
		}

		propIndex, err := propExpr.resolveStructAccess(val.SymbolTableEntry.Ref, propSym.Value)
		if err != nil {
			return err, nil, nil, nil
		}

		astType := val.SymbolTableEntry.Ref.Metadata.Types[propIndex]

		switch coltype := astType.(type) {
		case PointerType:
			isPointerType = true

			err, elementType = coltype.Underlying.LLVMType(ctx)
			if err != nil {
				return err, nil, nil, nil
			}

			arrayType = elementType
		case ArrayType:
			err, elementType = coltype.Underlying.LLVMType(ctx)
			if err != nil {
				return err, nil, nil, nil
			}

			err, arrayType = coltype.LLVMType(ctx)
			if err != nil {
				return err, nil, nil, nil
			}

			elementsCount = coltype.Size
		default:
			err := fmt.Errorf("Property %s is not an array", propSym.Value)

			return err, nil, nil, nil
		}

		if isPointerType {
			pointerValue := ctx.Builder.CreateLoad(*val.StuctPropertyValueType, *val.Value, "")
			array = &SymbolTableEntry{Value: pointerValue}
		} else {
			array = &SymbolTableEntry{Value: *val.Value}
		}

		entry = &ArraySymbolTableEntry{
			UnderlyingType: elementType,
			Type:           arrayType,
			ElementsCount:  elementsCount,
		}
	default:
		err := fmt.Errorf("ArrayAccessExpression not implemented")

		return err, nil, nil, nil
	}

	var indices []llvm.Value

	switch expr.Index.(type) {
	case NumberExpression:
		idx, _ := expr.Index.(NumberExpression)

		if int(idx.Value) < 0 {
			key := "ArrayAccessExpression.AccessedIndexIsNotANumber"

			return ctx.Dialect.Error(key, expr.Index), nil, nil, nil
		}

		// Skip bounds checking for pointer types (ElementsCount == -1)
		if entry.ElementsCount >= 0 && int(idx.Value) > entry.ElementsCount-1 {
			key := "ArrayAccessExpression.IndexOutOfBounds"

			return ctx.Dialect.Error(key, int(idx.Value), name), nil, nil, nil
		}

		indices = []llvm.Value{
			llvm.ConstInt(llvm.GlobalContext().Int32Type(), uint64(idx.Value), false),
		}
	case SymbolExpression, BinaryExpression, MemberExpression:
		err, res := expr.Index.CompileLLVM(ctx)
		if err != nil {
			return err, nil, nil, nil
		}

		if res.Value.Type().TypeKind() == llvm.PointerTypeKind {
			load := ctx.Builder.CreateLoad(res.Value.AllocatedType(), *res.Value, "")
			indices = []llvm.Value{load}
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

	return nil, &CompilerResult{Value: &itemPtr, ArraySymbolTableEntry: entry}
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
