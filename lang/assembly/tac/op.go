package tac

// op.go is intentionally minimal. The Op enum and opMap were removed when
// instructions were refactored into typed AsmOp structs (InstAdd, InstSub,
// etc.). Binary operator dispatch now lives in Triple.VisitBinaryExpression.
