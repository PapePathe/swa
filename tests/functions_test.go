package tests

import (
	"testing"
)

func TestFunctions(t *testing.T) {
	t.Parallel()

	t.Run("English", func(t *testing.T) {

		t.Run("Function taking struct as argument", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./functions/struct.source.english.swa",
				ExpectedExecutionOutput: "age: 40, height: 1.80, name: Pathe",
				T:                       t,
			}

			req.AssertCompileAndExecute()
		})

		t.Run("Substract", func(t *testing.T) {
			t.Run("Integer", func(t *testing.T) {
				req := CompileRequest{
					InputPath:               "./functions/substract.integer.english.swa",
					ExpectedExecutionOutput: "10 - 5 = 5",
					T:                       t,
				}

				req.AssertCompileAndExecute()
			})

			t.Run("Float", func(t *testing.T) {
				req := CompileRequest{
					InputPath:               "./functions/substract.float.english.swa",
					ExpectedExecutionOutput: "10.5 - 5.25 = 5.25",
					T:                       t,
				}

				req.AssertCompileAndExecute()
			})
		})

		t.Run("Add", func(t *testing.T) {
			t.Run("Float", func(t *testing.T) {
				req := CompileRequest{
					InputPath:               "./functions/add.float.english.swa",
					ExpectedExecutionOutput: "10.25 + 4.75 = 15.00",
					T:                       t,
				}

				req.AssertCompileAndExecute()
			})

			t.Run("Integer", func(t *testing.T) {
				req := CompileRequest{
					InputPath:               "./functions/add.integer.english.swa",
					ExpectedExecutionOutput: "10 + 5 = 15",
					T:                       t,
				}

				req.AssertCompileAndExecute()
			})
		})

		t.Run("Divide", func(t *testing.T) {
			t.Run("Integer", func(t *testing.T) {
				req := CompileRequest{
					InputPath:               "./functions/divide.integer.english.swa",
					ExpectedExecutionOutput: "10 / 5 = 2",
					T:                       t,
				}

				req.AssertCompileAndExecute()
			})

			t.Run("Float", func(t *testing.T) {
				req := CompileRequest{
					InputPath:               "./functions/divide.float.english.swa",
					ExpectedExecutionOutput: "10.0 / 5.0 = 2.000000",
					T:                       t,
				}

				req.AssertCompileAndExecute()
			})

		})

		t.Run("Multiply", func(t *testing.T) {
			t.Run("Float", func(t *testing.T) {
				req := CompileRequest{
					InputPath:               "./functions/multiply.float.english.swa",
					ExpectedExecutionOutput: "10.02, 5.05 = 50.601000",
					T:                       t,
				}

				req.AssertCompileAndExecute()
			})

			t.Run("Integer", func(t *testing.T) {
				req := CompileRequest{
					InputPath:               "./functions/multiply.integer.english.swa",
					ExpectedExecutionOutput: "10 * 5 = 50",
					T:                       t,
				}

				req.AssertCompileAndExecute()
			})
		})

	})

	t.Run("French", func(t *testing.T) {
		t.Run("Substract", func(t *testing.T) {
			t.Run("Integer", func(t *testing.T) {
				req := CompileRequest{
					InputPath:               "./functions/substract.integer.french.swa",
					ExpectedExecutionOutput: "10 - 5 = 5",
					T:                       t,
				}

				req.AssertCompileAndExecute()
			})

			t.Run("Float", func(t *testing.T) {
				req := CompileRequest{
					InputPath:               "./functions/substract.float.french.swa",
					ExpectedExecutionOutput: "10.5 - 5.25 = 5.25",
					T:                       t,
				}

				req.AssertCompileAndExecute()
			})
		})

		t.Run("Add", func(t *testing.T) {
			t.Run("Float", func(t *testing.T) {
				req := CompileRequest{
					InputPath:               "./functions/add.float.french.swa",
					ExpectedExecutionOutput: "10 + 5 = 15",
					T:                       t,
				}

				req.AssertCompileAndExecute()
			})
			t.Run("Integer", func(t *testing.T) {
				req := CompileRequest{
					InputPath:               "./functions/add.integer.french.swa",
					ExpectedExecutionOutput: "10 + 5 = 15",
					T:                       t,
				}

				req.AssertCompileAndExecute()
			})
		})

		t.Run("Divide", func(t *testing.T) {
			t.Run("Float", func(t *testing.T) {
				req := CompileRequest{
					InputPath:               "./functions/divide.float.french.swa",
					ExpectedExecutionOutput: "10.0 / 5.0 = 2.000000",
					T:                       t,
				}

				req.AssertCompileAndExecute()
			})
			t.Run("Integer", func(t *testing.T) {
				req := CompileRequest{
					InputPath:               "./functions/divide.integer.french.swa",
					ExpectedExecutionOutput: "10 / 5 = 2",
					T:                       t,
				}

				req.AssertCompileAndExecute()
			})
		})

		t.Run("Multiply", func(t *testing.T) {
			t.Run("Float", func(t *testing.T) {
				req := CompileRequest{
					InputPath:               "./functions/multiply.float.french.swa",
					ExpectedExecutionOutput: "10.02, 5.05 = 50.601000",
					T:                       t,
				}

				req.AssertCompileAndExecute()
			})

			t.Run("Integer", func(t *testing.T) {
				req := CompileRequest{
					InputPath:               "./functions/multiply.integer.french.swa",
					ExpectedExecutionOutput: "10 * 5 = 50",
					T:                       t,
				}

				req.AssertCompileAndExecute()
			})
		})
	})
}
