package values

type Value interface {
	GetValue() any
	String() string
}
