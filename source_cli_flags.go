package croconf

import (
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

func (cb *cliBinder) ToString() TypedBinding[string] {
	return ToBinding(cb.boundName(), cb.source, func() (string, error) {
		val, err := cb.lookup()
		if err != nil {
			return "", err
		}
		return val, nil
	})
}

func (cb *cliBinder) ToInt64() TypedBinding[int64] {
	return ToBinding(cb.boundName(), cb.source, func() (int64, error) {
		val, err := cb.lookup()
		if err != nil {
			return 0, err
		}
		parsedVal, bindErr := parseInt(val)
		if bindErr != nil {
			return 0, bindErr.withFuncName("ToInt64")
		}
		return parsedVal, nil
	})
}

func (cb *cliBinder) ToUint64() TypedBinding[uint64] {
	return ToBinding(cb.boundName(), cb.source, func() (uint64, error) {
		val, err := cb.lookup()
		if err != nil {
			return 0, err
		}
		parsedVal, bindErr := parseUint(val)
		if bindErr != nil {
			return 0, bindErr.withFuncName("ToUint64")
		}
		return parsedVal, nil
	})
}

func (cb *cliBinder) ToFloat64() TypedBinding[float64] {
	return ToBinding(cb.boundName(), cb.source, func() (float64, error) {
		val, err := cb.lookup()
		if err != nil {
			return 0, err
		}
		parsedVal, bindErr := parseFloat(val)
		if bindErr != nil {
			return 0, bindErr.withFuncName("ToFloat64")
		}
		return parsedVal, nil
	})
}

func (cb *cliBinder) ToBool() TypedBinding[bool] {
	return ToBinding(cb.boundName(), cb.source, func() (bool, error) {
		val, err := cb.lookup()
		if err != nil {
			return false, err
		}
		parsedVal, bindErr := strconv.ParseBool(val)
		if bindErr != nil {
			return false, err
		}
		return parsedVal, nil
	})
}

/*
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
*/
