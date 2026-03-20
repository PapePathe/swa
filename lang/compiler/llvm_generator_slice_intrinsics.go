package compiler

import (
	"fmt"
	"swahili/lang/ast"

	"tinygo.org/x/go-llvm"
)

func (g *LLVMGenerator) VisitTypeExpression(node *ast.TypeExpression) error {
	return node.Type.Accept(g)
}

func (g *LLVMGenerator) handleMakeIntrinsic(node *ast.FunctionCallExpression) error {
	if len(node.Args) < 2 {
		return fmt.Errorf("make expects at least 2 arguments (type, length, [capacity])")
	}

	// 1. Get Type
	// The first argument is a TypeExpression
	typeExpr, ok := node.Args[0].(*ast.TypeExpression)
	if !ok {
		return fmt.Errorf("first argument to make must be a type")
	}

	err := typeExpr.Accept(g)
	if err != nil {
		return err
	}
	sliceTypeRes := g.getLastTypeVisitResult()
	sliceType := sliceTypeRes.Type
	elemType := sliceTypeRes.SubType

	// 2. Evaluate Length
	err = node.Args[1].Accept(g)
	if err != nil {
		return err
	}
	lenVal := *g.getLastResult().Value

	// 3. Evaluate Capacity
	var capVal llvm.Value
	if len(node.Args) > 2 {
		err = node.Args[2].Accept(g)
		if err != nil {
			return err
		}
		capVal = *g.getLastResult().Value
	} else {
		capVal = lenVal
	}

	// 4. Malloc
	// cap might be i32, malloc expects i64
	cap64 := g.Ctx.Builder.CreateZExt(capVal, g.Ctx.Context.Int64Type(), "cap64")
	ptr := g.Ctx.Builder.CreateArrayMalloc(elemType, cap64, "slice_data")

	// 5. Build Slice Struct { T*, i32, i32 }
	slice := llvm.ConstNull(sliceType)
	slice = g.Ctx.Builder.CreateInsertValue(slice, ptr, 0, "slice.ptr")
	slice = g.Ctx.Builder.CreateInsertValue(slice, lenVal, 1, "slice.len")
	slice = g.Ctx.Builder.CreateInsertValue(slice, capVal, 2, "slice.cap")

	g.setLastResult(&CompilerResult{
		Value:   &slice,
		SwaType: typeExpr.Type,
	})

	return nil
}

func (g *LLVMGenerator) handleLenIntrinsic(node *ast.FunctionCallExpression) error {
	if len(node.Args) != 1 {
		return fmt.Errorf("len expects 1 argument")
	}

	err := node.Args[0].Accept(g)
	if err != nil {
		return err
	}

	res := g.getLastResult()
	sliceType := res.SwaType

	if sliceType.Value() != ast.DataTypeSlice {
		return fmt.Errorf("len only supported for slices")
	}

	// Extract length field (index 1)
	lenVal := g.Ctx.Builder.CreateExtractValue(*res.Value, 1, "len")

	g.setLastResult(&CompilerResult{
		Value:   &lenVal,
		SwaType: &ast.NumberType{},
	})

	return nil
}

func (g *LLVMGenerator) handleCapIntrinsic(node *ast.FunctionCallExpression) error {
	if len(node.Args) != 1 {
		return fmt.Errorf("cap expects 1 argument")
	}

	err := node.Args[0].Accept(g)
	if err != nil {
		return err
	}

	res := g.getLastResult()
	sliceType := res.SwaType

	if sliceType.Value() != ast.DataTypeSlice {
		return fmt.Errorf("cap only supported for slices")
	}

	// Extract capacity field (index 2)
	capVal := g.Ctx.Builder.CreateExtractValue(*res.Value, 2, "cap")

	g.setLastResult(&CompilerResult{
		Value:   &capVal,
		SwaType: &ast.NumberType{},
	})

	return nil
}

func (g *LLVMGenerator) handleAppendIntrinsic(node *ast.FunctionCallExpression) error {
	if len(node.Args) != 2 {
		return fmt.Errorf("append expects 2 arguments (slice, item)")
	}

	// 1. Evaluate Slice
	err := node.Args[0].Accept(g)
	if err != nil {
		return err
	}
	sliceRes := g.getLastResult()
	slice := *sliceRes.Value
	var sliceSwaType ast.SliceType
	if st, ok := sliceRes.SwaType.(ast.SliceType); ok {
		sliceSwaType = st
	} else if st, ok := sliceRes.SwaType.(*ast.SliceType); ok {
		sliceSwaType = *st
	} else {
		return fmt.Errorf("append expects a slice, got %T", sliceRes.SwaType)
	}

	sliceLLVMType := slice.Type()
	if sliceLLVMType.TypeKind() != llvm.StructTypeKind {
		return fmt.Errorf("append: slice must be a struct, got %s", sliceLLVMType.String())
	}
	err = sliceSwaType.Underlying.Accept(g)
	if err != nil {
		return err
	}
	elemLLVMType := g.getLastTypeVisitResult().Type

	structFields := sliceLLVMType.StructElementTypes()
	ptrType := structFields[0]
	lenType := structFields[1]
	capType := structFields[2]

	// 2. Evaluate Item
	err = node.Args[1].Accept(g)
	if err != nil {
		return err
	}
	item := *g.getLastResult().Value

	// 3. Logic: if len == cap { grow }
	// Extract fields early to avoid issues with blocks and segfaults
	ptrOld := g.Ctx.Builder.CreateExtractValue(slice, 0, "ptr.old")
	lenOld := g.Ctx.Builder.CreateExtractValue(slice, 1, "len.old")
	capOld := g.Ctx.Builder.CreateExtractValue(slice, 2, "cap.old")

	full := g.Ctx.Builder.CreateICmp(llvm.IntEQ, lenOld, capOld, "is_full")

	// We'll use a branch to handle growing if full
	currentBlock := g.Ctx.Builder.GetInsertBlock()
	function := currentBlock.Parent()
	growBlock := g.Ctx.Context.AddBasicBlock(function, "slice.grow")
	mergeBlock := g.Ctx.Context.AddBasicBlock(function, "slice.append")

	g.Ctx.Builder.CreateCondBr(full, growBlock, mergeBlock)

	// --- Grow Block ---
	g.Ctx.Builder.SetInsertPointAtEnd(growBlock)

	// new_cap = cap == 0 ? 1 : cap * 2
	isZero := g.Ctx.Builder.CreateICmp(llvm.IntEQ, capOld, llvm.ConstInt(capOld.Type(), 0, false), "cap_is_zero")
	newCap := g.Ctx.Builder.CreateSelect(isZero,
		llvm.ConstInt(capOld.Type(), 1, false),
		g.Ctx.Builder.CreateMul(capOld, llvm.ConstInt(capOld.Type(), 2, false), "cap_x2"),
		"new_cap")

	// Realloc
	targetData := llvm.NewTargetData(g.Ctx.Module.DataLayout())
	sizeOfElem := targetData.TypeAllocSize(elemLLVMType)
	targetData.Dispose()
	sizeOfElemVal := llvm.ConstInt(g.Ctx.Context.Int64Type(), sizeOfElem, false)

	newCap64 := g.Ctx.Builder.CreateZExt(newCap, g.Ctx.Context.Int64Type(), "new_cap64")
	newSize := g.Ctx.Builder.CreateMul(newCap64, sizeOfElemVal, "realloc_size")

	oldPtr := g.Ctx.Builder.CreateExtractValue(slice, 0, "old_ptr")
	oldPtrI8 := g.Ctx.Builder.CreateBitCast(oldPtr, llvm.PointerType(g.Ctx.Context.Int8Type(), 0), "old_ptr_i8")

	reallocFunc := g.getReallocFunc()
	newPtrI8 := g.Ctx.Builder.CreateCall(reallocFunc.lltype, reallocFunc.llval, []llvm.Value{oldPtrI8, newSize}, "new_ptr_i8")
	newPtr := g.Ctx.Builder.CreateBitCast(newPtrI8, llvm.PointerType(elemLLVMType, 0), "new_ptr")

	sliceGrown := g.Ctx.Builder.CreateInsertValue(slice, newPtr, 0, "slice.new_ptr")
	sliceGrown = g.Ctx.Builder.CreateInsertValue(sliceGrown, newCap, 2, "slice.new_cap")

	// Values from grown slice
	ptrGrown := g.Ctx.Builder.CreateExtractValue(sliceGrown, 0, "ptr.grown")
	lenGrown := g.Ctx.Builder.CreateExtractValue(sliceGrown, 1, "len.grown")
	capGrown := g.Ctx.Builder.CreateExtractValue(sliceGrown, 2, "cap.grown")

	g.Ctx.Builder.CreateBr(mergeBlock)

	// --- Merge Block ---
	g.Ctx.Builder.SetInsertPointAtEnd(mergeBlock)

	ptrPhi := g.Ctx.Builder.CreatePHI(ptrType, "slice.ptr.phi")
	lenPhi := g.Ctx.Builder.CreatePHI(lenType, "slice.len.phi")
	capPhi := g.Ctx.Builder.CreatePHI(capType, "slice.cap.phi")

	ptrPhi.AddIncoming([]llvm.Value{ptrOld, ptrGrown}, []llvm.BasicBlock{currentBlock, growBlock})
	lenPhi.AddIncoming([]llvm.Value{lenOld, lenGrown}, []llvm.BasicBlock{currentBlock, growBlock})
	capPhi.AddIncoming([]llvm.Value{capOld, capGrown}, []llvm.BasicBlock{currentBlock, growBlock})

	// Store item
	itemPtr := g.Ctx.Builder.CreateGEP(elemLLVMType, ptrPhi, []llvm.Value{lenPhi}, "item_ptr")
	g.Ctx.Builder.CreateStore(item, itemPtr)

	// Increment len
	newLen := g.Ctx.Builder.CreateAdd(lenPhi, llvm.ConstInt(lenType, 1, false), "new_len")

	// Build final slice
	finalSlice := llvm.ConstNull(sliceLLVMType)
	finalSlice = g.Ctx.Builder.CreateInsertValue(finalSlice, ptrPhi, 0, "slice.final.ptr")
	finalSlice = g.Ctx.Builder.CreateInsertValue(finalSlice, newLen, 1, "slice.final.len")
	finalSlice = g.Ctx.Builder.CreateInsertValue(finalSlice, capPhi, 2, "slice.final.cap")

	g.setLastResult(&CompilerResult{
		Value:   &finalSlice,
		SwaType: sliceSwaType,
	})

	return nil
}

type ExternalFunc struct {
	lltype llvm.Type
	llval  llvm.Value
}

func (g *LLVMGenerator) getReallocFunc() ExternalFunc {
	name := "realloc"
	f := g.Ctx.Module.NamedFunction(name)
	i64 := g.Ctx.Context.Int64Type()
	ptrI8 := llvm.PointerType(g.Ctx.Context.Int8Type(), 0)

	ftype := llvm.FunctionType(ptrI8, []llvm.Type{ptrI8, i64}, false)
	if f.IsNil() {
		f = llvm.AddFunction(*g.Ctx.Module, name, ftype)
	}

	return ExternalFunc{lltype: ftype, llval: f}
}
