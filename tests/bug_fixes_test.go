package tests

import (
	"testing"
)

func TestBugFixes(t *testing.T) {
	t.Run("00001", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./bug-fixes/00001-assign-array-element-to-variable.swa",
			"first elem: 15, assigned value: 15",
		)
	})

	t.Run("00002", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./bug-fixes/00002-assign-array-element-to-variable.swa",
			"last elem: 45",
		)
	})

	t.Run("00003", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./bug-fixes/00003-assign-array-element-to-variable.swa",
			"(first persons's age: 11)(last persons's age: 44)",
		)
	})

	t.Run("111-string-reassignment-produces-garbage-values", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./bug-fixes/111-string-reassignment-produces-garbage-values.english.swa",
			"Before: initial, After: changed",
		)
	})

	t.Run("95-logical-operators-and-are-not-implemented", func(t *testing.T) {
		t.Run("1", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./bug-fixes/95-logical-operators-and-are-not-implemented.1.english.swa",
				"1 == 1 && 0 == 0 evaluates to true",
			)
		})

		t.Run("2", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./bug-fixes/95-logical-operators-and-are-not-implemented.2.english.swa",
				"1 == 1 && 1 == 0 evaluates to false",
			)
		})

		t.Run("3", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./bug-fixes/95-logical-operators-and-are-not-implemented.3.english.swa",
				"1 == 0 && 1 == 1 evaluates to false",
			)
		})

		t.Run("4", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./bug-fixes/95-logical-operators-and-are-not-implemented.4.english.swa",
				"1 == 0 && 1 == 0 evaluates to false",
			)
		})

		t.Run("5", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./bug-fixes/95-logical-operators-and-are-not-implemented.5.english.swa",
				"1 == 1 || 1 == 1 evaluates to true",
			)
		})

		t.Run("6", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./bug-fixes/95-logical-operators-and-are-not-implemented.6.english.swa",
				"1 == 1 || 1 == 0 evaluates to true",
			)
		})

		t.Run("7", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./bug-fixes/95-logical-operators-and-are-not-implemented.7.english.swa",
				"1 == 0 || 1 == 1 evaluates to true",
			)
		})

		t.Run("8", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./bug-fixes/95-logical-operators-and-are-not-implemented.8.english.swa",
				"1 == 0 || 1 == 0 evaluates to false",
			)
		})

		t.Run("9", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./bug-fixes/95-logical-operators-and-are-not-implemented.9.english.swa",
				"a == 1 && b == 0 evaluates to true",
			)
		})

		t.Run("10", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./bug-fixes/95-logical-operators-and-are-not-implemented.10.english.swa",
				"a == 0 || b == 1 evaluates to true",
			)
		})

		t.Run("11", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./bug-fixes/95-logical-operators-and-are-not-implemented.11.english.swa",
				"x == 1.5 && y == 2.5 evaluates to true",
			)
		})

		t.Run("12", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./bug-fixes/95-logical-operators-and-are-not-implemented.12.english.swa",
				"x == 0.5 || y == 0.0 evaluates to true",
			)
		})

		t.Run("13", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./bug-fixes/95-logical-operators-and-are-not-implemented.13.english.swa",
				"a > 3 && x == 2.5 || 1 == 1 evaluates to true",
			)
		})

		t.Run("14", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./bug-fixes/95-logical-operators-and-are-not-implemented.14.english.swa",
				"(1 == 0 || 1 == 1) && 1 == 1 evaluates to true",
			)
		})

		t.Run("15", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./bug-fixes/95-logical-operators-and-are-not-implemented.15.english.swa",
				"1 == 1 || (1 == 0 && 1 == 0) evaluates to true",
			)
		})

		t.Run("16", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./bug-fixes/95-logical-operators-and-are-not-implemented.16.english.swa",
				"((1 == 1 || 1 == 0) && 1 == 1) || 1 == 0 evaluates to true",
			)
		})

		t.Run("17", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./bug-fixes/95-logical-operators-and-are-not-implemented.17.english.swa",
				"(a == 1 && b == 0) || (a == 0 && b == 1) evaluates to true",
			)
		})
	})

	t.Run("93-integer-comparison-uses-unsigned-instructions-for-signed-integers", func(t *testing.T) {
		t.Run("1", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./bug-fixes/93-integer-comparison-uses-unsigned-instructions-for-signed-integers.1.french.swa",
				"okok",
			)
		})

		t.Run("2", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./bug-fixes/93-integer-comparison-uses-unsigned-instructions-for-signed-integers.2.french.swa",
				"okok",
			)
		})
	})

	t.Run("106-unhandled-64-bit-integer-overflow", func(t *testing.T) {
		t.Run("1", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./bug-fixes/106-unhandled-64-bit-integer-overflow.max.english.swa",
				"9223372036854775808: value out of range while parsing number expression",
			)
		})

		t.Run("2", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./bug-fixes/106-unhandled-64-bit-integer-overflow.min.english.swa",
				"9223372036854775809: value out of range while parsing number expression",
			)
		})
	})
	t.Run("105-silent-32-bit-integer-overflow", func(t *testing.T) {
		t.Run("1", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./bug-fixes/105-silent-32-bit-integer-overflow-max.english.swa",
				"2147483648 is greater than max value for int32\n",
			)
		})

		t.Run("2", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./bug-fixes/105-silent-32-bit-integer-overflow-min.english.swa",
				"2147483649 is greater than max value for int32\n",
			)
		})
	})

	t.Run("missing-type-check-in-assignment-expression", func(t *testing.T) {
		t.Run("1", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./bug-fixes/missing-type-check-in-assignment-expression.1.english.swa",
				"Expected assignment of Float but got String\n",
			)
		})

		t.Run("2", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./bug-fixes/missing-type-check-in-assignment-expression.2.english.swa",
				"Expected assignment of Float but got Number\n",
			)
		})

		t.Run("3", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./bug-fixes/missing-type-check-in-assignment-expression.3.english.swa",
				"Expected assignment of Float but got Array\n",
			)
		})

		t.Run("4", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./bug-fixes/missing-type-check-in-assignment-expression.4.english.swa",
				"Expected assignment of Float but got Symbol\n",
			)
		})

		t.Run("5", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./bug-fixes/missing-type-check-in-assignment-expression.5.english.swa",
				"Expected assignment of Number but got Float\n",
			)
		})

		t.Run("6", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./bug-fixes/missing-type-check-in-assignment-expression.6.english.swa",
				"Expected assignment of Number but got Float\n",
			)
		})

		t.Run("7", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./bug-fixes/missing-type-check-in-assignment-expression.7.english.swa",
				"Expected assignment of Number but got Float\n",
			)
		})

		t.Run("8", func(t *testing.T) {
			t.Run("English", func(t *testing.T) {
				NewSuccessfulCompileRequest(t,
					"./bug-fixes/missing-type-check-in-assignment-expression.8.english.swa",
					"",
				)
			})
			t.Run("French", func(t *testing.T) {
				NewSuccessfulCompileRequest(t,
					"./bug-fixes/missing-type-check-in-assignment-expression.8.french.swa",
					"",
				)
			})
			t.Run("Soussou", func(t *testing.T) {
				NewSuccessfulCompileRequest(t,
					"./bug-fixes/missing-type-check-in-assignment-expression.8.soussou.swa",
					"",
				)
			})
		})
	})

	t.Run("99-missing-type-check-in-variable-declaration", func(t *testing.T) {
		t.Run("1", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./bug-fixes/99-missing-type-check-in-variable-declaration.english.1.swa",
				"expected String but got Float\n",
			)
		})

		t.Run("2", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./bug-fixes/99-missing-type-check-in-variable-declaration.english.2.swa",
				"expected String but got Number\n",
			)
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
		//			ExpectedOutput: "expected Number got pointer of %!v(PANIC=String method: unreachable)\n",
		//			T:              t,
		//		}

		//		assert.Error(t, req.Compile())
		//	})

		t.Run("5", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./bug-fixes/99-missing-type-check-in-variable-declaration.english.5.swa",
				"expected Number but got Float\n",
			)
		})

		t.Run("6", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./bug-fixes/99-missing-type-check-in-variable-declaration.english.6.swa",
				"expected Number but got String\n",
			)
		})
	})

	//	t.Run("103-missing-argument-type-check-in-function-calls", func(t *testing.T) {
	//		t.Run("English", func(t *testing.T) {
	//			req := CompileRequest{
	//				InputPath:      "./bug-fixes/103-missing-argument-type-check-in-function-calls.english.swa",
	//				ExpectedOutput: "expected argument of type IntegerType(32 bits) but got PointerType(Reference)\n",
	//				T:              t,
	//			}
	//
	//			assert.Error(t, req.Compile())
	//		})
	//	})

	t.Run("102-missing-arity-check-in-function-calls", func(t *testing.T) {
		t.Run("English", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./bug-fixes/102-missing-arity-check-in-function-calls.english.swa",
				"function add expect 2 arguments but was given 1\n",
			)
		})
	})

	t.Run("90-function-calls-with-variables-pass-pointers-instead-of-values", func(t *testing.T) {
		t.Run("English", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./bug-fixes/90-function-calls-with-variables-pass-pointers-instead-of-values.english.swa",
				"z: 30",
			)
		})

		t.Run("French", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./bug-fixes/90-function-calls-with-variables-pass-pointers-instead-of-values.french.swa",
				"z: 30",
			)
		})

		t.Run("Soussou", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./bug-fixes/90-function-calls-with-variables-pass-pointers-instead-of-values.soussou.swa",
				"z: 30",
			)
		})
	})

	t.Run("94-block-statements-do-not-create-new-scopes-shadowing-fails", func(t *testing.T) {
		t.Run("English", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./bug-fixes/94-block-statements-do-not-create-new-scopes-shadowing-fails.english.swa",
				"Inner x: 20, Outer x: 10",
			)
		})

		t.Run("French", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./bug-fixes/94-block-statements-do-not-create-new-scopes-shadowing-fails.french.swa",
				"Inner x: 20, Outer x: 10",
			)
		})

		t.Run("Soussou", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./bug-fixes/94-block-statements-do-not-create-new-scopes-shadowing-fails.soussou.swa",
				"Inner x: 20, Outer x: 10",
			)
		})
	})

	t.Run("114-variable-redeclaration-allowed-in-same-scope", func(t *testing.T) {
		t.Run("English", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./bug-fixes/114-variable-redeclaration-allowed-in-same-scope.english.swa",
				"variable x is already defined\n",
			)
		})

		t.Run("French", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./bug-fixes/114-variable-redeclaration-allowed-in-same-scope.french.swa",
				"variable x is already defined\n",
			)
		})

		t.Run("Soussou", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./bug-fixes/114-variable-redeclaration-allowed-in-same-scope.soussou.swa",
				"variable x na na yi khorun\n",
			)
		})
	})

	t.Run("91-array-of-structs-indexing-does-not-support-variables", func(t *testing.T) {
		t.Run("English", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./bug-fixes/91-array-of-structs-indexing-does-not-support-variables.english.swa",
				"Name: Alice",
			)
		})

		t.Run("French", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./bug-fixes/91-array-of-structs-indexing-does-not-support-variables.french.swa",
				"Name: Alice",
			)
		})

		t.Run("Soussou", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./bug-fixes/91-array-of-structs-indexing-does-not-support-variables.soussou.swa",
				"Name: Alice",
			)
		})
	})

	t.Run("116-keyword-regex-patterns-lack-word-boundaries", func(t *testing.T) {
		t.Run("English", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./bug-fixes/116-keyword-regex-patterns-lack-word-boundaries.english.swa",
				"",
			)
		})

		t.Run("French", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./bug-fixes/116-keyword-regex-patterns-lack-word-boundaries.french.swa",
				"",
			)
		})

		t.Run("Soussou", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./bug-fixes/116-keyword-regex-patterns-lack-word-boundaries.soussou.swa",
				"",
			)
		})
	})
}
