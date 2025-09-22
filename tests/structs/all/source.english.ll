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
   ret i32 100
}
