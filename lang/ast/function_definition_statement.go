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
	switch t.Value() {
	case DataTypeNumber:
		return nil, extractedType{typ: llvm.GlobalContext().Int32Type()}
	case DataTypeFloat:
		return nil, extractedType{typ: llvm.GlobalContext().DoubleType()}
	case DataTypeIntType:
		return nil, extractedType{typ: llvm.GlobalContext().Int32Type()}
	case DataTypeString:
		return nil, extractedType{typ: llvm.PointerType(llvm.GlobalContext().Int8Type(), 0)}
	case DataTypeSymbol:
		sym, _ := t.(SymbolType)

		err, entry := ctx.FindStructSymbol(sym.Name)
		if err != nil {
			return err, extractedType{typ: llvm.Type{}}
		}

		return nil, extractedType{typ: llvm.PointerType(entry.LLVMType, 0), sEntry: entry}
	case DataTypeArray:
		var sEntry *StructSymbolTableEntry

		var innerType llvm.Type

		arr, _ := t.(ArrayType)

		switch arr.Underlying.Value() {
		case DataTypeSymbol:
			arrSym, _ := t.(ArrayType)
			sym, _ := arrSym.Underlying.(SymbolType)

			err, entry := ctx.FindStructSymbol(sym.Name)
			if err != nil {
				return err, extractedType{typ: llvm.Type{}}
			}

			innerType = entry.LLVMType
			sEntry = entry
		case DataTypeNumber:
			innerType = llvm.GlobalContext().Int32Type()
		case DataTypeFloat:
			innerType = llvm.GlobalContext().DoubleType()
		default:
			return fmt.Errorf("FuncDeclStatement Type %s not supported as array element", arr.Underlying), extractedType{typ: llvm.Type{}}
		}

		etype := extractedType{
			typ: llvm.PointerType(innerType, 0),
			aEntry: &ArraySymbolTableEntry{
				UnderlyingType:    innerType,
				UnderlyingTypeDef: sEntry,
				ElementsCount:     arr.Size,
				Type:              llvm.ArrayType(innerType, arr.Size),
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
