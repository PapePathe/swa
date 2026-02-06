package tests

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	t.Run("json", func(t *testing.T) {
		expected, _ := os.ReadFile("./parser/1.english.json")
		req := CompileRequest{
			InputPath:      "./parser/1.english.swa",
			ExpectedOutput: string(expected),
			T:              t,
		}

		assert.NoError(t, req.Parse("json"))
	})

	t.Run("tree", func(t *testing.T) {
		expected, _ := os.ReadFile("./parser/1.english.tree")
		req := CompileRequest{
			InputPath:      "./parser/1.english.swa",
			ExpectedOutput: string(expected),
			T:              t,
		}

		assert.NoError(t, req.Parse("tree"))
	})
}
