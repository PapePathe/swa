package tests

import (
	"testing"
)

func TestTuples(t *testing.T) {
	t.Run("01", func(t *testing.T) {
		NewSuccessfulCompileRequest(t, "./tuples/01.swa",
			"Alice is 30 years old, active = 1")
	})
	t.Run("02", func(t *testing.T) {
		NewFailedCompileRequest(t, "./tuples/02.swa",
			"Expected assignment of Number but got Tuple\n")
	})
	t.Run("03", func(t *testing.T) {
		NewFailedCompileRequest(t, "./tuples/03.swa",
			"length of tuples does not match expected 2 got 3\n")
	})
	t.Run("04", func(t *testing.T) {
		NewFailedCompileRequest(t, "./tuples/04.swa",
			"Values mismatch in tuple expected String got Number at index 1\n")
	})
	t.Run("05", func(t *testing.T) {
		NewFailedCompileRequest(t, "./tuples/05.swa",
			"cannot coerce StructTypeKind and IntegerTypeKind\n")
	})
	t.Run("06", func(t *testing.T) {
		NewFailedCompileRequest(t, "./tuples/06.swa",
			"function add expect 2 arguments but was given 1\n")
	})
	t.Run("07", func(t *testing.T) {
		NewFailedCompileRequest(t, "./tuples/07.swa",
			"expected argument of type IntegerType(32 bits) but got StructType\n")
	})
	t.Run("08", func(t *testing.T) {
		NewSuccessfulCompileRequest(t, "./tuples/08.swa", "x=3, y=-1")
	})
	// FIXME developer should be able to ignore one
	// or more values of a tuple
	//	t.Run("09", func(t *testing.T) {
	//		NewSuccessfulCompileRequest(t, "./tuples/09.swa", "")
	//	})

	//	FIXME fix this tst
	//	t.Run("11", func(t *testing.T) {
	//		NewSuccessfulCompileRequest(t, "./tuples/11.swa", "")
	//	})

	//	FIXME: fix the below
	// Basic Block in function 'minmax' does not have terminator!\nlabel %merge\n\n
	//	t.Run("12", func(t *testing.T) {
	//		NewSuccessfulCompileRequest(t, "./tuples/12.swa", "")
	//	})
	t.Run("13", func(t *testing.T) {
		NewSuccessfulCompileRequest(t, "./tuples/13.swa",
			"i=0: (0,0)i=1: (1,2)i=2: (2,4)")
	})

	t.Run("14", func(t *testing.T) {
		NewSuccessfulCompileRequest(t, "./tuples/14.swa", "fib(5)=5, fib(6)=8")
	})

	// FIXME
	// Function return type does not match operand type of return
	// inst!\n  ret { ptr, ptr } %6\n { [3 x i32], [2 x i32] }\n"
	//
	//	t.Run("15", func(t *testing.T) {
	//		NewSuccessfulCompileRequest(t, "./tuples/15.swa",
	//			"")
	//	})

	// FIXME we do not support returning structs from functions
	//	t.Run("16", func(t *testing.T) {
	//		NewSuccessfulCompileRequest(t, "./tuples/16.swa",
	//			"")
	//	})

	t.Run("17", func(t *testing.T) {
		NewFailedCompileRequest(t, "./tuples/17.swa",
			"cannot insert Tuple in array of Number\n")
	})

	t.Run("18", func(t *testing.T) {
		NewFailedCompileRequest(t, "./tuples/18.swa",
			"Cannot assign tuple to struct property\n")
	})

	t.Run("19", func(t *testing.T) {
		NewSuccessfulCompileRequest(t, "./tuples/19.swa",
			"first called (1)\nsecond called (2)\na=1, b=2\n")
	})
	t.Run("20", func(t *testing.T) {
		NewFailedCompileRequest(t, "./tuples/20.swa",
			"expected return value of function to be Tuple(Number,Number), but got Number\n")
	})
	t.Run("21", func(t *testing.T) {
		NewFailedCompileRequest(t, "./tuples/21.swa",
			"function has no return statement\n")
	})
	t.Run("22", func(t *testing.T) {
		NewSuccessfulCompileRequest(t, "./tuples/22.swa",
			"(-1): 0,0\n(0): 1,1\n(5): 5,10\n(-3): 0,0\n")
	})
	t.Run("23", func(t *testing.T) {
		NewSuccessfulCompileRequest(t, "./tuples/23.swa",
			"loop\nloop\nloop\nloop\n")
	})
	t.Run("24", func(t *testing.T) {
		NewFailedCompileRequest(t, "./tuples/24.swa",
			"\x1b[33mexpected ASSIGNMENT, but got OPEN_PAREN at line 7\x1b[0m\n\n\x1b[34m7\x1b[0m \x1b[32mx\x1b[0m \x1b[32m:\x1b[0m \x1b[31m(\x1b[0m \n")
	})
	t.Run("25", func(t *testing.T) {
		NewFailedCompileRequest(t, "./tuples/25.swa",
			"Values mismatch in tuple expected Number got Float at index 0\n")
	})
	t.Run("26", func(t *testing.T) {
		NewSuccessfulCompileRequest(t, "./tuples/26.swa",
			"done")
	})
	t.Run("27", func(t *testing.T) {
		NewFailedCompileRequest(t, "./tuples/27.swa",
			"function inner expect 2 arguments but was given 1\n")
	})
	t.Run("28", func(t *testing.T) {
		NewSuccessfulCompileRequest(t, "./tuples/28.swa",
			"1 2 3 4 5 6 7 8 9 10")
	})
	t.Run("29", func(t *testing.T) {
		NewFailedCompileRequest(t, "./tuples/29.swa",
			"expected return value of function to be Tuple(Number,Number), but got Tuple(Number,Number,Number)\n")
	})
	t.Run("30", func(t *testing.T) {
		NewSuccessfulCompileRequest(t, "./tuples/30.swa",
			"a=5, b=6")
	})
	t.Run("31", func(t *testing.T) {
		NewFailedCompileRequest(t, "./tuples/31.swa",
			"Expected right side of tuple assignment to be tuple\n")
	})
	// FIXME the typechecker should catch this error
	t.Run("32", func(t *testing.T) {
		NewFailedCompileRequest(t, "./tuples/32.swa",
			"Branch condition is not 'i1' type!\n  br { i1, i32 } %0, label %if, label %else\n  %0 = call { i1, i32 } @cond()\n\n")
	})
	t.Run("33", func(t *testing.T) {
		NewSuccessfulCompileRequest(t, "./tuples/33.swa",
			"1 2")
	})
	t.Run("34", func(t *testing.T) {
		NewFailedCompileRequest(t, "./tuples/34.swa",
			"expected return value of function to be Number, but got Tuple(Number,Number)\n")
	})
	t.Run("35", func(t *testing.T) {
		NewSuccessfulCompileRequest(t, "./tuples/35.swa",
			"7 8")
	})
	t.Run("36", func(t *testing.T) {
		NewFailedCompileRequest(t, "./tuples/36.swa",
			"ArrayAccessExpression not implemented for *ast.FunctionCallExpression\n")
	})
	t.Run("37", func(t *testing.T) {
		NewFailedCompileRequest(t, "./tuples/37.swa",
			"\x1b[33mexpected CLOSE_PAREN, but got RETURN at line 8\x1b[0m\n\n\x1b[34m0\x1b[0m \n\x1b[34m8\x1b[0m \x1b[31mreturn\x1b[0m \n")
	})
	// FIXME the error message is not user friendly
	t.Run("38", func(t *testing.T) {
		NewFailedCompileRequest(t, "./tuples/38.swa",
			"nud handler expected for token COMMA and binding power 0 \n 7")
	})
	t.Run("39", func(t *testing.T) {
		NewSuccessfulCompileRequest(t, "./tuples/39.swa",
			"b=1, c=2a=0")
	})
	t.Run("40", func(t *testing.T) {
		NewSuccessfulCompileRequest(t, "./tuples/40.swa",
			"g1=100, g2=200")
	})
	t.Run("41", func(t *testing.T) {
		NewSuccessfulCompileRequest(t, "./tuples/41.swa",
			"arr[0]=2, arr[1]=3")
	})
	t.Run("42", func(t *testing.T) {
		NewSuccessfulCompileRequest(t, "./tuples/42.swa",
			"p.x=2, p.y=3")
	})
	t.Run("43", func(t *testing.T) {
		NewFailedCompileRequest(t, "./tuples/43.swa",
			"struct with pointer reference to self not supported, property: next\n")
	})
	t.Run("44", func(t *testing.T) {
		NewFailedCompileRequest(t, "./tuples/44.swa",
			"Only numbers are supported as array index, current: (&{nextIndex [] [IDENTIFIER (nextIndex) (12:8) OPEN_PAREN (12:17) CLOSE_PAREN (12:18)] <nil>})\n")
	})
	t.Run("45", func(t *testing.T) {
		NewFailedCompileRequest(t, "./tuples/45.swa",
			"nud handler expected for token COMMA and binding power 0 \n 9")
	})
	t.Run("46", func(t *testing.T) {
		NewSuccessfulCompileRequest(t, "./tuples/46.swa",
			"")
	})
	t.Run("47", func(t *testing.T) {
		NewSuccessfulCompileRequest(t, "./tuples/47.swa",
			"a =2 b =3c =2 d =3")
	})
	t.Run("48", func(t *testing.T) {
		NewFailedCompileRequest(t, "./tuples/48.swa",
			"function takeTwo expect 2 arguments but was given 1\n")
	})
}

func TestTupleReturns(t *testing.T) {
	t.Run("Return Integers", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./tuples/return_ints.swa",
			"x: 10, y: 20\n")
	})

	t.Run("Return Floats", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./tuples/return_floats.swa",
			"x: 1.5, y: 2.5\n")
	})

	t.Run("Return Error", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./tuples/return_error.swa",
			"val: 10, err: \nval: 0, err: value is negative\n")
	})

	t.Run("Swap", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./tuples/swap.swa",
			"Before: a=10, b=20\nAfter: a=20, b=10\n")
	})

	t.Run("Swap typecheck", func(t *testing.T) {
		NewFailedCompileRequest(t,
			"./tuples/swap-typecheck.swa",
			"Expected assignment of Number but got Tuple\n")
	})

	t.Run("Mixed Types", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./tuples/mixed.swa",
			"i=1, f=2.5\n")
	})

	t.Run("Recursion", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./tuples/recursion.swa",
			"x=5, y=10\n")
	})

	t.Run("Assign to array access expression", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./tuples/assign-to-array-index.swa",
			"tab[0]=15, tab[1]=150\n")
	})

	t.Run("Assign to array access expression in struct", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./tuples/assign-to-array-index-in-struct.swa",
			"tab[0]=15, tab[1]=150\n")
	})

	t.Run("Assign to array of structs access", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./tuples/assign-to-array-of-structs-index.swa",
			"tab[0]=15, tab[1]=150\n")
	})

	t.Run("Assign to member expression", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./tuples/assign-to-member-expression.swa",
			"p.div = 4, p.mod = 8")
	})

	t.Run("Assign to member typecheck", func(t *testing.T) {
		NewFailedCompileRequest(t,
			"./tuples/assign-to-member-expression-typecheck.swa",
			"Values mismatch in tuple expected Number got Error at index 0\n")
	})

	t.Run("Assign to member typecheck length", func(t *testing.T) {
		NewFailedCompileRequest(t,
			"./tuples/assign-to-member-expression-typecheck-length.swa",
			"length of tuples does not match expected 3 got 4\n")
	})
}
