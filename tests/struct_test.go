package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStructsEnglish(t *testing.T) {
	t.Parallel()
	expectedIR := `; ModuleID = 'swa-main'
source_filename = "swa-main"

%Engineer = type { i32, ptr, ptr }

declare i32 @printf(ptr, ...)

define i32 @main() {
func-body:
  %0 = alloca %Engineer, align 8
  %1 = getelementptr inbounds %Engineer, ptr %0, i32 0, i32 0
  store i32 40, ptr %1, align 4
  %2 = getelementptr inbounds %Engineer, ptr %0, i32 0, i32 1
  store [6 x i8] c"Pathe\00", ptr %2, align 1
  %3 = getelementptr inbounds %Engineer, ptr %0, i32 0, i32 2
  store [15 x i8] c"Ruby, Rust, Go\00", ptr %3, align 1
  ret i32 100
}
`
	output, err := CompileSwaCode(t, "./examples/struct.english.swa", "struct.english")

	assert.NoError(t, err)
	assert.Equal(t, string(output), "")

	assertCodeGenerated(t, "struct.english")
	assertFileContent(t, "./struct.english.ll", expectedIR)
	cleanupSwaCode(t, "struct.english")
}

func TestStructsFrench(t *testing.T) {
	t.Parallel()
	expectedIR := `; ModuleID = 'swa-main'
source_filename = "swa-main"

%Engineer = type { i32, ptr, ptr }

declare i32 @printf(ptr, ...)

define i32 @main() {
func-body:
  %0 = alloca %Engineer, align 8
  %1 = getelementptr inbounds %Engineer, ptr %0, i32 0, i32 0
  store i32 40, ptr %1, align 4
  %2 = getelementptr inbounds %Engineer, ptr %0, i32 0, i32 1
  store [6 x i8] c"Pathe\00", ptr %2, align 1
  %3 = getelementptr inbounds %Engineer, ptr %0, i32 0, i32 2
  store [15 x i8] c"Ruby, Rust, Go\00", ptr %3, align 1
  %4 = getelementptr inbounds %Engineer, ptr %0, i32 0, i32 0
  %5 = alloca [15 x i8], align 1
  store [15 x i8] c"Struct created\00", ptr %5, align 1
  %6 = alloca [29 x i8], align 1
  store [29 x i8] c"TODO display value of r1.Age\00", ptr %6, align 1
  %7 = call i32 (ptr, ...) @printf(ptr %5, ptr %6)
  ret i32 100
}
`
	output, err := CompileSwaCode(t, "./examples/struct.french.swa", "struct.french")

	assert.NoError(t, err)
	assert.Equal(t, string(output), "")

	assertCodeGenerated(t, "struct.french")
	assertFileContent(t, "./struct.french.ll", expectedIR)
	cleanupSwaCode(t, "struct.french")
}

func TestStructWithUnknownType(t *testing.T) {
	t.Parallel()

	sourceCode := `dialect:english;
start() int {
	struct Test {
		Name: chaine_de_caractere,
	}
  return 100;
} `
	output, err := CompileSwaSourceCode(
		t,
		"./examples/struct.unknown-data-type.swa",
		"struct.unknown-data-type",
		[]byte(sourceCode),
	)
	expectedOutput := "struct proprerty type ({chaine_de_caractere}) not supported\n"

	assert.Error(t, err)
	assert.Equal(t, string(output), expectedOutput)
}

func TestStructPropertyAssignment(t *testing.T) {
	t.Parallel()

	expectedIR := `; ModuleID = 'swa-main'
source_filename = "swa-main"

%Vec = type { i32 }

declare i32 @printf(ptr, ...)

define i32 @main() {
func-body:
  %0 = alloca %Vec, align 8
  %1 = getelementptr inbounds %Vec, ptr %0, i32 0, i32 0
  store i32 0, ptr %1, align 4
  %2 = getelementptr inbounds %Vec, ptr %0, i32 0, i32 0
  store i32 100, ptr %2, align 4
  ret i32 100
}
`
	sourceCode := `dialect:english;	
start() int {
	struct Vec {
		x: int,
	}
	let v:Vec = Vec{x: 0};
	v.x = 100;
  return 100;
} 
`
	output, err := CompileSwaSourceCode(
		t,
		"./examples/struct.property-assignment.swa",
		"struct.property-assignment",
		[]byte(sourceCode),
	)
	expectedOutput := ""

	assert.NoError(t, err)
	assert.Equal(t, string(output), expectedOutput)

	assertFileContent(t, "./struct.property-assignment.ll", expectedIR)
	cleanupSwaCode(t, "./struct.property-assignment")
}

func TestStructPropertyInReturnExpression(t *testing.T) {
	t.Parallel()

	sourceCode := `dialect:french;	
demarrer() entier {
	structure Test { 
		Message: chaine, 
		ExitCode: entier, 
	}
	variable ab1: Test = Test{ 
		Message: "Done", 
  	ExitCode: 3
  }; 

  retourner ab1.ExitCode;
}	
`
	output, err := CompileSwaSourceCode(
		t,
		"./examples/struct.property-access-in-return.swa",
		"struct.property-access-in-return.swa",
		[]byte(sourceCode),
	)
	expectedOutput := ""
	expectedIR := `; ModuleID = 'swa-main'
source_filename = "swa-main"

%Test = type { ptr, i32 }

declare i32 @printf(ptr, ...)

define i32 @main() {
func-body:
  %0 = alloca %Test, align 8
  %1 = getelementptr inbounds %Test, ptr %0, i32 0, i32 0
  store [5 x i8] c"Done\00", ptr %1, align 1
  %2 = getelementptr inbounds %Test, ptr %0, i32 0, i32 1
  store i32 3, ptr %2, align 4
  %3 = getelementptr inbounds %Test, ptr %0, i32 0, i32 1
  %4 = load i32, ptr %3, align 4
  ret i32 %4
}
`

	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, string(output))

	assertFileContent(t, "./struct.property-access-in-return.swa.ll", expectedIR)

	cleanupSwaCode(t, "./struct.property-access-in-return.swa")
}
