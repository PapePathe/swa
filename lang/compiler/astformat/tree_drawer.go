package astformat

import (
	"fmt"
	"io"
	"strings"
	"swahili/lang/ast"
)

type TreeDrawer struct {
	w      io.Writer
	isLast []bool
}

var _ ast.CodeGenerator = (*TreeDrawer)(nil)

func NewTreeDrawer(w io.Writer) *TreeDrawer {
	return &TreeDrawer{
		w: w,
	}
}

func Draw(node ast.Node) string {
	var sb strings.Builder

	drawer := NewTreeDrawer(&sb)

	_ = node.Accept(drawer)

	return sb.String()
}

func (t *TreeDrawer) print(format string, a ...interface{}) {
	fmt.Fprintf(t.w, format, a...)
}

func (t *TreeDrawer) writeLine(label string) {
	for i := 0; i < len(t.isLast)-1; i++ {
		if t.isLast[i] {
			t.print("    ")
		} else {
			t.print("│   ")
		}
	}
	if len(t.isLast) > 0 {
		if t.isLast[len(t.isLast)-1] {
			t.print("└── ")
		} else {
			t.print("├── ")
		}
	}
	t.print("%s\n", label)
}

func (t *TreeDrawer) visitChild(node ast.Node, last bool) error {
	if node == nil {
		return nil
	}
	t.isLast = append(t.isLast, last)
	err := node.Accept(t)
	t.isLast = t.isLast[:len(t.isLast)-1]
	return err
}

func (t *TreeDrawer) visitType(typ ast.Type, last bool) error {
	if typ == nil {
		return nil
	}
	t.isLast = append(t.isLast, last)
	err := typ.Accept(t)
	t.isLast = t.isLast[:len(t.isLast)-1]
	return err
}

func (t *TreeDrawer) VisitSymbolAdressExpression(node *ast.SymbolAdressExpression) error {
	t.writeLine("SymbolAdressExpression")

	err := t.visitChild(node.Exp, true)
	if err != nil {
		return err
	}

	return nil
}

func (t *TreeDrawer) VisitSymbolValueExpression(node *ast.SymbolValueExpression) error {
	t.writeLine("SymbolValueExpression")

	err := t.visitChild(node.Exp, true)
	if err != nil {
		return err
	}

	return nil
}

func (t *TreeDrawer) VisitBlockStatement(node *ast.BlockStatement) error {
	t.writeLine("BlockStatement")

	for i, stmt := range node.Body {
		if err := t.visitChild(stmt, i == len(node.Body)-1); err != nil {
			return err
		}
	}

	return nil
}

func (t *TreeDrawer) VisitVarDeclaration(node *ast.VarDeclarationStatement) error {
	t.writeLine(fmt.Sprintf("VarDeclarationStatement (Name: %s, Constant: %v)", node.Name, node.IsConstant))
	if node.ExplicitType != nil {
		if err := t.visitType(node.ExplicitType, node.Value == nil); err != nil {
			return err
		}
	}
	if node.Value != nil {
		if err := t.visitChild(node.Value, true); err != nil {
			return err
		}
	}
	return nil
}

func (t *TreeDrawer) VisitReturnStatement(node *ast.ReturnStatement) error {
	t.writeLine("ReturnStatement")
	if node.Value != nil {
		return t.visitChild(node.Value, true)
	}
	return nil
}

func (t *TreeDrawer) VisitExpressionStatement(node *ast.ExpressionStatement) error {
	t.writeLine("ExpressionStatement")
	if node.Exp != nil {
		return t.visitChild(node.Exp, true)
	}
	return nil
}

func (t *TreeDrawer) VisitConditionalStatement(node *ast.ConditionalStatetement) error {
	t.writeLine("IfStatement")

	t.isLast = append(t.isLast, false)
	t.writeLine("Condition")

	if err := t.visitChild(node.Condition, false); err != nil {
		return err
	}

	t.isLast = t.isLast[:len(t.isLast)-1]
	t.isLast = append(t.isLast, false)
	t.writeLine("Success")

	if err := t.visitChild(&node.Success, false); err != nil {
		return err
	}

	t.isLast = t.isLast[:len(t.isLast)-1]
	t.isLast = append(t.isLast, false)

	t.writeLine("Failure")
	if err := t.visitChild(&node.Failure, true); err != nil {
		return err
	}

	return nil
}

func (t *TreeDrawer) VisitWhileStatement(node *ast.WhileStatement) error {
	t.writeLine("WhileStatement")
	t.isLast = append(t.isLast, false)
	t.isLast = t.isLast[:len(t.isLast)-1]
	t.writeLine("Condition")

	if err := t.visitChild(node.Condition, false); err != nil {
		return err
	}

	t.isLast = t.isLast[:len(t.isLast)-1]
	t.isLast = append(t.isLast, false)
	t.writeLine("Body")

	if err := t.visitChild(&node.Body, true); err != nil {
		return err
	}

	return nil
}

func (t *TreeDrawer) VisitFunctionDefinition(node *ast.FuncDeclStatement) error {
	t.writeLine(fmt.Sprintf("FunctionDefinition (Name: %s)", node.Name))
	for i, arg := range node.Args {
		// Simplified last child logic for args
		lastArg := i == len(node.Args)-1 && node.ReturnType == nil && len(node.Body.Body) == 0
		t.isLast = append(t.isLast, lastArg)
		t.writeLine(fmt.Sprintf("Arg: %s", arg.Name))
		if arg.ArgType != nil {
			t.visitType(arg.ArgType, true)
		}
		t.isLast = t.isLast[:len(t.isLast)-1]
	}

	if node.ReturnType != nil {
		if err := t.visitType(node.ReturnType, len(node.Body.Body) == 0); err != nil {
			return err
		}
	}
	return t.visitChild(&node.Body, true)
}

func (t *TreeDrawer) VisitStructDeclaration(node *ast.StructDeclarationStatement) error {
	t.writeLine(fmt.Sprintf("StructDeclaration (Name: %s)", node.Name))
	for i, propName := range node.Properties {
		last := i == len(node.Properties)-1
		t.isLast = append(t.isLast, last)
		t.writeLine(fmt.Sprintf("Property: %s", propName))
		if i < len(node.Types) {
			t.visitType(node.Types[i], true)
		}
		t.isLast = t.isLast[:len(t.isLast)-1]
	}
	return nil
}

func (t *TreeDrawer) VisitPrintStatement(node *ast.PrintStatetement) error {
	t.writeLine("PrintStatement")
	for i, val := range node.Values {
		if err := t.visitChild(val, i == len(node.Values)-1); err != nil {
			return err
		}
	}
	return nil
}

func (t *TreeDrawer) VisitMainStatement(node *ast.MainStatement) error {
	t.writeLine("MainStatement")
	return t.visitChild(&node.Body, true)
}

// Expressions

func (t *TreeDrawer) VisitTupleExpression(node *ast.TupleExpression) error {
	t.writeLine("TupleExpression")
	for i, expr := range node.Expressions {
		if err := t.visitChild(expr, i == len(node.Expressions)-1); err != nil {
			return err
		}
	}
	return nil
}

func (t *TreeDrawer) VisitTupleAssignmentExpression(node *ast.TupleAssignmentExpression) error {
	t.writeLine("TupleAssignmentExpression")
	t.isLast = append(t.isLast, false)
	t.visitChild(node.Assignees, false)
	t.isLast[len(t.isLast)-1] = true
	t.visitChild(node.Value, true)
	t.isLast = t.isLast[:len(t.isLast)-1]
	return nil
}

func (t *TreeDrawer) VisitErrorExpression(node *ast.ErrorExpression) error {
	t.writeLine("ErrorExpression")
	if err := t.visitChild(node.Exp, false); err != nil {
		return err
	}
	return nil
}

func (t *TreeDrawer) VisitBinaryExpression(node *ast.BinaryExpression) error {
	t.writeLine(fmt.Sprintf("BinaryExpression (%s)", node.Operator.Value))
	t.isLast = append(t.isLast, false)
	t.writeLine("Left")

	if err := t.visitChild(node.Left, false); err != nil {
		return err
	}

	t.isLast[len(t.isLast)-1] = true

	t.writeLine("Right")
	if err := t.visitChild(node.Right, true); err != nil {
		return err
	}

	t.isLast = t.isLast[:len(t.isLast)-1]

	return nil
}

func (t *TreeDrawer) VisitFunctionCall(node *ast.FunctionCallExpression) error {
	t.writeLine(fmt.Sprintf("FunctionCall %s", node.Name))

	for i, arg := range node.Args {
		if err := t.visitChild(arg, i == len(node.Args)-1); err != nil {
			return err
		}
	}
	return nil
}

func (t *TreeDrawer) VisitCallExpression(node *ast.CallExpression) error {
	t.writeLine("CallExpression")
	if err := t.visitChild(node.Caller, len(node.Arguments) == 0); err != nil {
		return err
	}
	for i, arg := range node.Arguments {
		if err := t.visitChild(arg, i == len(node.Arguments)-1); err != nil {
			return err
		}
	}
	return nil
}

func (t *TreeDrawer) VisitStringExpression(node *ast.StringExpression) error {
	t.writeLine(fmt.Sprintf("StringExpression (%q)", node.Value))
	return nil
}

func (t *TreeDrawer) VisitNumberExpression(node *ast.NumberExpression) error {
	t.writeLine(fmt.Sprintf("NumberExpression (%d)", node.Value))
	return nil
}

func (t *TreeDrawer) VisitFloatExpression(node *ast.FloatExpression) error {
	t.writeLine(fmt.Sprintf("FloatExpression (%f)", node.Value))
	return nil
}

func (t *TreeDrawer) VisitFloatingBlockExpression(node *ast.FloatingBlockExpression) error {
	t.writeLine("FloatingBlockExpression")
	if err := t.visitChild(node.Stmt, false); err != nil {
		return err
	}

	return nil
}

func (t *TreeDrawer) VisitSymbolExpression(node *ast.SymbolExpression) error {
	t.writeLine(fmt.Sprintf("SymbolExpression (%s)", node.Value))
	return nil
}

func (t *TreeDrawer) VisitAssignmentExpression(node *ast.AssignmentExpression) error {
	t.writeLine(fmt.Sprintf("AssignmentExpression (%s)", node.Operator.Value))
	if err := t.visitChild(node.Assignee, false); err != nil {
		return err
	}
	return t.visitChild(node.Value, true)
}

func (t *TreeDrawer) VisitArrayAccessExpression(node *ast.ArrayAccessExpression) error {
	t.writeLine("ArrayAccessExpression")
	t.isLast = append(t.isLast, false)
	t.writeLine("Name")

	if err := t.visitChild(node.Name, false); err != nil {
		return err
	}

	t.isLast[len(t.isLast)-1] = true
	t.writeLine("Index")

	if err := t.visitChild(node.Index, true); err != nil {
		return err
	}

	t.isLast = t.isLast[:len(t.isLast)-1]

	return nil
}

func (t *TreeDrawer) VisitMemberExpression(node *ast.MemberExpression) error {
	t.writeLine("MemberExpression")
	if err := t.visitChild(node.Object, false); err != nil {
		return err
	}
	return t.visitChild(node.Property, true)
}

func (t *TreeDrawer) VisitStructInitializationExpression(node *ast.StructInitializationExpression) error {
	t.writeLine(fmt.Sprintf("StructInitialization (%s)", node.Name))
	for i, propName := range node.Properties {
		t.isLast = append(t.isLast, i == len(node.Properties)-1)
		t.writeLine(fmt.Sprintf("Property: %s", propName))
		if i < len(node.Values) {
			t.visitChild(node.Values[i], true)
		}
		t.isLast = t.isLast[:len(t.isLast)-1]
	}
	return nil
}

func (t *TreeDrawer) VisitArrayInitializationExpression(node *ast.ArrayInitializationExpression) error {
	t.writeLine("ArrayInitialization")
	for i, value := range node.Contents {
		if err := t.visitChild(value, i == len(node.Contents)-1); err != nil {
			return err
		}
	}
	return nil
}

func (t *TreeDrawer) VisitArrayOfStructsAccessExpression(node *ast.ArrayOfStructsAccessExpression) error {
	t.writeLine("ArrayOfStructsAccess")
	if err := t.visitChild(node.Name, false); err != nil {
		return err
	}
	if err := t.visitChild(node.Index, false); err != nil {
		return err
	}
	return t.visitChild(node.Property, true)
}

func (t *TreeDrawer) VisitPrefixExpression(node *ast.PrefixExpression) error {
	t.writeLine(fmt.Sprintf("PrefixExpression (%s)", node.Operator.Value))
	return t.visitChild(node.RightExpression, true)
}

// Types

func (t *TreeDrawer) VisitSymbolType(node *ast.SymbolType) error {
	t.writeLine(fmt.Sprintf("Type: Symbol (%s)", node.Name))
	return nil
}

func (t *TreeDrawer) VisitTupleType(node *ast.TupleType) error {
	t.writeLine("Type: Tuple")
	for i, typ := range node.Types {
		if err := t.visitType(typ, i == len(node.Types)-1); err != nil {
			return err
		}
	}
	return nil
}

func (t *TreeDrawer) VisitNumberType(node *ast.NumberType) error {
	t.writeLine("Type: Number")
	return nil
}

func (t *TreeDrawer) VisitNumber64Type(node *ast.Number64Type) error {
	t.writeLine("Type: Number64")
	return nil
}

func (t *TreeDrawer) VisitFloatType(node *ast.FloatType) error {
	t.writeLine("Type: Float")
	return nil
}

func (t *TreeDrawer) VisitPointerType(node *ast.PointerType) error {
	t.writeLine("Type: Pointer")
	return t.visitType(node.Underlying, true)
}

func (t *TreeDrawer) VisitStringType(node *ast.StringType) error {
	t.writeLine("Type: String")
	return nil
}

func (t *TreeDrawer) VisitErrorType(node *ast.ErrorType) error {
	t.writeLine("Type: Error")
	return nil
}

func (t *TreeDrawer) VisitArrayType(node *ast.ArrayType) error {
	t.writeLine("Type: Array")
	return t.visitType(node.Underlying, true)
}

func (t *TreeDrawer) VisitVoidType(node *ast.VoidType) error {
	t.writeLine("Type: Void")
	return nil
}

func (t *TreeDrawer) VisitZeroExpression(node *ast.ZeroExpression) error {
	t.writeLine("ZeroExpression")

	return t.visitType(node.T, true)
}

func (t *TreeDrawer) ZeroOfSymbolType(node *ast.SymbolType) error     { return nil }
func (t *TreeDrawer) ZeroOfTupleType(node *ast.TupleType) error       { return nil }
func (t *TreeDrawer) ZeroOfNumberType(node *ast.NumberType) error     { return nil }
func (t *TreeDrawer) ZeroOfNumber64Type(node *ast.Number64Type) error { return nil }
func (t *TreeDrawer) ZeroOfFloatType(node *ast.FloatType) error       { return nil }
func (t *TreeDrawer) ZeroOfPointerType(node *ast.PointerType) error   { return nil }
func (t *TreeDrawer) ZeroOfStringType(node *ast.StringType) error     { return nil }
func (t *TreeDrawer) ZeroOfErrorType(node *ast.ErrorType) error       { return nil }
func (t *TreeDrawer) ZeroOfArrayType(node *ast.ArrayType) error       { return nil }
func (t *TreeDrawer) ZeroOfVoidType(node *ast.VoidType) error         { return nil }
