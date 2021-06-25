package croconf

import "encoding"

type Field interface {
	Consolidate() []error
	ValueSource() Source // nil for default value
	Destination() interface{}
}

type FieldOption func(field Field) // TODO: do we need this?

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
	BindStringValueTo(dest *string) func() error
}

type IntValueBinding interface {
	SourceGetter
	BindIntValue() func(bitSize int) (int64, error)
}

type UintValueBinding interface {
	SourceGetter
	BindUintValue() func(bitSize int) (uint64, error)
}

type FloatValueBinding interface {
	SourceGetter
	BindFloatValue() func(bitSize int) (float64, error)
}

type TextBasedValueBinding interface {
	SourceGetter
	BindTextBasedValueTo(dest encoding.TextUnmarshaler) func() error
}

type CustomValueBinding interface {
	SourceGetter
	BindValue() func() error
}
