package tests

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMissingFile(t *testing.T) {
	t.Parallel()

	cmd := exec.Command("./swa", "compile", "-s", "./examples/not-found.swa")
	_, err := cmd.CombinedOutput()

	assert.Equal(t, err.Error(), "exit status 1")
}

func TestMissingDialect(t *testing.T) {
	t.Parallel()

	cmd := exec.Command("./swa", "compile", "-s", "./examples/missing_dialect.swa")
	_, err := cmd.CombinedOutput()

	assert.Equal(t, err.Error(), "exit status 1")
}

func TestConditionals(t *testing.T) {
	t.Parallel()

	cmd := exec.Command("./swa", "compile", "-s", "./examples/conditional.swa")
	output, err := cmd.CombinedOutput()

	assert.NoError(t, err)
	assert.Equal(t, string(output), "")

	assertFileExists(t, "./start.ll")
	assertFileExists(t, "./start.s")
	assertFileExists(t, "./start.o")
	assertFileExists(t, "./start.exe")
}

func assertFileExists(t *testing.T, path string) {
	_, err := os.Stat(path)

	assert.True(t, !os.IsNotExist(err))
}
