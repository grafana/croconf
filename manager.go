package croconf

import (
	"errors"
	"fmt"
	"strings"
)

type Manager struct {
	sources      []Source
	seenSources  map[Source]struct{}
	fields       []*ManagedField
	fieldsByDest map[interface{}]*ManagedField

	defaultSourceOfFieldNames Source
}

type ManagerOption func(*Manager)

func NewManager(options ...ManagerOption) *Manager {
	m := &Manager{
		fields:       make([]*ManagedField, 0),
		fieldsByDest: make(map[interface{}]*ManagedField),
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
	var firstCanonicalBinding, firstNonDefaultBinding BindingFromSource
	for _, binding := range field.Bindings() {
		if bindingFromSource, ok := binding.(BindingFromSource); ok {
			source := bindingFromSource.Source()
			if source != nil && firstCanonicalBinding == nil && source == m.defaultSourceOfFieldNames {
				firstCanonicalBinding = bindingFromSource
			}
			if source != nil && firstNonDefaultBinding == nil {
				firstNonDefaultBinding = bindingFromSource
			}
		}
	}

	if firstCanonicalBinding != nil {
		return firstCanonicalBinding.BoundName()
	}

	if firstNonDefaultBinding != nil {
		return firstNonDefaultBinding.BoundName()
	}

	// TODO: some better fallback? for example, a way to get field names from
	// the config struct that is being managed, through reflect. See
	// https://golang.org/pkg/reflect/#Value.NumField and
	// https://golang.org/pkg/reflect/#Value.Field. Something like this:
	//
	// func WithFieldNamesFromGoStruct(config interface{}) ManagerOption {}
	return fmt.Sprintf("field-%d", fieldIndex)
}

func (m *Manager) AddField(field Field, options ...ManagedFieldOption) *ManagedField {
	mf := &ManagedField{
		Field: field,
	}

	for _, opt := range options {
		opt(mf)
	}

	m.fields = append(m.fields, mf)
	m.fieldsByDest[field.Destination()] = mf

	m.addSources(mf)

	if mf.Name == "" {
		mf.Name = m.deriveFieldName(len(m.fields) - 1)
	}

	return mf
}

func (m *Manager) addSources(field Field) {
	for _, b := range field.Bindings() {
		if fromSource, ok := b.(BindingFromSource); ok {
			s := fromSource.Source()
			if _, seen := m.seenSources[s]; s != nil && !seen {
				m.seenSources[s] = struct{}{}
				m.sources = append(m.sources, s)
			}
		}
	}
}

func (m *Manager) Field(dest interface{}) *ManagedField {
	return m.fieldsByDest[dest]
}

func (m *Manager) Fields() []*ManagedField {
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

// WithDefaultSourceOfFieldNames designates a specific Source as the canonical
// source of field names. For example, that way all validation errors and help
// texts will always use the JSON property names or the CLI flag names.
func WithDefaultSourceOfFieldNames(source Source) ManagerOption {
	return func(m *Manager) {
		m.defaultSourceOfFieldNames = source
	}
}
