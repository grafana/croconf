package croconf

type SourceJSON struct {
	// TODO: I'm thinking that when this receives a JSON file, it should parse
	// it to a map[string]json.RawMessage. Then, it can parse every
	// json.RawMessage on demand, to the type specified by its `name` (set in
	// `From()` below) and
}

func (sj *SourceJSON) ParseAndApply() error {
	return nil // TODO
}

func (sj *SourceJSON) From(name string) MultiSingleValueSource {
	// TODO: this actually returns a closure
	return nil
}
