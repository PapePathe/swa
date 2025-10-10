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

	switch arrayType.Underlying.(type) {
	case NumberType:
		innerType = ctx.Context.Int32Type()
	default:
		// FIX: error messages should be translated
		return fmt.Errorf("Type (%s) not implemented in array expression", expr.Underlying), nil
	}

	//	arrayType := llvm.ArrayType(innerType, len(expr.Contents))
	contents := []llvm.Value{}

	for _, value := range expr.Contents {
		err, content := value.CompileLLVM(ctx)
		if err != nil {
			return err, nil
		}
		contents = append(contents, *content)
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
