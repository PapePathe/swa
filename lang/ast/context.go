package ast

import (
	"fmt"
	"swahili/lang/lexer"

	"tinygo.org/x/go-llvm"
)

type StructSymbolTableEntry struct {
	LLVMType      llvm.Type
	PropertyTypes []llvm.Type
	Metadata      StructDeclarationStatement
}

type SymbolTableEntry struct {
	Value llvm.Value
	Ref   *StructSymbolTableEntry
}

type ArraySymbolTableEntry struct {
	UnderlyingType llvm.Type
	Type           llvm.Type
	ElementsCount  int
}

type CompilerCtx struct {
	Context           *llvm.Context
	Builder           *llvm.Builder
	Module            *llvm.Module
	Dialect           lexer.Dialect
	SymbolTable       map[string]SymbolTableEntry
	StructSymbolTable map[string]StructSymbolTableEntry
	ArraysSymbolTable map[string]ArraySymbolTableEntry
	FuncSymbolTable   map[string]llvm.Type
}

type CompilerResult struct {
	// Some expressions return results
	Value *llvm.Value

	// Some expressions return a value that will be added to the structs symbol table
	StructSymbolTableEntry *StructSymbolTableEntry

	// Some expressions return a value that will be added to the arrays symbol table
	ArraySymbolTableEntry *ArraySymbolTableEntry
}

func (ctx CompilerCtx) PrintVarNames() {
	for k := range ctx.SymbolTable {
		fmt.Println("Variable name: ", k)
	}
}
