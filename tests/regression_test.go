package tests

import (
	"os"
	"testing"
)

func TestCollatz(t *testing.T) {
	NewSuccessfulCompileRequest(
		t,
		"./regression/collatz.swa",
		"Testing Collatz Conjecture...\nCollatz(27) took 111 steps.\nSuccess: Collatz(27) is correct.\n",
	)
}

func TestBrainFuck(t *testing.T) {
	NewSuccessfulCompileRequest(
		t,
		"./regression/brainfuck.swa",
		"Testing Brainfuck Simulation...\nSimulating BF: +++[->+<]\nResult in cell 1: 3\nSuccess: BF simulation worked.\n",
	)
}

func TestAckerman(t *testing.T) {
	NewSuccessfulCompileRequest(
		t,
		"./regression/ackermann.swa",
		"Testing Ackermann function...\nackermann(3, 2) = 29\nSuccess: Ackermann A(3, 2) is correct.\n",
	)
}

func TestGenetic(t *testing.T) {
	NewSuccessfulCompileRequest(
		t,
		"./regression/genetic.swa",
		"",
	)
}

func TestIntegrales(t *testing.T) {
	NewSuccessfulCompileRequest(
		t,
		"./regression/integrales.swa",
		"",
	)
}

func TestRsa(t *testing.T) {
	NewSuccessfulCompileRequest(t, "./regression/rsa.swa", "")
}

func TestParticles(t *testing.T) {
	NewSuccessfulCompileRequest(t, "./regression/particles.swa", "")
}

func TestStack(t *testing.T) {
	NewSuccessfulCompileRequest(t, "./regression/stack.swa", "")
}

func TestRipProtocol(t *testing.T) {
	NewSuccessfulCompileRequest(t, "./regression/rip-protocol.swa", "")
}

func TestProjectile(t *testing.T) {
	NewSuccessfulCompileRequest(t, "./regression/projectile.swa", "")
}

func TestOscillateur(t *testing.T) {
	NewSuccessfulCompileRequest(t, "./regression/oscillateur.swa", "")
}

func TestOrbital(t *testing.T) {
	NewSuccessfulCompileRequest(t, "./regression/orbital.swa", "")
}

func TestNeurons(t *testing.T) {
	NewSuccessfulCompileRequest(t, "./regression/neurones.swa", "")
}

func TestMutualRecursion(t *testing.T) {
	NewSuccessfulCompileRequest(t, "./regression/mutual-recursion.swa", "")
}

func TestTextAnalysis(t *testing.T) {
	NewSuccessfulCompileRequest(t, "./regression/matrices.swa", "")
}

func TestMatrices(t *testing.T) {
	NewSuccessfulCompileRequest(t, "./regression/matrices.swa", "")
}

func TestMatrices2(t *testing.T) {
	NewSuccessfulCompileRequest(t, "./regression/matrices.swa", "")
}

func TestIsotopes(t *testing.T) {
	NewSuccessfulCompileRequest(t, "./regression/isotopes.swa", "")
}

func TestGravitation(t *testing.T) {
	NewSuccessfulCompileRequest(t, "./regression/gravitation.swa", "")
}

func TestGaz(t *testing.T) {
	NewSuccessfulCompileRequest(t, "./regression/gaz.swa", "")
}

func TestDiffie(t *testing.T) {
	NewSuccessfulCompileRequest(t, "./regression/diffie-hellman.swa", "")
}

func TestCircuit(t *testing.T) {
	NewSuccessfulCompileRequest(t, "./regression/circuit.swa", "")

}

func TestChiffrement(t *testing.T) {
	NewSuccessfulCompileRequest(t, "./regression/chiffrement.swa", "")
}

func TestChiffrementCesar(t *testing.T) {
	NewSuccessfulCompileRequest(t, "./regression/chiffrement-cesar.swa", "")
}

func TestBibliotheque(t *testing.T) {
	NewSuccessfulCompileRequest(t, "./regression/bibliotheque.swa", "")
}

func TestArpProtocol(t *testing.T) {
	NewSuccessfulCompileRequest(t, "./regression/arp-protocol.swa", "")
}

func TestArithmtic(t *testing.T) {
	NewSuccessfulCompileRequest(t, "./regression/arithmetic.swa", "")
}

func TcpPacketCheckSum(t *testing.T) {
	NewSuccessfulCompileRequest(t, "./regression/tcp-packet-checksum.swa", "")
}

func TcpHandshake(t *testing.T) {
	NewSuccessfulCompileRequest(t, "./regression/tcp-handshake.swa", "")
}

func TestStringEdgeCases(t *testing.T) {
	NewSuccessfulCompileRequest(t, "./regression/string_edge_cases.swa", "s1: Hello\nWorld\ns2: Tab\tCharacter\ns3: Quote \" inside\ns4: Backslash \\ inside\nLong string: This is a very long string designed to test how the compiler handles larger literals in the assembly generation phase. It should be able to handle strings that are significantly longer than short identifiers or simple messages.\n")
}

func TestMultiDimensionalArrayAccess(t *testing.T) {
	NewFailedCompileRequest(t,
		"./regression/multidim_array.swa",
		"TODO: nested array access not yet supported\n")
}

func TestLargeBoolean(t *testing.T) {
	NewSuccessfulCompileRequest(t,
		"./regression/large_boolean.swa",
		"Success: Complex boolean logic evaluated correctly.\n")
}

func TestSymbolValue(t *testing.T) {
	NewFailedCompileRequest(t,
		"./regression/symbol-value.swa",
		"TODO implement VisitSymbolValueExpression\n")
}

func TestSymbolAddress(t *testing.T) {
	NewFailedCompileRequest(t,
		"./regression/symbol-address.swa",
		"TODO implement VisitSymbolAdressExpression\n")
}

func TestDoublePointer(t *testing.T) {
	t.Run("french", func(t *testing.T) {
		t.Run("compile", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./regression/double_pointer.french.swa",
				"TODO implement VisitSymbolAdressExpression\n")
		})

		t.Run("parse as json", func(t *testing.T) {
			expected, _ := os.ReadFile("./regression/double_pointer.french.json")
			NewSuccessfulCompileRequest(t,
				"./regression/double_pointer.french.swa",
				string(expected))
		})

		t.Run("parse as tree", func(t *testing.T) {
			expected, _ := os.ReadFile("./regression/double_pointer.french.tree")
			NewSuccessfulCompileRequest(t,
				"./regression/double_pointer.french.swa",
				string(expected))
		})
	})

	t.Run("english", func(t *testing.T) {
		NewFailedCompileRequest(t,
			"./regression/double_pointer.swa",
			"TODO implement VisitSymbolAdressExpression\n")
	})

	t.Run("soussou", func(t *testing.T) {
		NewFailedCompileRequest(t,
			"./regression/double_pointer.soussou.swa",
			"TODO implement VisitSymbolAdressExpression\n")
	})
}

func TestFloatingBlockStatement(t *testing.T) {
	t.Run("english", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./regression/shadowing.swa",
			"Outer x: 10\nInner x: 20\nOuter x again: 10\nSuccess: Variable shadowing works.\n")

		t.Run("parse as json", func(t *testing.T) {
			expected, _ := os.ReadFile("./regression/shadowing.json")
			NewSuccessfulCompileRequest(t,
				"./regression/shadowing.swa",
				string(expected))
		})

		t.Run("parse as tree", func(t *testing.T) {
			expected, _ := os.ReadFile("./regression/shadowing.tree")
			NewSuccessfulCompileRequest(t,
				"./regression/shadowing.swa",
				string(expected))
		})
	})

	t.Run("soussou", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./regression/shadowing.soussou.swa",
			"Outer x: 10\nInner x: 20\nOuter x again: 10\nSuccess: Variable shadowing works.\n")
	})

	t.Run("french", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./regression/shadowing.french.swa",
			"Outer x: 10\nInner x: 20\nOuter x again: 10\nSuccess: Variable shadowing works.\n")
	})
}
