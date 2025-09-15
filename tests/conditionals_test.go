package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConditionalsEnglish(t *testing.T) {
	t.Parallel()

	output, err := CompileSwaCode(t, "./examples/conditional.swa", "conditional.english")

	assert.NoError(t, err)
	assert.Equal(t, string(output), "")

	expectedIR := `; ModuleID = 'swa-main'
source_filename = "swa-main"

@salary_string = global [11 x i8] c"100000 EUR\00"
@salary = global i32 500000

declare i32 @printf(ptr, ...)

define i32 @main() {
func-body:
  ret i32 0
}
`

	assertCodeGenerated(t, "conditional.english")
	assertFileContent(t, "./conditional.english.ll", expectedIR)
	cleanupSwaCode(t, "conditional.english")
}

func TestConditionalsFrench(t *testing.T) {
	t.Parallel()

	output, err := CompileSwaCode(t, "./examples/conditional.french.swa", "conditional.french")

	assert.NoError(t, err)
	assert.Equal(t, string(output), "")

	expectedIR := `; ModuleID = 'swa-main'
source_filename = "swa-main"

@salary_string = global [11 x i8] c"100000 EUR\00"
@salary = global i32 500000
@new_salary = global i32 1009

declare i32 @printf(ptr, ...)

define i32 @main() {
func-body:
  store i32 999999, ptr @salary, align 4
  %0 = icmp uge ptr @salary, @new_salary
  br i1 %0, label %if, label %else

if:                                               ; preds = %func-body
  %1 = alloca [15 x i8], align 1
  store [15 x i8] c"%s %s %s %s %d\00", ptr %1, align 1
  %2 = alloca [22 x i8], align 1
  store [22 x i8] c"new salary is greater\00", ptr %2, align 1
  %3 = alloca [5 x i8], align 1
  store [5 x i8] c"than\00", ptr %3, align 1
  %4 = alloca [7 x i8], align 1
  store [7 x i8] c"salary\00", ptr %4, align 1
  %5 = call i32 (ptr, ...) @printf(ptr %1, ptr %2, ptr %3, ptr %4, ptr @salary_string, ptr @salary)
  ret i32 10

merge:                                            ; preds = %else
  ret i32 100

else:                                             ; preds = %func-body
  %6 = alloca [6 x i8], align 1
  store [6 x i8] c"%s %d\00", ptr %6, align 1
  %7 = alloca [39 x i8], align 1
  store [39 x i8] c"salary is greater ---- from else block\00", ptr %7, align 1
  %8 = call i32 (ptr, ...) @printf(ptr %6, ptr %7, ptr @new_salary)
  br label %merge
}
`

	assertCodeGenerated(t, "conditional.french")
	assertFileContent(t, "./conditional.french.ll", expectedIR)
	cleanupSwaCode(t, "conditional.french")
}
