package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProgramEntryPoint(t *testing.T) {
	t.Run("No main function", func(t *testing.T) {
		tests := []MultiDialectTest{
			{
				name:           "French",
				inputPath:      "./program_entrypoint/no-main.french.swa",
				expectedOutput: "Vous devez definir le programme principal\n",
			},
			{
				name:           "English",
				inputPath:      "./program_entrypoint/no-main.english.swa",
				expectedOutput: "Your program is missing a main function\n",
			},

			{
				name:           "Wolof",
				inputPath:      "./program_entrypoint/no-main.wolof.swa",
				expectedOutput: "Your program is missing a main function\n",
			},
			{
				name:           "Soussou",
				inputPath:      "./program_entrypoint/no-main.soussou.swa",
				expectedOutput: "Your program is missing a main function\n",
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				req := CompileRequest{
					T:              t,
					InputPath:      test.inputPath,
					ExpectedOutput: test.expectedOutput,
				}

				assert.Error(t, req.Compile())
			})
		}
	})

	t.Run("Many main functions", func(t *testing.T) {
		tests := []MultiDialectTest{
			{
				name:           "French",
				inputPath:      "./program_entrypoint/many-main.french.swa",
				expectedOutput: "Vous devez definir un seul programme principal, nombre (2)\n",
			},
			{
				name:           "English",
				inputPath:      "./program_entrypoint/many-main.english.swa",
				expectedOutput: "Your program must have exactly one main function, count (2)\n",
			},

			{
				name:           "Wolof",
				inputPath:      "./program_entrypoint/many-main.wolof.swa",
				expectedOutput: "Your program must have exactly one main function, count (2)\n",
			},
			{
				name:           "Soussou",
				inputPath:      "./program_entrypoint/many-main.soussou.swa",
				expectedOutput: "Your program must have exactly one main function, count (2)\n",
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				req := CompileRequest{
					T:              t,
					InputPath:      test.inputPath,
					ExpectedOutput: test.expectedOutput,
				}

				assert.Error(t, req.Compile())
			})
		}
	})
}
