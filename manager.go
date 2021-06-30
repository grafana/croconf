package croconf

import (
	"errors"
	"fmt"
	"strings"
)

type Manager struct {
	sources      []Source
	fields       []Field
	fieldsByDest map[interface{}]Field
	// TODO: internal data structure for tracking things
}

func NewManager() *Manager {
	return &Manager{
		fields:       make([]Field, 0),
		fieldsByDest: make(map[interface{}]Field),
	}
}

func (m *Manager) GetManager() *Manager {
	return m
}

func (m *Manager) AddField(field Field, options ...FieldOption) {
	// TODO: apply options?
	m.fields = append(m.fields, field)
	m.fieldsByDest[field.Destination()] = field
}

func (m *Manager) AddSource(source Source) {
	m.sources = append(m.sources, source)
}

func (m *Manager) Field(dest interface{}) Field {
	return m.fieldsByDest[dest]
}

func (m *Manager) Consolidate() error {
	var errs []error

	for _, s := range m.sources {
		err := s.Initialize()
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return consolidateErrorMessage(errs, "Config errors: ")
	}

	for _, f := range m.fields {
		errs = append(errs, f.Consolidate()...)
	}

	return consolidateErrorMessage(errs, "Config value errors: ")
}

func consolidateErrorMessage(errList []error, title string) error {
	if len(errList) == 0 {
		return nil
	}

	errMsgParts := []string{title}
	for _, err := range errList {
		errMsgParts = append(errMsgParts, fmt.Sprintf("\t- %s", err.Error()))
	}

	return errors.New(strings.Join(errMsgParts, "\n"))
}
