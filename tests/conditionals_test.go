package tests

import (
	"testing"
)

func TestGreaterThanEquals(t *testing.T) {
	t.Parallel()

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./conditionals/greater-than-equals/source.english.swa",
			ExpectedExecutionOutput: "okok",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})

	t.Run("French", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./conditionals/greater-than-equals/source.french.swa",
			ExpectedExecutionOutput: "okok",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})
}

func TestGreaterThanEqualsWithPointerAndInt(t *testing.T) {
	t.Parallel()

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./conditionals/greater-than-equals-pointer-and-int/source.english.swa",
			ExpectedExecutionOutput: "okok",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})

	t.Run("French", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./conditionals/greater-than-equals-pointer-and-int/source.french.swa",
			ExpectedExecutionOutput: "okok",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})
}

func TestLessThanEquals(t *testing.T) {
	t.Parallel()

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./conditionals/less-than-equals/source.english.swa",
			ExpectedExecutionOutput: "okokokokokokokok",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})

	t.Run("French", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./conditionals/less-than-equals/source.french.swa",
			ExpectedExecutionOutput: "okokokokokokok",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})
}

func TestEquals(t *testing.T) {
	t.Parallel()

	t.Run("Soussou", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./conditionals/equals/source.soussou.swa",
			ExpectedExecutionOutput: "okokokokokok",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./conditionals/equals/source.english.swa",
			ExpectedExecutionOutput: "okokokokokok",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})

	t.Run("French", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./conditionals/equals/source.french.swa",
			ExpectedExecutionOutput: "okokokokokok",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})
}
