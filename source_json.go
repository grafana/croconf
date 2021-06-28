package croconf

import (
	"encoding/json"
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

func (sj *SourceJSON) From(name string) LazySingleValueBinding {
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

func (jb *jsonBinding) SaveStringTo(dest *string) error {
	raw, ok := jb.source.fields[jb.name]
	if !ok {
		return ErrorMissing // TODO: better error message, e.g. 'field %s is not present in %s'?
	}

	return json.Unmarshal(raw, dest) // TODO: less reflection, better error messages
}

func (jb *jsonBinding) SaveInt64To(dest *int64) error {
	raw, ok := jb.source.fields[jb.name]
	if !ok {
		return ErrorMissing // TODO: better error message, e.g. 'field %s is not present in %s'?
	}

	return json.Unmarshal(raw, dest) // TODO: less reflection, better error messages
}

func (jb *jsonBinding) SaveCustomValueTo(dest CustomValue) error {
	raw, ok := jb.source.fields[jb.name]
	if !ok {
		return ErrorMissing // TODO: better error message, e.g. 'field %s is not present in %s'?
	}

	return json.Unmarshal(raw, dest) // TODO: less reflection, better error messages
}
