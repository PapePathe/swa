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

		// Get the AST type from struct metadata to determine element type
		var elementType llvm.Type
		var arrayType llvm.Type
		var elementsCount int = -1
		var isPointerType bool = false

		if val.SymbolTableEntry != nil && val.SymbolTableEntry.Ref != nil {
			// Find which property this is
			propExpr, ok := expr.Name.(MemberExpression)
			if ok {
				if propSym, ok := propExpr.Property.(SymbolExpression); ok {
					propIndex, err := propExpr.resolveStructAccess(val.SymbolTableEntry.Ref, propSym.Value)
					if err == nil {
						// Get the AST type from metadata
						astType := val.SymbolTableEntry.Ref.Metadata.Types[propIndex]

						// Check if it's a pointer type
						if ptrType, ok := astType.(PointerType); ok {
							isPointerType = true
							elementType = ptrType.Underlying.LLVMType()
							arrayType = elementType
						} else if arrType, ok := astType.(ArrayType); ok {
							// Embedded array type
							elementType = arrType.Underlying.LLVMType()
							arrayType = arrType.LLVMType()
							elementsCount = arrType.Size
						}
					}
				}
			}
		}

		// Fallback if we couldn't determine from metadata
		if elementType.C == nil {
			elementType = llvm.GlobalContext().Int32Type()
			arrayType = elementType
		}

		// Handle pointer vs embedded array differently
		if isPointerType {
			// Load the pointer value from the struct field
			pointerValue := ctx.Builder.CreateLoad(*val.StuctPropertyValueType, *val.Value, "")
			array = &SymbolTableEntry{Value: pointerValue}
		} else {
			// Embedded array - use the field address directly
			array = &SymbolTableEntry{Value: *val.Value}
		}

		entry = &ArraySymbolTableEntry{
			UnderlyingType: elementType,
			Type:           arrayType,
			ElementsCount:  elementsCount,
		}
	default:
		panic("ArrayAccessExpression not implemented")
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
	case SymbolExpression, BinaryExpression:
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
