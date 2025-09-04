; ModuleID = 'my_module'
source_filename = "my_module"

@salary_string = global [11 x i8] c"100000 EUR\00"
@salary = global i32 500000
@new_salary = global i32 1009

declare i32 @printf(ptr, ...)

define i32 @main() {
func-body:
  br i1 true, label %if, label %else

if:                                               ; preds = %func-body
  %0 = alloca [15 x i8], align 1
  store [15 x i8] c"%s %s %s %s %d\00", ptr %0, align 1
  %1 = alloca [22 x i8], align 1
  store [22 x i8] c"new salary is greater\00", ptr %1, align 1
  %2 = alloca [5 x i8], align 1
  store [5 x i8] c"than\00", ptr %2, align 1
  %3 = alloca [7 x i8], align 1
  store [7 x i8] c"salary\00", ptr %3, align 1
  %4 = call i32 (ptr, ...) @printf(ptr %0, ptr %1, ptr %2, ptr %3, ptr @salary_string, ptr @salary)
  ret i32 10

merge:                                            ; preds = %else
  ret i32 100

else:                                             ; preds = %func-body
  %5 = alloca [6 x i8], align 1
  store [6 x i8] c"%s %d\00", ptr %5, align 1
  %6 = alloca [39 x i8], align 1
  store [39 x i8] c"salary is greater ---- from else block\00", ptr %6, align 1
  %7 = call i32 (ptr, ...) @printf(ptr %5, ptr %6, ptr @new_salary)
  br label %merge
}
