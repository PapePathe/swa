package helpers

import (
	"fmt"
	"reflect"
)

func ExpectType[T any](r any) T {
	expectedType := reflect.TypeOf((*T)(nil)).Elem()
	receivedType := reflect.TypeOf(r)

	if expectedType == receivedType {
		result, ok := r.(T)

		if !ok {
			panic(fmt.Sprintf("Type assertion failed for %s", r))
		}

		return result
	}

	panic(fmt.Sprintf("Expected %s but got %s", expectedType, receivedType))
}
