package tests

import (
	"testing"
)

// TestSlicePrimitives tests basic slice operations with int elements:
// make, append, len, cap, and indexed access.
func TestSlicePrimitives(t *testing.T) {
	t.Run("append and index", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./slices/primitives/ints.swa",
			"len=3 cap=4 [10 20 30]")
	})

	t.Run("zero-initialized slice", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./slices/primitives/zero.swa",
			"len=0len=1 val=42")
	})
}

// TestSliceGrow verifies that the slice backing array is reallocated when
// capacity is exhausted (analogous to Go's append doubling behaviour).
func TestSliceGrow(t *testing.T) {
	t.Run("int slice doubles capacity on overflow", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./slices/grow/int_grow.swa",
			"full: len=2 cap=2grown: len=3 cap=4values: 1 2 3")
	})
}

// TestSliceOfStructs tests that slices whose element type is a struct
// correctly track length and capacity after successive appends.
func TestSliceOfStructs(t *testing.T) {
	t.Run("len and cap after appending Point structs", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./slices/structs/points.swa",
			"len=3 cap=3")
	})
}

// TestSliceOfFloats tests make/append/index on a []float slice.
func TestSliceOfFloats(t *testing.T) {
	t.Run("append and index contents", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./slices/primitives/floats.swa",
			"len=3 [1.5 2.5 3.5]")
	})
}

// TestSliceOfStrings tests make/append/index on a []string slice.
func TestSliceOfStrings(t *testing.T) {
	t.Run("append and index contents", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./slices/primitives/strings.swa",
			"len=3 [hello world swa]")
	})
}

// TestSliceOfBools tests make/append/index on a []bool slice.
// Booleans are printed as 1 (true) and 0 (false).
func TestSliceOfBools(t *testing.T) {
	t.Run("append and index contents", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./slices/primitives/bools.swa",
			"len=3 [1 0 1]")
	})
}
