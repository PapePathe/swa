package tests

import (
	"testing"
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
			ExpectedExecutionOutput: "val: 10, err: 0\nval: 0, err: 1\n",
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
}
