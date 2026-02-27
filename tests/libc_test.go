package tests

import "testing"

func TestLIBCReadfile(t *testing.T) {
	t.Run("read hello world with `read`", func(t *testing.T) {
		NewSuccessfulCompileRequest(
			t, "./libc/read-hello-world.swa",
			"Hello World\n")
	})
	t.Run("read hello with `pread`", func(t *testing.T) {
		NewSuccessfulCompileRequest(
			t, "./libc/read-hello-world-as-bytes.swa",
			"Hello World\n")
	})
}
