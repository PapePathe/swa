package tests

import (
	"testing"
)

func TestMissingDialect(t *testing.T) {
	NewFailedCompileRequest(t,
		"./examples/missing_dialect.swa",
		"you must define your dialect\n")
}
