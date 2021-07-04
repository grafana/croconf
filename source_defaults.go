package croconf

import (
	"encoding"
)

// TODO: make more flexible with callbacks, so that besides defaut values, we
// can use these for custom wrappers as well?

type defaultStringValue string

func (dsv defaultStringValue) BindStringValueTo(dest *string) func() error {
	return func() error {
		*dest = string(dsv)
		return nil
	}
}

func (dsv defaultStringValue) BindTextBasedValueTo(dest encoding.TextUnmarshaler) func() error {
	return func() error {
		return dest.UnmarshalText([]byte(dsv))
	}
}

func (dsv defaultStringValue) Source() Source {
	return nil
}

func DefaultStringValue(s string) interface {
	StringValueBinder
	TextBasedValueBinder
	BindingFromSource
} {
	return defaultStringValue(s)
}

type defaultIntValue int64

func (div defaultIntValue) BindIntValueTo(dest *int64) func() error {
	return func() error {
		*dest = int64(div)
		return nil
	}
}

func (div defaultIntValue) Source() Source {
	return nil
}

func DefaultIntValue(i int64) interface {
	IntValueBinder
	BindingFromSource
} {
	return defaultIntValue(i)
}

type DefaultCustomValue func()

var _ interface {
	CustomValueBinder
	BindingFromSource
} = DefaultCustomValue(nil)

func (dcv DefaultCustomValue) Source() Source {
	return nil
}

func (dcv DefaultCustomValue) BindValue() func() error {
	return func() error {
		dcv()
		return nil
	}
}

// TODO: sources for the rest of the types
