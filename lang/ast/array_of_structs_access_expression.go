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
	var array *SymbolTableEntry

	var entry *ArraySymbolTableEntry

	var err error

	var name string

	switch typedExpr := expr.Name.(type) {
	case SymbolExpression:
		err, array = ctx.FindSymbol(typedExpr.Value)
		if err != nil {
			key := "ArrayAccessExpression.NotFoundInSymbolTable"

			return ctx.Dialect.Error(key, typedExpr.Value), nil, nil, nil
		}

		err, entry = ctx.FindArraySymbol(typedExpr.Value)
		if err != nil {
			key := "ArrayAccessExpression.NotFoundInArraysSymbolTable"

			return ctx.Dialect.Error(key, typedExpr.Value), nil, nil, nil
		}

		name = typedExpr.Value
	case MemberExpression:
		err, val := expr.Name.CompileLLVM(ctx)
		if err != nil {
			return err, nil, nil, nil
		}

		var elementType llvm.Type

		var arrayType llvm.Type

		var typedef *StructSymbolTableEntry

		if val.SymbolTableEntry == nil || val.SymbolTableEntry.Ref == nil {
			return fmt.Errorf("ArrayOfStructsAccessExpression: Missing SymbolTableEntry"), nil, nil, nil
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

		coltype, ok := astType.(ArrayType)
		if !ok {
			err := fmt.Errorf("ArrayOfStructsAccessExpression type is not an array")

			return err, nil, nil, nil
		}

		err, elementType = coltype.Underlying.LLVMType(ctx)
		if err != nil {
			return err, nil, nil, nil
		}

		err, arrayType = coltype.LLVMType(ctx)
		if err != nil {
			return err, nil, nil, nil
		}

		underlying, _ := coltype.Underlying.(SymbolType)

		err, typedef = ctx.FindStructSymbol(underlying.Name)
		if err != nil {
			return err, nil, nil, nil
		}

		array = &SymbolTableEntry{Value: *val.Value}

		entry = &ArraySymbolTableEntry{
			UnderlyingTypeDef: typedef,
			UnderlyingType:    elementType,
			Type:              arrayType,
			ElementsCount:     coltype.Size,
		}
	default:
		format := "Expression %s not supported in ArrayOfStructsAccessExpression"

		return fmt.Errorf(format, typedExpr), nil, nil, nil
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

			return ctx.Dialect.Error(key, int(idx.Value), name), nil, nil, nil
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
		err, ref = ctx.FindStructSymbol(symbolType.Name)
		if err != nil {
			return err, nil
		}
	}

	return nil, &CompilerResult{
		Value: &load,
		SymbolTableEntry: &SymbolTableEntry{
			Address: &structPtr,
			Ref:     ref,
		},
	}
}

func (expr ArrayOfStructsAccessExpression) Accept(g CodeGenerator) error {
	return g.VisitArrayOfStructsAccessExpression(&expr)
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
