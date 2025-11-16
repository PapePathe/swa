package values

import "fmt"

type FloatValue struct {
	Value float64
}

func (fv FloatValue) GetValue() any {
	return fv.Value
}

func (fv FloatValue) String() string {
	return fmt.Sprintf("%f", fv.Value)
}
