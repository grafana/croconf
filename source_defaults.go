package croconf

/*
import (
	"encoding"
)

const defaultsBoundName = "default"

// TODO: setting Source=nil for defaults is confusing, should we just not set
// any Source at all for them? By convention, if the first binding doesn't have
// a source, we can consider it the default one?

type defaultStringValue string

func (dsv defaultStringValue) BindStringValueTo(dest *string) Binding {
	return NewCallbackBindingFromSource(nil, defaultsBoundName, func() error {
		*dest = string(dsv)
		return nil
	})
}

func (dsv defaultStringValue) BindTextBasedValueTo(dest encoding.TextUnmarshaler) Binding {
	return NewCallbackBindingFromSource(nil, defaultsBoundName, func() error {
		return dest.UnmarshalText([]byte(dsv))
	})
}

func DefaultStringValue(s string) interface {
	StringValueBinder
	TextBasedValueBinder
} {
	return defaultStringValue(s)
}

type defaultIntValue int64

func (div defaultIntValue) BindIntValueTo(dest *int64) Binding {
	return NewCallbackBindingFromSource(nil, defaultsBoundName, func() error {
		*dest = int64(div)
		return nil
	})
}

func (div defaultIntValue) Source() Source {
	return nil
}

func DefaultIntValue(i int64) interface {
	IntValueBinder
} {
	return defaultIntValue(i)
}

type DefaultCustomValue func()

var _ interface {
	CustomValueBinder
} = DefaultCustomValue(nil)

func (dcv DefaultCustomValue) Source() Source {
	return nil
}

func (dcv DefaultCustomValue) BindValue() Binding {
	return NewCallbackBindingFromSource(nil, defaultsBoundName, func() error {
		dcv()
		return nil
	})
}

// TODO: sources for the rest of the types

*/
