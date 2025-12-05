package ast

import (
	"encoding/json"
	"fmt"
	"swahili/lang/lexer"

	"tinygo.org/x/go-llvm"
)

type MemberExpression struct {
	Object   Expression
	Property Expression
	Computed bool
	Tokens   []lexer.Token
}

var _ Expression = (*MemberExpression)(nil)

func (expr MemberExpression) CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult) {
	if nestedMember, ok := expr.Object.(MemberExpression); ok {
		err, nestedAddr := expr.getNestedMemberAddress(ctx, nestedMember)
		if err != nil {
			return err, nil
		}

		baseObj, err := expr.findBaseSymbol(nestedMember.Object)
		if err != nil {
			return err, nil
		}

		err, varDef := ctx.FindSymbol(baseObj.Value)
		if err != nil {
			return err, nil
		}

		nestedStructType, err := expr.getNestedStructType(ctx, nestedMember, varDef.Ref)
		if err != nil {
			return err, nil
		}

		propName, err := expr.getProperty()
		if err != nil {
			return err, nil
		}

		propIndex, err := expr.resolveStructAccess(nestedStructType, propName)
		if err != nil {
			return err, nil
		}

		addr := ctx.Builder.CreateStructGEP(nestedStructType.LLVMType, nestedAddr, propIndex, "")
		propType := nestedStructType.PropertyTypes[propIndex]
		return nil, &CompilerResult{
			Value:                  &addr,
			SymbolTableEntry:       &SymbolTableEntry{Ref: nestedStructType},
			StuctPropertyValueType: &propType,
		}
	}

	obj, ok := expr.Object.(SymbolExpression)
	if !ok {
		return fmt.Errorf("struct object should be a symbol"), nil
	}

	err, varDef := ctx.FindSymbol(obj.Value)
	if err != nil {
		return fmt.Errorf("variable %s is not defined", obj.Value), nil
	}

	propName, err := expr.getProperty()
	if err != nil {
		return err, nil
	}

	propIndex, err := expr.resolveStructAccess(varDef.Ref, propName)
	if err != nil {
		return err, nil
	}

	addr := ctx.Builder.CreateStructGEP(varDef.Ref.LLVMType, varDef.Value, propIndex, "")
	propType := varDef.Ref.PropertyTypes[propIndex]
	return nil, &CompilerResult{
		Value:                  &addr,
		SymbolTableEntry:       varDef,
		StuctPropertyValueType: &propType,
	}
}

func (expr MemberExpression) getNestedMemberAddress(ctx *CompilerCtx, member MemberExpression) (error, llvm.Value) {
	propName, err := member.getProperty()
	if err != nil {
		return err, llvm.Value{}
	}

	if nestedMember, ok := member.Object.(MemberExpression); ok {
		err, nestedAddr := expr.getNestedMemberAddress(ctx, nestedMember)
		if err != nil {
			return err, llvm.Value{}
		}

		baseObj, err := expr.findBaseSymbol(nestedMember.Object)
		if err != nil {
			return err, llvm.Value{}
		}

		err, varDef := ctx.FindSymbol(baseObj.Value)
		if err != nil {
			return err, llvm.Value{}
		}

		nestedStructType, err := expr.getNestedStructType(ctx, nestedMember, varDef.Ref)
		if err != nil {
			return err, llvm.Value{}
		}

		propIndex, err := expr.resolveStructAccess(nestedStructType, propName)
		if err != nil {
			return err, llvm.Value{}
		}

		addr := ctx.Builder.CreateStructGEP(nestedStructType.LLVMType, nestedAddr, propIndex, "")
		return nil, addr
	}

	obj, ok := member.Object.(SymbolExpression)
	if !ok {
		return fmt.Errorf("struct object should be a symbol"), llvm.Value{}
	}

	err, varDef := ctx.FindSymbol(obj.Value)
	if err != nil {
		return err, llvm.Value{}
	}

	propIndex, err := expr.resolveStructAccess(varDef.Ref, propName)
	if err != nil {
		return err, llvm.Value{}
	}

	addr := ctx.Builder.CreateStructGEP(varDef.Ref.LLVMType, varDef.Value, propIndex, "")
	return nil, addr
}

func (expr MemberExpression) getNestedStructType(
	ctx *CompilerCtx,
	member MemberExpression,
	baseStructType *StructSymbolTableEntry,
) (*StructSymbolTableEntry, error) {
	propName, err := member.getProperty()
	if err != nil {
		return nil, err
	}

	currentStructType := baseStructType
	if nestedMember, ok := member.Object.(MemberExpression); ok {
		currentStructType, err = expr.getNestedStructType(ctx, nestedMember, baseStructType)
		if err != nil {
			return nil, err
		}
	}

	propIndex, err := expr.resolveStructAccess(currentStructType, propName)
	if err != nil {
		return nil, err
	}

	propType := currentStructType.Metadata.Types[propIndex]
	if symbolType, ok := propType.(SymbolType); ok {
		err, structDef := ctx.FindStructSymbol(symbolType.Name)
		if err != nil {
			return nil, fmt.Errorf("cannot find struct type %s", symbolType.Name)
		}
		return structDef, nil
	}

	return nil, fmt.Errorf("property %s is not a struct type", propName)
}

func (expr MemberExpression) TokenStream() []lexer.Token {
	return expr.Tokens
}

func (expr MemberExpression) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["Object"] = expr.Object
	m["Property"] = expr.Property
	m["Computed"] = expr.Computed

	res := make(map[string]any)
	res["ast.MemberExpression"] = m

	return json.Marshal(res)
}

func (expr MemberExpression) findBaseSymbol(obj Expression) (SymbolExpression, error) {
	switch v := obj.(type) {
	case SymbolExpression:
		return v, nil
	case MemberExpression:
		return expr.findBaseSymbol(v.Object)
	default:
		return SymbolExpression{}, fmt.Errorf("cannot resolve base object")
	}
}

func (expr MemberExpression) getProperty() (string, error) {
	prop, ok := expr.Property.(SymbolExpression)
	if !ok {
		return "", fmt.Errorf("struct property should be a symbol")
	}
	return prop.Value, nil
}

func (expr MemberExpression) resolveStructAccess(
	structType *StructSymbolTableEntry,
	propName string,
) (int, error) {
	err, propIndex := structType.Metadata.PropertyIndex(propName)
	if err != nil {
		return 0, fmt.Errorf("struct %s has no field %s", structType.Metadata.Name, propName)
	}
	return propIndex, nil
}
