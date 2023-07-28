package croconf

import (
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

func (eb *envBinder) ToString() TypedBinding[string] {
	return ToBinding(eb.name, eb.source, func() (string, error) {
		val, err := eb.lookup()
		if err != nil {
			return "", err
		}
		return val, nil
	})
}

func (eb *envBinder) ToInt64() TypedBinding[int64] {
	return ToBinding(eb.name, eb.source, func() (int64, error) {
		val, err := eb.lookup()
		if err != nil {
			return 0, NewBindFieldMissingError(eb.source.GetName(), eb.name)
		}
		parsedVal, bindErr := parseInt(val)
		if bindErr != nil {
			return 0, bindErr.withFuncName("ToInt64")
		}
		return parsedVal, nil
	})
}

func (eb *envBinder) ToUint64() TypedBinding[uint64] {
	return ToBinding(eb.name, eb.source, func() (uint64, error) {
		val, err := eb.lookup()
		if err != nil {
			return 0, NewBindFieldMissingError(eb.source.GetName(), eb.name)
		}
		parsedVal, bindErr := parseUint(val)
		if bindErr != nil {
			return 0, bindErr.withFuncName("ToUint64")
		}
		return parsedVal, nil
	})
}

func (eb *envBinder) ToFloat64() TypedBinding[float64] {
	return ToBinding(eb.name, eb.source, func() (float64, error) {
		val, err := eb.lookup()
		if err != nil {
			return 0, NewBindFieldMissingError(eb.source.GetName(), eb.name)
		}
		parsedVal, bindErr := parseFloat(val)
		if bindErr != nil {
			return 0, bindErr.withFuncName("ToFloat64")
		}
		return parsedVal, nil
	})
}

func (eb *envBinder) ToBool() TypedBinding[bool] {
	return ToBinding(eb.name, eb.source, func() (bool, error) {
		val, err := eb.lookup()
		if err != nil {
			return false, NewBindFieldMissingError(eb.source.GetName(), eb.name)
		}
		parsedVal, bindErr := strconv.ParseBool(val)
		if bindErr != nil {
			return false, err
		}
		return parsedVal, nil
	})
}

/*
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
*/

func parseEnvKeyValue(kv string) (string, string) {
	if idx := strings.IndexRune(kv, '='); idx != -1 {
		return kv[:idx], kv[idx+1:]
	}
	return kv, ""
}
