; ModuleID = 'swa-main'
source_filename = "swa-main"

%Engineer = type { i32, ptr, ptr }

declare i32 @printf(ptr, ...)

define i32 @main() {
func-body:
  %0 = alloca %Engineer, align 8
  %Age = getelementptr inbounds %Engineer, ptr %0, i32 0, i32 0
  store i32 40, ptr %Age, align 4
  %Name = getelementptr inbounds %Engineer, ptr %0, i32 0, i32 1
  store [6 x i8] c"Pathe\00", ptr %Name, align 1
  %TechStack = getelementptr inbounds %Engineer, ptr %0, i32 0, i32 2
  store [15 x i8] c"Ruby, Rust, Go\00", ptr %TechStack, align 1
  %1 = alloca [15 x i8], align 1
  store [15 x i8] c"Struct created\00", ptr %1, align 1
  %2 = alloca [29 x i8], align 1
  store [29 x i8] c"TODO display value of r1.Age\00", ptr %2, align 1
  %3 = call i32 (ptr, ...) @printf(ptr %1, ptr %2)
  ret i32 100
}
