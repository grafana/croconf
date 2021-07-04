package croconf

import (
	"encoding"
	"strconv"
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

func intValHelper(sources []IntValueBinder, bitSize int, saveToDest func(int64)) func(sourceNum int) Binding {
	return func(sourceNum int) Binding {
		var val int64
		binding := sources[sourceNum].BindIntValueTo(&val)

		return wrapBinding(binding, func() error {
			if err := binding.Apply(); err != nil {
				return err
			}
			if err := checkIntBitsize(val, bitSize); err != nil {
				return err
			}
			saveToDest(val)
			return nil
		})
	}
}

func NewIntField(dest *int, sources ...IntValueBinder) Field {
	return newField(dest, len(sources), intValHelper(sources, strconv.IntSize, func(val int64) {
		*dest = int(val) // this is safe, intValHelper checks val against bitSize
	}))
}

func NewInt8Field(dest *int8, sources ...IntValueBinder) Field {
	return newField(dest, len(sources), intValHelper(sources, 8, func(val int64) {
		*dest = int8(val) // this is safe, intValHelper checks val against bitSize
	}))
}

func NewInt16Field(dest *int16, sources ...IntValueBinder) Field {
	return newField(dest, len(sources), intValHelper(sources, 16, func(val int64) {
		*dest = int16(val) // this is safe, intValHelper checks val against bitSize
	}))
}

func NewInt32Field(dest *int32, sources ...IntValueBinder) Field {
	return newField(dest, len(sources), intValHelper(sources, 32, func(val int64) {
		*dest = int32(val) // this is safe, intValHelper checks val against bitSize
	}))
}

func NewInt64Field(dest *int64, sources ...IntValueBinder) Field {
	return newField(dest, len(sources), func(sourceNum int) Binding {
		return sources[sourceNum].BindIntValueTo(dest)
	})
}

func NewInt8SliceField(dest *[]int8, sources ...ArrayValueBinder) Field {
	// TODO: figure out some way to avoid the boilerplate?
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

			newArr := make([]int8, arrLength)
			for i := 0; i < arrLength; i++ {
				var val int64
				elBinding := getElement(i).BindIntValueTo(&val)
				if err := elBinding.Apply(); err != nil {
					return err
				}
				if err := checkIntBitsize(val, 8); err != nil {
					return err
				}
				newArr[i] = int8(val) // this is safe
			}
			*dest = newArr
			return nil
		})
	})
}

func NewInt64SliceField(dest *[]int64, sources ...ArrayValueBinder) Field {
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

			newArr := make([]int64, arrLength)
			for i := 0; i < arrLength; i++ {
				var val int64
				elBinding := getElement(i).BindIntValueTo(&val)
				if err := elBinding.Apply(); err != nil {
					return err
				}
				newArr[i] = val
			}
			*dest = newArr
			return nil
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
