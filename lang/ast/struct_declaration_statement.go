package ast

import (
	"fmt"
	"os"

	"tinygo.org/x/go-llvm"
)

type StructDeclarationStatement struct {
	Name       string
	Properties []string
	Types      []Type
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

func (sd StructDeclarationStatement) CompileLLVM(ctx *CompilerCtx) (error, *llvm.Value) {
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
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}

	newtype := ctx.Context.StructCreateNamed(sd.Name)
	newtype.StructSetBody(attrs, false)

	ctx.StructSymbolTable[sd.Name] = StructSymbolTableEntry{
		LLVMType: newtype,
		Metadata: sd,
	}

	return nil, nil
}
