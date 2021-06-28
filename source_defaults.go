package croconf

type defaultStringValue string

func (dsv defaultStringValue) SaveStringTo(dest *string) error {
	*dest = string(dsv)
	return nil
}

func (dsv defaultStringValue) GetSource() Source {
	return nil
}

func DefaultStringValue(s string) StringValueSource {
	return defaultStringValue(s)
}

type defaultInt64Value int64

func (div defaultInt64Value) SaveInt64To(dest *int64) error {
	*dest = int64(div)
	return nil
}
func (div defaultInt64Value) GetSource() Source {
	return nil
}

func DefaultInt64Value(i int64) Int64ValueSource {
	return defaultInt64Value(i)
}

type defaultCustomValue struct {
	value CustomValue
}

func (dcv defaultCustomValue) SaveCustomValueTo(dest CustomValue) error {
	dest = dcv.value
	return nil
}

func (dcv defaultCustomValue) GetSource() Source {
	return nil
}

func DefaultCustomValue(val CustomValue) CustomValueSource {
	return defaultCustomValue{value: val}
}

//TODO: sources for the rest of the types
