package ast

import (
	"fmt"

	"tinygo.org/x/go-llvm"
)

type StructProperty struct {
	PropType Type
}
type StructDeclarationStatement struct {
	Name       string
	Properties map[string]StructProperty
}

var _ Statement = (*StructDeclarationStatement)(nil)

func (sd StructDeclarationStatement) CompileLLVM(ctx *CompilerCtx) (error, *llvm.Value) {
	attrs := []llvm.Type{}

	for _, attr := range sd.Properties {
		switch v := attr.PropType.(type) {
		case SymbolType:
			switch v.Name {
			case "Chaine":
				attr := llvm.PointerType(ctx.Context.Int8Type(), 0)
				attrs = append(attrs, attr)
			case "Nombre":
				attr := ctx.Context.Int32Type()
				attrs = append(attrs, attr)
			default:
				err := fmt.Errorf("struct proprerty type (%s) not supported", v.Name)
				panic(err)
			}
		default:
			err := fmt.Errorf("struct proprerty does not support type (%v)", v)
			panic(err)
		}
	}

	newtype := ctx.Context.StructCreateNamed(sd.Name)
	newtype.StructSetBody(attrs, false)

	ctx.StructSymbolTable[sd.Name] = newtype

	return nil, nil
}
