package tests

import (
	"testing"
)

func TestMissingDialect(t *testing.T) {
	t.Parallel()

	req := CompileRequest{
		InputPath:               "./examples/missing_dialect.swa",
		OutputPath:              "55109834-f1de-425f-9e5a-562e8defd4d7",
		ExpectedExecutionOutput: "you must define your dialect\n",
		T:                       t,
	}

	defer req.Cleanup()

	req.AssertCompileAndExecute()
}
