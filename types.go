package croconf

type ConfigSource interface {
	ParseAndApply() error
}

type MultiSingleValueSource interface {
	StringValueSource
	Int64ValueSource
	// ... TODO: all all other types
}

type StringValueSource interface {
	SaveStringTo(dest *string) error
}

type Int64ValueSource interface {
	SaveInt64To(dest *int64) error
}
