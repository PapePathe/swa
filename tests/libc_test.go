package tests

import "testing"

func TestLIBCReadfile(t *testing.T) {
	t.Run("read hello world with `read`", func(t *testing.T) {
		t.Run("English", func(t *testing.T) {
			NewSuccessfulCompileRequest(
				t, "./libc/read-hello-world.swa",
				"Hello World\n")
		})
		t.Run("Igbo", func(t *testing.T) {
			NewSuccessfulCompileRequest(
				t, "./libc/read-hello-world.igbo.swa",
				"Hello World\n")
		})
	})
	t.Run("read hello with `pread`", func(t *testing.T) {
		t.Run("English", func(t *testing.T) {
			NewSuccessfulCompileRequest(
				t, "./libc/read-hello-world-as-bytes.swa",
				"Hello World\n")
		})
		t.Run("Igbo", func(t *testing.T) {
			NewSuccessfulCompileRequest(
				t, "./libc/read-hello-world-as-bytes.igbo.swa",
				"Hello World\n")
		})
	})

	t.Run("lib math", func(t *testing.T) {
		t.Run("sqrt", func(t *testing.T) {
			t.Run("English", func(t *testing.T) {
				NewSuccessfulCompileRequest(
					t, "./libc/sqrt.swa",
					"sqrt(16.0) = 4.000000\nsqrt(2.0) = 1.414214\nsqrt(-16.0) = -nan\n")
			})
			t.Run("Igbo", func(t *testing.T) {
				NewSuccessfulCompileRequest(
					t, "./libc/sqrt.igbo.swa",
					"sqrt(16.0) = 4.000000\nsqrt(2.0) = 1.414214\nsqrt(-16.0) = -nan\n")
			})
		})
	})
}
