package croconf

type ConfigSource interface {
	ParseAndApply() error
}

type MultiSingleValueSource interface {
	StringValueSource
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
	*/
	Int64ValueSource
	// ... TODO: all all other types
}

type StringValueSource interface {
	SaveStringTo(dest *string) error
}

type Int64ValueSource interface {
	SaveInt64To(dest *int64) error
}
