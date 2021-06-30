package croconf

type SourceGoMap struct {
	values map[string]interface{}
}

func NewGoMapSource(values map[string]interface{}) (*SourceGoMap, error) {
	return &SourceGoMap{values: values}, nil
}

// TODO: implement something like https://github.com/mitchellh/mapstructure
