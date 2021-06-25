package croconf

type Manager struct {
	//TODO: internal data structure for tracking things
}

func (m *Manager) GetManager() *Manager {
	return m
}

func (m *Manager) Int64Field(field *int64, defaultVal int64, sources ...Int64ValueSource) {
	//TODO: actually track this field in some internal data structure and set its default value
}

func (m *Manager) StringField(field *string, defaultVal string, sources ...StringValueSource) {
	//TODO: actually track this field in some internal data structure and set its default value
}

func (m *Manager) HasChanged(field *interface{}) bool {
	return false //TODO
}

func NewManager(...ConfigSource) *Manager {
	return &Manager{
		//TODO: initialize internal data structs
	}
}
