package croconf

import (
	"encoding/json"
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

func (sj *SourceJSON) From(name string) *jsonBinder {
	return &jsonBinder{
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
type jsonBinder struct {
	source Source
	lookup func() (json.RawMessage, error)
	name   string
}

func (jb *jsonBinder) From(name string) *jsonBinder {
	return &jsonBinder{
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

func toTypedJSONBinding[T any](jb *jsonBinder) TypedBinding[T] {
	conv := func() (T, error) {
		var result T
		val, err := jb.lookup()
		if err != nil {
			return result, err
		}

		// TODO: use something else, incl. custom functions for common types, to
		// have less reflection and better error messages
		err = json.Unmarshal(val, &result)
		return result, err
	}

	return ToBinding(jb.name, jb.source, conv)
}

func (jb *jsonBinder) ToString() TypedBinding[string] {
	return toTypedJSONBinding[string](jb)
}

func (jb *jsonBinder) ToInt64() TypedBinding[int64] {
	return toTypedJSONBinding[int64](jb)
}

func (jb *jsonBinder) ToUint64() TypedBinding[uint64] {
	return toTypedJSONBinding[uint64](jb)
}

func (jb *jsonBinder) ToFloat64() TypedBinding[float64] {
	return toTypedJSONBinding[float64](jb)
}

func (jb *jsonBinder) ToBool() TypedBinding[bool] {
	return toTypedJSONBinding[bool](jb)
}

/*
func (jb *jsonBinder) BindTextBasedValueTo(dest encoding.TextUnmarshaler) Binding {
	return jb.newBinding(func() error {
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
	})
}

func (jb *jsonBinder) BindArrayValueTo(length *int, element *func(int) LazySingleValueBinder) Binding {
	return jb.newBinding(func() error {
		raw, err := jb.lookup()
		if err != nil {
			return err
		}

		var rawArr []json.RawMessage
		if err := json.Unmarshal(raw, &rawArr); err != nil { // TODO: better error message
			return err
		}

		*length = len(rawArr)
		*element = func(elNum int) LazySingleValueBinder {
			name := fmt.Sprintf("%s[%d]", jb.name, elNum)
			return &jsonBinder{
				source: jb.source,
				name:   name,
				lookup: func() (json.RawMessage, error) {
					if elNum >= len(rawArr) {
						return nil, fmt.Errorf("tried to access invalid element %s, array only has %d elements", name, elNum)
					}
					return rawArr[elNum], nil
				},
			}
		}
		return nil
	})
}
*/
