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

func (dsv defaultStringValue) GetSource() Source {
	return nil
}

func DefaultStringValue(s string) interface {
	StringValueBinding
	TextBasedValueBinding
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

func (div defaultIntValue) GetSource() Source {
	return nil
}

func DefaultIntValue(i int64) IntValueBinding {
	return defaultIntValue(i)
}

type DefaultCustomValue func()

func (dcv DefaultCustomValue) GetSource() Source {
	return nil
}

func (dcv DefaultCustomValue) BindValue() func() error {
	return func() error {
		dcv()
		return nil
	}
}

// TODO: sources for the rest of the types
