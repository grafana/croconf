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

type LazySingleValueBinder interface {
	StringValueBinder
	IntValueBinder
	UintValueBinder
	FloatValueBinder
	BoolValueBinder
	TextBasedValueBinder
}

type ArrayBinder interface {
	BindArray() func() (Array, error)
}

type Array interface { // TODO: rename to List and ListBinding?
	Len() int
	Element(int) LazySingleValueBinder
}

type StringValueBinder interface {
	BindStringValueTo(*string) func() error
}

type IntValueBinder interface {
	BindIntValueTo(*int64) func() error
}

type UintValueBinder interface {
	BindUintValueTo(*uint64) func() error
}

type FloatValueBinder interface {
	BindFloatValueTo(*float64) func() error
}

type BoolValueBinder interface {
	BindBoolValueTo(dest *bool) func() error
}

type TextBasedValueBinder interface {
	BindTextBasedValueTo(dest encoding.TextUnmarshaler) func() error
}

type CustomValueBinder interface {
	BindValue() func() error
}
