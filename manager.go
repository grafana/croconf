package croconf

type Field interface {
}

func NewInt64Field(field *int64, sources ...Int64ValueSource) Field {
	return nil
}

func NewStringField(field *string, sources ...StringValueSource) Field {
	return nil
}

type Manager struct {
	fields []Field
	// TODO: internal data structure for tracking things
}

func (m *Manager) MarshalJSON() ([]byte, error) {
	// TODO:
	// do you emit defaut value (0, "") or null or nothing for unset fields
	return []byte(""), nil
}

func (m *Manager) GetManager() *Manager {
	return m
}

type FieldOption func(field Field)

func (m *Manager) AddField(field Field, options ...FieldOption) {
	// TODO: actually track this field in some internal data structure and set its default value
	m.fields = append(m.fields, field)
}

func (m *Manager) HasChanged(field *interface{}) bool {
	return false //TODO
}

func NewManager(...ConfigSource) *Manager {
	return &Manager{
		//TODO: initialize internal data structs
	}
}
