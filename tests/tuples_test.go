package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTupleReturns(t *testing.T) {
	t.Run("Return Integers", func(t *testing.T) {
		t.Parallel()
		req := CompileRequest{
			InputPath:               "./tuples/return_ints.swa",
			ExpectedExecutionOutput: "x: 10, y: 20\n",
			T:                       t,
		}
		req.AssertCompileAndExecute()
	})

	t.Run("Return Floats", func(t *testing.T) {
		t.Parallel()
		req := CompileRequest{
			InputPath:               "./tuples/return_floats.swa",
			ExpectedExecutionOutput: "x: 1.5, y: 2.5\n",
			T:                       t,
		}
		req.AssertCompileAndExecute()
	})

	t.Run("Return Error", func(t *testing.T) {
		t.Parallel()
		req := CompileRequest{
			InputPath:               "./tuples/return_error.swa",
			ExpectedExecutionOutput: "val: 10, err: \nval: 0, err: value is negative\n",
			T:                       t,
		}
		req.AssertCompileAndExecute()
	})

	t.Run("Swap", func(t *testing.T) {
		t.Parallel()
		req := CompileRequest{
			InputPath:               "./tuples/swap.swa",
			ExpectedExecutionOutput: "Before: a=10, b=20\nAfter: a=20, b=10\n",
			T:                       t,
		}
		req.AssertCompileAndExecute()
	})

	t.Run("Swap typecheck", func(t *testing.T) {
		t.Parallel()
		req := CompileRequest{
			InputPath:      "./tuples/swap-typecheck.swa",
			ExpectedOutput: "Expected assignment of Number but got Tuple\n",
			T:              t,
		}

		assert.Error(t, req.Compile())
	})

	t.Run("Mixed Types", func(t *testing.T) {
		t.Parallel()
		req := CompileRequest{
			InputPath:               "./tuples/mixed.swa",
			ExpectedExecutionOutput: "i=1, f=2.5\n",
			T:                       t,
		}
		req.AssertCompileAndExecute()
	})

	t.Run("Recursion", func(t *testing.T) {
		t.Parallel()
		req := CompileRequest{
			InputPath:               "./tuples/recursion.swa",
			ExpectedExecutionOutput: "x=5, y=10\n",
			T:                       t,
		}
		req.AssertCompileAndExecute()
	})

	t.Run("Assign to array access expression", func(t *testing.T) {
		t.Parallel()
		req := CompileRequest{
			InputPath:               "./tuples/assign-to-array-index.swa",
			ExpectedExecutionOutput: "tab[0]=15, tab[1]=150\n",
			T:                       t,
		}
		req.AssertCompileAndExecute()
	})

	t.Run("Assign to array access expression in struct", func(t *testing.T) {
		t.Parallel()
		req := CompileRequest{
			InputPath:               "./tuples/assign-to-array-index-in-struct.swa",
			ExpectedExecutionOutput: "tab[0]=15, tab[1]=150\n",
			T:                       t,
		}
		req.AssertCompileAndExecute()
	})

	t.Run("Assign to array of structs access", func(t *testing.T) {
		t.Parallel()
		req := CompileRequest{
			InputPath:               "./tuples/assign-to-array-of-structs-index.swa",
			ExpectedExecutionOutput: "tab[0]=15, tab[1]=150\n",
			T:                       t,
		}
		req.AssertCompileAndExecute()
	})

	t.Run("Assign to member expression", func(t *testing.T) {
		t.Parallel()
		req := CompileRequest{
			InputPath:               "./tuples/assign-to-member-expression.swa",
			ExpectedExecutionOutput: "p.div = 4, p.mod = 8",
			T:                       t,
		}
		req.AssertCompileAndExecute()
	})

	t.Run("Assign to member typecheck", func(t *testing.T) {
		t.Parallel()
		req := CompileRequest{
			InputPath:      "./tuples/assign-to-member-expression-typecheck.swa",
			ExpectedOutput: "Values mismatch in tuple expected Number got Error at index 0\n",
			T:              t,
		}

		assert.Error(t, req.Compile())
	})

	t.Run("Assign to member typecheck length", func(t *testing.T) {
		t.Parallel()
		req := CompileRequest{
			InputPath:      "./tuples/assign-to-member-expression-typecheck-length.swa",
			ExpectedOutput: "length of tuples does not match expected 3 got 4\n",
			T:              t,
		}

		assert.Error(t, req.Compile())
	})
}
