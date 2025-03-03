package values

import "fmt"

type StringValue struct {
	Value string
}

func (sv StringValue) GetValue() any {
	return sv.Value
}

func (sv StringValue) String() string {
	return fmt.Sprintf("\"%s\"", sv.Value)
}
