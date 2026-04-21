; ModuleID = 'test.swa'
source_filename = "test.swa"

%Outer = type { { ptr, i32, i32 } }
%Inner = type { { ptr, i32, i32 } }

declare i32 @printf(ptr, ...)

declare i32 @strcmp(ptr, ptr)

declare void @exit(i64)

declare i32 @open(ptr, i32, ...)

declare i64 @read(i32, ptr, i64)

declare i32 @close(i32)

declare i64 @write(i32, ptr, i64)

declare i64 @pread(i32, ptr, i32, i64)

declare i64 @pwrite(i32, ptr, i32, i64)

declare double @sqrt(double)

declare ptr @realloc(ptr, i64)

declare void @free(ptr)

define i32 @main() {
entry:
  %Outer.instance = alloca %Outer, align 8
  %slice_data = tail call ptr @malloc(i32 mul (i32 ptrtoint (ptr getelementptr (%Inner, ptr null, i32 1) to i32), i32 5))
  %slice.ptr = insertvalue { ptr, i32, i32 } zeroinitializer, ptr %slice_data, 0
  %slice.len = insertvalue { ptr, i32, i32 } %slice.ptr, i32 0, 1
  %slice.cap = insertvalue { ptr, i32, i32 } %slice.len, i32 5, 2
  %Outer.inners = getelementptr inbounds %Outer, ptr %Outer.instance, i32 0, i32 0
  store { ptr, i32, i32 } %slice.cap, ptr %Outer.inners, align 8
  %Inner.instance = alloca %Inner, align 8
  %slice_data1 = tail call ptr @malloc(i32 mul (i32 ptrtoint (ptr getelementptr (i32, ptr null, i32 1) to i32), i32 5))
  %slice.ptr2 = insertvalue { ptr, i32, i32 } zeroinitializer, ptr %slice_data1, 0
  %slice.len3 = insertvalue { ptr, i32, i32 } %slice.ptr2, i32 0, 1
  %slice.cap4 = insertvalue { ptr, i32, i32 } %slice.len3, i32 5, 2
  %Inner.vals = getelementptr inbounds %Inner, ptr %Inner.instance, i32 0, i32 0
  store { ptr, i32, i32 } %slice.cap4, ptr %Inner.vals, align 8
  ret i32 0
}

declare noalias ptr @malloc(i32)
