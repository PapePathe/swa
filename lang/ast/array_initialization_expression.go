package ast

import (
	"encoding/json"
	"fmt"

	"tinygo.org/x/go-llvm"
)

type ArrayInitializationExpression struct {
	Underlying Type
	Contents   []Expression
}

var _ Expression = (*ArrayInitializationExpression)(nil)

func (expr ArrayInitializationExpression) extractArrayType(ctx *CompilerCtx) llvm.Type {
	arrayType, ok := expr.Underlying.(ArrayType)

	if !ok {
		panic(fmt.Errorf("Type (%s) cannot be casted to array type", expr.Underlying))
	}

	switch arrayType.Underlying.Value() {
	case DataTypeNumber:
		return ctx.Context.Int32Type()
	case DataTypeSymbol:
		sym, _ := arrayType.Underlying.(SymbolType)

		sdef, ok := ctx.StructSymbolTable[sym.Name]
		if !ok {
			panic(fmt.Errorf("Type (%s) is not a valid struct", expr.Underlying))
		}
		return sdef.LLVMType
	case DataTypeString:
		return llvm.PointerType(ctx.Context.Int32Type(), 0)
	default:
		panic(fmt.Errorf("Type (%v) not implemented in array expression", arrayType.Underlying.Value()))
	}
}

func (expr ArrayInitializationExpression) CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult) {
	if len(expr.Contents) == 0 {
		// FIX: error messages should be translated
		return fmt.Errorf("Static arrays must be initialized"), nil
	}

	innerType := expr.extractArrayType(ctx)
	arrayType := llvm.ArrayType(innerType, len(expr.Contents))
	arrayPtr := ctx.Builder.CreateAlloca(arrayType, "")
	for i, value := range expr.Contents {
		err, content := value.CompileLLVM(ctx)
		if err != nil {
			return err, nil
		}
		gep := ctx.Builder.CreateGEP(
			innerType,
			arrayPtr,
			[]llvm.Value{llvm.ConstInt(llvm.GlobalContext().Int32Type(), uint64(i), false)}, "",
		)
		ctx.Builder.CreateStore(*content.Value, gep)
	}

	res := CompilerResult{
		Value: &arrayPtr,
		ArraySymbolTableEntry: &ArraySymbolTableEntry{
			ElementsCount:  arrayType.ArrayLength(),
			UnderlyingType: arrayType.ElementType(),
			Type:           arrayType,
		},
	}
	return nil, &res
}

func (cs ArrayInitializationExpression) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["UnderlyingType"] = cs.Underlying
	m["Contents"] = cs.Contents

	res := make(map[string]any)
	res["ast.ArrayInitializationExpression"] = m

	return json.Marshal(res)
}
