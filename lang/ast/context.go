package ast

import (
	"fmt"
	"swahili/lang/lexer"

	"tinygo.org/x/go-llvm"
)

type CompilerResult struct {
	// Some expressions return results
	Value                  *llvm.Value
	StuctPropertyValueType *llvm.Type
	SymbolTableEntry       *SymbolTableEntry
	// Some expressions return a value that will be added to the structs symbol table
	StructSymbolTableEntry *StructSymbolTableEntry
	// Some expressions return a value that will be added to the arrays symbol table
	ArraySymbolTableEntry *ArraySymbolTableEntry
}
type StructSymbolTableEntry struct {
	LLVMType      llvm.Type
	PropertyTypes []llvm.Type
	Metadata      StructDeclarationStatement
}

type SymbolTableEntry struct {
	Value        llvm.Value
	Address      *llvm.Value
	Ref          *StructSymbolTableEntry
	DeclaredType Type
}

type ArraySymbolTableEntry struct {
	UnderlyingType    llvm.Type
	UnderlyingTypeDef *StructSymbolTableEntry
	Type              llvm.Type
	ElementsCount     int
}

type CompilerCtx struct {
	parent            *CompilerCtx
	Context           *llvm.Context
	Builder           *llvm.Builder
	Module            *llvm.Module
	Dialect           lexer.Dialect
	symbolTable       map[string]SymbolTableEntry
	structSymbolTable map[string]StructSymbolTableEntry
	arraysSymbolTable map[string]ArraySymbolTableEntry
	funcSymbolTable   map[string]llvm.Type
	Debugging         bool
}

func NewCompilerContext(
	c *llvm.Context,
	b *llvm.Builder,
	m *llvm.Module,
	d lexer.Dialect,
	p *CompilerCtx,
) *CompilerCtx {
	return &CompilerCtx{
		parent:            p,
		Context:           c,
		Builder:           b,
		Module:            m,
		Dialect:           d,
		Debugging:         true,
		symbolTable:       map[string]SymbolTableEntry{},
		structSymbolTable: map[string]StructSymbolTableEntry{},
		funcSymbolTable:   map[string]llvm.Type{},
		arraysSymbolTable: map[string]ArraySymbolTableEntry{},
	}
}

func (ctx CompilerCtx) AddStructSymbol(name string, value *StructSymbolTableEntry) error {
	if _, exists := ctx.structSymbolTable[name]; exists {
		return fmt.Errorf("struct named %s already exists in symbol table", name)
	}

	ctx.structSymbolTable[name] = *value

	return nil
}

func (ctx CompilerCtx) FindStructSymbol(name string) (error, *StructSymbolTableEntry) {
	entry, exists := ctx.structSymbolTable[name]

	if !exists {
		if ctx.parent == nil {
			return fmt.Errorf("struct named %s does not exist in symbol table", name), nil
		}

		return ctx.parent.FindStructSymbol(name)
	}

	return nil, &entry
}

func (ctx CompilerCtx) AddArraySymbol(name string, value *ArraySymbolTableEntry) error {
	if _, exists := ctx.arraysSymbolTable[name]; exists {
		return fmt.Errorf("array named %s already exists in symbol table", name)
	}

	ctx.arraysSymbolTable[name] = *value

	return nil
}

func (ctx CompilerCtx) FindArraySymbol(name string) (error, *ArraySymbolTableEntry) {
	entry, exists := ctx.arraysSymbolTable[name]

	if !exists {
		if ctx.parent == nil {
			return fmt.Errorf("array named %s does not exist in symbol table", name), nil
		}

		return ctx.parent.FindArraySymbol(name)
	}

	return nil, &entry
}

func (ctx CompilerCtx) AddFuncSymbol(name string, value *llvm.Type) error {
	if _, exists := ctx.funcSymbolTable[name]; exists {
		return fmt.Errorf("function named %s already exists in symbol table", name)
	}

	ctx.funcSymbolTable[name] = *value

	return nil
}

func (ctx CompilerCtx) SymbolExistsInCurrentScope(name string) bool {
	_, exists := ctx.symbolTable[name]

	return exists
}

func (ctx CompilerCtx) FindFuncSymbol(name string) (error, *llvm.Type) {
	entry, exists := ctx.funcSymbolTable[name]

	if !exists {
		if ctx.parent == nil {
			return fmt.Errorf("function named %s does not exist in symbol table", name), nil
		}

		return ctx.parent.FindFuncSymbol(name)
	}

	return nil, &entry
}

func (ctx CompilerCtx) AddSymbol(name string, value *SymbolTableEntry) error {
	if _, exists := ctx.symbolTable[name]; exists {
		return fmt.Errorf("variable named %s already exists in symbol table", name)
	}

	ctx.symbolTable[name] = *value

	return nil
}

func (ctx CompilerCtx) FindSymbol(name string) (error, *SymbolTableEntry) {
	entry, exists := ctx.symbolTable[name]

	if !exists {
		if ctx.parent == nil {
			return fmt.Errorf("variable named %s does not exist in symbol table", name), nil
		}

		return ctx.parent.FindSymbol(name)
	}

	return nil, &entry
}

func (ctx CompilerCtx) PrintVarNames() {
	for k := range ctx.symbolTable {
		fmt.Println("Variable name: ", k)
	}
}
