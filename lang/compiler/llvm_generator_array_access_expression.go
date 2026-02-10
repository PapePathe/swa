package compiler

import (
	"fmt"
	"swahili/lang/ast"

	"tinygo.org/x/go-llvm"
)

func (g *LLVMGenerator) VisitArrayAccessExpression(node *ast.ArrayAccessExpression) error {
	old := g.logger.Step("ArrayAccessExpr")

	defer g.logger.Restore(old)

	g.Debugf("(%d) %s[%s]", node.Tokens[0].Line, node.Name, node.Index)

	err, entry, array, indices := g.findArraySymbolTableEntry(node)
	if err != nil {
		return err
	}

	var arrayAddr llvm.Value
	if array.Address == nil {
		arrayAddr = array.Value
	} else {
		arrayAddr = *array.Address
	}

	if entry.UnderlyingType.IsNil() {
		key := "LLVMGenerator.VisitArrayAccessExpression.UnderlyingTypeNotSet"

		return g.Ctx.Dialect.Error(key)
	}

	itemPtr := g.Ctx.Builder.CreateInBoundsGEP(
		entry.UnderlyingType,
		arrayAddr,
		indices,
		"",
	)

	array.Address = &itemPtr

	g.setLastResult(&CompilerResult{
		Value:                 &itemPtr,
		ArraySymbolTableEntry: entry,
		SymbolTableEntry:      array,
		SwaType:               node.SwaType,
	})

	return nil
}

func (g *LLVMGenerator) findArraySymbolTableEntry(
	expr *ast.ArrayAccessExpression,
) (error, *ArraySymbolTableEntry, *SymbolTableEntry, []llvm.Value) {
	old := g.logger.Step("findArraySymbolTableEntry")

	defer g.logger.Restore(old)

	var name string

	var array *SymbolTableEntry

	var entry *ArraySymbolTableEntry

	switch typedExpr := expr.Name.(type) {
	case *ast.ArrayOfStructsAccessExpression:
		g.Debugf("array access expr (%d) %s ", expr.Tokens[0].Line, expr.Name)

		err := expr.Name.Accept(g)
		if err != nil {
			return err, nil, nil, nil
		}

		lastesult := g.getLastResult()

		symbol, _ := typedExpr.Property.(*ast.SymbolExpression)

		embed, _ := lastesult.SymbolTableEntry.Ref.ArrayEmbeds[symbol.Value]

		propIndex, err := g.resolveStructAccess(lastesult.SymbolTableEntry.Ref, symbol.Value)
		propType := lastesult.SymbolTableEntry.Ref.PropertyTypes[propIndex]
		g.Debugf(
			"struct types for prop(%s) idx(%d): %s %s",
			typedExpr.Property,
			propIndex,
			lastesult.SymbolTableEntry.Address,
			propType.TypeKind(),
		)

		arrayItemPtr := g.Ctx.Builder.CreateInBoundsGEP(
			propType,
			*lastesult.SymbolTableEntry.Address,
			[]llvm.Value{
				llvm.ConstInt(llvm.GlobalContext().Int32Type(), uint64(propIndex), false),
			},
			fmt.Sprintf("%s[%s]", expr.Name, expr.Index),
		)

		array = lastesult.SymbolTableEntry
		array.Address = &arrayItemPtr
		array.Ref = embed.UnderlyingTypeDef
		entry = &embed

	case *ast.SymbolExpression:
		varName, _ := expr.Name.(*ast.SymbolExpression)

		err, arrayEntry := g.Ctx.FindSymbol(varName.Value)
		if err != nil {
			key := "ArrayAccessExpression.NotFoundInSymbolTable"

			return g.Ctx.Dialect.Error(key, varName.Value), nil, nil, nil
		}

		array = arrayEntry

		err, arraySymEntry := g.Ctx.FindArraySymbol(varName.Value)
		if err != nil {
			key := "ArrayAccessExpression.NotFoundInArraySymbolTable"

			return g.Ctx.Dialect.Error(key, varName.Value), nil, nil, nil
		}

		entry = arraySymEntry
		name = varName.Value

		arrtyp, ok := arrayEntry.DeclaredType.(ast.ArrayType)
		if !ok {
			key := "ArrayAccessExpression.DeclaredTypeIsNotArray"

			return g.Ctx.Dialect.Error(key), nil, nil, nil
		}
		expr.SwaType = arrtyp.Underlying

	case *ast.MemberExpression:
		err := expr.Name.Accept(g)
		if err != nil {
			return err, nil, nil, nil
		}

		val := g.getLastResult()

		var elementsCount int = -1

		var isPointerType bool = false

		if val.SymbolTableEntry == nil && val.SymbolTableEntry.Ref == nil {
			key := "LLVMGenerator.VisitArrayAccessExpression.MissingSymbolTableEntry"

			return g.Ctx.Dialect.Error(key), nil, nil, nil
		}

		propExpr, _ := expr.Name.(*ast.MemberExpression)
		propSym, _ := propExpr.Property.(*ast.SymbolExpression)

		if val.SymbolTableEntry.Ref == nil {
			key := "LLVMGenerator.VisitArrayAccessExpression.PropertyIsnotAnArray"

			return g.Ctx.Dialect.Error(key, propSym.Value), nil, nil, nil
		}

		propIndex, err := g.resolveStructAccess(val.SymbolTableEntry.Ref, propSym.Value)
		if err != nil {
			return err, nil, nil, nil
		}

		astType := val.SymbolTableEntry.Ref.Metadata.Types[propIndex]

		expr.SwaType = astType

		err = astType.Accept(g)
		if err != nil {
			return err, nil, nil, nil
		}

		etype := g.getLastTypeVisitResult()

		switch coltype := astType.(type) {
		case ast.PointerType:
			isPointerType = true
		case ast.ArrayType:
			elementsCount = coltype.Size
			expr.SwaType = coltype.Underlying
		default:
			key := "LLVMGenerator.VisitArrayAccessExpression.PropertyIsnotAnArray"
			err := g.Ctx.Dialect.Error(key, propSym.Value)

			return err, nil, nil, nil
		}

		if isPointerType {
			pointerValue := g.Ctx.Builder.CreateLoad(*val.StuctPropertyValueType, *val.Value, "")
			// TODO declared type should be set
			array = &SymbolTableEntry{Value: pointerValue, DeclaredType: astType}
		} else {
			// TODO declared type should be set
			array = &SymbolTableEntry{Value: *val.Value, DeclaredType: astType}
		}

		entry = &ArraySymbolTableEntry{
			UnderlyingType: etype.SubType,
			Type:           etype.Type,
			ElementsCount:  elementsCount,
		}
	default:
		key := "LLVMGenerator.VisitArrayAccessExpression.NotImplementedFor"
		err := g.Ctx.Dialect.Error(key, expr.Name)

		return err, nil, nil, nil
	}

	var indices []llvm.Value

	switch expr.Index.(type) {
	case *ast.NumberExpression:
		idx, _ := expr.Index.(*ast.NumberExpression)

		if int(idx.Value) < 0 {
			key := "ArrayAccessExpression.AccessedIndexIsNotANumber"

			return g.Ctx.Dialect.Error(key, expr.Index), nil, nil, nil
		}

		if entry.ElementsCount >= 0 && int(idx.Value) > entry.ElementsCount-1 {
			key := "ArrayAccessExpression.IndexOutOfBounds"

			return g.Ctx.Dialect.Error(key, int(idx.Value), name), nil, nil, nil
		}

		indices = []llvm.Value{
			llvm.ConstInt((*g.Ctx.Context).Int32Type(), uint64(idx.Value), false),
		}
	case *ast.SymbolExpression, *ast.BinaryExpression, *ast.MemberExpression:
		err := expr.Index.Accept(g)
		if err != nil {
			return err, nil, nil, nil
		}

		res := g.getLastResult()

		if res.Value.Type().TypeKind() == llvm.PointerTypeKind {
			load := g.Ctx.Builder.CreateLoad(res.Value.AllocatedType(), *res.Value, "")
			indices = []llvm.Value{load}
		} else {
			indices = []llvm.Value{*res.Value}
		}
	default:
		key := "ArrayAccessExpression.AccessedIndexIsNotANumber"

		return g.Ctx.Dialect.Error(key, expr.Index), nil, nil, nil
	}

	return nil, entry, array, indices
}

func (g *LLVMGenerator) resolveStructAccess(
	structType *StructSymbolTableEntry,
	propName string,
) (int, error) {
	err, propIndex := structType.Metadata.PropertyIndex(propName)
	if err != nil {
		key := "LLVMGenerator.VisitArrayAccessExpression.FieldDoesNotExistInStruct"

		return 0, g.Ctx.Dialect.Error(key, structType.Metadata.Name, propName)
	}

	return propIndex, nil
}
