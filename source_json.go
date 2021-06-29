package croconf

import (
	"encoding"
	"encoding/json"
	"fmt"
)

type SourceJSON struct {
	fields map[string]json.RawMessage
	// TODO: I'm thinking that when this receives a JSON file, it should parse
	// it to a map[string]json.RawMessage. Then, it can parse every
	// json.RawMessage on demand (i.e. lazily), to the type specified by its
	// `name` (set in `From()` below)
}

func NewJSONSource(data []byte) (*SourceJSON, error) {
	fields := make(map[string]json.RawMessage)

	if len(data) > 0 {
		// TODO: differentiate between an empty data and no data (nil)?
		if err := json.Unmarshal(data, &fields); err != nil {
			return nil, err
		}
	}
	return &SourceJSON{fields: fields}, nil
}

func (sj *SourceJSON) GetName() string {
	return "json"
}

func (sj *SourceJSON) From(name string) *jsonBinding {
	return &jsonBinding{
		source: sj,
		name:   name,
	}
}

type jsonBinding struct {
	source *SourceJSON
	name   string
}

func (jb *jsonBinding) GetSource() Source {
	return jb.source
}

func (jb *jsonBinding) BindStringValueTo(dest *string) func() error {
	return func() error {
		raw, ok := jb.source.fields[jb.name]
		if !ok {
			return ErrorMissing // TODO: better error message, e.g. 'field %s is not present in %s'?
		}

		return json.Unmarshal(raw, dest) // TODO: less reflection, better error messages
	}
}

func (jb *jsonBinding) BindInt64ValueTo(dest *int64) func() error {
	return func() error {
		raw, ok := jb.source.fields[jb.name]
		if !ok {
			return ErrorMissing // TODO: better error message, e.g. 'field %s is not present in %s'?
		}

		return json.Unmarshal(raw, dest) // TODO: less reflection, better error messages
	}
}

func (jb *jsonBinding) BindTextBasedValueTo(dest encoding.TextUnmarshaler) func() error {
	return func() error {
		raw, ok := jb.source.fields[jb.name]
		if !ok {
			return ErrorMissing // TODO: better error message, e.g. 'field %s is not present in %s'?
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

func (jb *jsonBinding) BindValue(dest interface{}) func() error {
	return func() error {
		raw, ok := jb.source.fields[jb.name]
		if !ok {
			return ErrorMissing // TODO: better error message, e.g. 'field %s is not present in %s'?
		}
		tdest, ok := dest.(json.Unmarshaler)
		if !ok {
			return json.Unmarshal(raw, dest)
		}
		return tdest.UnmarshalJSON(raw)
	}
}

//func (jb *jsonBinding) To(dest json.Unmarshaler) *jsonBindingWithDest {
	//return &jsonBindingWithDest{jsonBinding: jb, dest: dest}
//}

//type jsonBindingWithDest struct {
	//*jsonBinding
	//dest json.Unmarshaler
//}

//func (jbd *jsonBindingWithDest) BindValue(dest interface{}) func() error {
	//return func() error {
		//raw, ok := jbd.source.fields[jbd.name]
		//if !ok {
			//return ErrorMissing // TODO: better error message, e.g. 'field %s is not present in %s'?
		//}
		//tdest, ok := dest.(json.Unmarshaler)
		//if !ok {
			//return json.Unmarshal(raw, dest)
		//}
		//return tdest.UnmarshalJSON(raw)
	//}
//}
