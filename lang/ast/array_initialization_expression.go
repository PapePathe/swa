package ast

import (
	"encoding/json"
	"fmt"
	"swahili/lang/lexer"

	"tinygo.org/x/go-llvm"
)

type ArrayInitializationExpression struct {
	Underlying Type
	Contents   []Expression
	Tokens     []lexer.Token
}

var _ Expression = (*ArrayInitializationExpression)(nil)

func (expr ArrayInitializationExpression) extractArrayType(ctx *CompilerCtx) (*llvm.Type, *StructSymbolTableEntry, error) {
	arrayType, ok := expr.Underlying.(ArrayType)

	if !ok {
		panic(fmt.Errorf("Type (%s) cannot be casted to array type", expr.Underlying))
	}

	switch arrayType.Underlying.Value() {
	case DataTypeNumber:
		typ := ctx.Context.Int32Type()
		return &typ, nil, nil
	case DataTypeSymbol:
		sym, _ := arrayType.Underlying.(SymbolType)

		err, sdef := ctx.FindStructSymbol(sym.Name)
		if err != nil {
			return nil, nil, fmt.Errorf("Type (%s) is not a valid struct", expr.Underlying)
		}
		return &sdef.LLVMType, sdef, nil
	case DataTypeString:
		typ := llvm.PointerType(ctx.Context.Int32Type(), 0)
		return &typ, nil, nil
	default:
		return nil, nil, fmt.Errorf("Type (%v) not implemented in array expression", arrayType.Underlying.Value())
	}
}

func (expr ArrayInitializationExpression) CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult) {
	if len(expr.Contents) == 0 {
		// FIX: error messages should be translated
		return fmt.Errorf("Static arrays must be initialized"), nil
	}

	innerType, sdef, err := expr.extractArrayType(ctx)
	if err != nil {
		return err, nil
	}
	arrayType := llvm.ArrayType(*innerType, len(expr.Contents))
	arrayPtr := ctx.Builder.CreateAlloca(arrayType, "")
	for i, value := range expr.Contents {
		itemGep := ctx.Builder.CreateGEP(
			arrayType,
			arrayPtr,
			[]llvm.Value{
				llvm.ConstInt(llvm.GlobalContext().Int32Type(), uint64(0), false),
				llvm.ConstInt(llvm.GlobalContext().Int32Type(), uint64(i), false),
			}, "",
		)

		switch value.(type) {
		case NumberExpression, StringExpression:
			err, content := value.CompileLLVM(ctx)
			if err != nil {
				return err, nil
			}

			ctx.Builder.CreateStore(*content.Value, itemGep)
		case StructInitializationExpression:
			structExpr, _ := value.(StructInitializationExpression)
			err, structFields := structExpr.InitValues(ctx)
			if err != nil {
				return err, nil
			}

			for _, field := range structFields {
				gep := ctx.Builder.CreateGEP(
					*innerType,
					itemGep,
					[]llvm.Value{
						llvm.ConstInt(llvm.GlobalContext().Int32Type(), uint64(0), false),
						llvm.ConstInt(llvm.GlobalContext().Int32Type(), uint64(field.Position), false),
					}, "",
				)
				ctx.Builder.CreateStore(*field.Value, gep)
			}

		default:
			panic(fmt.Sprintf("ArrayInitializationExpression: Expression %s not supported", value))
		}
	}

	res := CompilerResult{
		Value: &arrayPtr,
		ArraySymbolTableEntry: &ArraySymbolTableEntry{
			ElementsCount:     arrayType.ArrayLength(),
			UnderlyingType:    arrayType.ElementType(),
			UnderlyingTypeDef: sdef,
			Type:              arrayType,
		},
	}
	return nil, &res
}

func (expr ArrayInitializationExpression) TokenStream() []lexer.Token {
	return expr.Tokens
}

func (cs ArrayInitializationExpression) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["UnderlyingType"] = cs.Underlying
	m["Contents"] = cs.Contents
	m["Tokens"] = cs.Tokens

	res := make(map[string]any)
	res["ast.ArrayInitializationExpression"] = m

	return json.Marshal(res)
}
