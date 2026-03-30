package tests

// import (
// 	"fmt"
// 	"math/rand"
// 	"os"
// 	"os/exec"
// 	"strings"
// 	"testing"
//
// 	"github.com/google/uuid"
// )
//
// // ---------------------------------------------------------------------------
// // Helper: compile source text through the `./swa` binary.
// // Returns compiler stdout/stderr and any exec error.
// // ---------------------------------------------------------------------------
//
// func compileSWA(t *testing.T, source string) (string, error) {
// 	t.Helper()
//
// 	// Write source to a temp file
// 	tmpFile := fmt.Sprintf("/tmp/fuzz_%s.swa", uuid.New().String())
// 	if err := os.WriteFile(tmpFile, []byte(source), 0o644); err != nil {
// 		return "", fmt.Errorf("failed to write temp file: %w", err)
// 	}
// 	defer os.Remove(tmpFile)
//
// 	outName := uuid.New().String()
// 	cmd := exec.Command("./swa", "compile", "-s", tmpFile, "-o", outName)
// 	output, err := cmd.CombinedOutput()
//
// 	// Cleanup artifacts regardless of outcome
// 	for _, ext := range []string{".ll", ".s", ".o", ".exe"} {
// 		os.Remove(outName + ext)
// 	}
//
// 	return string(output), err
// }
//
// // ---------------------------------------------------------------------------
// // FuzzCompilerNocrash – feed arbitrary bytes to the compiler and ensure it
// // never panics / segfaults.  The compiler may return errors – that is fine.
// // ---------------------------------------------------------------------------
//
// func FuzzCompilerNoCrash(f *testing.F) {
// 	// ---- Seed corpus: valid programs covering supported features ----------
//
// 	// Minimal valid program
// 	f.Add(`dialect:english;
// start() int {
//   return 0;
// }
// `)
//
// 	// Variable declarations and print
// 	f.Add(`dialect:english;
// start() int {
//   let x: int = 42;
//   let y: float = 3.14;
//   let s: string = "hello";
//   print("%d %.2f %s", x, y, s);
//   return 0;
// }
// `)
//
// 	// Basic arithmetic
// 	f.Add(`dialect:english;
// start() int {
//   let a: int = 10;
//   let b: int = 20;
//   let c: int = a + b;
//   let d: int = a - b;
//   let e: int = a * b;
//   let f: int = b / a;
//   print("%d %d %d %d", c, d, e, f);
//   return 0;
// }
// `)
//
// 	// Function definitions and calls
// 	f.Add(`dialect:english;
// func add(a: int, b: int) int {
//   return a + b;
// }
// start() int {
//   let r: int = add(10, 20);
//   print("%d", r);
//   return 0;
// }
// `)
//
// 	// Struct definition and access
// 	f.Add(`dialect:english;
// start() int {
//   struct Point {
//     X: int,
//     Y: int,
//   }
//   let p: Point = Point{X: 10, Y: 20};
//   print("x=%d y=%d", p.X, p.Y);
//   return 0;
// }
// `)
//
// 	// Array initialization and access
// 	f.Add(`dialect:english;
// start() int {
//   let arr: [5]int = [5]int{1, 2, 3, 4, 5};
//   print("%d %d", arr[0], arr[4]);
//   return 0;
// }
// `)
//
// 	// If-else conditionals
// 	f.Add(`dialect:english;
// start() int {
//   let x: int = 10;
//   if (x == 10) {
//     print("ok");
//   } else {
//     print("fail");
//   }
//   return 0;
// }
// `)
//
// 	// While loop
// 	f.Add(`dialect:english;
// start() int {
//   let i: int = 0;
//   while (i < 5) {
//     print("%d ", i);
//     i = i + 1;
//   }
//   return 0;
// }
// `)
//
// 	// Tuples / multiple return values
// 	f.Add(`dialect:english;
// func get_pair() (int, int) {
//   return 10, 20;
// }
// start() int {
//   let a: int = 0;
//   let b: int = 0;
//   (a, b) = get_pair();
//   print("a=%d b=%d", a, b);
//   return 0;
// }
// `)
//
// 	// Error expression
// 	f.Add(`dialect:english;
// func divide(a: int, b: int) (int, error) {
//   if (b == 0) {
//     return 0, error("division by zero");
//   }
//   return a / b, zero(error);
// }
// start() int {
//   let res: int = 0;
//   let err: error;
//   (res, err) = divide(100, 0);
//   if (err != zero(error)) {
//     print("error: %s", err);
//     return 0;
//   }
//   print("%d", res);
//   return 0;
// }
// `)
//
// 	// Struct with array field
// 	f.Add(`dialect:english;
// start() int {
//   struct Container {
//     Size: int,
//     Data: *int,
//   }
//   let backing: [3]int = [3]int{10, 20, 30};
//   let c: Container = Container{Size: 3, Data: backing};
//   print("size=%d", c.Size);
//   return 0;
// }
// `)
//
// 	// Nested conditionals and boolean operators
// 	f.Add(`dialect:english;
// start() int {
//   let a: int = 1;
//   let b: int = 0;
//   if (a == 1 && b == 0) {
//     print("both true");
//   }
//   if (a == 1 || b == 1) {
//     print("or true");
//   }
//   return 0;
// }
// `)
//
// 	// Float arithmetic
// 	f.Add(`dialect:english;
// func addf(a: float, b: float) float {
//   return a + b;
// }
// start() int {
//   let x: float = 1.5;
//   let y: float = 2.5;
//   let z: float = addf(x, y);
//   print("%.1f", z);
//   return 0;
// }
// `)
//
// 	// String variable
// 	f.Add(`dialect:english;
// start() int {
//   let name: string = "world";
//   print("hello %s", name);
//   return 0;
// }
// `)
//
// 	// Prefix expression with negation
// 	f.Add(`dialect:english;
// start() int {
//   let x: int = 10;
//   let y: int = -x;
//   print("%d", y);
//   return 0;
// }
// `)
//
// 	// French dialect
// 	f.Add(`dialect:french;
// entree() entier {
//   variable x: entier = 42;
//   ecrire("%d", x);
//   retourner 0;
// }
// `)
//
// 	// Soussou dialect
// 	f.Add(`dialect:soussou;
// fɔlɛ() kɔntin {
//   a sa x: kɔntin = 42;
//   masen("%d", x);
//   ragiri 0;
// }
// `)
//
// 	// ---- Seed corpus: programs that should produce compiler errors ----------
//
// 	// Self-referencing struct
// 	f.Add(`dialect:english;
// struct Node {
//   id: int,
//   parent: Node,
// }
// start() int {
//   return 0;
// }
// `)
//
// 	// Pointer self-reference
// 	f.Add(`dialect:english;
// struct Node {
//   id: int,
//   parent: *Node,
// }
// start() int {
//   return 0;
// }
// `)
//
// 	// Type mismatch in variable declaration
// 	f.Add(`dialect:english;
// start() int {
//   let x: int = 3.14;
//   return 0;
// }
// `)
//
// 	// Type mismatch in assignment
// 	f.Add(`dialect:english;
// start() int {
//   let x: float = 1.0;
//   x = "hello";
//   return 0;
// }
// `)
//
// 	// Missing function arity
// 	f.Add(`dialect:english;
// func add(a: int, b: int) int {
//   return a + b;
// }
// start() int {
//   let x: int = add(1);
//   return 0;
// }
// `)
//
// 	// Undefined variable reference
// 	f.Add(`dialect:english;
// start() int {
//   print("%d", undefined_var);
//   return 0;
// }
// `)
//
// 	// Variable redeclaration
// 	f.Add(`dialect:english;
// start() int {
//   let x: int = 10;
//   let x: int = 20;
//   return 0;
// }
// `)
//
// 	// Integer overflow (int32)
// 	f.Add(`dialect:english;
// start() int {
//   let x: int = 2147483648;
//   return 0;
// }
// `)
//
// 	// Negative integer overflow (int32)
// 	f.Add(`dialect:english;
// start() int {
//   let x: int = -2147483649;
//   return 0;
// }
// `)
//
// 	// Field access on primitive
// 	f.Add(`dialect:english;
// start() int {
//   let i: int = 10;
//   print("%d", i.field);
//   return 0;
// }
// `)
//
// 	// ---- Seed corpus: edge cases and boundary inputs ----------
//
// 	// Empty source
// 	f.Add(``)
//
// 	// Only dialect
// 	f.Add(`dialect:english;`)
//
// 	// Missing semicolon
// 	f.Add(`dialect:english
// start() int { return 0; }
// `)
//
// 	// Incomplete function
// 	f.Add(`dialect:english;
// func add(a: int, b: int) int {
// `)
//
// 	// Unclosed brace
// 	f.Add(`dialect:english;
// start() int {
//   let x: int = 10;
// `)
//
// 	// Unclosed string
// 	f.Add(`dialect:english;
// start() int {
//   let s: string = "unterminated;
//   return 0;
// }
// `)
//
// 	// Random junk
// 	f.Add(`@#$%^&*()!~`)
//
// 	// Very long identifier
// 	f.Add(`dialect:english;
// start() int {
//   let aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa: int = 0;
//   return 0;
// }
// `)
//
// 	// Deeply nested if
// 	f.Add(`dialect:english;
// start() int {
//   let x: int = 1;
//   if (x == 1) { if (x == 1) { if (x == 1) { if (x == 1) { print("deep"); } } } }
//   return 0;
// }
// `)
//
// 	// Large array
// 	f.Add(`dialect:english;
// start() int {
//   let arr: [10]int = [10]int{0,1,2,3,4,5,6,7,8,9};
//   print("%d", arr[9]);
//   return 0;
// }
// `)
//
// 	// Struct with all field types
// 	f.Add(`dialect:english;
// start() int {
//   struct Mixed {
//     I: int,
//     F: float,
//     S: string,
//   }
//   let m: Mixed = Mixed{I: 42, F: 3.14, S: "test"};
//   print("%d %.2f %s", m.I, m.F, m.S);
//   return 0;
// }
// `)
//
// 	// Multiple functions
// 	f.Add(`dialect:english;
// func square(n: int) int {
//   return n * n;
// }
// func cube(n: int) int {
//   return n * n * n;
// }
// start() int {
//   print("sq=%d cu=%d", square(3), cube(3));
//   return 0;
// }
// `)
//
// 	// While loop with array
// 	f.Add(`dialect:english;
// start() int {
//   let arr: [5]int = [5]int{10, 20, 30, 40, 50};
//   let i: int = 0;
//   let sum: int = 0;
//   while (i < 5) {
//     sum = sum + arr[i];
//     i = i + 1;
//   }
//   print("sum=%d", sum);
//   return 0;
// }
// `)
//
// 	// Comparison operators
// 	f.Add(`dialect:english;
// start() int {
//   let a: int = 5;
//   let b: int = 10;
//   if (a < b) { print("lt "); }
//   if (b > a) { print("gt "); }
//   if (a <= 5) { print("le "); }
//   if (b >= 10) { print("ge "); }
//   if (a != b) { print("ne "); }
//   if (a == 5) { print("eq"); }
//   return 0;
// }
// `)
//
// 	// Struct passed to function
// 	f.Add(`dialect:english;
// struct Employee {
//   Age: int,
//   Name: string,
// }
// func get_age(e: Employee) int {
//   return e.Age;
// }
// start() int {
//   let e: Employee = Employee{Age: 30, Name: "Alice"};
//   print("age=%d", get_age(e));
//   return 0;
// }
// `)
//
// 	// Array of structs
// 	f.Add(`dialect:english;
// struct Point {
//   X: int,
//   Y: int,
// }
// start() int {
//   let pts: [2]Point = [2]Point{Point{X:1, Y:2}, Point{X:3, Y:4}};
//   print("x0=%d y1=%d", pts[0].X, pts[1].Y);
//   return 0;
// }
// `)
//
// 	// Tuple assign to struct field
// 	f.Add(`dialect:english;
// func swap(a: int, b: int) (int, int) {
//   return b, a;
// }
// start() int {
//   let x: int = 1;
//   let y: int = 2;
//   (x, y) = swap(x, y);
//   print("x=%d y=%d", x, y);
//   return 0;
// }
// `)
//
// 	// Zero value initialization
// 	f.Add(`dialect:english;
// start() int {
//   let x: int = zero(int);
//   let y: float = zero(float);
//   print("x=%d y=%.1f", x, y);
//   return 0;
// }
// `)
//
// 	// Scoping / shadowing
// 	f.Add(`dialect:english;
// start() int {
//   let x: int = 10;
//   if (x == 10) {
//     let x: int = 20;
//     print("inner=%d ", x);
//   }
//   print("outer=%d", x);
//   return 0;
// }
// `)
//
// 	// Early return
// 	f.Add(`dialect:english;
// func check(x: int) int {
//   if (x > 10) {
//     return 1;
//   }
//   return 0;
// }
// start() int {
//   print("%d %d", check(5), check(15));
//   return 0;
// }
// `)
//
// 	// Modern print (variadic)
// 	f.Add(`dialect:english;
// start() int {
//   let a: int = 42;
//   let b: float = 3.14;
//   print("a=", a, " b=", b, "\n");
//   return 0;
// }
// `)
//
// 	// ---- The actual fuzzer ----
// 	f.Fuzz(func(t *testing.T, data string) {
// 		// The main invariant: the compiler must never panic or crash (segfault).
// 		// Compilation errors are expected and acceptable.
// 		out, err := compileSWA(t, data)
// 		if err != nil {
// 			t.Logf("compiler returned error (expected for fuzz): %v, output: %s", err, out)
// 		}
// 	})
// }
//
// // ---------------------------------------------------------------------------
// // FuzzStructuredPrograms – generate structured, syntactically plausible
// // SWA programs by mutating building blocks to explore deeper compiler paths.
// // ---------------------------------------------------------------------------
//
// func FuzzStructuredPrograms(f *testing.F) {
// 	// Seeds are minimal working programs; Go's fuzzer will mutate the
// 	// individual fields to produce interesting variants.
//
// 	f.Add(
// 		"int", // return type
// 		"x",   // var name
// 		"int", // var type
// 		"42",  // var value
// 		"%d",  // format string
// 		"==",  // comparison operator
// 		"42",  // comparison value
// 		5,     // loop iterations
// 	)
//
// 	f.Add(
// 		"int",
// 		"counter",
// 		"int",
// 		"0",
// 		"%d",
// 		"<",
// 		"10",
// 		10,
// 	)
//
// 	f.Add(
// 		"int",
// 		"val",
// 		"float",
// 		"3.14",
// 		"%.2f",
// 		">",
// 		"0.0",
// 		3,
// 	)
//
// 	f.Add(
// 		"int",
// 		"msg",
// 		"string",
// 		`"hello"`,
// 		"%s",
// 		"!=",
// 		`""`,
// 		0,
// 	)
//
// 	f.Fuzz(func(t *testing.T, retType, varName, varType, varValue, fmtStr, cmpOp, cmpVal string, loopCount int) {
// 		// Sanitize inputs to prevent trivially enormous programs
// 		if len(varName) > 64 || len(varValue) > 128 || len(fmtStr) > 64 || len(cmpVal) > 64 {
// 			return
// 		}
// 		if loopCount < 0 {
// 			loopCount = 0
// 		}
// 		if loopCount > 20 {
// 			loopCount = 20
// 		}
//
// 		// Sanitize varName to contain only alphanumeric and underscores
// 		sanitized := strings.Builder{}
// 		for _, c := range varName {
// 			if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' {
// 				sanitized.WriteRune(c)
// 			}
// 		}
// 		varName = sanitized.String()
// 		if len(varName) == 0 {
// 			varName = "x"
// 		}
// 		// Ensure variable name starts with a letter
// 		if varName[0] >= '0' && varName[0] <= '9' {
// 			varName = "v" + varName
// 		}
//
// 		var sb strings.Builder
// 		sb.WriteString("dialect:english;\n\n")
// 		sb.WriteString(fmt.Sprintf("start() %s {\n", retType))
// 		sb.WriteString(fmt.Sprintf("  let %s: %s = %s;\n", varName, varType, varValue))
//
// 		// Add a conditional
// 		sb.WriteString(fmt.Sprintf("  if (%s %s %s) {\n", varName, cmpOp, cmpVal))
// 		sb.WriteString(fmt.Sprintf("    print(\"%s\", %s);\n", fmtStr, varName))
// 		sb.WriteString("  }\n")
//
// 		// Add a loop if requested
// 		if loopCount > 0 {
// 			sb.WriteString("  let _i: int = 0;\n")
// 			sb.WriteString(fmt.Sprintf("  while (_i < %d) {\n", loopCount))
// 			sb.WriteString("    _i = _i + 1;\n")
// 			sb.WriteString("  }\n")
// 		}
//
// 		sb.WriteString("  return 0;\n")
// 		sb.WriteString("}\n")
//
// 		if out, err := compileSWA(t, sb.String()); err != nil {
// 			t.Logf("compiler error (expected for fuzz): %v, output: %s", err, out)
// 		}
// 	})
// }
//
// // ---------------------------------------------------------------------------
// // FuzzStructDefinitions – fuzzes struct definitions with varying field
// // counts, names and types to explore struct-related compiler paths.
// // ---------------------------------------------------------------------------
//
// func FuzzStructDefinitions(f *testing.F) {
// 	f.Add(1, "id", "int", false)
// 	f.Add(2, "name", "string", false)
// 	f.Add(3, "value", "float", true) // true = try self-reference
// 	f.Add(1, "data", "int", false)
//
// 	fieldTypes := []string{"int", "float", "string"}
//
// 	f.Fuzz(func(t *testing.T, numFields int, baseName string, baseType string, selfRef bool) {
// 		if numFields < 1 {
// 			numFields = 1
// 		}
// 		if numFields > 10 {
// 			numFields = 10
// 		}
// 		if len(baseName) > 32 {
// 			baseName = baseName[:32]
// 		}
//
// 		// Sanitize base name
// 		sanitized := strings.Builder{}
// 		for _, c := range baseName {
// 			if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' {
// 				sanitized.WriteRune(c)
// 			}
// 		}
// 		baseName = sanitized.String()
// 		if len(baseName) == 0 {
// 			baseName = "field"
// 		}
// 		if baseName[0] >= '0' && baseName[0] <= '9' {
// 			baseName = "f" + baseName
// 		}
//
// 		// Capitalize first letter for field name
// 		fieldBase := strings.ToUpper(baseName[:1]) + baseName[1:]
//
// 		var sb strings.Builder
// 		sb.WriteString("dialect:english;\n\n")
// 		sb.WriteString("start() int {\n")
// 		sb.WriteString("  struct TestStruct {\n")
//
// 		rng := rand.New(rand.NewSource(int64(numFields + len(baseName))))
// 		for i := 0; i < numFields; i++ {
// 			fname := fmt.Sprintf("%s%d", fieldBase, i)
// 			ftype := fieldTypes[rng.Intn(len(fieldTypes))]
// 			if i == 0 {
// 				ftype = baseType
// 			}
// 			sb.WriteString(fmt.Sprintf("    %s: %s,\n", fname, ftype))
// 		}
//
// 		if selfRef {
// 			sb.WriteString("    Parent: TestStruct,\n")
// 		}
//
// 		sb.WriteString("  }\n\n")
//
// 		// Don't try to initialize self-referencing structs
// 		if !selfRef {
// 			sb.WriteString("  let ts: TestStruct = TestStruct{\n")
// 			for i := 0; i < numFields; i++ {
// 				fname := fmt.Sprintf("%s%d", fieldBase, i)
// 				ftype := fieldTypes[rng.Intn(len(fieldTypes))]
// 				if i == 0 {
// 					ftype = baseType
// 				}
// 				switch ftype {
// 				case "int":
// 					sb.WriteString(fmt.Sprintf("    %s: %d,\n", fname, i))
// 				case "float":
// 					sb.WriteString(fmt.Sprintf("    %s: %d.%d,\n", fname, i, i))
// 				case "string":
// 					sb.WriteString(fmt.Sprintf("    %s: \"%s_%d\",\n", fname, baseName, i))
// 				default:
// 					sb.WriteString(fmt.Sprintf("    %s: %d,\n", fname, i))
// 				}
// 			}
// 			sb.WriteString("  };\n")
// 		}
//
// 		sb.WriteString("  return 0;\n")
// 		sb.WriteString("}\n")
//
// 		if out, err := compileSWA(t, sb.String()); err != nil {
// 			t.Logf("compiler error (expected for fuzz): %v, output: %s", err, out)
// 		}
// 	})
// }
//
// // ---------------------------------------------------------------------------
// // FuzzArrayOperations – fuzzes array declarations, accesses, and
// // operations to stress-test array-related compiler paths.
// // ---------------------------------------------------------------------------
//
// func FuzzArrayOperations(f *testing.F) {
// 	f.Add(3, "int", 0, true) // small int array, valid access, with loop
// 	f.Add(5, "float", 2, false)
// 	f.Add(1, "int", 0, true)  // single element
// 	f.Add(10, "int", 9, true) // boundary access at last index
//
// 	f.Fuzz(func(t *testing.T, size int, elemType string, accessIdx int, withLoop bool) {
// 		if size < 1 {
// 			size = 1
// 		}
// 		if size > 20 {
// 			size = 20
// 		}
// 		if accessIdx < 0 {
// 			accessIdx = 0
// 		}
//
// 		if elemType != "int" && elemType != "float" {
// 			elemType = "int"
// 		}
//
// 		var sb strings.Builder
// 		sb.WriteString("dialect:english;\n\n")
// 		sb.WriteString("start() int {\n")
//
// 		// Generate array init
// 		sb.WriteString(fmt.Sprintf("  let arr: [%d]%s = [%d]%s{", size, elemType, size, elemType))
// 		for i := 0; i < size; i++ {
// 			if i > 0 {
// 				sb.WriteString(", ")
// 			}
// 			switch elemType {
// 			case "int":
// 				sb.WriteString(fmt.Sprintf("%d", i*10))
// 			case "float":
// 				sb.WriteString(fmt.Sprintf("%d.%d", i, i))
// 			}
// 		}
// 		sb.WriteString("};\n")
//
// 		// Access element (may be out of bounds — that's intentional)
// 		if elemType == "int" {
// 			sb.WriteString(fmt.Sprintf("  print(\"%%d\", arr[%d]);\n", accessIdx))
// 		} else {
// 			sb.WriteString(fmt.Sprintf("  print(\"%%.2f\", arr[%d]);\n", accessIdx))
// 		}
//
// 		// Optional loop over array
// 		if withLoop {
// 			sb.WriteString("  let i: int = 0;\n")
// 			sb.WriteString(fmt.Sprintf("  while (i < %d) {\n", size))
// 			sb.WriteString("    i = i + 1;\n")
// 			sb.WriteString("  }\n")
// 		}
//
// 		sb.WriteString("  return 0;\n")
// 		sb.WriteString("}\n")
//
// 		if out, err := compileSWA(t, sb.String()); err != nil {
// 			t.Logf("compiler error (expected for fuzz): %v, output: %s", err, out)
// 		}
// 	})
// }
//
// // ---------------------------------------------------------------------------
// // FuzzFunctionSignatures – fuzzes function definitions with varying
// // parameter counts, types, and return types.
// // ---------------------------------------------------------------------------
//
// func FuzzFunctionSignatures(f *testing.F) {
// 	f.Add(1, "int", "int", false)     // single param, int return
// 	f.Add(2, "float", "float", false) // two params, float return
// 	f.Add(0, "int", "int", true)      // no params, tuple return
// 	f.Add(3, "int", "int", false)
//
// 	paramTypes := []string{"int", "float", "string"}
//
// 	f.Fuzz(func(t *testing.T, numParams int, paramType string, returnType string, tupleReturn bool) {
// 		if numParams < 0 {
// 			numParams = 0
// 		}
// 		if numParams > 10 {
// 			numParams = 10
// 		}
//
// 		if paramType != "int" && paramType != "float" && paramType != "string" {
// 			paramType = "int"
// 		}
// 		if returnType != "int" && returnType != "float" {
// 			returnType = "int"
// 		}
//
// 		var sb strings.Builder
// 		sb.WriteString("dialect:english;\n\n")
//
// 		// Function definition
// 		sb.WriteString("func test_fn(")
// 		rng := rand.New(rand.NewSource(int64(numParams)))
// 		for i := 0; i < numParams; i++ {
// 			if i > 0 {
// 				sb.WriteString(", ")
// 			}
// 			pt := paramTypes[rng.Intn(len(paramTypes))]
// 			if i == 0 {
// 				pt = paramType
// 			}
// 			sb.WriteString(fmt.Sprintf("p%d: %s", i, pt))
// 		}
// 		sb.WriteString(") ")
//
// 		if tupleReturn {
// 			sb.WriteString(fmt.Sprintf("(%s, %s)", returnType, returnType))
// 		} else {
// 			sb.WriteString(returnType)
// 		}
// 		sb.WriteString(" {\n")
//
// 		if tupleReturn {
// 			switch returnType {
// 			case "int":
// 				sb.WriteString("  return 1, 2;\n")
// 			case "float":
// 				sb.WriteString("  return 1.0, 2.0;\n")
// 			}
// 		} else {
// 			switch returnType {
// 			case "int":
// 				sb.WriteString("  return 0;\n")
// 			case "float":
// 				sb.WriteString("  return 0.0;\n")
// 			}
// 		}
// 		sb.WriteString("}\n\n")
//
// 		// Main function that calls test_fn
// 		sb.WriteString("start() int {\n")
// 		if numParams > 0 {
// 			sb.WriteString("  let r: int = 0;\n")
//
// 			// Build call args
// 			sb.WriteString("  test_fn(")
// 			for i := 0; i < numParams; i++ {
// 				if i > 0 {
// 					sb.WriteString(", ")
// 				}
// 				pt := paramTypes[rng.Intn(len(paramTypes))]
// 				if i == 0 {
// 					pt = paramType
// 				}
// 				switch pt {
// 				case "int":
// 					sb.WriteString(fmt.Sprintf("%d", i+1))
// 				case "float":
// 					sb.WriteString(fmt.Sprintf("%d.0", i+1))
// 				case "string":
// 					sb.WriteString(fmt.Sprintf("\"%s\"", "arg"))
// 				}
// 			}
// 			sb.WriteString(");\n")
// 		}
// 		sb.WriteString("  return 0;\n")
// 		sb.WriteString("}\n")
//
// 		if out, err := compileSWA(t, sb.String()); err != nil {
// 			t.Logf("compiler error (expected for fuzz): %v, output: %s", err, out)
// 		}
// 	})
// }
//
// // ---------------------------------------------------------------------------
// // FuzzDialectVariations – tests programs with different dialect declarations
// // to explore dialect-specific lexing and parsing paths.
// // ---------------------------------------------------------------------------
//
// func FuzzDialectVariations(f *testing.F) {
// 	f.Add("english", "start", "let", "func", "return", "if", "else", "while", "print", "int", "float", "string", "struct")
// 	f.Add("french", "entree", "variable", "fonction", "retourner", "si", "sinon", "tantque", "ecrire", "entier", "decimal", "chaine", "structure")
// 	f.Add("soussou", "fɔlɛ", "a sa", "wali", "ragiri", "xa", "xa mu", "han ma", "masen", "kɔntin", "kɔntin dɔxɔ", "sɛbɛli", "sɛnbɛ")
//
// 	f.Fuzz(func(t *testing.T, dialect, mainKw, letKw, funcKw, returnKw, ifKw, elseKw, whileKw, printKw, intKw, floatKw, stringKw, structKw string) {
// 		if len(dialect) > 32 {
// 			return
// 		}
//
// 		var sb strings.Builder
// 		sb.WriteString(fmt.Sprintf("dialect:%s;\n\n", dialect))
// 		sb.WriteString(fmt.Sprintf("%s() %s {\n", mainKw, intKw))
// 		sb.WriteString(fmt.Sprintf("  %s x: %s = 42;\n", letKw, intKw))
// 		sb.WriteString(fmt.Sprintf("  %s(\"%%d\", x);\n", printKw))
// 		sb.WriteString(fmt.Sprintf("  %s 0;\n", returnKw))
// 		sb.WriteString("}\n")
//
// 		if out, err := compileSWA(t, sb.String()); err != nil {
// 			t.Logf("compiler error (expected for fuzz): %v, output: %s", err, out)
// 		}
// 	})
// }
//
// // ---------------------------------------------------------------------------
// // FuzzBinaryExpressions – fuzzes binary expressions with various operators,
// // operand types, and nesting levels.
// // ---------------------------------------------------------------------------
//
// func FuzzBinaryExpressions(f *testing.F) {
// 	f.Add("int", "+", "10", "20", false)     // int addition, no parens
// 	f.Add("int", "-", "30", "15", false)     // int subtraction
// 	f.Add("int", "*", "5", "6", false)       // int multiplication
// 	f.Add("int", "/", "100", "10", false)    // int division
// 	f.Add("float", "+", "1.5", "2.5", false) // float addition
// 	f.Add("int", "==", "10", "10", true)     // comparison in conditional
// 	f.Add("int", "!=", "10", "20", true)
// 	f.Add("int", "<", "5", "10", true)
// 	f.Add("int", ">", "10", "5", true)
// 	f.Add("int", "<=", "5", "5", true)
// 	f.Add("int", ">=", "10", "5", true)
//
// 	f.Fuzz(func(t *testing.T, opType string, operator string, left string, right string, inConditional bool) {
// 		if len(left) > 32 || len(right) > 32 || len(operator) > 4 {
// 			return
// 		}
// 		if opType != "int" && opType != "float" {
// 			opType = "int"
// 		}
//
// 		var sb strings.Builder
// 		sb.WriteString("dialect:english;\n\n")
// 		sb.WriteString("start() int {\n")
//
// 		if inConditional {
// 			sb.WriteString(fmt.Sprintf("  if (%s %s %s) {\n", left, operator, right))
// 			sb.WriteString("    print(\"yes\");\n")
// 			sb.WriteString("  }\n")
// 		} else {
// 			sb.WriteString(fmt.Sprintf("  let result: %s = %s %s %s;\n", opType, left, operator, right))
// 			if opType == "int" {
// 				sb.WriteString("  print(\"%d\", result);\n")
// 			} else {
// 				sb.WriteString("  print(\"%.2f\", result);\n")
// 			}
// 		}
//
// 		sb.WriteString("  return 0;\n")
// 		sb.WriteString("}\n")
//
// 		if out, err := compileSWA(t, sb.String()); err != nil {
// 			t.Logf("compiler error (expected for fuzz): %v, output: %s", err, out)
// 		}
// 	})
// }
//
// // ---------------------------------------------------------------------------
// // FuzzGlobalDeclarations – fuzzes global variable declarations.
// // ---------------------------------------------------------------------------
//
// func FuzzGlobalDeclarations(f *testing.F) {
// 	f.Add("BITS", "int", 2, true)
// 	f.Add("MAX", "int", 1, false)
// 	f.Add("PI", "float", 1, false)
//
// 	f.Fuzz(func(t *testing.T, name string, elemType string, size int, isArray bool) {
// 		if len(name) > 32 || size < 1 || size > 10 {
// 			return
// 		}
//
// 		// Sanitize name
// 		sanitized := strings.Builder{}
// 		for _, c := range name {
// 			if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' {
// 				sanitized.WriteRune(c)
// 			}
// 		}
// 		name = sanitized.String()
// 		if len(name) == 0 {
// 			name = "G"
// 		}
// 		if name[0] >= '0' && name[0] <= '9' {
// 			name = "G" + name
// 		}
//
// 		if elemType != "int" && elemType != "float" {
// 			elemType = "int"
// 		}
//
// 		var sb strings.Builder
// 		sb.WriteString("dialect:english;\n\n")
//
// 		if isArray {
// 			sb.WriteString(fmt.Sprintf("const %s:[%d]%s = [%d]%s{", name, size, elemType, size, elemType))
// 			for i := 0; i < size; i++ {
// 				if i > 0 {
// 					sb.WriteString(",")
// 				}
// 				sb.WriteString(fmt.Sprintf("%d", i))
// 			}
// 			sb.WriteString("};\n\n")
// 		} else {
// 			switch elemType {
// 			case "int":
// 				sb.WriteString(fmt.Sprintf("const %s:%s = 42;\n\n", name, elemType))
// 			case "float":
// 				sb.WriteString(fmt.Sprintf("const %s:%s = 3.14;\n\n", name, elemType))
// 			}
// 		}
//
// 		sb.WriteString("start() int {\n")
// 		sb.WriteString("  return 0;\n")
// 		sb.WriteString("}\n")
//
// 		if out, err := compileSWA(t, sb.String()); err != nil {
// 			t.Logf("compiler error (expected for fuzz): %v, output: %s", err, out)
// 		}
// 	})
// }
