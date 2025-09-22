ModuleID = 'swa-main'
 source_filename = "swa-main"

 %Vec = type { i32 }

 declare i32 @printf(ptr, ...)

 define i32 @main() {
 func-body:
   %0 = alloca %Vec, align 8
   %1 = getelementptr inbounds %Vec, ptr %0, i32 0, i32 0
   store i32 0, ptr %1, align 4
   %2 = getelementptr inbounds %Vec, ptr %0, i32 0, i32 0
   store i32 100, ptr %2, align 4
   ret i32 100
}
