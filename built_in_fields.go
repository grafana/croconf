package croconf

import (
	"encoding"
	"errors"
	"strconv"
)

type valueBinding struct {
	sourceGetter SourceGetter
	apply        func() error
}

func vb(sg SourceGetter, apply func() error) valueBinding {
	return valueBinding{sourceGetter: sg, apply: apply}
}

type field struct {
	source        Source
	destination   interface{}
	valueBindings []valueBinding
}

func (f *field) Consolidate() []error {
	var errs []error
	for _, vb := range f.valueBindings {
		err := vb.apply()
		if err == nil {
			f.source = vb.sourceGetter.GetSource()
			continue
		}
		var bindErr *BindFieldMissingError
		if !errors.Is(ErrorMissing, err) && !errors.As(err, &bindErr) {
			errs = append(errs, err)
		}
	}
	return errs
}

func (f *field) ValueSource() Source {
	return f.source
}

func (f *field) Destination() interface{} {
	return f.destination
}

func newField(dest interface{}, sourcesLen int, callback func(sourceNum int) valueBinding) *field {
	f := &field{
		destination:   dest,
		valueBindings: make([]valueBinding, sourcesLen),
	}
	for i := 0; i < sourcesLen; i++ {
		f.valueBindings[i] = callback(i)
	}

	return f
}

func intValHelper(sources []IntValueBinding, bitSize int, saveToDest func(int64)) func(sourceNum int) valueBinding {
	return func(sourceNum int) valueBinding {
		var val int64
		bind := sources[sourceNum].BindIntValueTo(&val)
		return vb(sources[sourceNum], func() error {
			if err := bind(); err != nil {
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

func NewIntField(dest *int, sources ...IntValueBinding) Field {
	return newField(dest, len(sources), intValHelper(sources, strconv.IntSize, func(val int64) {
		*dest = int(val) // this is safe, intValHelper checks val against bitSize
	}))
}

func NewInt8Field(dest *int8, sources ...IntValueBinding) Field {
	return newField(dest, len(sources), intValHelper(sources, 8, func(val int64) {
		*dest = int8(val) // this is safe, intValHelper checks val against bitSize
	}))
}

func NewInt16Field(dest *int16, sources ...IntValueBinding) Field {
	return newField(dest, len(sources), intValHelper(sources, 16, func(val int64) {
		*dest = int16(val) // this is safe, intValHelper checks val against bitSize
	}))
}

func NewInt32Field(dest *int32, sources ...IntValueBinding) Field {
	return newField(dest, len(sources), intValHelper(sources, 32, func(val int64) {
		*dest = int32(val) // this is safe, intValHelper checks val against bitSize
	}))
}

func NewInt64Field(dest *int64, sources ...IntValueBinding) Field {
	return newField(dest, len(sources), func(sourceNum int) valueBinding {
		return vb(sources[sourceNum], sources[sourceNum].BindIntValueTo(dest))
	})
}

func NewInt8SliceField(dest *[]int8, sources ...ArrayBinding) Field {
	// TODO: figure out some way to avoid the boilerplate?
	return newField(dest, len(sources), func(sourceNum int) valueBinding {
		source := sources[sourceNum]
		arrBind := source.BindArray()
		return vb(source, func() error {
			sourceArr, err := arrBind()
			if err != nil {
				return err
			}

			arrLen := sourceArr.Len()
			newArr := make([]int8, arrLen)
			for i := 0; i < arrLen; i++ {
				var val int64
				bind := sourceArr.Element(i).BindIntValueTo(&val)
				if err := bind(); err != nil {
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

func NewInt64SliceField(dest *[]int64, sources ...ArrayBinding) Field {
	return newField(dest, len(sources), func(sourceNum int) valueBinding {
		source := sources[sourceNum]
		arrBind := source.BindArray()
		return vb(source, func() error {
			sourceArr, err := arrBind()
			if err != nil {
				return err
			}

			arrLen := sourceArr.Len()
			newArr := make([]int64, arrLen)
			for i := 0; i < arrLen; i++ {
				var val int64
				bind := sourceArr.Element(i).BindIntValueTo(&val)
				if err := bind(); err != nil {
					return err
				}
				newArr[i] = val
			}
			*dest = newArr
			return nil
		})
	})
}

func NewStringField(dest *string, sources ...StringValueBinding) Field {
	return newField(dest, len(sources), func(sourceNum int) valueBinding {
		return vb(sources[sourceNum], sources[sourceNum].BindStringValueTo(dest))
	})
}

func NewTextBasedField(dest encoding.TextUnmarshaler, sources ...TextBasedValueBinding) Field {
	return newField(dest, len(sources), func(sourceNum int) valueBinding {
		return vb(sources[sourceNum], sources[sourceNum].BindTextBasedValueTo(dest))
	})
}

func NewBoolField(dest *bool, sources ...BoolValueBinding) Field {
	return newField(dest, len(sources), func(sourceNum int) valueBinding {
		return vb(sources[sourceNum], sources[sourceNum].BindBoolValueTo(dest))
	})
}

func NewCustomField(dest interface{}, sources ...CustomValueBinding) Field {
	return newField(dest, len(sources), func(sourceNum int) valueBinding {
		return vb(sources[sourceNum], sources[sourceNum].BindValue())
	})
}
