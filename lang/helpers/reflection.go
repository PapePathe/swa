package helpers

import (
	"fmt"
	"reflect"
)

func ExpectType[T any](r any) T {
	expectedType := reflect.TypeOf((*T)(nil)).Elem()
	receivedType := reflect.TypeOf(r)

	if expectedType == receivedType {
		return r.(T)
	}

	panic(fmt.Sprintf("Expected %s but got %s", expectedType, receivedType))
}
