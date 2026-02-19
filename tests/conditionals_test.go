package tests

import (
	"testing"
)

func TestEarlyReturnStatements(t *testing.T) {
	t.Run("English", func(t *testing.T) {
		t.Run("1", func(t *testing.T) {
			t.Run("Condition is true", func(t *testing.T) {
				req := CompileRequest{
					InputPath:               "./conditionals/early-return/1-condition-true.english.swa",
					ExpectedExecutionOutput: "program interrupted by early return",
					T:                       t,
				}

				req.AssertCompileAndExecute()
			})
			//	t.Run("Condition is true.X", func(t *testing.T) {
			//		req := CompileRequest{
			//			InputPath: "./conditionals/early-return/1-condition-true.english.swa",
			//			T:         t,
			//		}

			//		req.AssertCompileAndExecuteX()
			//	})

			t.Run("Condition is false", func(t *testing.T) {
				req := CompileRequest{
					InputPath:               "./conditionals/early-return/1-condition-false.english.swa",
					ExpectedExecutionOutput: "program not interrupted by early return",
					T:                       t,
				}

				req.AssertCompileAndExecute()
			})

			//	t.Run("Condition is false.X", func(t *testing.T) {
			//		req := CompileRequest{
			//			InputPath: "./conditionals/early-return/1-condition-false.english.swa",
			//			T:         t,
			//		}

			//		req.AssertCompileAndExecuteX()
			//	})
		})
	})

	t.Run("Wolof", func(t *testing.T) {
		t.Run("1", func(t *testing.T) {
			t.Run("Condition is true", func(t *testing.T) {
				req := CompileRequest{
					InputPath:               "./conditionals/early-return/1-condition-true.wolof.swa",
					ExpectedExecutionOutput: "program interrupted by early return",
					T:                       t,
				}

				req.AssertCompileAndExecute()
			})
			//	t.Run("Condition is true.X", func(t *testing.T) {
			//		req := CompileRequest{
			//			InputPath: "./conditionals/early-return/1-condition-true.wolof.swa",
			//			T:         t,
			//		}

			//		req.AssertCompileAndExecuteX()
			//	})

			t.Run("Condition is false", func(t *testing.T) {
				req := CompileRequest{
					InputPath:               "./conditionals/early-return/1-condition-false.wolof.swa",
					ExpectedExecutionOutput: "program not interrupted by early return",
					T:                       t,
				}

				req.AssertCompileAndExecute()
			})
			//	t.Run("Condition is false.X", func(t *testing.T) {
			//		req := CompileRequest{
			//			InputPath: "./conditionals/early-return/1-condition-false.wolof.swa",
			//			T:         t,
			//		}

			//		req.AssertCompileAndExecuteX()
			//	})
		})
	})

	t.Run("French", func(t *testing.T) {
		t.Run("1", func(t *testing.T) {
			t.Run("Condition is true", func(t *testing.T) {
				req := CompileRequest{
					InputPath:               "./conditionals/early-return/1-condition-true.french.swa",
					ExpectedExecutionOutput: "program interrupted by early retourner",
					T:                       t,
				}

				req.AssertCompileAndExecute()
			})
			t.Run("Condition is false", func(t *testing.T) {
				req := CompileRequest{
					InputPath:               "./conditionals/early-return/1-condition-false.french.swa",
					ExpectedExecutionOutput: "program not interrupted by early retourner",
					T:                       t,
				}

				req.AssertCompileAndExecute()
			})
		})
	})

	t.Run("Soussou", func(t *testing.T) {
		t.Run("1", func(t *testing.T) {
			t.Run("Condition is true", func(t *testing.T) {
				req := CompileRequest{
					InputPath:               "./conditionals/early-return/1-condition-true.soussou.swa",
					ExpectedExecutionOutput: "program interrupted by early return",
					T:                       t,
				}

				req.AssertCompileAndExecute()
			})
			t.Run("Condition is false", func(t *testing.T) {
				req := CompileRequest{
					InputPath:               "./conditionals/early-return/1-condition-false.soussou.swa",
					ExpectedExecutionOutput: "program not interrupted by early return",
					T:                       t,
				}

				req.AssertCompileAndExecute()
			})
		})
	})
}

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

	t.Run("English.X", func(t *testing.T) {
		req := CompileRequest{
			InputPath: "./conditionals/greater-than-equals/source.english.swa",
			T:         t,
		}

		req.AssertCompileAndExecuteX()
	})

	t.Run("French", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./conditionals/greater-than-equals/source.french.swa",
			ExpectedExecutionOutput: "okok",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})

	t.Run("French.X", func(t *testing.T) {
		req := CompileRequest{
			InputPath: "./conditionals/greater-than-equals/source.french.swa",
			T:         t,
		}

		req.AssertCompileAndExecuteX()
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

	t.Run("English.X", func(t *testing.T) {
		req := CompileRequest{
			InputPath: "./conditionals/greater-than-equals-pointer-and-int/source.english.swa",
			T:         t,
		}

		req.AssertCompileAndExecuteX()
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

	t.Run("English.X", func(t *testing.T) {
		req := CompileRequest{
			InputPath: "./conditionals/less-than-equals/source.english.swa",
			T:         t,
		}

		req.AssertCompileAndExecuteX()
	})

	t.Run("French", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./conditionals/less-than-equals/source.french.swa",
			ExpectedExecutionOutput: "okokokokokokok",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})

	t.Run("French.X", func(t *testing.T) {
		req := CompileRequest{
			InputPath: "./conditionals/less-than-equals/source.french.swa",
			T:         t,
		}

		req.AssertCompileAndExecuteX()
	})
}

func TestEquals(t *testing.T) {
	t.Parallel()

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./conditionals/equals/source.english.swa",
			ExpectedExecutionOutput: "okokokokokok",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})

	t.Run("English.X", func(t *testing.T) {
		req := CompileRequest{
			InputPath: "./conditionals/equals/source.english.swa",
			T:         t,
		}

		req.AssertCompileAndExecuteX()
	})

	t.Run("French", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./conditionals/equals/source.french.swa",
			ExpectedExecutionOutput: "okokokokokok",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})

	t.Run("French.X", func(t *testing.T) {
		req := CompileRequest{
			InputPath: "./conditionals/equals/source.french.swa",
			T:         t,
		}

		req.AssertCompileAndExecuteX()
	})
}
