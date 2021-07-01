package croconf

import (
	"encoding"
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

func (sev *SourceEnvVars) From(name string) LazySingleValueBinding {
	return &envBinding{
		source: sev,
		name:   name,
	}
}

type envBinding struct {
	source *SourceEnvVars
	name   string
}

func (eb *envBinding) GetSource() Source {
	return eb.source
}

func (eb *envBinding) BindStringValueTo(dest *string) func() error {
	return func() error {
		val, ok := eb.source.env[eb.name]
		if !ok {
			return ErrorMissing // TODO: better error message, e.g. 'field %s is not present in %s'?
		}
		*dest = val
		return nil
	}
}

func (eb *envBinding) BindIntValue() func(bitSize int) (int64, error) {
	return func(bitSize int) (int64, error) {
		val, ok := eb.source.env[eb.name]
		if !ok {
			return 0, ErrorMissing // TODO: better error message, e.g. 'field %s is not present in %s'?
		}
		return strconv.ParseInt(val, 10, bitSize) // TODO: use a custom function with better error messages
	}
}

func (eb *envBinding) BindTextBasedValueTo(dest encoding.TextUnmarshaler) func() error {
	return func() error {
		val, ok := eb.source.env[eb.name]
		if !ok {
			return ErrorMissing // TODO: better error message, e.g. 'field %s is not present in %s'?
		}

		return dest.UnmarshalText([]byte(val))
	}
}

func parseEnvKeyValue(kv string) (string, string) {
	if idx := strings.IndexRune(kv, '='); idx != -1 {
		return kv[:idx], kv[idx+1:]
	}
	return kv, ""
}
