package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArrayBounds(t *testing.T) {
	t.Parallel()

	t.Run("ValidBounds", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./arrays/bounds/valid_bounds.swa",
			ExpectedExecutionOutput: "First: 10, Last: 30",
			T:                       t,
		}
		req.AssertCompileAndExecute()
	})

	t.Run("InvalidNegative", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./arrays/bounds/invalid_negative.swa",
			ExpectedExecutionOutput: "Runtime Error: Index -1 out of bounds (0-2)",
			T:                       t,
		}

		assert.NoError(t, req.Compile())
		assert.Error(t, req.RunProgram())
	})

	t.Run("InvalidEqualLength", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./arrays/bounds/invalid_equal_length.swa",
			ExpectedExecutionOutput: "Runtime Error: Index 3 out of bounds (0-2)",
			T:                       t,
		}

		assert.NoError(t, req.Compile())
		assert.Error(t, req.RunProgram())
	})

	t.Run("InvalidLarge", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./arrays/bounds/invalid_large.swa",
			ExpectedExecutionOutput: "Runtime Error: Index 100 out of bounds (0-2)",
			T:                       t,
		}

		assert.NoError(t, req.Compile())
		assert.Error(t, req.RunProgram())
	})

	t.Run("StructArrayBounds", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./arrays/bounds/struct_array_bounds.swa",
			ExpectedExecutionOutput: "Runtime Error: Index 5 out of bounds (0-1)",
			T:                       t,
		}

		assert.NoError(t, req.Compile())
		assert.Error(t, req.RunProgram())
	})
}
