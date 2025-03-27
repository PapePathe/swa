package ast

import (
	"encoding/json"
	"swahili/lang/values"

	"github.com/llir/llvm/ir"
)

type FuncArg struct {
	Name    string
	ArgType string
}

type FuncDeclStatement struct {
	Body       BlockStatement
	Name       string
	ReturnType string
	Args       []FuncArg
}

func (fd FuncDeclStatement) Evaluate(s *Scope) (error, values.Value) {
	return nil, nil
}

func (FuncDeclStatement) Compile(m *ir.Module, b *ir.Block) error {
	return nil
}

func (fd FuncDeclStatement) statement() {}

func (fd FuncDeclStatement) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["Args"] = fd.Args
	m["Body"] = fd.Body
	m["Name"] = fd.Name
	m["ReturnType"] = fd.ReturnType

	res := make(map[string]any)
	res["ast.FuncDeclStatement"] = m

	return json.Marshal(res)
}
