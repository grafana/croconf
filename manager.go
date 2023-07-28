package croconf

import (
	"errors"
	"fmt"
	"strings"
)

type Manager struct {
	sources            []Source
	seenSources        map[Source]struct{}
	fields             []Field
	consolidatedFields []ConsolidatedField
	fieldsByDest       map[interface{}]int

	defaultSourceOfFieldNames Source
}

type ManagerOption func(*Manager)

func NewManager(options ...ManagerOption) *Manager {
	m := &Manager{
		fields:       make([]Field, 0),
		fieldsByDest: make(map[interface{}]int),
		sources:      make([]Source, 0),
		seenSources:  make(map[Source]struct{}),
	}

	for _, opt := range options {
		opt(m)
	}

	return m
}

func (m *Manager) deriveFieldName(fieldIndex int) string {
	field := m.fields[fieldIndex]
	var firstCanonicalBinding, firstNonDefaultBinding Binding
	for _, binding := range field.Bindings() {
		source := binding.Source()
		if source != nil && firstCanonicalBinding == nil && source == m.defaultSourceOfFieldNames {
			firstCanonicalBinding = binding
		}
		if source != nil && firstNonDefaultBinding == nil {
			firstNonDefaultBinding = binding
		}
	}

	if firstCanonicalBinding != nil {
		return firstCanonicalBinding.Identifier()
	}

	if firstNonDefaultBinding != nil {
		return firstNonDefaultBinding.Identifier()
	}

	// TODO: some better fallback? for example, a way to get field names from
	// the config struct that is being managed, through reflect. See
	// https://golang.org/pkg/reflect/#Value.NumField and
	// https://golang.org/pkg/reflect/#Value.Field. Something like this:
	//
	// func WithFieldNamesFromGoStruct(config interface{}) ManagerOption {}
	return fmt.Sprintf("field-%d", fieldIndex)
}

func (m *Manager) AddField(field Field) {
	idx := len(m.fields)
	m.fieldsByDest[field.Destination()] = idx
	m.fields = append(m.fields, field)

	m.addSources(field)
}

func (m *Manager) addSources(field Field) {
	for _, b := range field.Bindings() {
		s := b.Source()
		if _, seen := m.seenSources[s]; s != nil && !seen {
			m.seenSources[s] = struct{}{}
			m.sources = append(m.sources, s)
		}
	}
}

func (m *Manager) Field(dest interface{}) Field {
	if idx, ok := m.fieldsByDest[dest]; ok {
		return m.fields[idx]
	}
	return nil
}

func (m *Manager) Fields() []Field {
	return m.fields
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

	m.consolidatedFields = make([]ConsolidatedField, len(m.fields))
	for i, f := range m.fields {
		cf, err := f.Consolidate()
		m.consolidatedFields[i] = cf

		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return consolidateErrorMessage(errs, "Config value errors: ")
	}

	for _, cf := range m.consolidatedFields {
		fieldErr := cf.Validate()
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

// WithDefaultSourceOfFieldNames designates a specific Source as the canonical
// source of field names. For example, that way all validation errors and help
// texts will always use the JSON property names or the CLI flag names.
func WithDefaultSourceOfFieldNames(source Source) ManagerOption {
	return func(m *Manager) {
		m.defaultSourceOfFieldNames = source
	}
}
