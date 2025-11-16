package values

import "fmt"

type IntegerValue struct {
	Value int64
}

func (iv IntegerValue) GetValue() any {
	return iv.Value
}

func (iv IntegerValue) String() string {
	return fmt.Sprintf("%d", iv.Value)
}
