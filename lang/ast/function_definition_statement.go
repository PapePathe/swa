package ast

import (
	"encoding/json"
	"fmt"
	"swahili/lang/lexer"

	"tinygo.org/x/go-llvm"
)

type FuncArg struct {
	Name    string
	ArgType Type
}

type FuncDeclStatement struct {
	Body       BlockStatement
	Name       string
	ReturnType Type
	Args       []FuncArg
	Tokens     []lexer.Token
}

var _ Statement = (*FuncDeclStatement)(nil)

func (fd FuncDeclStatement) CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult) {
	newCtx := NewCompilerContext(
		ctx.Context,
		ctx.Builder,
		ctx.Module,
		ctx.Dialect,
		ctx,
	)

	err, params := fd.funcParams(newCtx)
	if err != nil {
		return err, nil
	}

	err, returnType := fd.extractType(newCtx, fd.ReturnType)
	if err != nil {
		return err, nil
	}

	newfuncType := llvm.FunctionType(returnType.typ, params, false)
	newFunc := llvm.AddFunction(*newCtx.Module, fd.Name, newfuncType)

	err = ctx.AddFuncSymbol(fd.Name, &newfuncType)
	if err != nil {
		return err, nil
	}

	for i, p := range newFunc.Params() {
		argType := fd.Args[i].ArgType
		name := fd.Args[i].Name
		p.SetName(name)

		entry := SymbolTableEntry{Value: p}

		err, eType := fd.extractType(newCtx, argType)
		if err != nil {
			return err, nil
		}

		if eType.sEntry != nil {
			entry.Ref = eType.sEntry
		}

		err = newCtx.AddSymbol(name, &entry)
		if err != nil {
			return fmt.Errorf("failed to add parameter %s to symbol table: %w", name, err), nil
		}

		if eType.aEntry != nil {
			err := newCtx.AddArraySymbol(name, eType.aEntry)
			if err != nil {
				return fmt.Errorf("failed to add parameter %s to arrays symbol table: %w", name, err), nil
			}
		}
	}

	block := ctx.Context.AddBasicBlock(newFunc, "body")
	ctx.Builder.SetInsertPointAtEnd(block)

	err, _ = fd.Body.CompileLLVM(newCtx)
	if err != nil {
		return err, nil
	}

	return nil, nil
}

func (expr FuncDeclStatement) TokenStream() []lexer.Token {
	return expr.Tokens
}

func (fd FuncDeclStatement) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["Args"] = fd.Args
	m["Body"] = fd.Body
	m["Name"] = fd.Name
	m["ReturnType"] = fd.ReturnType

	res := make(map[string]any)
	res["ast.FuncDeclStatement"] = m

	return json.Marshal(res)
}

func (fd FuncDeclStatement) extractType(ctx *CompilerCtx, t Type) (error, extractedType) {
	err, compiledType := t.LLVMType(ctx)
	if err != nil {
		return err, extractedType{}
	}

	switch typ := t.(type) {
	case NumberType, FloatType, StringType:
		return nil, extractedType{typ: compiledType}
	case SymbolType:
		err, entry := ctx.FindStructSymbol(typ.Name)
		if err != nil {
			return err, extractedType{typ: llvm.Type{}}
		}

		etyp := extractedType{
			// TODO: need to dinstinguish between passing a struct as value and as a pointer
			typ:    llvm.PointerType(compiledType, 0),
			sEntry: entry,
		}

		return nil, etyp
	case ArrayType:
		var sEntry *StructSymbolTableEntry

		switch undType := typ.Underlying.(type) {
		case SymbolType:
			err, entry := ctx.FindStructSymbol(undType.Name)
			if err != nil {
				return err, extractedType{}
			}

			sEntry = entry
		case PointerType:
		// TODO handle pointer type here
		default:
		}

		etype := extractedType{
			typ: llvm.PointerType(compiledType.ElementType(), 0),
			aEntry: &ArraySymbolTableEntry{
				UnderlyingType:    compiledType.ElementType(),
				UnderlyingTypeDef: sEntry,
				ElementsCount:     compiledType.ArrayLength(),
				Type:              llvm.ArrayType(compiledType.ElementType(), typ.Size),
			},
		}

		return nil, etype
	default:
		return fmt.Errorf("FuncDeclStatement argument type %v not supported", t), extractedType{}
	}
}

type extractedType struct {
	typ    llvm.Type
	entry  *SymbolTableEntry
	sEntry *StructSymbolTableEntry
	aEntry *ArraySymbolTableEntry
}

func (fd FuncDeclStatement) funcParams(ctx *CompilerCtx) (error, []llvm.Type) {
	params := []llvm.Type{}

	for _, arg := range fd.Args {
		err, typ := fd.extractType(ctx, arg.ArgType)
		if err != nil {
			return err, nil
		}

		params = append(params, typ.typ)
	}

	return nil, params
}
