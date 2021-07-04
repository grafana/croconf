package croconf

import (
	"encoding"
	"fmt"
	"strconv"
	"strings"
)

type SourceEnvVars struct {
	env map[string]string
	// TODO
}

func NewSourceFromEnv(environ []string) *SourceEnvVars {
	env := make(map[string]string, len(environ))
	for _, kv := range environ {
		k, v := parseEnvKeyValue(kv)
		env[k] = v
	}
	return &SourceEnvVars{env: env}
}

func (sev *SourceEnvVars) Initialize() error {
	return nil // TODO? maybe prefix handling?
}

func (sev *SourceEnvVars) GetName() string {
	return "environment variables" // TODO
}

func (sev *SourceEnvVars) From(name string) *envBinder {
	return &envBinder{
		source: sev,
		name:   name,
		lookup: func() (string, error) {
			val, ok := sev.env[name]
			if !ok {
				return "", ErrorMissing // TODO: better error message, e.g. 'field %s is not present in %s'?
			}
			return val, nil
		},
	}
}

type envBinder struct {
	source Source
	name   string
	lookup func() (string, error)
}

func (eb *envBinder) newBinding(apply func() error) *envBinding {
	return &envBinding{
		binder: eb,
		apply:  apply,
	}
}

func (eb *envBinder) BindStringValueTo(dest *string) Binding {
	return eb.newBinding(func() error {
		val, err := eb.lookup()
		if err != nil {
			return err
		}
		*dest = val
		return nil
	})
}

func (eb *envBinder) BindIntValueTo(dest *int64) Binding {
	return eb.newBinding(func() error {
		val, err := eb.lookup()
		if err != nil {
			// TODO: we might want to integrate custom error into lookup() method
			return NewBindFieldMissingError(eb.source.GetName(), eb.name)
		}
		intVal, bindErr := parseInt(val)
		if bindErr != nil {
			return bindErr.withFuncName("BindIntValue")
		}
		*dest = intVal
		return nil
	})
}

func (eb *envBinder) BindUintValueTo(dest *uint64) Binding {
	return eb.newBinding(func() error {
		val, err := eb.lookup()
		if err != nil {
			// TODO: we might want to integrate custom error into lookup() method
			return NewBindFieldMissingError(eb.source.GetName(), eb.name)
		}
		uintVal, bindErr := parseUint(val)
		if bindErr != nil {
			return bindErr.withFuncName("BindUintValue")
		}
		*dest = uintVal
		return nil
	})
}

func (eb *envBinder) BindFloatValueTo(dest *float64) Binding {
	return eb.newBinding(func() error {
		strVal, err := eb.lookup()
		if err != nil {
			// TODO: we might want to integrate custom error into lookup() method
			return NewBindFieldMissingError(eb.source.GetName(), eb.name)
		}
		val, bindErr := parseFloat(strVal)
		if bindErr != nil {
			return bindErr.withFuncName("BindFloatValue")
		}
		*dest = val
		return nil
	})
}

func (eb *envBinder) BindBoolValueTo(dest *bool) Binding {
	return eb.newBinding(func() error {
		val, err := eb.lookup()
		if err != nil {
			return err
		}
		b, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}
		*dest = b
		return nil
	})
}

func (eb *envBinder) BindTextBasedValueTo(dest encoding.TextUnmarshaler) Binding {
	return eb.newBinding(func() error {
		val, err := eb.lookup()
		if err != nil {
			return NewBindFieldMissingError(eb.source.GetName(), eb.name)
		}

		return dest.UnmarshalText([]byte(val))
	})
}

func (eb *envBinder) BindArrayValueTo(length *int, element *func(int) LazySingleValueBinder) Binding {
	return eb.newBinding(func() error {
		val, err := eb.lookup()
		if err != nil {
			return NewBindFieldMissingError(eb.source.GetName(), eb.name)
		}

		arr := strings.Split(val, ",") // TODO: figure out how to make the delimiter configurable

		*length = len(arr)
		*element = func(elNum int) LazySingleValueBinder {
			name := fmt.Sprintf("%s[%d]", eb.name, elNum)
			return &envBinder{
				source: eb.source,
				name:   name,
				lookup: func() (string, error) {
					if elNum >= len(arr) {
						return "", fmt.Errorf("tried to access invalid element %s, array only has %d elements", name, elNum)
					}
					return arr[elNum], nil
				},
			}
		}
		return nil
	})
}

func parseEnvKeyValue(kv string) (string, string) {
	if idx := strings.IndexRune(kv, '='); idx != -1 {
		return kv[:idx], kv[idx+1:]
	}
	return kv, ""
}

type envBinding struct {
	binder *envBinder
	apply  func() error
}

var _ interface {
	Binding
	BindingFromSource
} = &envBinding{}

func (eb *envBinding) Apply() error {
	return eb.apply()
}

func (eb *envBinding) Source() Source {
	return eb.binder.source
}

func (eb *envBinding) BoundName() string {
	return eb.binder.name
}
