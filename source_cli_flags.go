package croconf

import (
	"encoding"
	"fmt"
	"strconv"

	"github.com/spf13/pflag"
	"go.k6.io/croconf/flag"
)

type SourceCLI struct {
	flags   []string
	flagSet *pflag.FlagSet // TODO: replace pflag, it's a very poor fit for this architecture

	fs     *flag.Set
	parser *flag.Parser
}

func NewSourceFromCLIFlags(flags []string) *SourceCLI {
	flagSet := pflag.NewFlagSet("this is only temporary", pflag.ContinueOnError)
	flagSet.SortFlags = false
	flagSet.ParseErrorsWhitelist = pflag.ParseErrorsWhitelist{UnknownFlags: true}

	return &SourceCLI{
		flags:   flags,
		flagSet: flagSet,
	}
}

func (sc *SourceCLI) Initialize() error {
	//err := sc.flagSet.Parse(sc.flags)
	sc.parser = flag.NewParser()
	fs, err := sc.parser.Parse(sc.flags)
	if err != nil {
		return err
	}
	sc.fs = fs
	return err
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
	unary     bool // TODO: figure out what we should do about boolean CLI flags
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
func (cb *cliBinding) BindIntValue() func(int) (int64, error) {
	return func(bitSize int) (int64, error) {
		v, err := cb.lookup()
		if err != nil {
			// TODO what to use? cb.shorthand or cb.longhand or smth joint
			return 0, NewBindFieldMissingError(cb.source.GetName(), cb.longhand)
		}
		val, bindErr := parseInt(v, 10, bitSize)
		if bindErr != nil {
			return 0, bindErr.withFuncName("BindIntValue")
		}
		return val, nil
	}
}

func (cb *cliBinding) BindUintValue() func(bitSize int) (uint64, error) {
	if cb.position > 0 {
		return func(bitSize int) (uint64, error) {
			if cb.source.flagSet.NArg() < cb.position {
				return 0, ErrorMissing
			}
			// TODO: use a custom function with better error message
			return parseUint(cb.source.flagSet.Arg(cb.position-1), 10, bitSize)
		}
	}
	s := cb.source.flagSet.StringP(cb.longhand, cb.shorthand, "", "")
	return func(bitSize int) (uint64, error) {
		if f := cb.source.flagSet.Lookup(cb.longhand); f.Changed {
			// TODO: use a custom function with better error message
			return parseUint(*s, 10, bitSize)
		}
		return 0, ErrorMissing
	}
}

func (cb *cliBinding) BindFloatValue() func(bitSize int) (float64, error) {
	if cb.position > 0 {
		return func(bitSize int) (float64, error) {
			if cb.source.flagSet.NArg() < cb.position {
				return 0, ErrorMissing
			}
			// TODO: use a custom function with better error message
			return parseFloat(cb.source.flagSet.Arg(cb.position-1), bitSize)
		}
	}
	s := cb.source.flagSet.StringP(cb.longhand, cb.shorthand, "", "")
	return func(bitSize int) (float64, error) {
		if f := cb.source.flagSet.Lookup(cb.longhand); f.Changed {
			// TODO: use a custom function with better error message
			return parseFloat(*s, bitSize)
		}
		return 0, ErrorMissing
	}
}

func (cb *cliBinding) BindBoolValueTo(dest *bool) func() error {
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
	err := cb.source.parser.RegisterSlice(cb.longhand, cb.shorthand)
	return func() (Array, error) {
		if err != nil {
			return nil, fmt.Errorf("slice binding failed")
		}

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
