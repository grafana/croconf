package croconf

import (
	"encoding"
	"encoding/json"
	"fmt"
)

// TODO: use json.Decoder for this? json.Unmarshal() is a bit too magical

// TODO: rename this to something else? e.g. JSONDocument?
type SourceJSON struct {
	fields map[string]json.RawMessage
	init   func() error
}

func NewJSONSource(data []byte) *SourceJSON {
	fields := make(map[string]json.RawMessage)
	return &SourceJSON{
		fields: fields,
		init: func() error {
			// TODO: differentiate between an empty data and no data (nil)?
			if len(data) == 0 {
				return nil
			}
			err := json.Unmarshal(data, &fields)
			if err != nil {
				return NewJSONSourceInitError(data, err)
			}
			return nil
		},
	}
}

func (sj *SourceJSON) Initialize() error {
	return sj.init()
}

func (sj *SourceJSON) GetName() string {
	return "json"
}

func (sj *SourceJSON) Lookup(name string) (json.RawMessage, bool) {
	res, ok := sj.fields[name]
	return res, ok
}

func (sj *SourceJSON) From(name string) *jsonBinding {
	return &jsonBinding{
		source: sj,
		name:   name,
		lookup: func() (json.RawMessage, error) {
			raw, ok := sj.Lookup(name)
			if !ok {
				return nil, NewBindFieldMissingError(sj.GetName(), name)
			}
			return raw, nil
		},
	}
}

// TODO: export and rename? e.g. to JSONProperty?
type jsonBinding struct {
	source Source
	lookup func() (json.RawMessage, error)
	name   string
}

func (jb *jsonBinding) GetSource() Source {
	return jb.source
}

func (jb *jsonBinding) From(name string) *jsonBinding {
	return &jsonBinding{
		source: jb.source,
		name:   jb.name + "." + name,
		lookup: func() (json.RawMessage, error) {
			raw, err := jb.lookup()
			if err != nil {
				return nil, err
			}

			// TODO: cache this, so we don't parse sub-configs multiple times
			subdoc := NewJSONSource(raw)
			if err := subdoc.init(); err != nil {
				return nil, NewJSONSourceInitError(raw, err)
			}

			rawEl, ok := subdoc.Lookup(name)
			if !ok {
				return nil, NewBindFieldMissingError(subdoc.GetName(), name)
			}
			return rawEl, nil
		},
	}
}

func (jb *jsonBinding) BindStringValueTo(dest *string) func() error {
	return func() error {
		raw, err := jb.lookup()
		if err != nil {
			return err
		}

		return json.Unmarshal(raw, dest) // TODO: less reflection, better error messages
	}
}

func (jb *jsonBinding) BindIntValue() func(bitSize int) (int64, error) {
	return func(bitSize int) (int64, error) {
		raw, err := jb.lookup()
		if err != nil {
			// TODO: we might want to integrate custom error into lookup() method
			return 0, NewBindFieldMissingError(jb.source.GetName(), jb.name)
		}
		intVal, bindErr := parseInt(string(raw), 10, bitSize)
		if bindErr != nil {
			return 0, bindErr.withFuncName("BindIntValue")
		}
		return intVal, nil
	}
}

func (jb *jsonBinding) BindUintValue() func(bitSize int) (uint64, error) {
	return func(bitSize int) (uint64, error) {
		raw, err := jb.lookup()
		if err != nil {
			// TODO: we might want to integrate custom error into lookup() method
			return 0, NewBindFieldMissingError(jb.source.GetName(), jb.name)
		}
		intVal, bindErr := parseUint(string(raw), 10, bitSize)
		if bindErr != nil {
			return 0, bindErr.withFuncName("BindIntValue")
		}
		return intVal, nil
	}
}

func (jb *jsonBinding) BindFloatValue() func(bitSize int) (float64, error) {
	return func(bitSize int) (float64, error) {
		raw, err := jb.lookup()
		if err != nil {
			// TODO: we might want to integrate custom error into lookup() method
			return 0, NewBindFieldMissingError(jb.source.GetName(), jb.name)
		}
		intVal, bindErr := parseFloat(string(raw), bitSize)
		if bindErr != nil {
			return 0, bindErr.withFuncName("BindIntValue")
		}
		return intVal, nil
	}
}

func (jb *jsonBinding) BindBoolValueTo(dest *bool) func() error {
	return func() error {
		raw, err := jb.lookup()
		if err != nil {
			return err
		}

		return json.Unmarshal(raw, dest) // TODO: less reflection, better error messages
	}
}

func (jb *jsonBinding) BindTextBasedValueTo(dest encoding.TextUnmarshaler) func() error {
	return func() error {
		raw, err := jb.lookup()
		if err != nil {
			return err
		}

		// Progressive enhancement ¯\_(ツ)_/¯ If the destination supports directly
		// unmarshaling JSON, we should use that. Otherwise, we will fall back to
		// the simple text unmarshaling we know we can rely on.
		if jum, ok := dest.(json.Unmarshaler); ok {
			return jum.UnmarshalJSON(raw)
		}

		rawLen := len(raw)
		if rawLen < 2 || raw[0] != '"' || raw[rawLen-1] != '"' {
			return fmt.Errorf("expected a string when parsing JSON value for %s, got '%s'", jb.name, string(raw))
		}

		return dest.UnmarshalText(raw[1 : rawLen-1])
	}
}

func (jb *jsonBinding) To(dest json.Unmarshaler) *jsonBindingWithDest {
	return &jsonBindingWithDest{jsonBinding: jb, dest: dest}
}

type jsonBindingWithDest struct {
	*jsonBinding
	dest json.Unmarshaler
}

func (jbd *jsonBindingWithDest) BindValue() func() error {
	return func() error {
		raw, err := jbd.lookup()
		if err != nil {
			return err
		}

		return jbd.dest.UnmarshalJSON(raw)
	}
}

func (jb *jsonBinding) BindArray() func() (Array, error) {
	return func() (Array, error) {
		raw, err := jb.lookup()
		if err != nil {
			return nil, err
		}

		var rawArr []json.RawMessage
		if err := json.Unmarshal(raw, &rawArr); err != nil { // TODO: better error message
			return nil, err
		}

		return &jsonArrBinding{jb: jb, array: rawArr}, nil
	}
}

type jsonArrBinding struct {
	jb    *jsonBinding
	array []json.RawMessage
}

func (jba *jsonArrBinding) Len() int {
	return len(jba.array)
}

func (jba *jsonArrBinding) Element(elNum int) LazySingleValueBinding {
	name := fmt.Sprintf("%s[%d]", jba.jb.name, elNum)
	return &jsonBinding{
		source: jba.jb.source,
		name:   name,
		lookup: func() (json.RawMessage, error) {
			if elNum >= len(jba.array) {
				return nil, fmt.Errorf("tried to access invalid element %s, array only has %d elements", name, elNum)
			}
			return jba.array[elNum], nil
		},
	}
}
