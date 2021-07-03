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

type SourceGetter interface {
	GetSource() Source
}

type LazySingleValueBinding interface {
	SourceGetter
	StringValueBinding
	IntValueBinding
	UintValueBinding
	FloatValueBinding
	BoolValueBinding
	TextBasedValueBinding
}

type ArrayBinding interface {
	SourceGetter
	BindArray() func() (Array, error)
}

type Array interface { // TODO: rename to List and ListBinding?
	Len() int
	Element(int) LazySingleValueBinding
}

type StringValueBinding interface {
	SourceGetter
	BindStringValueTo(*string) func() error
}

type IntValueBinding interface {
	SourceGetter
	BindIntValueTo(*int64) func() error
}

type UintValueBinding interface {
	SourceGetter
	BindUintValueTo(*uint64) func() error
}

type FloatValueBinding interface {
	SourceGetter
	BindFloatValueTo(*float64) func() error
}

type BoolValueBinding interface {
	SourceGetter
	BindBoolValueTo(dest *bool) func() error
}

type TextBasedValueBinding interface {
	SourceGetter
	BindTextBasedValueTo(dest encoding.TextUnmarshaler) func() error
}

type CustomValueBinding interface {
	SourceGetter
	BindValue() func() error
}
