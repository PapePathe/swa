package values

import "fmt"

type NumberValue struct {
	Value float64
}

func (iv NumberValue) GetValue() any {
	return iv.Value
}

func (iv NumberValue) String() string {
	return fmt.Sprintf("%f", iv.Value)
}
