; ModuleID = 'swa-main'
source_filename = "swa-main"

@salaire_de_lucien = global i32 500000
@salaire_de_pathe = global i32 50000

declare i32 @printf(ptr, ...)

define i32 @main() {
func-body:
  %0 = load ptr, ptr @salaire_de_lucien, align 8
  %1 = load ptr, ptr @salaire_de_pathe, align 8
  %2 = icmp uge ptr %0, %1
  br i1 %2, label %if, label %else

if:                                               ; preds = %func-body
  %3 = alloca [3 x i8], align 1
  store [3 x i8] c"%s\00", ptr %3, align 1
  %4 = alloca [3 x i8], align 1
  store [3 x i8] c"ok\00", ptr %4, align 1
  %5 = call i32 (ptr, ...) @printf(ptr %3, ptr %4)
  br label %merge

merge:                                            ; preds = %else, %if
  %6 = load ptr, ptr @salaire_de_pathe, align 8
  %7 = load ptr, ptr @salaire_de_lucien, align 8
  %8 = icmp uge ptr %6, %7
  br i1 %8, label %if2, label %else3

if2:                                              ; preds = %merge
  %9 = alloca [3 x i8], align 1
  store [3 x i8] c"%s\00", ptr %9, align 1
  %10 = alloca [104 x i8], align 1
  store [104 x i8] c"Ast BinaryExpression GreaterThanEqualsImplementationError cond: (salaire_de_pathe >= salaire_de_lucien)\00", ptr %10, align 1
  %11 = call i32 (ptr, ...) @printf(ptr %9, ptr %10)
  br label %merge1

merge1:                                           ; preds = %else3, %if2
  ret i32 0

else3:                                            ; preds = %merge
  %12 = alloca [3 x i8], align 1
  store [3 x i8] c"%s\00", ptr %12, align 1
  %13 = alloca [3 x i8], align 1
  store [3 x i8] c"ok\00", ptr %13, align 1
  %14 = call i32 (ptr, ...) @printf(ptr %12, ptr %13)
  br label %merge1

else:                                             ; preds = %func-body
  %15 = alloca [3 x i8], align 1
  store [3 x i8] c"%s\00", ptr %15, align 1
  %16 = alloca [58 x i8], align 1
  store [58 x i8] c"Ast BinaryExpression GreaterThanEqualsImplementationError\00", ptr %16, align 1
  %17 = call i32 (ptr, ...) @printf(ptr %15, ptr %16)
  br label %merge
}
