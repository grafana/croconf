package croconf

import (
	"encoding"
	"errors"
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
		} else if !errors.Is(ErrorMissing, err) {
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

func NewInt64Field(dest *int64, sources ...Int64ValueBinding) Field {
	return newField(dest, len(sources), func(sourceNum int) valueBinding {
		return vb(sources[sourceNum], sources[sourceNum].BindInt64ValueTo(dest))
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

func NewCustomField(dest interface{}, sources ...CustomValueBinding) Field {
	return newField(dest, len(sources), func(sourceNum int) valueBinding {
		return vb(sources[sourceNum], sources[sourceNum].BindValue())
	})
}
