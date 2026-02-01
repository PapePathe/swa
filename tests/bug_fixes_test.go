package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBugFixes(t *testing.T) {
	t.Parallel()

	t.Run("00001", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./bug-fixes/00001-assign-array-element-to-variable.swa",
			ExpectedExecutionOutput: "first elem: 15, assigned value: 15",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})

	t.Run("00002", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./bug-fixes/00002-assign-array-element-to-variable.swa",
			ExpectedExecutionOutput: "last elem: 45",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})

	t.Run("00003", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./bug-fixes/00003-assign-array-element-to-variable.swa",
			ExpectedExecutionOutput: "(first persons's age: 11)(last persons's age: 44)",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})

	t.Run("111-string-reassignment-produces-garbage-values", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./bug-fixes/111-string-reassignment-produces-garbage-values.english.swa",
			ExpectedExecutionOutput: "Before: initial, After: changed",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})

	t.Run("95-logical-operators-and-are-not-implemented", func(t *testing.T) {
		t.Run("1", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./bug-fixes/95-logical-operators-and-are-not-implemented.1.english.swa",
				ExpectedExecutionOutput: "1 == 1 && 0 == 0 evaluates to true",
				T:                       t,
			}

			req.AssertCompileAndExecute()
		})

		t.Run("2", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./bug-fixes/95-logical-operators-and-are-not-implemented.2.english.swa",
				ExpectedExecutionOutput: "1 == 1 && 1 == 0 evaluates to false",
				T:                       t,
			}

			req.AssertCompileAndExecute()
		})

		t.Run("3", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./bug-fixes/95-logical-operators-and-are-not-implemented.3.english.swa",
				ExpectedExecutionOutput: "1 == 0 && 1 == 1 evaluates to false",
				T:                       t,
			}

			req.AssertCompileAndExecute()
		})

		t.Run("4", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./bug-fixes/95-logical-operators-and-are-not-implemented.4.english.swa",
				ExpectedExecutionOutput: "1 == 0 && 1 == 0 evaluates to false",
				T:                       t,
			}

			req.AssertCompileAndExecute()
		})

		t.Run("5", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./bug-fixes/95-logical-operators-and-are-not-implemented.5.english.swa",
				ExpectedExecutionOutput: "1 == 1 || 1 == 1 evaluates to true",
				T:                       t,
			}

			req.AssertCompileAndExecute()
		})

		t.Run("6", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./bug-fixes/95-logical-operators-and-are-not-implemented.6.english.swa",
				ExpectedExecutionOutput: "1 == 1 || 1 == 0 evaluates to true",
				T:                       t,
			}

			req.AssertCompileAndExecute()
		})

		t.Run("7", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./bug-fixes/95-logical-operators-and-are-not-implemented.7.english.swa",
				ExpectedExecutionOutput: "1 == 0 || 1 == 1 evaluates to true",
				T:                       t,
			}

			req.AssertCompileAndExecute()
		})

		t.Run("8", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./bug-fixes/95-logical-operators-and-are-not-implemented.8.english.swa",
				ExpectedExecutionOutput: "1 == 0 || 1 == 0 evaluates to false",
				T:                       t,
			}

			req.AssertCompileAndExecute()
		})

		t.Run("9", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./bug-fixes/95-logical-operators-and-are-not-implemented.9.english.swa",
				ExpectedExecutionOutput: "a == 1 && b == 0 evaluates to true",
				T:                       t,
			}

			req.AssertCompileAndExecute()
		})

		t.Run("10", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./bug-fixes/95-logical-operators-and-are-not-implemented.10.english.swa",
				ExpectedExecutionOutput: "a == 0 || b == 1 evaluates to true",
				T:                       t,
			}

			req.AssertCompileAndExecute()
		})

		t.Run("11", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./bug-fixes/95-logical-operators-and-are-not-implemented.11.english.swa",
				ExpectedExecutionOutput: "x == 1.5 && y == 2.5 evaluates to true",
				T:                       t,
			}

			req.AssertCompileAndExecute()
		})

		t.Run("12", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./bug-fixes/95-logical-operators-and-are-not-implemented.12.english.swa",
				ExpectedExecutionOutput: "x == 0.5 || y == 0.0 evaluates to true",
				T:                       t,
			}

			req.AssertCompileAndExecute()
		})

		t.Run("13", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./bug-fixes/95-logical-operators-and-are-not-implemented.13.english.swa",
				ExpectedExecutionOutput: "a > 3 && x == 2.5 || 1 == 1 evaluates to true",
				T:                       t,
			}

			req.AssertCompileAndExecute()
		})

		t.Run("14", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./bug-fixes/95-logical-operators-and-are-not-implemented.14.english.swa",
				ExpectedExecutionOutput: "(1 == 0 || 1 == 1) && 1 == 1 evaluates to true",
				T:                       t,
			}

			req.AssertCompileAndExecute()
		})

		t.Run("15", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./bug-fixes/95-logical-operators-and-are-not-implemented.15.english.swa",
				ExpectedExecutionOutput: "1 == 1 || (1 == 0 && 1 == 0) evaluates to true",
				T:                       t,
			}

			req.AssertCompileAndExecute()
		})

		t.Run("16", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./bug-fixes/95-logical-operators-and-are-not-implemented.16.english.swa",
				ExpectedExecutionOutput: "((1 == 1 || 1 == 0) && 1 == 1) || 1 == 0 evaluates to true",
				T:                       t,
			}

			req.AssertCompileAndExecute()
		})

		t.Run("17", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./bug-fixes/95-logical-operators-and-are-not-implemented.17.english.swa",
				ExpectedExecutionOutput: "(a == 1 && b == 0) || (a == 0 && b == 1) evaluates to true",
				T:                       t,
			}

			req.AssertCompileAndExecute()
		})
	})

	t.Run("93-integer-comparison-uses-unsigned-instructions-for-signed-integers", func(t *testing.T) {
		t.Run("1", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./bug-fixes/93-integer-comparison-uses-unsigned-instructions-for-signed-integers.1.french.swa",
				ExpectedExecutionOutput: "okok",
				T:                       t,
			}

			req.AssertCompileAndExecute()
		})

		t.Run("2", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./bug-fixes/93-integer-comparison-uses-unsigned-instructions-for-signed-integers.2.french.swa",
				ExpectedExecutionOutput: "okok",
				T:                       t,
			}

			req.AssertCompileAndExecute()
		})
	})

	t.Run("106-unhandled-64-bit-integer-overflow", func(t *testing.T) {
		t.Run("1", func(t *testing.T) {
			req := CompileRequest{
				InputPath:      "./bug-fixes/106-unhandled-64-bit-integer-overflow.max.english.swa",
				ExpectedOutput: "9223372036854775808: value out of range while parsing number expression",
				T:              t,
			}

			assert.Error(t, req.Compile())
		})

		t.Run("2", func(t *testing.T) {
			req := CompileRequest{
				InputPath:      "./bug-fixes/106-unhandled-64-bit-integer-overflow.min.english.swa",
				ExpectedOutput: "-9223372036854775809: value out of range while parsing number expression",
				T:              t,
			}

			assert.Error(t, req.Compile())
		})
	})
	t.Run("105-silent-32-bit-integer-overflow", func(t *testing.T) {
		t.Run("1", func(t *testing.T) {
			req := CompileRequest{
				InputPath:      "./bug-fixes/105-silent-32-bit-integer-overflow-max.english.swa",
				ExpectedOutput: "2147483648 is greater than max value for int32\n",
				T:              t,
			}

			assert.Error(t, req.Compile())
		})

		t.Run("2", func(t *testing.T) {
			req := CompileRequest{
				InputPath:      "./bug-fixes/105-silent-32-bit-integer-overflow-min.english.swa",
				ExpectedOutput: "-2147483649 is smaller than min value for int32\n",
				T:              t,
			}

			assert.Error(t, req.Compile())
		})
	})

	t.Run("missing-type-check-in-assignment-expression", func(t *testing.T) {
		t.Run("1", func(t *testing.T) {
			req := CompileRequest{
				InputPath:      "./bug-fixes/missing-type-check-in-assignment-expression.1.english.swa",
				ExpectedOutput: "Expected assignment of DataTypeFloat but got DataTypeString\n",
				T:              t,
			}

			assert.Error(t, req.Compile())
		})

		t.Run("2", func(t *testing.T) {
			req := CompileRequest{
				InputPath:      "./bug-fixes/missing-type-check-in-assignment-expression.2.english.swa",
				ExpectedOutput: "Expected assignment of DataTypeFloat but got DataTypeNumber\n",
				T:              t,
			}

			assert.Error(t, req.Compile())
		})

		t.Run("3", func(t *testing.T) {
			req := CompileRequest{
				InputPath:      "./bug-fixes/missing-type-check-in-assignment-expression.3.english.swa",
				ExpectedOutput: "Expected assignment of DataTypeFloat but got DataTypeArray\n",
				T:              t,
			}

			assert.Error(t, req.Compile())
		})

		t.Run("4", func(t *testing.T) {
			req := CompileRequest{
				InputPath:      "./bug-fixes/missing-type-check-in-assignment-expression.4.english.swa",
				ExpectedOutput: "Expected assignment of DataTypeFloat but got DataTypeSymbol\n",
				T:              t,
			}

			assert.Error(t, req.Compile())
		})

		t.Run("5", func(t *testing.T) {
			req := CompileRequest{
				InputPath:      "./bug-fixes/missing-type-check-in-assignment-expression.5.english.swa",
				ExpectedOutput: "Expected assignment of DataTypeNumber but got DataTypeFloat\n",
				T:              t,
			}

			assert.Error(t, req.Compile())
		})

		t.Run("6", func(t *testing.T) {
			req := CompileRequest{
				InputPath:      "./bug-fixes/missing-type-check-in-assignment-expression.6.english.swa",
				ExpectedOutput: "Expected assignment of DataTypeNumber but got DataTypeFloat\n",
				T:              t,
			}

			assert.Error(t, req.Compile())
		})

		t.Run("7", func(t *testing.T) {
			req := CompileRequest{
				InputPath:      "./bug-fixes/missing-type-check-in-assignment-expression.7.english.swa",
				ExpectedOutput: "Expected assignment of DataTypeNumber but got DataTypeFloat\n",
				T:              t,
			}

			assert.Error(t, req.Compile())
		})
	})

	t.Run("99-missing-type-check-in-variable-declaration", func(t *testing.T) {
		t.Run("1", func(t *testing.T) {
			req := CompileRequest{
				InputPath:      "./bug-fixes/99-missing-type-check-in-variable-declaration.english.1.swa",
				ExpectedOutput: "expected DataTypeString but got DataTypeFloat\n",
				T:              t,
			}

			assert.Error(t, req.Compile())
		})

		t.Run("2", func(t *testing.T) {
			req := CompileRequest{
				InputPath:      "./bug-fixes/99-missing-type-check-in-variable-declaration.english.2.swa",
				ExpectedOutput: "expected DataTypeString but got DataTypeNumber\n",
				T:              t,
			}

			assert.Error(t, req.Compile())
		})

		//	TODO: implemenent typechecker for var decl with ArrayAccessExpression
		//	t.Run("3", func(t *testing.T) {
		//		req := CompileRequest{
		//			InputPath:      "./bug-fixes/99-missing-type-check-in-variable-declaration.english.3.swa",
		//			ExpectedOutput: "",
		//			T:              t,
		//		}

		//		assert.Error(t, req.Compile())
		//	})

		//	TODO: implemenent typechecker for var decl with ArrayAccessExpression
		//	t.Run("4", func(t *testing.T) {
		//		req := CompileRequest{
		//			InputPath:      "./bug-fixes/99-missing-type-check-in-variable-declaration.english.4.swa",
		//			ExpectedOutput: "expected DataTypeNumber got pointer of %!v(PANIC=String method: unreachable)\n",
		//			T:              t,
		//		}

		//		assert.Error(t, req.Compile())
		//	})

		t.Run("5", func(t *testing.T) {
			req := CompileRequest{
				InputPath:      "./bug-fixes/99-missing-type-check-in-variable-declaration.english.5.swa",
				ExpectedOutput: "expected DataTypeNumber but got DataTypeFloat\n",
				T:              t,
			}

			assert.Error(t, req.Compile())
		})

		t.Run("6", func(t *testing.T) {
			req := CompileRequest{
				InputPath:      "./bug-fixes/99-missing-type-check-in-variable-declaration.english.6.swa",
				ExpectedOutput: "expected DataTypeNumber but got DataTypeString\n",
				T:              t,
			}

			assert.Error(t, req.Compile())
		})
	})

	t.Run("103-missing-argument-type-check-in-function-calls", func(t *testing.T) {
		t.Run("English", func(t *testing.T) {
			req := CompileRequest{
				InputPath:      "./bug-fixes/103-missing-argument-type-check-in-function-calls.english.swa",
				ExpectedOutput: "expected argument of type IntegerType(32 bits) but got PointerType(Reference)\n",
				T:              t,
			}

			assert.Error(t, req.Compile())
		})
	})

	t.Run("102-missing-arity-check-in-function-calls", func(t *testing.T) {
		t.Run("English", func(t *testing.T) {
			req := CompileRequest{
				InputPath:      "./bug-fixes/102-missing-arity-check-in-function-calls.english.swa",
				ExpectedOutput: "function add expect 2 arguments but was given 1\n",
				T:              t,
			}

			assert.Error(t, req.Compile())
		})
	})

	t.Run("90-function-calls-with-variables-pass-pointers-instead-of-values", func(t *testing.T) {
		t.Run("English", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./bug-fixes/90-function-calls-with-variables-pass-pointers-instead-of-values.english.swa",
				ExpectedExecutionOutput: "z: 30",
				T:                       t,
			}

			req.AssertCompileAndExecute()
		})

		t.Run("French", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./bug-fixes/90-function-calls-with-variables-pass-pointers-instead-of-values.french.swa",
				ExpectedExecutionOutput: "z: 30",
				T:                       t,
			}
			req.AssertCompileAndExecute()
		})
	})

	t.Run("94-block-statements-do-not-create-new-scopes-shadowing-fails", func(t *testing.T) {
		t.Run("English", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./bug-fixes/94-block-statements-do-not-create-new-scopes-shadowing-fails.english.swa",
				ExpectedExecutionOutput: "Inner x: 20, Outer x: 10",
				T:                       t,
			}

			req.AssertCompileAndExecute()
		})

		t.Run("French", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./bug-fixes/94-block-statements-do-not-create-new-scopes-shadowing-fails.french.swa",
				ExpectedExecutionOutput: "Inner x: 20, Outer x: 10",
				T:                       t,
			}
			req.AssertCompileAndExecute()
		})
	})

	t.Run("114-variable-redeclaration-allowed-in-same-scope", func(t *testing.T) {
		t.Run("English", func(t *testing.T) {
			req := CompileRequest{
				InputPath:      "./bug-fixes/114-variable-redeclaration-allowed-in-same-scope.english.swa",
				ExpectedOutput: "variable x is already defined\n",
				T:              t,
			}

			assert.Error(t, req.Compile())
		})

		t.Run("French", func(t *testing.T) {
			req := CompileRequest{
				InputPath:      "./bug-fixes/114-variable-redeclaration-allowed-in-same-scope.french.swa",
				ExpectedOutput: "variable x is already defined\n",
				T:              t,
			}

			assert.Error(t, req.Compile())
		})
	})

	t.Run("91-array-of-structs-indexing-does-not-support-variables", func(t *testing.T) {
		t.Run("English", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./bug-fixes/91-array-of-structs-indexing-does-not-support-variables.english.swa",
				ExpectedExecutionOutput: "Name: Alice",
				T:                       t,
			}

			req.AssertCompileAndExecute()
		})
	})

	t.Run("116-keyword-regex-patterns-lack-word-boundaries", func(t *testing.T) {
		t.Run("English", func(t *testing.T) {
			req := CompileRequest{
				InputPath: "./bug-fixes/116-keyword-regex-patterns-lack-word-boundaries.english.swa",
				T:         t,
			}

			req.AssertCompileAndExecute()
		})

		t.Run("French", func(t *testing.T) {
			req := CompileRequest{
				InputPath: "./bug-fixes/116-keyword-regex-patterns-lack-word-boundaries.french.swa",
				T:         t,
			}

			req.AssertCompileAndExecute()
		})
	})
}
