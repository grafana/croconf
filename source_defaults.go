package croconf

import (
	"encoding"
	"fmt"
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

func (div defaultIntValue) BindIntValue() func(bitSize int) (int64, error) {
	return func(bitSize int) (int64, error) {
		val := int64(div)
		// See https://golang.org/pkg/math/#pkg-constants
		min, max := int64(-1<<(bitSize-1)), int64(1<<(bitSize-1)-1)
		if val < min || val > max {
			return 0, fmt.Errorf("invalid value %d, has to be between %d and %d", val, min, max)
		}
		return val, nil
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
