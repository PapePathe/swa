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
		format := "Type (%s) cannot be casted to array type"

		return nil, nil, fmt.Errorf(format, expr.Underlying)
	}

	switch arrayType.Underlying.Value() {
	case DataTypeFloat:
		typ := llvm.GlobalContext().DoubleType()

		return &typ, nil, nil
	case DataTypeNumber:
		typ := llvm.GlobalContext().Int32Type()

		return &typ, nil, nil
	case DataTypeSymbol:
		sym, _ := arrayType.Underlying.(SymbolType)

		err, sdef := ctx.FindStructSymbol(sym.Name)
		if err != nil {
			return nil, nil, fmt.Errorf("Type (%s) is not a valid struct", sym.Name)
		}

		return &sdef.LLVMType, sdef, nil
	case DataTypeString:
		typ := llvm.PointerType(llvm.GlobalContext().Int32Type(), 0)

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

	arrTyp, _ := expr.Underlying.(ArrayType)

	if len(expr.Contents) != arrTyp.Size {
		format := "Static arrays must be initialized with exact size: want: %d, has: %d"

		return fmt.Errorf(format, arrTyp.Size, len(expr.Contents)), nil
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
		case NumberExpression, StringExpression, FloatExpression:
			err, content := value.CompileLLVM(ctx)
			if err != nil {
				return err, nil
			}

			ctx.Builder.CreateStore(*content.Value, itemGep)
		case SymbolExpression:
			err, content := value.CompileLLVM(ctx)
			if err != nil {
				return err, nil
			}

			if content.SymbolTableEntry != nil && content.SymbolTableEntry.Ref != nil {
				load := ctx.Builder.CreateLoad(content.SymbolTableEntry.Ref.LLVMType, *content.Value, "")
				ctx.Builder.CreateStore(load, itemGep)
			} else {
				load := ctx.Builder.CreateLoad(content.Value.Type(), *content.Value, "")
				ctx.Builder.CreateStore(load, itemGep)
			}
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
			format := "ArrayInitializationExpression: Expression %s not supported"

			return fmt.Errorf(format, value), nil
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

func (expr ArrayInitializationExpression) Accept(g CodeGenerator) error {
	return g.VisitArrayInitializationExpression(&expr)
}

func (expr ArrayInitializationExpression) TokenStream() []lexer.Token {
	return expr.Tokens
}

func (cs ArrayInitializationExpression) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["UnderlyingType"] = cs.Underlying
	m["Contents"] = cs.Contents

	res := make(map[string]any)
	res["ast.ArrayInitializationExpression"] = m

	return json.Marshal(res)
}
