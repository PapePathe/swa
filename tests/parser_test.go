package tests

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseNiceNumbers(t *testing.T) {
	t.Run("French", func(t *testing.T) {
		t.Run("json", func(t *testing.T) {
			expected, _ := os.ReadFile("./parser/2.french.json")
			req := CompileRequest{
				InputPath:      "./parser/2.french.swa",
				ExpectedOutput: string(expected),
				T:              t,
			}

			assert.NoError(t, req.Parse("json"))
		})

		t.Run("tree", func(t *testing.T) {
			expected, _ := os.ReadFile("./parser/2.french.tree")
			req := CompileRequest{
				InputPath:      "./parser/2.french.swa",
				ExpectedOutput: string(expected),
				T:              t,
			}

			assert.NoError(t, req.Parse("tree"))
		})
	})
}

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
