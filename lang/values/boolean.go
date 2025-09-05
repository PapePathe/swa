package values

import "fmt"

var TrueBooleanValue = BooleaValue{Value: true}

type BooleaValue struct {
	Value bool
}

func (iv BooleaValue) GetValue() any {
	return iv.Value
}

func (iv BooleaValue) String() string {
	return fmt.Sprintf("%v", iv.Value)
}
