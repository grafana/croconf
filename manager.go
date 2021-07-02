package croconf

import (
	"errors"
	"fmt"
	"strings"
)

type Manager struct {
	sources      []Source
	fields       []*ManagedField
	fieldsByDest map[interface{}]*ManagedField
	// TODO: internal data structure for tracking things
}

func NewManager() *Manager {
	return &Manager{
		fields:       make([]*ManagedField, 0),
		fieldsByDest: make(map[interface{}]*ManagedField),
	}
}

func (m *Manager) GetManager() *Manager {
	return m
}

func (m *Manager) AddField(field Field, options ...FieldOption) *ManagedField {
	mf := &ManagedField{
		Field: field,
	}

	for _, opt := range options {
		opt(mf)
	}

	if mf.Name == "" {
		// TODO: add a way to designate a specific source as the canonical
		// source of field names, e.g. so that all validation errors contain the
		// JSON or CLI flag names
		mf.Name = fmt.Sprintf("field-%d", len(m.fields)+1)
	}

	m.fields = append(m.fields, mf)
	m.fieldsByDest[field.Destination()] = mf

	return mf
}

func (m *Manager) AddSource(source Source) {
	m.sources = append(m.sources, source)
}

func (m *Manager) Field(dest interface{}) *ManagedField {
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

	if len(errs) > 0 {
		return consolidateErrorMessage(errs, "Config value errors: ")
	}

	for _, f := range m.fields {
		fieldErr := f.Validate()
		if fieldErr != nil {
			errs = append(errs, fieldErr)
		}
	}

	return consolidateErrorMessage(errs, "Validation errors: ")
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
