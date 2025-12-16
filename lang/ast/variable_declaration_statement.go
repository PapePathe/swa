package ast

import (
	"encoding/json"
	"fmt"
	"swahili/lang/lexer"

	"tinygo.org/x/go-llvm"
)

// VarDeclarationStatement ...
type VarDeclarationStatement struct {
	// The name of the variable
	Name string
	// Wether or not the variable is a constant
	IsConstant bool
	// The value assigned to the variable
	Value Expression
	// The explicit type of the variable
	ExplicitType Type
	Tokens       []lexer.Token
}

var _ Statement = (*VarDeclarationStatement)(nil)

func (vd VarDeclarationStatement) CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult) {
	if ctx.SymbolExistsInCurrentScope(vd.Name) {
		return fmt.Errorf("variable %s is aleady defined", vd.Name), nil
	}

	var val *CompilerResult

	var err error

	if vd.Value != nil {
		err, val = vd.Value.CompileLLVM(ctx)
		if err != nil {
			return err, nil
		}
	}

	if val == nil {
		if vd.ExplicitType == nil {
			return fmt.Errorf("VarDeclarationStatement: return value is nil and no explicit type <%s>", vd.Name), nil
		}

		err, llvmType := vd.ExplicitType.LLVMType(ctx)
		if err != nil {
			return err, nil
		}

		alloc := ctx.Builder.CreateAlloca(llvmType, fmt.Sprintf("alloc.%s", vd.Name))
		ctx.Builder.CreateStore(llvm.ConstNull(llvmType), alloc)

		entry := &SymbolTableEntry{
			Value:        alloc,
			Address:      &alloc,
			DeclaredType: vd.ExplicitType,
		}

		if arrayType, ok := vd.ExplicitType.(ArrayType); ok {
			err, underlyingLLVMType := arrayType.Underlying.LLVMType(ctx)
			if err != nil {
				return err, nil
			}

			arrayEntry := &ArraySymbolTableEntry{
				UnderlyingType: underlyingLLVMType,
				Type:           llvmType,
				ElementsCount:  arrayType.Size,
			}

			err = ctx.AddArraySymbol(vd.Name, arrayEntry)
			if err != nil {
				return err, nil
			}
		} else if symType, ok := vd.ExplicitType.(SymbolType); ok {
			err, typeDef := ctx.FindStructSymbol(symType.Name)
			if err != nil {
				return err, nil
			}

			entry.Ref = typeDef
		}

		err = ctx.AddSymbol(vd.Name, entry)
		if err != nil {
			return err, nil
		}

		return nil, nil
	}

	if _, ok := vd.Value.(MemberExpression); ok {
		if val.SymbolTableEntry != nil && val.SymbolTableEntry.Ref != nil {
			memberExpr, _ := vd.Value.(MemberExpression)
			propExpr, _ := memberExpr.Property.(SymbolExpression)

			err, propIndex := val.SymbolTableEntry.Ref.Metadata.PropertyIndex(propExpr.Value)
			if err != nil {
				return err, nil
			}
			// Load the value from the address so TypeCheck can work properly
			loadedValue := ctx.Builder.CreateLoad(val.SymbolTableEntry.Ref.PropertyTypes[propIndex], *val.Value, "")
			val.Value = &loadedValue
		} else {
			return fmt.Errorf("VarDeclarationStatement/MemberExpression: unable to determine type"), nil
		}
	}

	err = vd.TypeCheck(vd.ExplicitType.Value(), val.Value.Type())
	if err != nil {
		return err, nil
	}

	switch vd.Value.(type) {
	case ArrayOfStructsAccessExpression:
		alloc := ctx.Builder.CreateAlloca(val.Value.Type(), fmt.Sprintf("alloc.%s", vd.Name))
		ctx.Builder.CreateStore(*val.Value, alloc)

		err = ctx.AddSymbol(vd.Name, &SymbolTableEntry{Value: *val.Value, Address: &alloc, DeclaredType: vd.ExplicitType})
		if err != nil {
			return err, nil
		}
	case ArrayAccessExpression:
		load := ctx.Builder.CreateLoad(val.Value.AllocatedType(), *val.Value, "load.from-array")
		alloc := ctx.Builder.CreateAlloca(val.Value.AllocatedType(), fmt.Sprintf("alloc.%s", vd.Name))
		ctx.Builder.CreateStore(load, alloc)

		err = ctx.AddSymbol(vd.Name, &SymbolTableEntry{Value: load, Address: &alloc, DeclaredType: vd.ExplicitType})
		if err != nil {
			return err, nil
		}
	case StructInitializationExpression:
		return vd.compileStructInitializationExpression(ctx, val)
	case StringExpression:
		return vd.compileStringExpression(ctx, val)
	case SymbolExpression:
		alloc := ctx.Builder.CreateAlloca(val.Value.Type(), fmt.Sprintf("alloc.%s", vd.Name))
		ctx.Builder.CreateStore(*val.Value, alloc)

		entry := &SymbolTableEntry{
			Value:        *val.Value,
			Address:      &alloc,
			DeclaredType: vd.ExplicitType,
		}

		switch vd.ExplicitType.(type) {
		case SymbolType:
			entry.Ref = val.SymbolTableEntry.Ref
		default:
		}

		err = ctx.AddSymbol(vd.Name, entry)
		if err != nil {
			return err, nil
		}
	case NumberExpression, FloatExpression, BinaryExpression, FunctionCallExpression, MemberExpression:
		alloc := ctx.Builder.CreateAlloca(val.Value.Type(), fmt.Sprintf("alloc.%s", vd.Name))
		ctx.Builder.CreateStore(*val.Value, alloc)

		entry := &SymbolTableEntry{
			Value:        *val.Value,
			Address:      &alloc,
			DeclaredType: vd.ExplicitType,
		}
		if val.SymbolTableEntry != nil {
			entry.Ref = val.SymbolTableEntry.Ref
		}

		err = ctx.AddSymbol(vd.Name, entry)
		if err != nil {
			return err, nil
		}
	case ArrayInitializationExpression:
		return vd.compileArrArrayInitializationExpression(ctx, val)
	default:
		return fmt.Errorf("VarDeclarationStatement: Unhandled expression type (%v)", vd.Value), nil
	}

	return nil, nil
}

func (expr VarDeclarationStatement) TypeCheck(t DataType, k llvm.Type) error {
	switch k.TypeKind() {
	case llvm.ArrayTypeKind:
		switch t {
		case DataTypeSymbol:
		case DataTypeString, DataTypeArray:
			// we are ok
		default:
			return fmt.Errorf("expected %s got %s", t, k.TypeKind())
		}
	case llvm.DoubleTypeKind:
		switch t {
		case DataTypeSymbol:
		case DataTypeFloat:
			// we are ok
		default:
			return fmt.Errorf("expected %s got %s", t, k.TypeKind())
		}
	case llvm.IntegerTypeKind:
		switch t {
		case DataTypeSymbol:
		case DataTypeNumber:
			// we are ok
		default:
			return fmt.Errorf("expected %s got %s", t, k.TypeKind())
		}
	case llvm.PointerTypeKind:
		switch t {
		case DataTypeString:
			switch k.ElementType() {
			case llvm.GlobalContext().Int8Type():
				// we good
			default:
				// TODO fix this
				//	return fmt.Errorf("expected %s got pointer of unknown value %v", t, k.IsNil())
			}
		case DataTypeNumber:
			switch k.ElementType().TypeKind() {
			case llvm.IntegerTypeKind:
				// we good
			default:
				// return fmt.Errorf("expected %s got pointer of %v", t, k)
			}
		case DataTypeArray:
		case DataTypeSymbol:
		default:
			return fmt.Errorf("unsupported element %s", t)
		}
	}

	return nil
}

func (expr VarDeclarationStatement) TokenStream() []lexer.Token {
	return expr.Tokens
}

func (cs VarDeclarationStatement) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["Name"] = cs.Name
	m["IsConstant"] = cs.IsConstant
	m["Value"] = cs.Value
	m["ExplicitType"] = cs.ExplicitType

	res := make(map[string]any)
	res["ast.VarDeclarationStatement"] = m

	return json.Marshal(res)
}

func (vd VarDeclarationStatement) compileStructInitializationExpression(
	ctx *CompilerCtx,
	val *CompilerResult,
) (error, *CompilerResult) {
	explicitType, ok := vd.ExplicitType.(SymbolType)
	if !ok {
		return fmt.Errorf("explicit type is not a symbol %v", vd.ExplicitType), nil
	}

	err, typeDef := ctx.FindStructSymbol(explicitType.Name)
	if err != nil {
		return fmt.Errorf("Could not find typedef for %s in structs symbol table", explicitType.Name), nil
	}

	entry := &SymbolTableEntry{
		Value:        *val.Value,
		Ref:          typeDef,
		DeclaredType: vd.ExplicitType,
	}

	err = ctx.AddSymbol(vd.Name, entry)
	if err != nil {
		return err, nil
	}

	return nil, nil
}

func (vd VarDeclarationStatement) compileStringExpression(
	ctx *CompilerCtx,
	val *CompilerResult,
) (error, *CompilerResult) {
	glob := llvm.AddGlobal(*ctx.Module, val.Value.Type(), fmt.Sprintf("global.%s", vd.Name))
	glob.SetInitializer(*val.Value)

	alloc := ctx.Builder.CreateAlloca(llvm.PointerType(llvm.GlobalContext().Int8Type(), 0), "")

	ctx.Builder.CreateStore(glob, alloc)

	entry := &SymbolTableEntry{
		Value:        *val.Value,
		Address:      &alloc,
		DeclaredType: vd.ExplicitType,
	}

	err := ctx.AddSymbol(vd.Name, entry)
	if err != nil {
		return err, nil
	}
	return nil, nil
}

func (vd VarDeclarationStatement) compileArrArrayInitializationExpression(
	ctx *CompilerCtx,
	val *CompilerResult,
) (error, *CompilerResult) {
	entry := &SymbolTableEntry{
		Value:        *val.Value,
		Address:      val.Value,
		DeclaredType: vd.ExplicitType,
	}

	err := ctx.AddSymbol(vd.Name, entry)
	if err != nil {
		return err, nil
	}

	err = ctx.AddArraySymbol(vd.Name, val.ArraySymbolTableEntry)
	if err != nil {
		return err, nil
	}

	return nil, nil
}
