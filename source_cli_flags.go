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

func (sc *SourceCLI) FromName(name string) *cliBinding {
	return &cliBinding{source: sc, longhand: name}
}

func (sc *SourceCLI) FromNameAndShorthand(name, shorthand string) *cliBinding {
	return &cliBinding{
		source:    sc,
		longhand:  name,
		shorthand: shorthand,
	}
}

func (sc *SourceCLI) FromPositionalArg(position int) LazySingleValueBinding {
	return &cliBinding{source: sc, position: position}
}

type cliBinding struct {
	source    *SourceCLI
	shorthand string
	longhand  string
	position  int

	// lookupfn defines a custom lookup logic
	lookupfn func() (string, error)
}

func (cb *cliBinding) GetSource() Source {
	return cb.source
}

// TODO: refactor
func (cb *cliBinding) lookup() (string, error) {
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
func (cb *cliBinding) textValueHelper(callback func(string) error) func() error {
	return func() error {
		vtext, err := cb.lookup()
		if err != nil {
			return ErrorMissing
		}
		return callback(vtext)
	}
}

func (cb *cliBinding) BindStringValueTo(dest *string) func() error {
	return cb.textValueHelper(func(s string) error {
		*dest = s
		return nil
	})
}

func (cb *cliBinding) BindTextBasedValueTo(dest encoding.TextUnmarshaler) func() error {
	return cb.textValueHelper(func(s string) error {
		return dest.UnmarshalText([]byte(s))
	})
}

func (cb *cliBinding) BindIntValueTo(dest *int64) func() error {
	return func() error {
		v, err := cb.lookup()
		if err != nil {
			// TODO what to use? cb.shorthand or cb.longhand or pos or smth joint
			return NewBindFieldMissingError(cb.source.GetName(), cb.longhand)
		}
		val, bindErr := parseInt(v)
		if bindErr != nil {
			return bindErr.withFuncName("BindIntValue")
		}
		*dest = val
		return nil
	}
}

func (cb *cliBinding) BindUintValueTo(dest *uint64) func() error {
	return func() error {
		v, err := cb.lookup()
		if err != nil {
			// TODO what to use? cb.shorthand or cb.longhand or pos or smth joint
			return NewBindFieldMissingError(cb.source.GetName(), cb.longhand)
		}
		val, bindErr := parseUint(v)
		if bindErr != nil {
			return bindErr.withFuncName("BindIntValue")
		}
		*dest = val
		return nil
	}
}

func (cb *cliBinding) BindFloatValueTo(dest *float64) func() error {
	return func() error {
		v, err := cb.lookup()
		if err != nil {
			// TODO what to use? cb.shorthand or cb.longhand or pos or smth joint
			return NewBindFieldMissingError(cb.source.GetName(), cb.longhand)
		}
		val, bindErr := parseFloat(v)
		if bindErr != nil {
			return bindErr.withFuncName("BindIntValue")
		}
		*dest = val
		return nil
	}
}

func (cb *cliBinding) BindBoolValueTo(dest *bool) func() error {
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

func (cb *cliBinding) BindArray() func() (Array, error) {
	cb.source.parser.RegisterSlice(cb.longhand, cb.shorthand)
	return func() (Array, error) {
		opts := cb.source.fs.Options(cb.longhand, cb.shorthand)
		if len(opts) < 1 {
			return nil, ErrorMissing
		}
		return cliArrayBinding{cb: cb, arr: opts}, nil
	}
}

type cliArrayBinding struct {
	cb  *cliBinding
	arr []string
}

func (cab cliArrayBinding) Len() int {
	return len(cab.arr)
}

func (cab cliArrayBinding) Element(i int) LazySingleValueBinding {
	cbcopy := *cab.cb
	cbcopy.lookupfn = func() (string, error) {
		if i >= len(cab.arr) {
			return "", fmt.Errorf("tried to access invalid element %s, array only has %d elements", cab.cb.longhand, i)
		}
		return cab.arr[i], nil
	}
	return &cbcopy
}
