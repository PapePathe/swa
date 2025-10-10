package ast

import (
	"encoding/json"
	"fmt"

	"tinygo.org/x/go-llvm"

	"swahili/lang/errmsg"
)

type ArrayAccessExpression struct {
	Name  Expression
	Index Expression
}

var _ Expression = (*ArrayAccessExpression)(nil)

func (expr ArrayAccessExpression) CompileLLVM(ctx *CompilerCtx) (error, *llvm.Value) {
	varName, ok := expr.Name.(SymbolExpression)
	if !ok {
		err := errmsg.AstError{Message: fmt.Sprintf("%s.ArrayAccessExpression.NameNotASymbol", "en")}
		// FIX: error messages should be translated
		return err, nil
	}

	array, ok := ctx.SymbolTable[varName.Value]
	if !ok {
		// FIX: error messages should be translated
		return fmt.Errorf("Array Name is not a symbol"), nil
	}

	itemIndex, ok := expr.Index.(NumberExpression)
	if !ok {
		// FIX: error messages should be translated
		return fmt.Errorf("Only numbers are supported as array index, current: (%v)", expr.Index), nil
	}

	entry, ok := ctx.ArraysSymbolTable[varName.Value]
	if !ok {
		// FIX: error messages should be translated
		return fmt.Errorf("Array (%s) does not exist in symbol table", varName.Value), nil
	}

	if int(itemIndex.Value) > entry.ElementsCount-1 {
		return fmt.Errorf("Element at index (%d) does not exist in array (%s)", int(itemIndex.Value), varName.Value), nil
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
