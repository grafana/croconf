package croconf

import "encoding"

type Field interface {
	Consolidate() []error
	ValueSource() Source // nil for default value
	Destination() interface{}
}

type FieldOption func(field Field) // TODO: do we need this?

type Source interface {
	// TODO: figure out what to put here? and if we even need this :/
	GetName() string
}

type SourceGetter interface {
	GetSource() Source
}

type LazySingleValueBinding interface {
	SourceGetter
	StringValueBinding
	// TODO: all all other types
	/*
		UintValueBinding
		Uint8ValueBinding
		Uint16ValueBinding
		Uint32ValueBinding
		Uint64ValueBinding
		IntValueBinding
		Int8ValueBinding
		Int16ValueBinding
		Int32ValueBinding
		...
	*/
	Int64ValueBinding
	TextBasedValueBinding
}

type StringValueBinding interface {
	SourceGetter
	BindStringValueTo(dest *string) func() error
}

type Int64ValueBinding interface {
	SourceGetter
	BindInt64ValueTo(dest *int64) func() error
}

type TextBasedValueBinding interface {
	SourceGetter
	BindTextBasedValueTo(dest encoding.TextUnmarshaler) func() error
}

type CustomValueBinding interface {
	SourceGetter
	BindValue() func() error
}
