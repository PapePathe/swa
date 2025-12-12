package ast

import (
	"encoding/json"
	"fmt"
	"swahili/lang/lexer"

	"tinygo.org/x/go-llvm"
)

type StructDeclarationStatement struct {
	Name       string
	Properties []string
	Types      []Type
	Tokens     []lexer.Token
}

var _ Statement = (*StructDeclarationStatement)(nil)

func (sd StructDeclarationStatement) PropertyIndex(name string) (error, int) {
	for propIndex, propName := range sd.Properties {
		if propName == name {
			return nil, propIndex
		}
	}

	return fmt.Errorf("Property with name (%s) does not exist on struct %s", name, sd.Name), 0
}

func (sd StructDeclarationStatement) CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult) {
	attrs := []llvm.Type{}

	for idx := range sd.Properties {
		typ := sd.Types[idx]
		switch typ := typ.(type) {
		case SymbolType:
			err, t := typ.LLVMTypeDyn(ctx)
			if err != nil {
				return err, nil
			}
			attrs = append(attrs, t)
		default:
			attrs = append(attrs, typ.LLVMType())
		}
	}

	newtype := ctx.Context.StructCreateNamed(sd.Name)
	newtype.StructSetBody(attrs, false)

	entry := &StructSymbolTableEntry{
		LLVMType:      newtype,
		Metadata:      sd,
		PropertyTypes: attrs,
	}
	err := ctx.AddStructSymbol(sd.Name, entry)

	return err, nil
}

func (expr StructDeclarationStatement) TokenStream() []lexer.Token {
	return expr.Tokens
}

func (expr StructDeclarationStatement) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["Name"] = expr.Name
	m["Properties"] = expr.Properties
	m["Types"] = expr.Types

	res := make(map[string]any)
	res["ast.StructDeclarationStatement"] = m

	return json.Marshal(res)
}
