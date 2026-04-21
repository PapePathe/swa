package tests

import (
	"testing"
)

func TestSlicePrimitives(t *testing.T) {
	t.Run("append and index", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./slices/primitives/all.swa",
			"count=3  [10 20 30]\ncount=3 [1 0 1]\ncount=3 [1.5 2.5 3.5]\ncount=3 [hello world swa]\n")
	})

	t.Run("zero-initialized slice", func(t *testing.T) {
		NewFailedCompileRequest(t,
			"./slices/primitives/zero.swa",
			"List should be initialized. Use make([]datatype)\n")
	})
}

func TestSliceGrow(t *testing.T) {
	t.Run("int slice doubles capacity on overflow", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./slices/grow/int_grow.swa",
			"full: count=2full: count=3values: 1 2 3")
	})
}

func TestSliceOfStructs(t *testing.T) {
	t.Run("len and cap after appending Point structs", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./slices/structs/points.swa",
			"len=3")
	})
}
