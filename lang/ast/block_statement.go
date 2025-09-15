package ast

import (
	"encoding/json"

	"tinygo.org/x/go-llvm"
)

// BlockStatement ...
type BlockStatement struct {
	// The body of the block statement
	Body []Statement
}

var _ Statement = (*BlockStatement)(nil)

func (bs BlockStatement) CompileLLVM(ctx *CompilerCtx) (error, *llvm.Value) {
	for _, stmt := range bs.Body {
		err, _ := stmt.CompileLLVM(ctx)
		if err != nil {
			return err, nil
		}
	}

	return nil, nil
}

func (bs BlockStatement) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["ast.BlockStatement"] = bs.Body

	return json.Marshal(m)
}
