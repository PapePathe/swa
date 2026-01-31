package compiler

import (
	"fmt"
	"swahili/lang/ast"

	"tinygo.org/x/go-llvm"
)

func (g *LLVMGenerator) VisitArrayOfStructsAccessExpression(node *ast.ArrayOfStructsAccessExpression) error {
	old := g.logger.Step("ArrayOfStructsAccessExpr")

	defer g.logger.Restore(old)

	g.Debugf("%s[%s].%s", node.Name, node.Index, node.Property)

	err, array, arrayEntry, itemIndex := g.findArrayOfStructsSymbolTableEntry(node)
	if err != nil {
		return err
	}

	if array.Address == nil {
		g.Ctx.PrintVarNames()
		fmt.Println(arrayEntry.ElementsCount)
		g.NotImplemented("NIL POINTER array.Address should be set.")
	}

	itemPtr := g.Ctx.Builder.CreateInBoundsGEP(
		arrayEntry.Type,
		*array.Address,
		itemIndex,
		"",
	)

	propName, ok := node.Property.(*ast.SymbolExpression)
	if !ok {
		g.NotImplemented(fmt.Sprintf("Type %s not supported in ArrayOfStructsAccessExpression", node.Property))
	}

	g.Debugf("struct name: %s, props: %v", array.Ref.Metadata.Name, array.Ref.Metadata.Properties)

	err, propIndex := array.Ref.Metadata.PropertyIndex(propName.Value)
	if err != nil {
		g.Debugf("prop %s not found", propName.Value)

		format := "property (%s) not found in struct %s"

		return fmt.Errorf(format, propName.Value, array.Ref.Metadata.Name)
	}

	structPtr := g.Ctx.Builder.CreateStructGEP(
		array.Ref.LLVMType,
		itemPtr,
		propIndex,
		"",
	)

	proptype := array.Ref.PropertyTypes[propIndex]

	if proptype.TypeKind() == llvm.StructTypeKind {
		nestedstruct, _ := array.Ref.Embeds[propName.Value]
		array.Ref = &nestedstruct
	}

	if proptype.TypeKind() == llvm.ArrayTypeKind &&
		proptype.ElementType().TypeKind() == llvm.StructTypeKind {
		format := "processing nested array of struct(%s)"

		g.Debugf(format, arrayEntry.UnderlyingTypeDef.Metadata.Name)

		array.Ref = arrayEntry.UnderlyingTypeDef
	}

	res := &CompilerResult{
		Value:                  &structPtr,
		SymbolTableEntry:       array,
		StuctPropertyValueType: &proptype,
		ArraySymbolTableEntry:  arrayEntry,
	}

	debugf := "result object, value:%v, sdef: %s, adef: %s"

	g.Debugf(
		debugf,
		res.Value,
		array.Ref.Metadata.Name,
		arrayEntry.UnderlyingTypeDef.Metadata.Name,
	)

	g.setLastResult(res)

	return nil
}

func (g *LLVMGenerator) findArrayOfStructsSymbolTableEntry(
	expr *ast.ArrayOfStructsAccessExpression,
) (error, *SymbolTableEntry, *ArraySymbolTableEntry, []llvm.Value) {
	err, baseSym, arrayMeta := g.resolveArrayOfStructsBase(expr.Name)
	if err != nil {
		return err, nil, nil, nil
	}

	err, indices := g.resolveGepIndices(expr.Index)
	if err != nil {
		return err, nil, nil, nil
	}

	return nil, baseSym, arrayMeta, indices
}

func (g *LLVMGenerator) resolveArrayOfStructsBase(nameNode ast.Node) (error, *SymbolTableEntry, *ArraySymbolTableEntry) {
	old := g.logger.Step("resolveArrayOfStructsBase")
	defer g.logger.Restore(old)

	g.Debugf("%s", nameNode)

	switch typednode := nameNode.(type) {
	case *ast.ArrayOfStructsAccessExpression:
		g.Debugf("processing ast.ArrayOfStructsAccessExpression")

		err := typednode.Accept(g)
		if err != nil {
			return err, nil, nil
		}

		lastres := g.getLastResult()

		g.Debugf("%s", nameNode)

		debugf := "base result object, value:%v, sdef: %s, adef: %s"

		g.Debugf(
			debugf,
			lastres.Value,
			lastres.SymbolTableEntry.Ref.Metadata.Name,
			lastres.ArraySymbolTableEntry.UnderlyingTypeDef.Metadata.Name,
		)

		mem, _ := typednode.Property.(*ast.SymbolExpression)

		propIndex, err := g.resolveStructAccess(lastres.SymbolTableEntry.Ref, mem.Value)
		if err != nil {
			return err, nil, nil
		}

		propType := lastres.SymbolTableEntry.Ref.PropertyTypes[propIndex]

		if propType.TypeKind() == llvm.ArrayTypeKind &&
			propType.ElementType().TypeKind() == llvm.StructTypeKind {
			g.Debugf("proptype %v", propType.TypeKind())

			nested, ok := lastres.SymbolTableEntry.Ref.ArrayEmbeds[mem.Value]
			if ok {
				lastres.SymbolTableEntry.Ref = nested.UnderlyingTypeDef

				itemPtr := g.Ctx.Builder.CreateInBoundsGEP(
					propType,
					*lastres.Value,
					[]llvm.Value{
						llvm.ConstInt(llvm.GlobalContext().Int32Type(), uint64(propIndex), false),
					},
					"",
				)

				g.Debugf("item ptr of array sym entry %s", itemPtr)

				lastres.SymbolTableEntry.Address = &itemPtr
				lastres.ArraySymbolTableEntry = &nested
			}
		}

		return nil, lastres.SymbolTableEntry, lastres.ArraySymbolTableEntry
	case *ast.MemberExpression:
		g.Debugf("processing ast.MemberExpression %s.%s", typednode.Object, typednode.Property)

		err := typednode.Object.Accept(g)
		if err != nil {
			return err, nil, nil
		}

		lastres := g.getLastResult()

		propName, err := g.getProperty(typednode)
		if err != nil {
			return err, nil, nil
		}

		propIndex, err := g.resolveStructAccess(lastres.SymbolTableEntry.Ref, propName)
		if err != nil {
			return err, nil, nil
		}

		propType := lastres.SymbolTableEntry.Ref.PropertyTypes[propIndex]
		underlyingTypeDef := lastres.SymbolTableEntry.Ref

		if propType.TypeKind() == llvm.ArrayTypeKind {
			if propType.ElementType().TypeKind() == llvm.StructTypeKind {
				g.Debugf("ast.MemberExpression with array of structs %s", propName)

				proptype := lastres.SymbolTableEntry.Ref.Metadata.Types[propIndex]
				arrtype, _ := proptype.(ast.ArrayType)
				stype, _ := arrtype.Underlying.(ast.SymbolType)

				err, sym := g.Ctx.FindStructSymbol(stype.Name)
				if err != nil {
					return err, nil, nil
				}

				underlyingTypeDef = sym
				lastres.SymbolTableEntry.Ref = sym
			}
		}

		return nil, lastres.SymbolTableEntry, &ArraySymbolTableEntry{
			ElementsCount:     propType.ArrayLength(),
			UnderlyingType:    propType.ElementType(),
			Type:              propType,
			UnderlyingTypeDef: underlyingTypeDef,
		}
	case *ast.SymbolExpression:
		g.Debugf("processing ast.SymbolExpression")

		name := typednode.Value

		err, symEntry := g.Ctx.FindSymbol(name)
		if err != nil {
			return g.Ctx.Dialect.Error("ArrayOfStructsAccessExpression.NotFoundInSymbolTable", name), nil, nil
		}

		err, arrEntry := g.Ctx.FindArraySymbol(name)
		if err != nil {
			return err, nil, nil
		}

		return nil, symEntry, arrEntry
	default:
		return fmt.Errorf("ArrayOfStructsAccessExpression not implemented for %T", typednode), nil, nil
	}
}
