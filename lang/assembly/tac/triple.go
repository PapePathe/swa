package tac

import (
	"fmt"
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

type CustomType struct {
	Types []ast.Type
}

type Triple struct {
	types     map[string]CustomType
	Insts     []Inst
	lastValue InstArg
	Main      *Proc
	Procs     []*Proc
	currproc  *Proc
}

var _ AsmOp = (*Triple)(nil)

func (gen *Triple) Gen(g AssemblyOpGenerator) error {
	return g.VisitTriple(gen)
}

func (gen *Triple) Display() {
	fmt.Println("header:")

	for i, v := range gen.Insts {
		var ArgOne string
		if v.ArgOne != nil {
			ArgOne = v.ArgOne.InstructionArg()
		}

		var ArgTwo string
		if v.ArgTwo != nil {
			ArgTwo = v.ArgTwo.InstructionArg()
		}

		fmt.Printf(" %5d %20s  %10s  %10s\n", i, v.Operation, ArgOne, ArgTwo)
	}

	fmt.Println()

	for _, proc := range gen.Procs {
		fmt.Printf("%s @%s( ", proc.Ret.String(), proc.Name)

		for _, arg := range proc.Args {
			fmt.Printf("%s:%s ", arg.Name, arg.ArgType.Value())
		}

		fmt.Println(")")

		for _, label := range proc.Labels {
			fmt.Printf("  %s:\n", label.Name)
			for i, v := range label.Insts {
				var ArgOne string
				if v.ArgOne != nil {
					ArgOne = v.ArgOne.InstructionArg()
				}

				var ArgTwo string
				if v.ArgTwo != nil {
					ArgTwo = v.ArgTwo.InstructionArg()
				}

				fmt.Printf(" %5d %20s  %10s  %10s\n", i, v.Operation, ArgOne, ArgTwo)
			}
		}

		fmt.Println()
	}

	fmt.Println("main:")
	for _, label := range gen.Main.Labels {
		fmt.Printf("  %s:\n", label.Name)

		for i, v := range label.Insts {
			var ArgOne string
			if v.ArgOne != nil {
				ArgOne = v.ArgOne.InstructionArg()
			}

			var ArgTwo string
			if v.ArgTwo != nil {
				ArgTwo = v.ArgTwo.InstructionArg()
			}

			fmt.Printf(" %5d %20s  %10s  %10s\n", i, v.Operation, ArgOne, ArgTwo)
		}
	}
}

func (gen *Triple) VisitByteType(node *ast.ByteType) error {
	panic("unimplemented")
}

func (gen *Triple) ZeroOfByteType(node *ast.ByteType) error {
	panic("unimplemented")
}

var _ ast.CodeGenerator = (*Triple)(nil)

func NewTripleGenerator() *Triple {
	return &Triple{
		Insts: []Inst{},
		types: map[string]CustomType{},
	}
}

func (gen Triple) getLastValue() InstArg {
	val := gen.lastValue
	gen.lastValue = nil

	return val
}

func (gen *Triple) setLastValue(a InstArg) {
	gen.lastValue = a
}

func appendOpToCurrentProc(gen *Triple, op AsmOp) {
	gen.currproc.AppendOp(op)
}

func appendToCurrentProc(gen *Triple, i Inst) {
	gen.currproc.Append(i)
}

func getLastGlobalInstID(gen *Triple) GlobalId {
	if len(gen.Insts) == 0 {
		return GlobalId{id: uint32(0)}
	}

	return GlobalId{id: uint32(len(gen.Insts) - 1)}
}

func getLastInstID(gen *Triple) InstID {
	if len(gen.currproc.currentLabel.Insts) == 0 {
		return InstID{id: uint32(0)}
	}

	return InstID{id: uint32(len(gen.currproc.currentLabel.Insts) - 1)}
}

func (gen *Triple) VisitArrayAccessExpression(node *ast.ArrayAccessExpression) error {
	panic("unimplemented")
}

func (gen *Triple) VisitArrayInitializationExpression(node *ast.ArrayInitializationExpression) error {
	panic("unimplemented")
}

func (gen *Triple) VisitArrayOfStructsAccessExpression(node *ast.ArrayOfStructsAccessExpression) error {
	panic("unimplemented")
}

func (gen *Triple) VisitArrayType(node *ast.ArrayType) error {
	panic("unimplemented")
}

func (gen *Triple) VisitAssignmentExpression(node *ast.AssignmentExpression) error {
	err := node.Value.Accept(gen)
	if err != nil {
		return err
	}

	incomingValue := gen.getLastValue()

	err = node.Assignee.Accept(gen)
	if err != nil {
		return err
	}

	assignee := gen.getLastValue()

	appendToCurrentProc(gen, Inst{
		Operation: OpWrite,
		ArgOne:    incomingValue,
		ArgTwo:    assignee,
	})

	return nil
}

func (gen *Triple) VisitBinaryExpression(node *ast.BinaryExpression) error {
	err := node.Left.Accept(gen)
	if err != nil {
		return err
	}

	left := gen.getLastValue()

	err = node.Right.Accept(gen)
	if err != nil {
		return err
	}

	right := gen.getLastValue()

	op, err := opMap(node.Operator.Kind)
	if err != nil {
		return err
	}

	appendToCurrentProc(gen, Inst{
		Operation: op,
		ArgOne:    left,
		ArgTwo:    right,
	})

	lastInstID := getLastInstID(gen)
	gen.setLastValue(lastInstID)

	return nil
}

func (gen *Triple) VisitBlockStatement(node *ast.BlockStatement) error {
	for _, stmt := range node.Body {
		err := stmt.Accept(gen)
		if err != nil {
			return err
		}
	}

	return nil
}

func (gen *Triple) VisitBoolType(node *ast.BoolType) error {
	panic("unimplemented")
}

func (gen *Triple) VisitBooleanExpression(node *ast.BooleanExpression) error {
	if node.Value {
		gen.setLastValue(BoolVal{value: "true"})

		return nil
	}

	gen.setLastValue(BoolVal{value: "false"})

	return nil
}

func (gen *Triple) VisitCallExpression(node *ast.CallExpression) error {
	panic("unimplemented")
}

func (gen *Triple) VisitConditionalStatement(node *ast.ConditionalStatetement) error {
	err := node.Condition.Accept(gen)
	if err != nil {
		return err
	}

	ifblock := gen.currproc.addLabel("if.block")
	elseblock := gen.currproc.addLabel("else.block")
	mergeblock := gen.currproc.addLabel("merge.block")

	appendToCurrentProc(gen, Inst{
		Operation: OpJumpCond,
		ArgOne:    getLastInstID(gen),
		ArgTwo: &JumpCond{
			Success: ifblock.Name,
			Failure: elseblock.Name,
		},
	})

	gen.currproc.setCurrentLabel(ifblock)

	err = node.Success.Accept(gen)
	if err != nil {
		return err
	}

	appendToCurrentProc(gen, Inst{
		Operation: OpJump,
		ArgOne:    mergeblock,
	})

	gen.currproc.setCurrentLabel(elseblock)

	if node.Failure.Body != nil {
		err := node.Failure.Accept(gen)
		if err != nil {
			return err
		}
	}

	appendToCurrentProc(gen, Inst{
		Operation: OpJump,
		ArgOne:    mergeblock,
	})

	gen.currproc.setCurrentLabel(mergeblock)

	return nil
}

func (gen *Triple) VisitErrorExpression(node *ast.ErrorExpression) error {
	panic("unimplemented")
}

func (gen *Triple) VisitErrorType(node *ast.ErrorType) error {
	panic("unimplemented")
}

func (gen *Triple) VisitExpressionStatement(node *ast.ExpressionStatement) error {
	return node.Exp.Accept(gen)
}

func (gen *Triple) VisitFloatExpression(node *ast.FloatExpression) error {
	res := Float64Val{value: node.Value}
	gen.setLastValue(&res)

	return nil
}

func (gen *Triple) VisitFloatType(node *ast.FloatType) error {
	panic("unimplemented")
}

func (gen *Triple) VisitFloatingBlockExpression(node *ast.FloatingBlockExpression) error {
	panic("unimplemented")
}

func (gen *Triple) VisitFunctionCall(node *ast.FunctionCallExpression) error {
	for _, arg := range node.Args {
		err := arg.Accept(gen)
		if err != nil {
			return err
		}

		appendToCurrentProc(gen, Inst{
			Operation: OpFunCallArg,
			ArgOne:    gen.getLastValue(),
		})
	}

	sym, _ := node.Name.(*ast.SymbolExpression)

	appendToCurrentProc(gen, Inst{
		Operation: OpFunCall,
		ArgOne:    &SymbolVal{value: sym.Value},
	})

	gen.setLastValue(getLastInstID(gen))

	return nil
}

func (gen *Triple) VisitFunctionDefinition(node *ast.FuncDeclStatement) error {
	lab := &Label{Name: "default", Insts: []Inst{}}
	proc := &Proc{
		Name:         node.Name,
		Ret:          node.ReturnType,
		Args:         node.Args,
		Labels:       []*Label{lab},
		currentLabel: lab,
		lmap:         map[string]int{"default": 0},
	}
	gen.Procs = append(gen.Procs, proc)
	gen.currproc = proc

	err := node.Body.Accept(gen)
	if err != nil {
		return err
	}

	gen.currproc = nil

	return nil
}

func (gen *Triple) VisitMainStatement(node *ast.MainStatement) error {
	lab := &Label{Name: "default", Insts: []Inst{}}
	proc := &Proc{
		Name:         "main",
		Labels:       []*Label{lab},
		currentLabel: lab,
		lmap:         map[string]int{"default": 0},
	}
	gen.Main = proc

	gen.currproc = gen.Main

	err := node.Body.Accept(gen)
	if err != nil {
		return err
	}

	gen.currproc = nil

	return nil
}

func (gen *Triple) VisitMemberExpression(node *ast.MemberExpression) error {
	panic("unimplemented")
}

func (gen *Triple) VisitNumber64Type(node *ast.Number64Type) error { return nil }

func (gen *Triple) VisitNumberExpression(node *ast.NumberExpression) error {
	gen.setLastValue(&Number32Val{value: int(node.Value)})

	return nil
}

func (gen *Triple) VisitNumberType(node *ast.NumberType) error   { return nil }
func (gen *Triple) VisitPointerType(node *ast.PointerType) error { return nil }

func (gen *Triple) VisitPrefixExpression(node *ast.PrefixExpression) error {
	err := node.RightExpression.Accept(gen)
	if err != nil {
		return err
	}

	res := gen.getLastValue()

	switch node.Operator.Kind {
	case lexer.Minus:
		appendToCurrentProc(gen, Inst{
			Operation: OpMinus,
			ArgOne:    res,
		})

		gen.setLastValue(getLastInstID(gen))

	case lexer.Not:
		appendToCurrentProc(gen, Inst{
			Operation: OpNegation,
			ArgOne:    res,
		})

		gen.setLastValue(getLastInstID(gen))
	default:
		return fmt.Errorf("Unsupported prefix %s", node.Operator.Kind)
	}

	return nil
}

func (gen *Triple) VisitPrintStatement(node *ast.PrintStatetement) error {
	for _, expr := range node.Values {
		err := expr.Accept(gen)
		if err != nil {
			return err
		}

		val := gen.getLastValue()

		appendToCurrentProc(gen, Inst{
			Operation: OpFunCallArg,
			ArgOne:    val,
		})
	}

	appendToCurrentProc(gen, Inst{
		Operation: OpFunCall,
		ArgOne:    &SymbolVal{value: "@printf"},
	})

	return nil
}

func (gen *Triple) VisitReturnStatement(node *ast.ReturnStatement) error {
	err := node.Value.Accept(gen)
	if err != nil {
		return err
	}

	last := gen.getLastValue()

	//	appendToCurrentProc(gen, Inst{
	//		Operation: OpReturn,
	//		ArgOne:    last,
	//	})

	appendOpToCurrentProc(gen, &Ret{
		Val: last,
	})

	return nil
}

func (gen *Triple) VisitStringExpression(node *ast.StringExpression) error {
	gen.Insts = append(gen.Insts, Inst{
		Operation: OpGlobal,
		ArgOne: &GlobalString{
			Value:  node.Value,
			Length: len(node.Value),
		},
	})

	id := getLastGlobalInstID(gen)

	gen.setLastValue(id)

	return nil
}

func (gen *Triple) VisitStringType(node *ast.StringType) error { return nil }

func (gen *Triple) VisitStructDeclaration(node *ast.StructDeclarationStatement) error {
	gen.types[node.Name] = CustomType{Types: node.Types}

	return nil
}

func (gen *Triple) VisitStructInitializationExpression(node *ast.StructInitializationExpression) error {
	panic("unimplemented")
}

func (gen *Triple) VisitSymbolAdressExpression(node *ast.SymbolAdressExpression) error {
	panic("unimplemented")
}

func (gen *Triple) VisitSymbolExpression(node *ast.SymbolExpression) error {
	gen.setLastValue(&SymbolVal{value: node.Value})

	return nil
}

func (gen *Triple) VisitSymbolType(node *ast.SymbolType) error { return nil }

func (gen *Triple) VisitSymbolValueExpression(node *ast.SymbolValueExpression) error {
	panic("unimplemented")
}

func (gen *Triple) VisitTupleAssignmentExpression(node *ast.TupleAssignmentExpression) error {
	panic("unimplemented")
}

func (gen *Triple) VisitTupleExpression(node *ast.TupleExpression) error {
	// gen.setLastValue(node)

	return nil
}

func (gen *Triple) VisitTupleType(node *ast.TupleType) error {
	panic("unimplemented")
}

func (gen *Triple) VisitVarDeclaration(node *ast.VarDeclarationStatement) error {
	if node.Value == nil {
		node.Value = &ast.ZeroExpression{
			T: node.ExplicitType,
		}
	}

	err := node.Value.Accept(gen)
	if err != nil {
		return err
	}

	value := gen.getLastValue()

	appendToCurrentProc(gen, Inst{
		Operation: OpAlloc,
		ArgOne:    &TypeID{T: node.ExplicitType},
		ArgTwo:    &SymbolVal{value: node.Name},
	})

	allocInstId := getLastInstID(gen)

	appendToCurrentProc(gen, Inst{
		Operation: OpWrite,
		ArgOne:    value,
		ArgTwo:    allocInstId,
	})

	return nil
}

func (gen *Triple) VisitVoidType(node *ast.VoidType) error { return nil }

func (gen *Triple) VisitWhileStatement(node *ast.WhileStatement) error {
	panic("unimplemented")
}

func (gen *Triple) VisitZeroExpression(node *ast.ZeroExpression) error {
	err := node.T.AcceptZero(gen)
	if err != nil {
		return err
	}

	return nil
}

func (gen *Triple) ZeroOfArrayType(node *ast.ArrayType) error {
	panic("unimplemented")
}

func (gen *Triple) ZeroOfBoolType(node *ast.BoolType) error {
	gen.setLastValue(BoolVal{value: "false"})

	return nil
}

func (gen *Triple) ZeroOfErrorType(node *ast.ErrorType) error {
	panic("unimplemented")
}

func (gen *Triple) ZeroOfFloatType(node *ast.FloatType) error {
	gen.setLastValue(&Float64Val{value: 0.0})

	return nil
}

func (gen *Triple) ZeroOfNumber64Type(node *ast.Number64Type) error {
	gen.setLastValue(&Number32Val{value: 0.0})

	return nil
}

func (gen *Triple) ZeroOfNumberType(node *ast.NumberType) error {
	gen.setLastValue(&Number32Val{value: 0})

	return nil
}

func (gen *Triple) ZeroOfPointerType(node *ast.PointerType) error {
	panic("unimplemented")
}

func (gen *Triple) ZeroOfStringType(node *ast.StringType) error {
	panic("unimplemented")
}

func (gen *Triple) ZeroOfSymbolType(node *ast.SymbolType) error {
	panic("unimplemented")
}

func (gen *Triple) ZeroOfTupleType(node *ast.TupleType) error {
	panic("unimplemented")
}

func (gen *Triple) ZeroOfVoidType(node *ast.VoidType) error {
	panic("unimplemented")
}
