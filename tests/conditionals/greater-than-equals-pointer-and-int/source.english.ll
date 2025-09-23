; ModuleID = 'swa-main'
source_filename = "swa-main"

@salaire_de_lucien = global i32 500000

declare i32 @printf(ptr, ...)

define i32 @main() {
func-body:
  %0 = load i32, ptr @salaire_de_lucien, align 4
  %1 = icmp uge i32 %0, 500000
  br i1 %1, label %if, label %else

if:                                               ; preds = %func-body
  %2 = alloca [3 x i8], align 1
  store [3 x i8] c"%s\00", ptr %2, align 1
  %3 = alloca [3 x i8], align 1
  store [3 x i8] c"ok\00", ptr %3, align 1
  %4 = call i32 (ptr, ...) @printf(ptr %2, ptr %3)
  br label %merge

merge:                                            ; preds = %else, %if
  %5 = load i32, ptr @salaire_de_lucien, align 4
  %6 = icmp uge i32 500000, %5
  br i1 %6, label %if2, label %else3

if2:                                              ; preds = %merge
  %7 = alloca [3 x i8], align 1
  store [3 x i8] c"%s\00", ptr %7, align 1
  %8 = alloca [3 x i8], align 1
  store [3 x i8] c"ok\00", ptr %8, align 1
  %9 = call i32 (ptr, ...) @printf(ptr %7, ptr %8)
  br label %merge1

merge1:                                           ; preds = %else3, %if2
  ret i32 0

else3:                                            ; preds = %merge
  %10 = alloca [3 x i8], align 1
  store [3 x i8] c"%s\00", ptr %10, align 1
  %11 = alloca [95 x i8], align 1
  store [95 x i8] c"Ast BinaryExpression GreaterThanEqualsImplementationError cond: (500000 >= salaire_de_lucien) \00", ptr %11, align 1
  %12 = call i32 (ptr, ...) @printf(ptr %10, ptr %11)
  br label %merge1

else:                                             ; preds = %func-body
  %13 = alloca [3 x i8], align 1
  store [3 x i8] c"%s\00", ptr %13, align 1
  %14 = alloca [94 x i8], align 1
  store [94 x i8] c"Ast BinaryExpression GreaterThanEqualsImplementationError cond: (salaire_de_lucien >= 500000)\00", ptr %14, align 1
  %15 = call i32 (ptr, ...) @printf(ptr %13, ptr %14)
  br label %merge
}
