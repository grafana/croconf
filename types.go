package croconf

import "encoding"

type Field interface {
	Consolidate() []error
	ValueSource() Source // nil for default value
	Destination() interface{}
}

type Source interface {
	Initialize() error
	// TODO: figure out what to put here? and if we even need this :/
	GetName() string
}

type Binding interface {
	Apply() error
}

type BindingFromSource interface {
	// Binding
	Source() Source
}

type BindingWithName interface {
	// Binding
	BoundName() string
}

type LazySingleValueBinding interface {
	StringValueBinding
	IntValueBinding
	UintValueBinding
	FloatValueBinding
	BoolValueBinding
	TextBasedValueBinding
}

type ArrayBinding interface {
	BindArray() func() (Array, error)
}

type Array interface { // TODO: rename to List and ListBinding?
	Len() int
	Element(int) LazySingleValueBinding
}

type StringValueBinding interface {
	BindStringValueTo(*string) func() error
}

type IntValueBinding interface {
	BindIntValueTo(*int64) func() error
}

type UintValueBinding interface {
	BindUintValueTo(*uint64) func() error
}

type FloatValueBinding interface {
	BindFloatValueTo(*float64) func() error
}

type BoolValueBinding interface {
	BindBoolValueTo(dest *bool) func() error
}

type TextBasedValueBinding interface {
	BindTextBasedValueTo(dest encoding.TextUnmarshaler) func() error
}

type CustomValueBinding interface {
	BindValue() func() error
}
