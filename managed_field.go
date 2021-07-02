package croconf

import "fmt"

type ManagedField struct {
	Field

	Name        string
	Description string
	Required    bool
	Validator   func() error
	// TODO: other meta information? e.g. deprecation warnings, usage
	// information and examples, annotations, etc.
}

func (mf *ManagedField) Validate() error {
	if mf.Required && mf.Field.ValueSource() == nil {
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
