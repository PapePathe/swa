package tests

import (
	"testing"
)

func TestTupleReturns(t *testing.T) {
	t.Run("Return Integers", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./tuples/return_ints.swa",
			"x: 10, y: 20\n")
	})

	t.Run("Return Floats", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./tuples/return_floats.swa",
			"x: 1.5, y: 2.5\n")
	})

	t.Run("Return Error", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./tuples/return_error.swa",
			"val: 10, err: \nval: 0, err: value is negative\n")
	})

	t.Run("Swap", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./tuples/swap.swa",
			"Before: a=10, b=20\nAfter: a=20, b=10\n")
	})

	t.Run("Swap typecheck", func(t *testing.T) {
		NewFailedCompileRequest(t,
			"./tuples/swap-typecheck.swa",
			"Expected assignment of Number but got Tuple\n")
	})

	t.Run("Mixed Types", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./tuples/mixed.swa",
			"i=1, f=2.5\n")
	})

	t.Run("Recursion", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./tuples/recursion.swa",
			"x=5, y=10\n")
	})

	t.Run("Assign to array access expression", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./tuples/assign-to-array-index.swa",
			"tab[0]=15, tab[1]=150\n")
	})

	t.Run("Assign to array access expression in struct", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./tuples/assign-to-array-index-in-struct.swa",
			"tab[0]=15, tab[1]=150\n")
	})

	t.Run("Assign to array of structs access", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./tuples/assign-to-array-of-structs-index.swa",
			"tab[0]=15, tab[1]=150\n")
	})

	t.Run("Assign to member expression", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./tuples/assign-to-member-expression.swa",
			"p.div = 4, p.mod = 8")
	})

	t.Run("Assign to member typecheck", func(t *testing.T) {
		NewFailedCompileRequest(t,
			"./tuples/assign-to-member-expression-typecheck.swa",
			"Values mismatch in tuple expected Number got Error at index 0\n")
	})

	t.Run("Assign to member typecheck length", func(t *testing.T) {
		NewFailedCompileRequest(t,
			"./tuples/assign-to-member-expression-typecheck-length.swa",
			"length of tuples does not match expected 3 got 4\n")
	})
}
