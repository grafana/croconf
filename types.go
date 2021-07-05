package croconf

import "encoding"

type Field interface {
	Destination() interface{}
	Bindings() []Binding
}

type Source interface {
	Initialize() error
	GetName() string // TODO: remove?
}

type Binding interface {
	Apply() error
}

type BindingFromSource interface {
	Binding
	Source() Source
	BoundName() string
}

type LazySingleValueBinder interface {
	StringValueBinder
	IntValueBinder
	UintValueBinder
	FloatValueBinder
	BoolValueBinder
	TextBasedValueBinder
}

type StringValueBinder interface {
	BindStringValueTo(*string) Binding
}

type IntValueBinder interface {
	BindIntValueTo(*int64) Binding
}

type UintValueBinder interface {
	BindUintValueTo(*uint64) Binding
}

type FloatValueBinder interface {
	BindFloatValueTo(*float64) Binding
}

type BoolValueBinder interface {
	BindBoolValueTo(dest *bool) Binding
}

type TextBasedValueBinder interface {
	BindTextBasedValueTo(dest encoding.TextUnmarshaler) Binding
}

type CustomValueBinder interface {
	BindValue() Binding
}

// TODO: rename to List or Slice instead of Array?
type ArrayValueBinder interface {
	BindArrayValueTo(length *int, element *func(int) LazySingleValueBinder) Binding
}
