package croconf

import (
	"errors"
	"fmt"
)

type ManagedField struct {
	Field

	wasConsolidated       bool
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
	if mf.wasConsolidated {
		return nil
	}
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
