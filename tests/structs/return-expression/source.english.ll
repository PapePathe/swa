; ModuleID = 'swa-main'
source_filename = "swa-main"

%Test = type { ptr, i32 }

declare i32 @printf(ptr, ...)

define i32 @main() {
func-body:
  %0 = alloca %Test, align 8
  %1 = getelementptr inbounds %Test, ptr %0, i32 0, i32 0
  store [5 x i8] c"Done\00", ptr %1, align 1
  %2 = getelementptr inbounds %Test, ptr %0, i32 0, i32 1
  store i32 0, ptr %2, align 4
  %3 = getelementptr inbounds %Test, ptr %0, i32 0, i32 1
  %4 = load i32, ptr %3, align 4
  ret i32 %4
}
