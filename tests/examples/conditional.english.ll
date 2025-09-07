; ModuleID = 'swa-main'
source_filename = "swa-main"

@salary_string = global [11 x i8] c"100000 EUR\00"
@salary = global i32 500000

declare i32 @printf(ptr, ...)

define i32 @main() {
func-body:
  ret i32 0
}
