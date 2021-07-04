package croconf

import (
	"encoding"
)

type field struct {
	destination interface{}
	bindings    []Binding
}

func (f *field) Destination() interface{} {
	return f.destination
}

func (f *field) Bindings() []Binding {
	return f.bindings
}

func newField(dest interface{}, sourcesLen int, callback func(sourceNum int) Binding) *field {
	f := &field{
		destination: dest,
		bindings:    make([]Binding, sourcesLen),
	}
	for i := 0; i < sourcesLen; i++ {
		f.bindings[i] = callback(i)
	}

	return f
}

type arrayHandler func(arrLength int, getElement func(int) LazySingleValueBinder) error

func newArrayField(dest interface{}, sources []ArrayValueBinder, handler arrayHandler) Field {
	return newField(dest, len(sources), func(sourceNum int) Binding {
		source := sources[sourceNum]
		var arrLength int
		var getElement func(int) LazySingleValueBinder
		binding := source.BindArrayValueTo(&arrLength, &getElement)
		return wrapBinding(binding, func() error {
			err := binding.Apply()
			if err != nil {
				return err
			}
			return handler(arrLength, getElement)
		})
	})
}

func NewStringField(dest *string, sources ...StringValueBinder) Field {
	return newField(dest, len(sources), func(sourceNum int) Binding {
		return sources[sourceNum].BindStringValueTo(dest)
	})
}

func NewTextBasedField(dest encoding.TextUnmarshaler, sources ...TextBasedValueBinder) Field {
	return newField(dest, len(sources), func(sourceNum int) Binding {
		return sources[sourceNum].BindTextBasedValueTo(dest)
	})
}

func NewBoolField(dest *bool, sources ...BoolValueBinder) Field {
	return newField(dest, len(sources), func(sourceNum int) Binding {
		return sources[sourceNum].BindBoolValueTo(dest)
	})
}

func NewCustomField(dest interface{}, sources ...CustomValueBinder) Field {
	return newField(dest, len(sources), func(sourceNum int) Binding {
		return sources[sourceNum].BindValue()
	})
}
