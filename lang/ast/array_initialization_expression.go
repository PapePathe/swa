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

func (expr ArrayInitializationExpression) CompileLLVM(ctx *CompilerCtx) (error, *llvm.Value) {
	if len(expr.Contents) == 0 {
		// FIX: error messages should be translated
		return fmt.Errorf("Static arrays must be initialized"), nil
	}

	var innerType llvm.Type
	arrayType, ok := expr.Underlying.(ArrayType)

	if !ok {
		// FIX: error messages should be translated
		return fmt.Errorf("Type (%s) cannot be casted to array type", expr.Underlying), nil
	}

	switch arrayType.Underlying.Value() {
	case DataTypeNumber:
		innerType = ctx.Context.Int32Type()
	case DataTypeSymbol:
		innerType = llvm.PointerType(ctx.Context.Int32Type(), 0)
	case DataTypeString:
		innerType = llvm.PointerType(ctx.Context.Int32Type(), 0)
	default:
		// FIX: error messages should be translated
		return fmt.Errorf("Type (%v) not implemented in array expression", arrayType.Underlying.Value()), nil
	}

	//	arrayType := llvm.ArrayType(innerType, len(expr.Contents))
	contents := []llvm.Value{}

	for _, value := range expr.Contents {
		err, content := value.CompileLLVM(ctx)
		if err != nil {
			return err, nil
		}

		switch value.(type) {
		case NumberExpression:
			contents = append(contents, *content)
		case StringExpression:
			glob := llvm.AddGlobal(*ctx.Module, content.Type(), "")
			glob.SetInitializer(*content)
			contents = append(contents, glob)
		}
	}
	arrayConst := llvm.ConstArray(innerType, contents)

	return nil, &arrayConst
}

func (cs ArrayInitializationExpression) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["UnderlyingType"] = cs.Underlying
	m["Contents"] = cs.Contents

	res := make(map[string]any)
	res["ast.ArrayInitializationExpression"] = m

	return json.Marshal(res)
}
