package tests

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {

	t.Run("English", func(t *testing.T) {
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
	})

	t.Run("Wolof", func(t *testing.T) {
		t.Run("json", func(t *testing.T) {
			expected, _ := os.ReadFile("./parser/1.wolof.json")
			req := CompileRequest{
				InputPath:      "./parser/1.wolof.swa",
				ExpectedOutput: string(expected),
				T:              t,
			}

			assert.NoError(t, req.Parse("json"))
		})

		t.Run("tree", func(t *testing.T) {
			expected, _ := os.ReadFile("./parser/1.wolof.tree")
			req := CompileRequest{
				InputPath:      "./parser/1.wolof.swa",
				ExpectedOutput: string(expected),
				T:              t,
			}

			assert.NoError(t, req.Parse("tree"))
		})
	})

	t.Run("French", func(t *testing.T) {
		t.Run("json", func(t *testing.T) {
			expected, _ := os.ReadFile("./parser/1.french.json")
			req := CompileRequest{
				InputPath:      "./parser/1.french.swa",
				ExpectedOutput: string(expected),
				T:              t,
			}

			assert.NoError(t, req.Parse("json"))
		})

		t.Run("tree", func(t *testing.T) {
			expected, _ := os.ReadFile("./parser/1.french.tree")
			req := CompileRequest{
				InputPath:      "./parser/1.french.swa",
				ExpectedOutput: string(expected),
				T:              t,
			}

			assert.NoError(t, req.Parse("tree"))
		})
	})

	t.Run("Soussou", func(t *testing.T) {
		t.Run("json", func(t *testing.T) {
			expected, _ := os.ReadFile("./parser/1.soussou.json")
			req := CompileRequest{
				InputPath:      "./parser/1.soussou.swa",
				ExpectedOutput: string(expected),
				T:              t,
			}

			assert.NoError(t, req.Parse("json"))
		})

		t.Run("tree", func(t *testing.T) {
			expected, _ := os.ReadFile("./parser/1.soussou.tree")
			req := CompileRequest{
				InputPath:      "./parser/1.soussou.swa",
				ExpectedOutput: string(expected),
				T:              t,
			}

			assert.NoError(t, req.Parse("tree"))
		})
	})
}
