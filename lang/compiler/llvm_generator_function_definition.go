package compiler

import (
	"swahili/lang/ast"

	"tinygo.org/x/go-llvm"
)

func (g *LLVMGenerator) VisitFunctionDefinition(node *ast.FuncDeclStatement) error {
	old := g.logger.Step("FunDefStmt")

	defer g.logger.Restore(old)

	g.Debugf("%s", node.Name)

	oldCtx := g.Ctx
	g.Ctx = NewCompilerContext(
		g.Ctx.Context,
		g.Ctx.Builder,
		g.Ctx.Module,
		g.Ctx.Dialect,
		g.Ctx,
	)
	g.Ctx.InsideFunction = true

	defer func() {
		g.Ctx.InsideFunction = false
		g.Ctx = oldCtx
	}()

	err, params := g.funcParams(g.Ctx, node)
	if err != nil {
		return err
	}

	err, returnType := g.extractType(g.Ctx, node.ReturnType)
	if err != nil {
		return err
	}

	typetoReturn := returnType.typ
	newfuncType := llvm.FunctionType(typetoReturn, params, node.ArgsVariadic)
	newFunc := llvm.AddFunction(*g.Ctx.Module, node.Name, newfuncType)

	err = oldCtx.AddFuncSymbol(node.Name, &newfuncType, node)
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

			var entry SymbolTableEntry

			switch argType.(type) {
			case ast.FloatType, ast.NumberType, ast.Number64Type,
				ast.SymbolType, ast.PointerType, ast.ArrayType:
				entry = SymbolTableEntry{
					Value:        p,
					DeclaredType: argType,
					Address:      &p,
				}
			default:
				alloca := g.Ctx.Builder.CreateAlloca(p.Type(), name)
				g.Ctx.Builder.CreateStore(p, alloca)

				entry = SymbolTableEntry{
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
					return err
				}
			}

			err = g.Ctx.AddSymbol(name, &entry)
			if err != nil {
				return err
			}
		}

		err := node.Body.Accept(g)
		if err != nil {
			return err
		}
	} else {
		newFunc.SetLinkage(llvm.ExternalLinkage)
	}

	return nil
}

type extractedType struct {
	typ    llvm.Type
	entry  *SymbolTableEntry
	sEntry *StructSymbolTableEntry
	aEntry *ArraySymbolTableEntry
}

func (g *LLVMGenerator) extractType(ctx *CompilerCtx, t ast.Type) (error, extractedType) {
	err := t.Accept(g)
	if err != nil {
		return err, extractedType{}
	}

	compiledType := g.getLastTypeVisitResult()

	switch typ := t.(type) {
	case ast.NumberType, ast.Number64Type, ast.FloatType,
		ast.StringType, ast.VoidType, *ast.ErrorType, *ast.TupleType:
		return nil, extractedType{typ: compiledType.Type}
	case ast.SymbolType:
		err, entry := ctx.FindStructSymbol(typ.Name)
		if err != nil {
			return err, extractedType{typ: llvm.Type{}}
		}

		etyp := extractedType{
			// TODO: need to dinstinguish between passing a struct as value and as a pointer
			typ:    llvm.PointerType(compiledType.Type, 0),
			sEntry: entry,
		}

		return nil, etyp
	case ast.PointerType:
		var sEntry *StructSymbolTableEntry

		switch undType := typ.Underlying.(type) {
		case ast.SymbolType:
			err, entry := ctx.FindStructSymbol(undType.Name)
			if err != nil {
				return err, extractedType{}
			}

			sEntry = entry
		default:
		}

		etype := extractedType{typ: compiledType.Type, sEntry: sEntry}

		return nil, etype
	case ast.ArrayType:
		var sEntry *StructSymbolTableEntry

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
			typ: llvm.PointerType(compiledType.SubType, 0),
			aEntry: &ArraySymbolTableEntry{
				UnderlyingType:    compiledType.SubType,
				UnderlyingTypeDef: sEntry,
				ElementsCount:     typ.Size,
				Type:              llvm.ArrayType(compiledType.SubType, typ.Size),
			},
		}

		return nil, etype
	default:
		key := "LLVMGenerator.VisitFunctionDefinition.UnsupportedArgumentType"

		return g.Ctx.Dialect.Error(key, t), extractedType{}
	}
}

func (g *LLVMGenerator) funcParams(ctx *CompilerCtx, node *ast.FuncDeclStatement) (error, []llvm.Type) {
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
