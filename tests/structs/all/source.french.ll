; ModuleID = 'swa-main'

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
