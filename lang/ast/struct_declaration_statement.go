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
				attrs = append(attrs, llvm.PointerType(llvm.GlobalContext().Int8Type(), 0))
			case "Nombre":
				attrs = append(attrs, llvm.GlobalContext().Int32Type())
			default:
				err := fmt.Errorf("struct proprerty type %s not supported", v.Name)
				panic(err)
			}
		default:
			err := fmt.Errorf("struct proprerty does not support type")
			panic(err)
		}
	}
	stype := ctx.Context.StructCreateNamed(sd.Name)
	stype.StructSetBody([]llvm.Type{}, false)

	return nil, nil
}
