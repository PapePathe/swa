package ast

import (
	"encoding/json"
	"fmt"

	"tinygo.org/x/go-llvm"
)

type ArrayAccessExpression struct {
	Name  Expression
	Index Expression
}

var _ Expression = (*ArrayAccessExpression)(nil)

func (expr ArrayAccessExpression) CompileLLVM(ctx *CompilerCtx) (error, *llvm.Value) {
	varName, ok := expr.Name.(SymbolExpression)
	if !ok {
		// FIX: error messages should be translated
		return fmt.Errorf("Array Name is not a symbol"), nil
	}

	array, ok := ctx.SymbolTable[varName.Value]
	if !ok {
		// FIX: error messages should be translated
		return fmt.Errorf("Array Name is not a symbol"), nil
	}

	itemIndex, ok := expr.Index.(NumberExpression)
	if !ok {
		// FIX: error messages should be translated
		return fmt.Errorf("Value should be a number expression"), nil
	}

	itemPtr := ctx.Builder.CreateInBoundsGEP(
		ctx.Context.Int32Type(),
		array.Value,
		[]llvm.Value{
			llvm.ConstInt(llvm.GlobalContext().Int32Type(), uint64(itemIndex.Value), false),
		},
		"",
	)

	return nil, &itemPtr
}

func (cs ArrayAccessExpression) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["Name"] = cs.Name
	m["Index"] = cs.Index

	res := make(map[string]any)
	res["ast.ArrayAccessExpression"] = m

	return json.Marshal(res)
}

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
		panic(fmt.Sprintf("Type (%s) cannot be casted to array type", expr.Underlying))
	}

	switch arrayType.Underlying.(type) {
	case NumberType:
		innerType = ctx.Context.Int32Type()
	default:
		// FIX: error messages should be translated
		panic(fmt.Sprintf("Type (%s) not implemented in array expression", expr.Underlying))
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
