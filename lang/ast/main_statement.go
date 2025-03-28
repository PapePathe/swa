package ast

import (
	"encoding/json"
	"swahili/lang/values"

	"github.com/llir/llvm/ir/types"
)

type MainStatement struct {
	Body BlockStatement
}

func (ms MainStatement) Evaluate(s *Scope) (error, values.Value) {
	return nil, nil
}

func (MainStatement) statement() {}
func (ms MainStatement) Compile(ctx *Context) error {
	main := ctx.mod.NewFunc("main", types.I32)
	mainCtx := ctx.NewContext(main.NewBlock(""))

	return ms.Body.Compile(mainCtx)
}

func (cs MainStatement) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["body"] = cs.Body

	res := make(map[string]any)
	res["ast.MainProGram"] = m

	return json.Marshal(res)
}
