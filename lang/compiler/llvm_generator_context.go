package compiler

import (
	"fmt"
	"os"
	"swahili/lang/ast"
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
	// The Swa type of this value
	SwaType ast.Type
}
type StructSymbolTableEntry struct {
	LLVMType      llvm.Type
	PropertyTypes []llvm.Type
	Metadata      ast.StructDeclarationStatement
	Embeds        map[string]StructSymbolTableEntry
	ArrayEmbeds   map[string]ArraySymbolTableEntry
}

type SymbolTableEntry struct {
	Value        llvm.Value
	Address      *llvm.Value
	Ref          *StructSymbolTableEntry
	DeclaredType ast.Type
	Global       bool
}

type ArraySymbolTableEntry struct {
	UnderlyingType    llvm.Type
	UnderlyingTypeDef *StructSymbolTableEntry
	Type              llvm.Type
	ElementsCount     int
}

type FuncDetails struct {
	meta   *ast.FuncDeclStatement
	lltype llvm.Type
}

type CompilerCtx struct {
	parent              *CompilerCtx
	Context             *llvm.Context
	Builder             *llvm.Builder
	Module              *llvm.Module
	Dialect             lexer.Dialect
	symbolTable         map[string]SymbolTableEntry
	structSymbolTable   map[string]StructSymbolTableEntry
	arraysSymbolTable   map[string]ArraySymbolTableEntry
	funcSymbolTable     map[string]FuncDetails
	Debugging           bool
	InsideFunction      bool
	MainFuncOccurrences int
}

func NewCompilerContext(
	c *llvm.Context,
	b *llvm.Builder,
	m *llvm.Module,
	d lexer.Dialect,
	p *CompilerCtx,
) *CompilerCtx {
	var debugging bool

	dbg := os.Getenv("SWA_DEBUG")
	if dbg == "yes" {
		debugging = true
	}

	ctx := &CompilerCtx{
		parent:            p,
		Context:           c,
		Builder:           b,
		Module:            m,
		Dialect:           d,
		Debugging:         debugging,
		symbolTable:       map[string]SymbolTableEntry{},
		structSymbolTable: map[string]StructSymbolTableEntry{},
		funcSymbolTable:   map[string]FuncDetails{},
		arraysSymbolTable: map[string]ArraySymbolTableEntry{},
	}

	if p != nil {
		ctx.InsideFunction = p.InsideFunction
	}

	return ctx
}

func (ctx *CompilerCtx) IncrementMainOccurrences() {
	if ctx.parent != nil {
		ctx.parent.IncrementMainOccurrences()

		return
	}

	ctx.MainFuncOccurrences += 1
}

func (ctx CompilerCtx) UpdateStructSymbol(name string, value *StructSymbolTableEntry) error {
	if ctx.parent != nil {
		return ctx.parent.UpdateStructSymbol(name, value)
	}

	if _, exists := ctx.structSymbolTable[name]; !exists {
		key := "CompilerCtx.UpdateStructSymbol.StructDoesNotExist"

		return ctx.Dialect.Error(key, name)
	}

	ctx.structSymbolTable[name] = *value

	return nil
}

func (ctx CompilerCtx) AddStructSymbol(name string, value *StructSymbolTableEntry) error {
	if ctx.parent != nil {
		return ctx.parent.AddStructSymbol(name, value)
	}

	if _, exists := ctx.structSymbolTable[name]; exists {
		key := "CompilerCtx.AddStructSymbol.StructAlreadyExist"

		return ctx.Dialect.Error(key, name)
	}

	ctx.structSymbolTable[name] = *value

	return nil
}

func (ctx CompilerCtx) FindStructSymbol(name string) (error, *StructSymbolTableEntry) {
	entry, exists := ctx.structSymbolTable[name]

	if !exists {
		if ctx.parent == nil {
			key := "CompilerCtx.FindStructSymbol.StructDoesNotExist"

			return ctx.Dialect.Error(key, name), nil
		}

		return ctx.parent.FindStructSymbol(name)
	}

	return nil, &entry
}

func (ctx CompilerCtx) AddArraySymbol(name string, value *ArraySymbolTableEntry) error {
	if _, exists := ctx.arraysSymbolTable[name]; exists {
		key := "CompilerCtx.AddArraySymbol.AlreadyExisits"
		return ctx.Dialect.Error(key, name)
	}

	ctx.arraysSymbolTable[name] = *value

	return nil
}

func (ctx CompilerCtx) FindArraySymbol(name string) (error, *ArraySymbolTableEntry) {
	entry, exists := ctx.arraysSymbolTable[name]

	if !exists {
		if ctx.parent == nil {
			key := "CompilerCtx.FindArraySymbol.DoesNotExist"

			return ctx.Dialect.Error(key, name), nil
		}

		return ctx.parent.FindArraySymbol(name)
	}

	return nil, &entry
}

func (ctx CompilerCtx) AddFuncSymbol(name string, value *llvm.Type, meta *ast.FuncDeclStatement) error {
	if ctx.parent != nil {
		return ctx.parent.AddFuncSymbol(name, value, meta)
	}

	if _, exists := ctx.funcSymbolTable[name]; exists {
		key := "CompilerCtx.AddFuncSymbol.AlreadyExisits"

		return ctx.Dialect.Error(key, name)
	}

	ctx.funcSymbolTable[name] = FuncDetails{lltype: *value, meta: meta}

	return nil
}

func (ctx CompilerCtx) SymbolExistsInCurrentScope(name string) bool {
	_, exists := ctx.symbolTable[name]

	return exists
}

func (ctx CompilerCtx) FindFuncSymbol(name string) (error, *FuncDetails) {
	entry, exists := ctx.funcSymbolTable[name]

	if !exists {
		if ctx.parent == nil {
			key := "CompilerCtx.FindFuncSymbol.DoesNotExist"

			return ctx.Dialect.Error(key, name), nil
		}

		return ctx.parent.FindFuncSymbol(name)
	}

	return nil, &entry
}

func (ctx CompilerCtx) AddSymbol(name string, value *SymbolTableEntry) error {
	if _, exists := ctx.symbolTable[name]; exists {
		key := "CompilerCtx.AddSymbol.AlreadyExisits"

		return ctx.Dialect.Error(key, name)
	}

	ctx.symbolTable[name] = *value

	return nil
}

func (ctx CompilerCtx) FindSymbol(name string) (error, *SymbolTableEntry) {
	entry, exists := ctx.symbolTable[name]

	if !exists {
		if ctx.parent == nil {
			key := "CompilerCtx.FindSymbol.DoesNotExist"

			return ctx.Dialect.Error(key, name), nil
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
