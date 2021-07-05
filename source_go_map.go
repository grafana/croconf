package croconf

import (
	"fmt"
)

type SourceGoMap struct {
	values map[string]interface{}
}

func NewGoMapSource(values map[string]interface{}) *SourceGoMap {
	return &SourceGoMap{values: values}
}

// TODO: implement something like https://github.com/mitchellh/mapstructure

func (sm *SourceGoMap) GetName() string {
	return "go map"
}

func (sm *SourceGoMap) Initialize() error { return nil }

func (sm *SourceGoMap) Lookup(name string) (interface{}, bool) {
	res, ok := sm.values[name]
	return res, ok
}

func (sm *SourceGoMap) From(name string) *gomapBinding {
	return &gomapBinding{
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

type gomapBinding struct {
	source Source
	lookup func() (interface{}, error)
	name   string
}

func (mb *gomapBinding) GetSource() Source {
	return mb.source
}

func (mb *gomapBinding) From(name string) *gomapBinding {
	return &gomapBinding{
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

func (mb *gomapBinding) BindStringValueTo(dest *string) func() error {
	return func() error {
		raw, err := mb.lookup()
		if err != nil {
			return err
		}

		val, ok := raw.(string)
		if !ok {
			return fmt.Errorf("failed cast value for key=%s as a string", mb.name)
		}
		*dest = val
		return nil

	}
}

func (mb *gomapBinding) BindIntValue() func(bitSize int) (int64, error) {
	return func(bitSize int) (int64, error) {
		raw, err := mb.lookup()
		if err != nil {
			return 0, NewBindFieldMissingError(mb.source.GetName(), mb.name)
		}

		var val int64
		switch bitSize {
		case 0:
			v, ok := raw.(int)
			if !ok {
				return 0, fmt.Errorf("failed cast value for key=%s as a int", mb.name)
			}
			val = int64(v)

		case 8:
			v, ok := raw.(int8)
			if !ok {
				return 0, fmt.Errorf("failed cast value for key=%s as a int8", mb.name)
			}
			val = int64(v)

		case 16:
			v, ok := raw.(int16)
			if !ok {
				return 0, fmt.Errorf("failed cast value for key=%s as a int16", mb.name)
			}
			val = int64(v)

		case 32:
			v, ok := raw.(int32)
			if !ok {
				return 0, fmt.Errorf("failed cast value for key=%s as a int32", mb.name)
			}
			val = int64(v)

		case 64:
			v, ok := raw.(int64)
			if !ok {
				return 0, fmt.Errorf("failed cast value for key=%s as a int64", mb.name)
			}
			val = v
		}

		return val, nil
	}
}

func (mb *gomapBinding) BindUintValue() func(bitSize int) (uint64, error) {
	return func(bitSize int) (uint64, error) {
		raw, err := mb.lookup()
		if err != nil {
			return 0, NewBindFieldMissingError(mb.source.GetName(), mb.name)
		}

		var val uint64
		switch bitSize {
		case 0:
			v, ok := raw.(uint)
			if !ok {
				return 0, fmt.Errorf("failed cast value for key=%s as a uint", mb.name)
			}
			val = uint64(v)

		case 8:
			v, ok := raw.(uint8)
			if !ok {
				return 0, fmt.Errorf("failed cast value for key=%s as a uint8", mb.name)
			}
			val = uint64(v)

		case 16:
			v, ok := raw.(uint16)
			if !ok {
				return 0, fmt.Errorf("failed cast value for key=%s as a uint16", mb.name)
			}
			val = uint64(v)

		case 32:
			v, ok := raw.(uint32)
			if !ok {
				return 0, fmt.Errorf("failed cast value for key=%s as a uint32", mb.name)
			}
			val = uint64(v)

		case 64:
			v, ok := raw.(uint64)
			if !ok {
				return 0, fmt.Errorf("failed cast value for key=%s as a uint64", mb.name)
			}
			val = v
		}

		return val, nil
	}
}

func (mb *gomapBinding) BindFloatValue() func(bitSize int) (float64, error) {
	return func(bitSize int) (float64, error) {
		raw, err := mb.lookup()
		if err != nil {
			return 0, NewBindFieldMissingError(mb.source.GetName(), mb.name)
		}

		var val float64
		switch bitSize {
		case 32:
			v, ok := raw.(float32)
			if !ok {
				return 0, fmt.Errorf("failed cast value for key=%s as a float32", mb.name)
			}
			val = float64(v)

		case 64:
			v, ok := raw.(float64)
			if !ok {
				return 0, fmt.Errorf("failed cast value for key=%s as a float64", mb.name)
			}
			val = v
		}

		return val, nil
	}
}

func (mb *gomapBinding) BindBoolValueTo(dest *bool) func() error {
	return func() error {
		raw, err := mb.lookup()
		if err != nil {
			return err
		}

		v, ok := raw.(bool)
		if !ok {
			return fmt.Errorf("failed cast value for key=%s as a bool", mb.name)
		}
		*dest = v
		return nil
	}
}
