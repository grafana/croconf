package croconf

import (
	"errors"
	"fmt"
	"reflect"
)

type ManagedField struct {
	Field

	wasConsolidated       bool
	lastBindingFromSource BindingFromSource // nil for default value

	Name         string
	DefaultValue string
	Description  string
	Required     bool
	Validator    func() error
	// TODO: other meta information? e.g. deprecation warnings, usage
	// information and examples, possible values, annotations, etc.
}

func (mf *ManagedField) getCurrentValueAsString() string {
	dest := mf.Destination()
	if stringer, ok := dest.(fmt.Stringer); ok {
		return stringer.String()
	}

	// TODO: check for encoding.TextMarshaler?

	// Since the destination is likely a pointer, we dereference it here
	value := reflect.Indirect(reflect.ValueOf(dest)).Interface()
	if stringer, ok := value.(fmt.Stringer); ok {
		return stringer.String()
	}

	// TODO: check for encoding.TextMarshaler?

	return fmt.Sprintf("%v", value)
}

func (mf *ManagedField) Consolidate() []error {
	if mf.wasConsolidated {
		return nil
	}
	// TODO: verify that sources have been initialized

	mf.DefaultValue = mf.getCurrentValueAsString()

	var errs []error
	for _, binding := range mf.Field.Bindings() {
		err := binding.Apply()
		if err == nil {
			if fromSource, ok := binding.(BindingFromSource); ok {
				mf.lastBindingFromSource = fromSource

				if fromSource.Source() == nil {
					// This was a default value
					mf.DefaultValue = mf.getCurrentValueAsString()
				}
			}
			continue
		}
		var bindErr *BindFieldMissingError
		if !errors.Is(ErrorMissing, err) && !errors.As(err, &bindErr) {
			errs = append(errs, err)
		}
	}
	mf.wasConsolidated = true
	return errs
}

func (mf *ManagedField) LastBindingFromSource() BindingFromSource {
	return mf.lastBindingFromSource
}

func (mf *ManagedField) HasBeenSetFromSource() bool {
	return mf.lastBindingFromSource != nil && mf.lastBindingFromSource.Source() != nil
}

func (mf *ManagedField) Validate() error {
	if mf.Required && !mf.HasBeenSetFromSource() {
		return fmt.Errorf("Field %s is required, but no value was set", mf.Name)
	}

	if mf.Validator != nil {
		return mf.Validator()
	}
	return nil
}

type ManagedFieldOption func(*ManagedField)

func WithName(name string) ManagedFieldOption {
	return func(mfield *ManagedField) {
		mfield.Name = name
	}
}

func WithDescription(description string) ManagedFieldOption {
	return func(mfield *ManagedField) {
		mfield.Description = description
	}
}

func WithValidator(validator func() error) ManagedFieldOption {
	return func(mfield *ManagedField) {
		mfield.Validator = validator
	}
}

func IsRequired() ManagedFieldOption {
	return func(mfield *ManagedField) {
		mfield.Required = true
	}
}
