package croconf

import "encoding/json"

type SourceJSON struct {
	fields map[string]json.RawMessage
	// TODO: I'm thinking that when this receives a JSON file, it should parse
	// it to a map[string]json.RawMessage. Then, it can parse every
	// json.RawMessage on demand (i.e. lazily), to the type specified by its
	// `name` (set in `From()` below)
}

func NewJSONSource(data []byte) (*SourceJSON, error) {
	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &fields); err != nil {
		return nil, err
	}
	return &SourceJSON{fields: fields}, nil
}

func (sj *SourceJSON) ParseAndApply() error {
	return nil // TODO
}

func (sj *SourceJSON) From(name string) MultiSingleValueSource {
	// TODO: this actually returns a closure
	return nil
}
