package tests

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
