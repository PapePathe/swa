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

	err, val := vd.Value.CompileLLVM(ctx)
	if err != nil {
		return err, nil
	}

	if val == nil {
		err := fmt.Errorf("VarDeclarationStatement: return value is nil <%s> <%s>", vd.Name, vd.Value)

		return err, nil
	}

	if err := vd.TypeCheck(vd.ExplicitType.Value(), val.Value.Type()); err != nil {
		return err, nil
	}

	switch vd.Value.(type) {
	case StructInitializationExpression:
		explicitType, ok := vd.ExplicitType.(SymbolType)
		if !ok {
			return fmt.Errorf("explicit type is not a symbol %v", vd.ExplicitType), nil
		}

		err, typeDef := ctx.FindStructSymbol(explicitType.Name)
		if err != nil {
			return fmt.Errorf("Could not find typedef for %s in structs symbol table", explicitType.Name), nil
		}

		ctx.AddSymbol(vd.Name, &SymbolTableEntry{Value: *val.Value, Ref: typeDef})
	case StringExpression:
		glob := llvm.AddGlobal(*ctx.Module, val.Value.Type(), fmt.Sprintf("global.%s", vd.Name))
		glob.SetInitializer(*val.Value)
		alloc := ctx.Builder.CreateAlloca(llvm.PointerType(llvm.GlobalContext().Int8Type(), 0), "")
		ctx.Builder.CreateStore(glob, alloc)
		ctx.AddSymbol(vd.Name, &SymbolTableEntry{Value: *val.Value, Address: &alloc})
	case NumberExpression, FloatExpression, BinaryExpression, FunctionCallExpression:
		alloc := ctx.Builder.CreateAlloca(val.Value.Type(), fmt.Sprintf("alloc.%s", vd.Name))
		ctx.Builder.CreateStore(*val.Value, alloc)
		ctx.AddSymbol(vd.Name, &SymbolTableEntry{Value: *val.Value, Address: &alloc})
	case ArrayInitializationExpression:
		ctx.AddSymbol(vd.Name, &SymbolTableEntry{Value: *val.Value, Address: val.Value})
		ctx.AddArraySymbol(vd.Name, val.ArraySymbolTableEntry)
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
				//				return fmt.Errorf("expected %s got pointer of unknown value %v", t, k.IsNil())
			}
		case DataTypeNumber:
			switch k.ElementType().TypeKind() {
			case llvm.IntegerTypeKind:
				// we good
			default:
				return fmt.Errorf("expected %s got pointer of %v", t, k)
			}
		case DataTypeArray:
		case DataTypeSymbol:
		default:
			panic(fmt.Sprintf("unsupported element %s", t))
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
