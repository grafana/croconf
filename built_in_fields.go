package croconf

import (
	"errors"
)

type field struct {
	hasBeenSet  bool
	source      Source
	consolidate func() []error
	destination interface{}
}

func (f *field) Consolidate() []error {
	return f.consolidate()
}

func (f *field) HasBeenSet() bool {
	return f.hasBeenSet
}

func (f *field) SourceOfValue() Source {
	return f.source
}

func (f *field) Destination() interface{} {
	return f.destination
}

func newField(dest interface{}, sourcesLen int, callback func(sourceNum int) (SourceGetter, error)) *field {
	// TODO: figure out some way to improve this?
	f := &field{destination: dest}
	f.consolidate = func() []error {
		var errs []error
		for i := 0; i < sourcesLen; i++ {
			sourceGetter, err := callback(i)
			if err == nil {
				f.source = sourceGetter.GetSource()
				if f.source != nil {
					f.hasBeenSet = true
				}
			} else if !errors.Is(ErrorMissing, err) {
				errs = append(errs, err)
			}
		}
		return errs
	}
	return f
}

func NewInt64Field(dest *int64, sources ...Int64ValueSource) Field {
	return newField(dest, len(sources), func(sourceNum int) (SourceGetter, error) {
		return sources[sourceNum], sources[sourceNum].SaveInt64To(dest)
	})
}

func NewStringField(dest *string, sources ...StringValueSource) Field {
	return newField(dest, len(sources), func(sourceNum int) (SourceGetter, error) {
		return sources[sourceNum], sources[sourceNum].SaveStringTo(dest)
	})
}
