package compiler

import (
	"fmt"
	"swahili/lang/ast"

	"tinygo.org/x/go-llvm"
)

// VisitMemberExpression implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitMemberExpression(node *ast.MemberExpression) error {
	old := g.logger.Step("MemberExpr")

	defer g.logger.Restore(old)

	g.Debugf("VisitMemberExpression %s.%s", node.Object, node.Property)

	switch objectType := node.Object.(type) {
	case *ast.SymbolExpression:
		err, varDef := g.Ctx.FindSymbol(objectType.Value)
		if err != nil {
			key := "LLVMGenerator.VisitMemberExpression.NotDefined"

			return g.Ctx.Dialect.Error(key, objectType.Value)
		}

		if varDef.Ref == nil {
			key := "LLVMGenerator.VisitMemberExpression.NotAStructInstance"

			return g.Ctx.Dialect.Error(key, objectType.Value)
		}

		propName, err := g.getProperty(node)
		if err != nil {
			return err
		}

		propIndex, err := g.resolveStructAccess(varDef.Ref, propName)
		if err != nil {
			return err
		}

		var baseValue llvm.Value = varDef.Value
		if varDef.Address != nil {
			baseValue = *varDef.Address
		}

		addr := g.Ctx.Builder.CreateStructGEP(varDef.Ref.LLVMType, baseValue, propIndex, "")
		propType := varDef.Ref.PropertyTypes[propIndex]
		result := &CompilerResult{
			Value:                  &addr,
			SymbolTableEntry:       varDef,
			StuctPropertyValueType: &propType,
		}

		swatype := result.SymbolTableEntry.Ref.Metadata.Types[propIndex]

		if propType.TypeKind() == llvm.StructTypeKind {
			prop, _ := varDef.Ref.Embeds[propName]
			result.SymbolTableEntry.Ref = &prop
		}

		node.SwaType = swatype

		g.setLastResult(result)

		return nil
	case *ast.ArrayOfStructsAccessExpression:
		err := node.Object.Accept(g)
		if err != nil {
			return err
		}

		lastresult := g.getLastResult()

		propName, err := g.getProperty(node)
		if err != nil {
			return err
		}

		propIndex, err := g.resolveStructAccess(lastresult.SymbolTableEntry.Ref, propName)
		if err != nil {
			return err
		}

		propType := lastresult.SymbolTableEntry.Ref.PropertyTypes[propIndex]
		addr := g.Ctx.Builder.CreateStructGEP(
			lastresult.SymbolTableEntry.Ref.LLVMType,
			*lastresult.Value,
			propIndex,
			"")
		result := &CompilerResult{
			Value:                  &addr,
			SymbolTableEntry:       lastresult.SymbolTableEntry,
			StuctPropertyValueType: &propType,
		}

		if propType.TypeKind() == llvm.StructTypeKind {
			// TODO handle the case where propName is not in embeds
			prop, _ := lastresult.SymbolTableEntry.Ref.Embeds[propName]
			result.SymbolTableEntry.Ref = &prop
		}

		g.setLastResult(result)

		return nil
	case *ast.MemberExpression:
		err := node.Object.Accept(g)
		if err != nil {
			return err
		}

		prop, _ := node.Property.(*ast.SymbolExpression)
		result := g.getLastResult()

		propIndex, err := g.resolveStructAccess(result.SymbolTableEntry.Ref, prop.Value)
		if err != nil {
			return err
		}

		addr := g.Ctx.Builder.CreateStructGEP(*result.StuctPropertyValueType, *result.Value, propIndex, "")

		if propIndex > len(result.SymbolTableEntry.Ref.PropertyTypes) ||
			len(result.SymbolTableEntry.Ref.PropertyTypes) == 0 {
			key := "LLVMGenerator.VisitMemberExpression.PropertyNotFound"

			return g.Ctx.Dialect.Error(key, prop.Value, propIndex)
		}

		propType := result.SymbolTableEntry.Ref.PropertyTypes[propIndex]

		finalresult := &CompilerResult{
			Value:                  &addr,
			SymbolTableEntry:       result.SymbolTableEntry,
			StuctPropertyValueType: &propType,
		}

		swatype := result.SymbolTableEntry.Ref.Metadata.Types[propIndex]

		if propType.TypeKind() == llvm.StructTypeKind {
			sEntry, ok := result.SymbolTableEntry.Ref.Embeds[prop.Value]
			if ok {
				result.SymbolTableEntry.Ref = &sEntry
			}
		}

		node.SwaType = swatype

		g.setLastResult(finalresult)

		return nil

	default:
		g.NotImplemented(fmt.Sprintf("VisitMemberExpression not implemented for %T", objectType))
	}

	return nil
}
