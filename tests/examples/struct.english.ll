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
  ret i32 0
}
