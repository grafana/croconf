package croconf

import (
	"errors"
	"fmt"
)

type ManagedField struct {
	Field

	lastBindingFromSource BindingFromSource // nil for default value

	Name        string
	Description string
	Required    bool
	Validator   func() error
	// TODO: other meta information? e.g. deprecation warnings, usage
	// information and examples, annotations, etc.
}

func (mf *ManagedField) ApplyDefault() error {
	for _, binding := range mf.Field.Bindings() {
		if fromSource, ok := binding.(BindingFromSource); ok {
			// nil source means the default value
			if fromSource.Source() == nil {
				return fromSource.Apply()
			}
		}
	}
	return nil
}

func (mf *ManagedField) Consolidate() []error {
	// TODO: verify that sources have been initialized
	var errs []error
	for _, binding := range mf.Field.Bindings() {
		err := binding.Apply()
		if err == nil {
			if fromSource, ok := binding.(BindingFromSource); ok {
				mf.lastBindingFromSource = fromSource
			}
			continue
		}
		var bindErr *BindFieldMissingError
		if !errors.Is(ErrorMissing, err) && !errors.As(err, &bindErr) {
			errs = append(errs, err)
		}
	}
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

type FieldOption func(*ManagedField)

func WithName(name string) FieldOption {
	return func(mfield *ManagedField) {
		mfield.Name = name
	}
}

func WithDescription(description string) FieldOption {
	return func(mfield *ManagedField) {
		mfield.Description = description
	}
}

func WithValidator(validator func() error) FieldOption {
	return func(mfield *ManagedField) {
		mfield.Validator = validator
	}
}

func IsRequired() FieldOption {
	return func(mfield *ManagedField) {
		mfield.Required = true
	}
}
