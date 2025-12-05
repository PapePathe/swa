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

		var baseObj SymbolExpression
		currentExpr := nestedMember
		for {
			if sym, ok := currentExpr.Object.(SymbolExpression); ok {
				baseObj = sym
				break
			} else if mem, ok := currentExpr.Object.(MemberExpression); ok {
				currentExpr = mem
			} else {
				return fmt.Errorf("cannot resolve base object for nested member access"), nil
			}
		}

		err, varDef := ctx.FindSymbol(baseObj.Value)
		if err != nil {
			return err, nil
		}

		nestedStructType, err := expr.getNestedStructType(ctx, nestedMember, varDef.Ref)
		if err != nil {
			return err, nil
		}

		prop, ok := expr.Property.(SymbolExpression)
		if !ok {
			return fmt.Errorf("struct property should be a symbol %v", expr.Property), nil
		}

		err, propIndex := nestedStructType.Metadata.PropertyIndex(prop.Value)
		if err != nil {
			return fmt.Errorf("Struct %s does not have a field named %s", nestedStructType.Metadata.Name, prop.Value), nil
		}

		addr := ctx.Builder.CreateStructGEP(nestedStructType.LLVMType, nestedAddr, propIndex, "")
		propType := nestedStructType.PropertyTypes[propIndex]
		return nil, &CompilerResult{Value: &addr, SymbolTableEntry: &SymbolTableEntry{Ref: nestedStructType}, StuctPropertyValueType: &propType}
	}

	obj, ok := expr.Object.(SymbolExpression)
	if !ok {
		return fmt.Errorf("struct object should be a symbol %v", obj), nil
	}

	err, varDef := ctx.FindSymbol(obj.Value)
	if err != nil {
		return fmt.Errorf("Variable %s of type Struct is not defined", obj.Value), nil
	}

	prop, ok := expr.Property.(SymbolExpression)
	if !ok {
		return fmt.Errorf("struct property should be a symbol %v", prop), nil
	}

	err, propIndex := varDef.Ref.Metadata.PropertyIndex(prop.Value)
	if err != nil {
		return fmt.Errorf("Struct %s does not have a field named %s", varDef.Ref.Metadata.Name, prop), nil
	}

	addr := ctx.Builder.CreateStructGEP(varDef.Ref.LLVMType, varDef.Value, propIndex, "")
	propType := varDef.Ref.PropertyTypes[propIndex]

	return nil, &CompilerResult{Value: &addr, SymbolTableEntry: varDef, StuctPropertyValueType: &propType}
}

func (expr MemberExpression) getNestedMemberAddress(ctx *CompilerCtx, member MemberExpression) (error, llvm.Value) {
	if nestedMember, ok := member.Object.(MemberExpression); ok {
		err, nestedAddr := expr.getNestedMemberAddress(ctx, nestedMember)
		if err != nil {
			return err, llvm.Value{}
		}

		var baseObj SymbolExpression
		currentExpr := nestedMember
		for {
			if sym, ok := currentExpr.Object.(SymbolExpression); ok {
				baseObj = sym
				break
			} else if mem, ok := currentExpr.Object.(MemberExpression); ok {
				currentExpr = mem
			} else {
				return fmt.Errorf("cannot resolve base object"), llvm.Value{}
			}
		}

		err, varDef := ctx.FindSymbol(baseObj.Value)
		if err != nil {
			return err, llvm.Value{}
		}

		nestedStructType, err := expr.getNestedStructType(ctx, nestedMember, varDef.Ref)
		if err != nil {
			return err, llvm.Value{}
		}

		prop, ok := member.Property.(SymbolExpression)
		if !ok {
			return fmt.Errorf("struct property should be a symbol"), llvm.Value{}
		}

		err, propIndex := nestedStructType.Metadata.PropertyIndex(prop.Value)
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

	prop, ok := member.Property.(SymbolExpression)
	if !ok {
		return fmt.Errorf("struct property should be a symbol"), llvm.Value{}
	}

	err, propIndex := varDef.Ref.Metadata.PropertyIndex(prop.Value)
	if err != nil {
		return err, llvm.Value{}
	}

	addr := ctx.Builder.CreateStructGEP(varDef.Ref.LLVMType, varDef.Value, propIndex, "")
	return nil, addr
}

func (expr MemberExpression) CompileLLVMForPropertyAccess(ctx *CompilerCtx) (error, *llvm.Value) {
	if nestedMember, ok := expr.Object.(MemberExpression); ok {
		err, nestedAddr := expr.getNestedMemberAddress(ctx, nestedMember)
		if err != nil {
			return err, nil
		}

		var baseObj SymbolExpression
		currentExpr := nestedMember
		for {
			if sym, ok := currentExpr.Object.(SymbolExpression); ok {
				baseObj = sym
				break
			} else if mem, ok := currentExpr.Object.(MemberExpression); ok {
				currentExpr = mem
			} else {
				return fmt.Errorf("cannot resolve base object for nested member access"), nil
			}
		}

		err, varDef := ctx.FindSymbol(baseObj.Value)
		if err != nil {
			return err, nil
		}

		nestedStructType, err := expr.getNestedStructType(ctx, nestedMember, varDef.Ref)
		if err != nil {
			return err, nil
		}

		prop, ok := expr.Property.(SymbolExpression)
		if !ok {
			return fmt.Errorf("struct property should be a symbol %v", expr.Property), nil
		}

		err, propIndex := nestedStructType.Metadata.PropertyIndex(prop.Value)
		if err != nil {
			return fmt.Errorf("Struct %s does not have a field named %s", nestedStructType.Metadata.Name, prop.Value), nil
		}

		addr := ctx.Builder.CreateStructGEP(nestedStructType.LLVMType, nestedAddr, propIndex, "")
		loadedval := ctx.Builder.CreateLoad(nestedStructType.PropertyTypes[propIndex], addr, "")

		return nil, &loadedval
	}

	obj, ok := expr.Object.(SymbolExpression)
	if !ok {
		return fmt.Errorf("struct object should be a symbol %v", obj), nil
	}

	err, varDef := ctx.FindSymbol(obj.Value)
	if err != nil {
		return fmt.Errorf("Variable %s of type Struct is not defined", obj.Value), nil
	}

	prop, ok := expr.Property.(SymbolExpression)
	if !ok {
		return fmt.Errorf("struct property should be a symbol %v", prop), nil
	}

	err, propIndex := varDef.Ref.Metadata.PropertyIndex(prop.Value)
	if err != nil {
		return fmt.Errorf("Struct %s does not have a field named %s", varDef.Ref.Metadata.Name, prop), nil
	}

	addr := ctx.Builder.CreateStructGEP(varDef.Ref.LLVMType, varDef.Value, propIndex, "")
	loadedval := ctx.Builder.CreateLoad(varDef.Ref.PropertyTypes[propIndex], addr, "")

	return nil, &loadedval
}

func (expr MemberExpression) getNestedStructType(
	ctx *CompilerCtx,
	member MemberExpression,
	baseStructType *StructSymbolTableEntry,
) (*StructSymbolTableEntry, error) {
	prop, ok := member.Property.(SymbolExpression)
	if !ok {
		return nil, fmt.Errorf("struct property should be a symbol")
	}

	currentStructType := baseStructType
	if nestedMember, ok := member.Object.(MemberExpression); ok {
		var err error
		currentStructType, err = expr.getNestedStructType(ctx, nestedMember, baseStructType)
		if err != nil {
			return nil, err
		}
	}

	err, propIndex := currentStructType.Metadata.PropertyIndex(prop.Value)
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

	return nil, fmt.Errorf("property %s is not a struct type", prop.Value)
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
