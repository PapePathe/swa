package ast

import (
	"encoding/json"
	"swahili/lang/values"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

type MainStatement struct {
	Body BlockStatement
}

func (ms MainStatement) Evaluate(s *Scope) (error, values.Value) {
	return nil, nil
}

func (MainStatement) statement() {}
func (ms MainStatement) Compile(m *ir.Module, b *ir.Block) error {
	main := m.NewFunc("main", types.I32)
	entry := main.NewBlock("")

	for _, stmt := range ms.Body.Body {
		stmt.Compile(m, entry)
	}

	return nil
}

func (cs MainStatement) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["body"] = cs.Body

	res := make(map[string]any)
	res["ast.MainProGram"] = m

	return json.Marshal(res)
}
