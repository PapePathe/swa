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
		switch typ.(type) {
		case StringType:
			attrs = append(attrs, llvm.PointerType(ctx.Context.Int8Type(), 0))
		case NumberType:
			attrs = append(attrs, ctx.Context.Int32Type())
		default:
			err := fmt.Errorf("struct proprerty type (%s) not supported", typ)
			return err, nil
		}
	}

	newtype := ctx.Context.StructCreateNamed(sd.Name)
	newtype.StructSetBody(attrs, false)

	ctx.AddStructSymbol(
		sd.Name,
		&StructSymbolTableEntry{
			LLVMType:      newtype,
			Metadata:      sd,
			PropertyTypes: attrs,
		})

	return nil, nil
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
