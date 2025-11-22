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
	err, params := fd.funcParams(ctx)
	if err != nil {
		return err, nil
	}

	err, returnType, _ := fd.extractType(ctx, fd.ReturnType)
	if err != nil {
		return err, nil
	}

	newfuncType := llvm.FunctionType(returnType, params, false)
	newFunc := llvm.AddFunction(*ctx.Module, fd.Name, newfuncType)
	if err := ctx.AddFuncSymbol(fd.Name, &newfuncType); err != nil {
		return err, nil
	}

	for i, p := range newFunc.Params() {
		argType := fd.Args[i].ArgType
		name := fd.Args[i].Name
		p.SetName(name)

		entry := SymbolTableEntry{Value: p}
		err, _, strucSymbolEntry := fd.extractType(ctx, argType)
		if err != nil {
			return err, nil
		}

		if strucSymbolEntry != nil {
			entry.Ref = strucSymbolEntry
		}

		if err := ctx.AddSymbol(name, &entry); err != nil {
			return fmt.Errorf("failed to add parameter %s to symbol table: %w", name, err), nil
		}
	}

	block := ctx.Context.AddBasicBlock(newFunc, "body")
	ctx.Builder.SetInsertPointAtEnd(block)

	if err, _ := fd.Body.CompileLLVM(ctx); err != nil {
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

func (fd FuncDeclStatement) extractType(ctx *CompilerCtx, t Type) (error, llvm.Type, *StructSymbolTableEntry) {
	switch t.Value() {
	case DataTypeNumber:
		return nil, llvm.GlobalContext().Int32Type(), nil
	case DataTypeFloat:
		return nil, llvm.GlobalContext().DoubleType(), nil
	case DataTypeIntType:
		return nil, llvm.GlobalContext().Int32Type(), nil
	case DataTypeString:
		return nil, llvm.PointerType(llvm.GlobalContext().Int8Type(), 0), nil
	case DataTypeSymbol:
		sym, _ := t.(SymbolType)
		err, entry := ctx.FindStructSymbol(sym.Name)
		if err != nil {
			return err, llvm.Type{}, nil
		}
		return nil, llvm.PointerType(entry.LLVMType, 0), entry
	default:
		return fmt.Errorf("argument type %v not supported", t), llvm.Type{}, nil
	}
}

func (fd FuncDeclStatement) funcParams(ctx *CompilerCtx) (error, []llvm.Type) {
	params := []llvm.Type{}

	for _, arg := range fd.Args {
		err, typ, _ := fd.extractType(ctx, arg.ArgType)
		if err != nil {
			return err, nil
		}

		params = append(params, typ)
	}

	return nil, params
}
