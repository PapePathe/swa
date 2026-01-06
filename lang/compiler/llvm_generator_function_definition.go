package compiler

import (
	"fmt"
	"swahili/lang/ast"

	"tinygo.org/x/go-llvm"
)

func (g *LLVMGenerator) VisitFunctionDefinition(node *ast.FuncDeclStatement) error {
	oldCtx := g.Ctx
	g.Ctx = ast.NewCompilerContext(
		g.Ctx.Context,
		g.Ctx.Builder,
		g.Ctx.Module,
		g.Ctx.Dialect,
		g.Ctx,
	)

	err, params := g.funcParams(g.Ctx, node)
	if err != nil {
		return err
	}

	err, returnType := g.extractType(g.Ctx, node.ReturnType)
	if err != nil {
		return err
	}

	newfuncType := llvm.FunctionType(returnType.typ, params, node.ArgsVariadic)
	newFunc := llvm.AddFunction(*g.Ctx.Module, node.Name, newfuncType)

	err = oldCtx.AddFuncSymbol(node.Name, &newfuncType)
	if err != nil {
		return err
	}

	if len(node.Body.Body) > 0 {
		// Create entry block
		entryBlock := g.Ctx.Context.AddBasicBlock(newFunc, "entry")
		g.Ctx.Builder.SetInsertPointAtEnd(entryBlock)

		for i, p := range newFunc.Params() {
			argType := node.Args[i].ArgType
			name := node.Args[i].Name
			p.SetName(name)

			// Create alloca for the parameter to support address taking and array indexing
			var entry ast.SymbolTableEntry

			switch argType.(type) {
			case ast.ArrayType:
				param := p
				entry = ast.SymbolTableEntry{
					Value:        param,
					DeclaredType: argType,
					Address:      &param,
				}
			case ast.PointerType:
				param := p
				entry = ast.SymbolTableEntry{
					Value:        param,
					DeclaredType: argType,
					Address:      &param,
				}
			case ast.SymbolType:
				entry = ast.SymbolTableEntry{
					Value:        p,
					DeclaredType: argType,
					Address:      &p,
				}
			case ast.FloatType, ast.NumberType, ast.Number64Type:
				entry = ast.SymbolTableEntry{
					Value:        p,
					DeclaredType: argType,
					Address:      &p,
				}
			default:
				alloca := g.Ctx.Builder.CreateAlloca(p.Type(), name)
				g.Ctx.Builder.CreateStore(p, alloca)

				entry = ast.SymbolTableEntry{
					Value:        alloca,
					DeclaredType: argType,
					Address:      &alloca,
				}
			}

			err, eType := g.extractType(g.Ctx, argType)
			if err != nil {
				return err
			}

			if eType.sEntry != nil {
				entry.Ref = eType.sEntry
			}

			if eType.aEntry != nil {
				entry.Ref = eType.aEntry.UnderlyingTypeDef

				err := g.Ctx.AddArraySymbol(name, eType.aEntry)
				if err != nil {
					return fmt.Errorf("failed to add parameter %s to arrays symbol table: %w", name, err)
				}
			}

			err = g.Ctx.AddSymbol(name, &entry)
			if err != nil {
				return fmt.Errorf("failed to add parameter %s to symbol table: %w", name, err)
			}

		}

		err := node.Body.Accept(g)
		if err != nil {
			return err
		}
	} else {
		newFunc.SetLinkage(llvm.ExternalLinkage)
	}

	g.Ctx = oldCtx

	return nil
}

type extractedType struct {
	typ    llvm.Type
	entry  *ast.SymbolTableEntry
	sEntry *ast.StructSymbolTableEntry
	aEntry *ast.ArraySymbolTableEntry
}

func (g *LLVMGenerator) extractType(ctx *ast.CompilerCtx, t ast.Type) (error, extractedType) {
	err, compiledType := t.LLVMType(ctx)
	if err != nil {
		return err, extractedType{}
	}

	switch typ := t.(type) {
	case ast.NumberType, ast.Number64Type, ast.FloatType, ast.StringType, ast.VoidType:
		return nil, extractedType{typ: compiledType}
	case ast.SymbolType:
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
	case ast.PointerType:
		var sEntry *ast.StructSymbolTableEntry

		switch undType := typ.Underlying.(type) {
		case ast.SymbolType:
			err, entry := ctx.FindStructSymbol(undType.Name)
			if err != nil {
				return err, extractedType{}
			}

			sEntry = entry
		default:
		}

		etype := extractedType{typ: compiledType, sEntry: sEntry}

		return nil, etype
	case ast.ArrayType:
		var sEntry *ast.StructSymbolTableEntry

		switch undType := typ.Underlying.(type) {
		case ast.SymbolType:
			err, entry := ctx.FindStructSymbol(undType.Name)
			if err != nil {
				return err, extractedType{}
			}

			sEntry = entry
		default:
		}

		etype := extractedType{
			typ: llvm.PointerType(compiledType.ElementType(), 0),
			aEntry: &ast.ArraySymbolTableEntry{
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

func (g *LLVMGenerator) funcParams(ctx *ast.CompilerCtx, node *ast.FuncDeclStatement) (error, []llvm.Type) {
	params := []llvm.Type{}

	for _, arg := range node.Args {
		err, typ := g.extractType(ctx, arg.ArgType)
		if err != nil {
			return err, nil
		}

		// Pass arrays by reference (pointer)
		if _, ok := arg.ArgType.(ast.ArrayType); ok {
			ptrType := llvm.PointerType(typ.typ, 0)
			params = append(params, ptrType)
		} else {
			params = append(params, typ.typ)
		}
	}

	return nil, params
}

func (g *LLVMGenerator) setLastResult(res *ast.CompilerResult) {
	g.lastResult = res
}

func (g *LLVMGenerator) getLastResult() *ast.CompilerResult {
	res := g.lastResult
	g.lastResult = nil

	return res
}
