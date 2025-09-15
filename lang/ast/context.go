package ast

import (
	"fmt"

	"tinygo.org/x/go-llvm"
)

type StructSymbolTableEntry struct {
	LLVMType llvm.Type
	Metadata StructDeclarationStatement
}

type SymbolTableEntry struct {
	Value llvm.Value
	Ref   *StructSymbolTableEntry
}

type CompilerCtx struct {
	Context           *llvm.Context
	Builder           *llvm.Builder
	Module            *llvm.Module
	SymbolTable       map[string]SymbolTableEntry
	StructSymbolTable map[string]StructSymbolTableEntry
	FuncSymbolTable   map[string]llvm.Type
}

func (ctx CompilerCtx) PrintVarNames() {
	for k := range ctx.SymbolTable {
		fmt.Println("Variable name: ", k)
	}
}
