package croconf

import (
	"fmt"
)

type SourceGoMap struct {
	fields map[string]interface{}
}

func NewGoMapSource(fields map[string]interface{}) *SourceGoMap {
	return &SourceGoMap{fields: fields}
}

func (sm *SourceGoMap) GetName() string {
	return "go map"
}

func (sm *SourceGoMap) Initialize() error { return nil }

func (sm *SourceGoMap) Lookup(name string) (interface{}, bool) {
	res, ok := sm.fields[name]
	if !ok {
		return nil, false
	}
	return res, ok
}

func (sm *SourceGoMap) From(name string) *gomapBinder {
	return &gomapBinder{
		source: sm,
		name:   name,
		lookup: func() (interface{}, error) {
			raw, ok := sm.Lookup(name)
			if !ok {
				return nil, NewBindFieldMissingError(sm.GetName(), name)
			}
			return raw, nil
		},
	}
}

// TODO: export and rename? e.g. to JSONProperty?
type gomapBinder struct {
	source Source
	lookup func() (interface{}, error)
	name   string
}

func (mb *gomapBinder) From(name string) *gomapBinder {
	return &gomapBinder{
		source: mb.source,
		name:   mb.name + "." + name,
		lookup: func() (interface{}, error) {
			raw, err := mb.lookup()
			if err != nil {
				return nil, err
			}

			data, ok := raw.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("key=%s must be of a map[string]interface{} type", mb.name)
			}

			// TODO: cache this, so we don't parse sub-configs multiple times
			submap := NewGoMapSource(data)

			rawEl, ok := submap.Lookup(name)
			if !ok {
				return nil, NewBindFieldMissingError(submap.GetName(), name)
			}
			return rawEl, nil
		},
	}
}

func (mb *gomapBinder) newBinding(apply func() error) *gomapBinding {
	return &gomapBinding{
		binder: mb,
		apply:  apply,
	}
}

///
//
//

func (mb *gomapBinder) BindStringValueTo(dest *string) Binding {
	return mb.newBinding(func() error {
		raw, err := mb.lookup()
		if err != nil {
			return err
		}

		val, ok := raw.(string)
		if !ok {
			return NewBindValueError("BindStringValueTo", fmt.Sprintf("%v", raw), fmt.Errorf("casting string failed"))
		}
		*dest = val
		return nil
	})
}

func (mb *gomapBinder) BindIntValueTo(dest *int64) Binding {
	return mb.newBinding(func() error {
		raw, err := mb.lookup()
		if err != nil {
			return NewBindFieldMissingError(mb.source.GetName(), mb.name)
		}
		switch val := raw.(type) {
		case int:
			*dest = int64(val)
		case int8:
			*dest = int64(val)
		case int16:
			*dest = int64(val)
		case int32:
			*dest = int64(val)
		case int64:
			*dest = val
		default:
			return NewBindValueError("BindIntValueTo", fmt.Sprintf("%v", raw), fmt.Errorf("casting any int* type failed"))
		}
		return nil
	})
}

func (mb *gomapBinder) BindUintValueTo(dest *uint64) Binding {
	return mb.newBinding(func() error {
		raw, err := mb.lookup()
		if err != nil {
			return NewBindFieldMissingError(mb.source.GetName(), mb.name)
		}
		switch val := raw.(type) {
		case uint:
			*dest = uint64(val)
		case uint8:
			*dest = uint64(val)
		case uint16:
			*dest = uint64(val)
		case uint32:
			*dest = uint64(val)
		case uint64:
			*dest = val
		default:
			return NewBindValueError("BindUintValueTo", fmt.Sprintf("%v", raw), fmt.Errorf("casting any uint* type failed"))
		}
		return nil
	})
}

func (mb *gomapBinder) BindFloatValueTo(dest *float64) Binding {
	return mb.newBinding(func() error {
		raw, err := mb.lookup()
		if err != nil {
			return NewBindFieldMissingError(mb.source.GetName(), mb.name)
		}
		switch val := raw.(type) {
		case float32:
			*dest = float64(val)
		case float64:
			*dest = val
		default:
			return NewBindValueError("BindFloatValueTo", fmt.Sprintf("%v", raw), fmt.Errorf("casting any float* type failed"))
		}
		return nil
	})
}

func (mb *gomapBinder) BindBoolValueTo(dest *bool) Binding {
	return mb.newBinding(func() error {
		raw, err := mb.lookup()
		if err != nil {
			return err
		}

		v, ok := raw.(bool)
		if !ok {
			return NewBindValueError("BindBoolValueTo", fmt.Sprintf("%v", raw), fmt.Errorf("casting bool failed"))
		}
		*dest = v
		return nil
	})
}

// TODO: BindArrayValueTo
// func (mb *gomapBinder) BindArrayValueTo(length *int, element *func(int) LazySingleValueBinder) Binding {
// }

type gomapBinding struct {
	binder *gomapBinder
	apply  func() error
}

var _ interface {
	Binding
	BindingFromSource
} = &gomapBinding{}

func (gmb *gomapBinding) Apply() error {
	return gmb.apply()
}

func (gmb *gomapBinding) Source() Source {
	return gmb.binder.source
}

func (gmb *gomapBinding) BoundName() string {
	return gmb.binder.name
}
