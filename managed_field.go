package croconf

import (
	"errors"
	"fmt"
	"reflect"
)

//TODO: call this TypedField and cut it in half, have a separate and simpler
//Field interface and a ManagedField struct that can wrap either? :/

type ManagedField[T any] struct {
	destination  *T
	defaultValue T
	bindings     []TypedBinding[T]
	validators   []func(T) error

	// TODO: add validation and source strategies (e.g. validate every value
	// from every source, not just the final one)

	// TODO: split some of these in a separate non-generic struct?

	// TODO: other meta information? e.g. deprecation warnings, usage
	// information and examples, possible values, annotations, etc.
	name        string
	description string
	required    bool
}

type ConsolidatedManagedField[T any] struct {
	mf          *ManagedField[T]
	lastBinding TypedBinding[T]
}

func NewField[T any](dest *T) *ManagedField[T] {
	return &ManagedField[T]{destination: dest}
}

var _ Field = &ManagedField[any]{}
var _ ConsolidatedField = &ConsolidatedManagedField[any]{}

func (mf *ManagedField[T]) WithDefault(val T) *ManagedField[T] {
	mf.defaultValue = val
	return mf
}

func (mf *ManagedField[T]) WithBinding(source TypedBinding[T]) *ManagedField[T] {
	mf.bindings = append(mf.bindings, source)
	return mf
}

func (mf *ManagedField[T]) WithValidator(validator func(T) error) *ManagedField[T] {
	mf.validators = append(mf.validators, validator)
	return mf
}

func (mf *ManagedField[T]) WithName(name string) *ManagedField[T] {
	mf.name = name
	return mf
}

func (mf *ManagedField[T]) getCurrentValueAsString() string {
	dest := any(mf.destination)
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

func (mf *ManagedField[T]) Destination() any {
	return mf.destination
}

func (mf *ManagedField[T]) Name() string {
	return mf.name
}

func (mf *ManagedField[T]) Bindings() []Binding {
	res := make([]Binding, len(mf.bindings))
	for i, b := range mf.bindings {
		res[i] = b
	}
	return res
}

func (mf *ManagedField[T]) Consolidate() (ConsolidatedField, error) {
	// TODO: verify that sources have been initialized

	cf := &ConsolidatedManagedField[T]{mf: mf}
	var errs []error
	for _, binding := range mf.bindings {
		val, err := binding.GetValue()
		if err == nil {
			*mf.destination = val
			cf.lastBinding = binding
			continue
		}
		var bindErr *BindFieldMissingError
		if !errors.Is(ErrorMissing, err) && !errors.As(err, &bindErr) {
			errs = append(errs, err)
		}
	}
	return cf, errors.Join(errs...) // TODO: use a custom type
}

func (cf *ConsolidatedManagedField[T]) HasBeenSet() bool {
	return cf.lastBinding != nil
}

func (cf *ConsolidatedManagedField[T]) Source() Source {
	if cf.lastBinding == nil {
		return nil
	}
	return cf.lastBinding.Source()
}

func (cf *ConsolidatedManagedField[T]) Validate() error {
	if cf.mf.required && !cf.HasBeenSet() {
		return fmt.Errorf("Field %s is required, but no value was set", cf.mf.name)
	}

	var errs []error
	for _, validator := range cf.mf.validators {
		errs = append(errs, validator(*cf.mf.destination))
	}
	return errors.Join(errs...) // TODO: use a custom type
}
