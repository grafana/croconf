package croconf

import (
	"encoding"
	"fmt"
	"strconv"

	"go.k6.io/croconf/flag"
)

type SourceCLI struct {
	flags []string

	parser *flag.Parser
	fs     *flag.Set
}

func NewSourceFromCLIFlags(flags []string) *SourceCLI {
	return &SourceCLI{
		flags:  flags,
		parser: flag.NewParser(),
	}
}

func (sc *SourceCLI) Initialize() error {
	fs, err := sc.parser.Parse(sc.flags)
	if err != nil {
		return err
	}
	sc.fs = fs
	return nil
}

func (sc *SourceCLI) GetName() string {
	return "CLI flags" // TODO
}

func (sc *SourceCLI) FromName(name string) *cliBinder {
	return &cliBinder{source: sc, longhand: name}
}

func (sc *SourceCLI) FromNameAndShorthand(name, shorthand string) *cliBinder {
	return &cliBinder{
		source:    sc,
		longhand:  name,
		shorthand: shorthand,
	}
}

func (sc *SourceCLI) FromPositionalArg(position int) *cliBinder {
	return &cliBinder{source: sc, position: position}
}

type cliBinder struct {
	source    *SourceCLI
	shorthand string
	longhand  string
	position  int

	// lookupfn defines a custom lookup logic
	lookupfn func() (string, error)
}

func (cb *cliBinder) boundName() string {
	// TODO: improve?
	if cb.position > 0 {
		return fmt.Sprintf("argument #%d", cb.position)
	}
	if cb.shorthand != "" {
		return fmt.Sprintf("--%s / -%s", cb.shorthand, cb.longhand)
	}
	return fmt.Sprintf("--%s", cb.longhand)
}

func (cb *cliBinder) newBinding(apply func() error) *cliBinding {
	return &cliBinding{
		binder: cb,
		apply:  apply,
	}
}

// TODO: refactor
func (cb *cliBinder) lookup() (string, error) {
	// custom lookup
	if cb.lookupfn != nil {
		return cb.lookupfn()
	}
	// default
	if cb.position > 0 {
		arg, ok := cb.source.fs.Positional(uint(cb.position))
		if !ok {
			return "", ErrorMissing
		}
		return arg, nil
	}
	opt, ok := cb.source.fs.Option(cb.longhand, cb.shorthand)
	if !ok {
		return opt, ErrorMissing
	}
	return opt, nil
}

// we can use this function to get the string representation of any simple
// binary value and parse it ourselves
func (cb *cliBinder) textValueHelper(callback func(string) error) Binding {
	return cb.newBinding(func() error {
		vtext, err := cb.lookup()
		if err != nil {
			return ErrorMissing
		}
		return callback(vtext)
	})
}

func (cb *cliBinder) BindStringValueTo(dest *string) Binding {
	return cb.textValueHelper(func(s string) error {
		*dest = s
		return nil
	})
}

func (cb *cliBinder) BindTextBasedValueTo(dest encoding.TextUnmarshaler) Binding {
	return cb.textValueHelper(func(s string) error {
		return dest.UnmarshalText([]byte(s))
	})
}

func (cb *cliBinder) BindIntValueTo(dest *int64) Binding {
	return cb.newBinding(func() error {
		v, err := cb.lookup()
		if err != nil {
			return NewBindFieldMissingError(cb.source.GetName(), cb.boundName())
		}
		val, bindErr := parseInt(v)
		if bindErr != nil {
			return bindErr.withFuncName("BindIntValue")
		}
		*dest = val
		return nil
	})
}

func (cb *cliBinder) BindUintValueTo(dest *uint64) Binding {
	return cb.newBinding(func() error {
		v, err := cb.lookup()
		if err != nil {
			return NewBindFieldMissingError(cb.source.GetName(), cb.boundName())
		}
		val, bindErr := parseUint(v)
		if bindErr != nil {
			return bindErr.withFuncName("BindIntValue")
		}
		*dest = val
		return nil
	})
}

func (cb *cliBinder) BindFloatValueTo(dest *float64) Binding {
	return cb.newBinding(func() error {
		v, err := cb.lookup()
		if err != nil {
			return NewBindFieldMissingError(cb.source.GetName(), cb.boundName())
		}
		val, bindErr := parseFloat(v)
		if bindErr != nil {
			return bindErr.withFuncName("BindIntValue")
		}
		*dest = val
		return nil
	})
}

func (cb *cliBinder) BindBoolValueTo(dest *bool) Binding {
	cb.source.parser.RegisterUnary(cb.longhand, cb.shorthand)
	return cb.textValueHelper(func(v string) error {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return err
		}
		*dest = b
		return nil
	})
}

func (cb *cliBinder) BindArrayValueTo(length *int, element *func(int) LazySingleValueBinder) Binding {
	cb.source.parser.RegisterSlice(cb.longhand, cb.shorthand)
	return cb.newBinding(func() error {
		opts := cb.source.fs.Options(cb.longhand, cb.shorthand)
		if len(opts) < 1 {
			return ErrorMissing
		}

		*length = len(opts)
		*element = func(i int) LazySingleValueBinder {
			cbcopy := *cb
			cbcopy.lookupfn = func() (string, error) {
				if i >= len(opts) {
					return "", fmt.Errorf("tried to access invalid element %s, array only has %d elements", cbcopy.longhand, i)
				}
				return opts[i], nil
			}
			return &cbcopy
		}

		return nil
	})
}

type cliBinding struct {
	binder *cliBinder
	apply  func() error
}

var _ interface {
	Binding
	BindingFromSource
} = &cliBinding{}

func (cb *cliBinding) Apply() error {
	return cb.apply()
}

func (cb *cliBinding) Source() Source {
	return cb.binder.source
}

func (cb *cliBinding) BoundName() string {
	return cb.binder.boundName()
}
