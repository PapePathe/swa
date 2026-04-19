package tests

import (
	"testing"
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
				NewFailedCompileRequest(t, test.inputPath, test.expectedOutput)
			})
		}
	})

	t.Run("Many main functions", func(t *testing.T) {
		tests := []MultiDialectTest{
			{
				name:           "French",
				inputPath:      "./program_entrypoint/many-main.french.swa",
				expectedOutput: "function named main already exists in symbol table\n",
			},
			{
				name:           "English",
				inputPath:      "./program_entrypoint/many-main.english.swa",
				expectedOutput: "function named main already exists in symbol table\n",
			},

			{
				name:           "Wolof",
				inputPath:      "./program_entrypoint/many-main.wolof.swa",
				expectedOutput: "function named main already exists in symbol table\n",
			},
			{
				name:           "Soussou",
				inputPath:      "./program_entrypoint/many-main.soussou.swa",
				expectedOutput: "wali main na na yi khorun\n",
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				NewFailedCompileRequest(t, test.inputPath, test.expectedOutput)
			})
		}
	})
}
