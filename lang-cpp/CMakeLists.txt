cmake_minimum_required(VERSION 4.2.1)
project(SwaCompiler)

set(CMAKE_CXX_STANDARD 17)
set(LLVM_INSTALL_DIR "" CACHE_PATH "LLVM install dir")
find_package(LLVM REQUIRED CONFIG)
include_directories(${LLVM_INCLUDE_DIRS})
link_directories(${LLVM_LIBRARY_DIRS})
add_definitions(${LLVM_DEFINITIONS})

if(NOT LLVM_ENABLE_RTTI)
  set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -fno-rtti")
endif()

# 2. Include our headers
include_directories(include)

# 3. Add executable
add_executable(swa
    src/lexer.cpp 
    src/swa.cpp 
)

# 4. Link LLVM Libraries
llvm_map_components_to_libnames(llvm_libs support )
target_link_libraries(swa ${llvm_libs})
