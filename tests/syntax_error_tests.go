package tests

import (
	"testing"
)

func TestMissingDialect(t *testing.T) {
	t.Parallel()

	req := CompileRequest{
		InputPath:               "./examples/missing_dialect.swa",
		ExpectedExecutionOutput: "you must define your dialect\n",
		T:                       t,
	}

	defer req.Cleanup()

	req.AssertCompileAndExecute()
}
