package tests

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollatz(t *testing.T) {
	req := CompileRequest{
		InputPath:               "./regression/collatz.swa",
		ExpectedExecutionOutput: "Testing Collatz Conjecture...\nCollatz(27) took 111 steps.\nSuccess: Collatz(27) is correct.\n",
		T:                       t,
	}

	req.AssertCompileAndExecute()
}

func TestBrainFuck(t *testing.T) {
	req := CompileRequest{
		InputPath:               "./regression/brainfuck.swa",
		ExpectedExecutionOutput: "Testing Brainfuck Simulation...\nSimulating BF: +++[->+<]\nResult in cell 1: 3\nSuccess: BF simulation worked.\n",
		T:                       t,
	}

	req.AssertCompileAndExecute()
}

func TestAckerman(t *testing.T) {
	req := CompileRequest{
		InputPath:               "./regression/ackermann.swa",
		ExpectedExecutionOutput: "Testing Ackermann function...\nackermann(3, 2) = 29\nSuccess: Ackermann A(3, 2) is correct.\n",
		T:                       t,
	}

	req.AssertCompileAndExecute()
}

func TestStringEdgeCases(t *testing.T) {
	req := CompileRequest{
		InputPath:               "./regression/string_edge_cases.swa",
		ExpectedExecutionOutput: "s1: Hello\nWorld\ns2: Tab\tCharacter\ns3: Quote \" inside\ns4: Backslash \\ inside\nLong string: This is a very long string designed to test how the compiler handles larger literals in the assembly generation phase. It should be able to handle strings that are significantly longer than short identifiers or simple messages.\n",
		T:                       t,
	}

	req.AssertCompileAndExecute()
}

func TestMultiDimensionalArrayAccess(t *testing.T) {
	req := CompileRequest{
		InputPath:      "./regression/multidim_array.swa",
		ExpectedOutput: "TODO: nested array access not yet supported\n",
		T:              t,
	}

	assert.Error(t, req.Compile())
}

func TestLargeBoolean(t *testing.T) {
	req := CompileRequest{
		InputPath:               "./regression/large_boolean.swa",
		ExpectedExecutionOutput: "Success: Complex boolean logic evaluated correctly.\n",
		T:                       t,
	}

	req.AssertCompileAndExecute()
}

func TestSymbolValue(t *testing.T) {
	req := CompileRequest{
		InputPath:      "./regression/symbol-value.swa",
		ExpectedOutput: "TODO implement VisitSymbolValueExpression\n",
		T:              t,
	}

	assert.Error(t, req.Compile())
}

func TestSymbolAddress(t *testing.T) {
	req := CompileRequest{
		InputPath:      "./regression/symbol-address.swa",
		ExpectedOutput: "TODO implement VisitSymbolAdressExpression\n",
		T:              t,
	}

	assert.Error(t, req.Compile())
}

func TestDoublePointer(t *testing.T) {
	t.Run("french", func(t *testing.T) {
		t.Run("compile", func(t *testing.T) {
			req := CompileRequest{
				InputPath:      "./regression/double_pointer.french.swa",
				ExpectedOutput: "TODO implement VisitSymbolAdressExpression\n",
				T:              t,
			}

			assert.Error(t, req.Compile())
		})

		t.Run("parse as json", func(t *testing.T) {
			expected, _ := os.ReadFile("./regression/double_pointer.french.json")
			req := CompileRequest{
				InputPath:      "./regression/double_pointer.french.swa",
				ExpectedOutput: string(expected),
				T:              t,
			}
			assert.NoError(t, req.Parse("json"))
		})

		t.Run("parse as tree", func(t *testing.T) {
			expected, _ := os.ReadFile("./regression/double_pointer.french.tree")
			req := CompileRequest{
				InputPath:      "./regression/double_pointer.french.swa",
				ExpectedOutput: string(expected),
				T:              t,
			}
			assert.NoError(t, req.Parse("tree"))
		})
	})

	t.Run("english", func(t *testing.T) {
		req := CompileRequest{
			InputPath:      "./regression/double_pointer.swa",
			ExpectedOutput: "TODO implement VisitSymbolAdressExpression\n",
			T:              t,
		}

		assert.Error(t, req.Compile())
	})

	t.Run("soussou", func(t *testing.T) {
		req := CompileRequest{
			InputPath:      "./regression/double_pointer.soussou.swa",
			ExpectedOutput: "TODO implement VisitSymbolAdressExpression\n",
			T:              t,
		}

		assert.Error(t, req.Compile())
	})
}

func TestFloatingBlockStatement(t *testing.T) {
	t.Run("english", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./regression/shadowing.swa",
			ExpectedExecutionOutput: "Outer x: 10\nInner x: 20\nOuter x again: 10\nSuccess: Variable shadowing works.\n",
			T:                       t,
		}

		req.AssertCompileAndExecute()

		t.Run("parse as json", func(t *testing.T) {
			expected, _ := os.ReadFile("./regression/shadowing.json")
			req := CompileRequest{
				InputPath:      "./regression/shadowing.swa",
				ExpectedOutput: string(expected),
				T:              t,
			}
			assert.NoError(t, req.Parse("json"))
		})

		t.Run("parse as tree", func(t *testing.T) {
			expected, _ := os.ReadFile("./regression/shadowing.tree")
			req := CompileRequest{
				InputPath:      "./regression/shadowing.swa",
				ExpectedOutput: string(expected),
				T:              t,
			}
			assert.NoError(t, req.Parse("tree"))
		})
	})

	t.Run("french", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./regression/shadowing.soussou.swa",
			ExpectedExecutionOutput: "Outer x: 10\nInner x: 20\nOuter x again: 10\nSuccess: Variable shadowing works.\n",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})

	t.Run("french", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./regression/shadowing.french.swa",
			ExpectedExecutionOutput: "Outer x: 10\nInner x: 20\nOuter x again: 10\nSuccess: Variable shadowing works.\n",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})
}
