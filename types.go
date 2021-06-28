package croconf

type Field interface {
	Consolidate() []error
	HasBeenSet() bool
	SourceOfValue() Source
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
	StringValueSource
	// TODO: all all other types
	/*
		UintValueSource
		Uint8ValueSource
		Uint16ValueSource
		Uint32ValueSource
		Uint64ValueSource
		IntValueSource
		Int8ValueSource
		Int16ValueSource
		Int32ValueSource
		...
	*/
	Int64ValueSource
	CustomValueSource
}

type StringValueSource interface {
	SourceGetter
	SaveStringTo(dest *string) error
}

type Int64ValueSource interface {
	SourceGetter
	GetSource() Source
	SaveInt64To(dest *int64) error
}

type CustomValueSource interface {
	SourceGetter
	GetSource() Source
	SaveCustomTo(dest interface{}) error
}
