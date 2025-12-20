package ast

type CodeGenerator interface {
	// Statements
	VisitBlockStatement(node *BlockStatement) error
	VisitVarDeclaration(node *VarDeclarationStatement) error
	VisitReturnStatement(node *ReturnStatement) error
	VisitExpressionStatement(node *ExpressionStatement) error
	VisitConditionalStatement(node *ConditionalStatetement) error
	VisitWhileStatement(node *WhileStatement) error
	VisitFunctionDefinition(node *FuncDeclStatement) error
	VisitStructDeclaration(node *StructDeclarationStatement) error
	VisitPrintStatement(node *PrintStatetement) error
	VisitMainStatement(node *MainStatement) error

	// Expressions
	VisitBinaryExpression(node *BinaryExpression) error
	VisitFunctionCall(node *FunctionCallExpression) error
	VisitCallExpression(node *CallExpression) error
	VisitStringExpression(node *StringExpression) error
	VisitNumberExpression(node *NumberExpression) error
	VisitFloatExpression(node *FloatExpression) error
	VisitSymbolExpression(node *SymbolExpression) error
	VisitAssignmentExpression(node *AssignmentExpression) error
	VisitArrayAccessExpression(node *ArrayAccessExpression) error
	VisitMemberExpression(node *MemberExpression) error
	VisitStructInitializationExpression(node *StructInitializationExpression) error
	VisitArrayInitializationExpression(node *ArrayInitializationExpression) error
	VisitArrayOfStructsAccessExpression(node *ArrayOfStructsAccessExpression) error
	VisitPrefixExpression(node *PrefixExpression) error
}
