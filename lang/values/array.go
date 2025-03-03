package values

type ArrayValue struct {
	Values []Value
}

func (av ArrayValue) GetValue() interface{} {
	return av.Values
}

func (av ArrayValue) String() string {
	result := "["
	for i, v := range av.Values {
		result += v.String()
		if i < len(av.Values)-1 {
			result += ", "
		}
	}
	result += "]"
	return result
}
