package croconf

import (
	"encoding"

	"github.com/spf13/pflag"
)

type SourceCLI struct {
	flags   []string
	flagSet *pflag.FlagSet // TODO: replace pflag, it's a very poor fit for this architecture
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
	return sc.flagSet.Parse(sc.flags)
}

func (sc *SourceCLI) GetName() string {
	return "CLI flags" // TODO
}

func (sc *SourceCLI) FromName(name string) LazySingleValueBinding {
	return &cliBinding{source: sc, longhand: name}
}

func (sc *SourceCLI) FromNameAndShorthand(name, shorthand string) LazySingleValueBinding {
	return &cliBinding{source: sc, longhand: name, shorthand: shorthand}
}

func (sc *SourceCLI) FromPositionalArg(position int) LazySingleValueBinding {
	return &cliBinding{source: sc, position: position}
}

type cliBinding struct {
	source    *SourceCLI
	isUnary   bool // TODO: figure out what we should do about boolean CLI flags
	shorthand string
	longhand  string
	position  int
}

func (cb *cliBinding) GetSource() Source {
	return cb.source
}

// we can use this function to get the string representation of any simple
// binary value and parse it ourselves
func (cb *cliBinding) textValueHelper(callback func(string) error) func() error {
	if cb.position > 0 {
		return func() error {
			if cb.source.flagSet.NArg() < cb.position {
				return ErrorMissing
			}
			return callback(cb.source.flagSet.Arg(cb.position - 1))
		}
	}
	s := cb.source.flagSet.StringP(cb.longhand, cb.shorthand, "", "")
	return func() error {
		if f := cb.source.flagSet.Lookup(cb.longhand); f.Changed {
			return callback(*s)
		}
		return ErrorMissing
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

func (cb *cliBinding) BindIntValue() func(bitSize int) (int64, error) {
	if cb.position > 0 {
		return func(bitSize int) (int64, error) {
			if cb.source.flagSet.NArg() < cb.position {
				// TODO what to use? cb.shorthand or cb.longhand or smth joint
				return 0, NewBindFieldMissingError(cb.source.GetName(), cb.longhand)
			}
			// TODO: use a custom function with better error message
			val, bindErr := parseInt(cb.source.flagSet.Arg(cb.position-1), 10, bitSize)
			if bindErr != nil {
				return 0, bindErr.withFuncName("BindIntValue")
			}
			return val, nil
		}
	}
	s := cb.source.flagSet.StringP(cb.longhand, cb.shorthand, "", "")
	return func(bitSize int) (int64, error) {
		if f := cb.source.flagSet.Lookup(cb.longhand); f.Changed {
			// TODO: use a custom function with better error message
			val, bindErr := parseInt(*s, 10, bitSize)
			if bindErr != nil {
				return 0, bindErr.withFuncName("BindIntValue")
			}
			return val, nil
		}
		return 0, NewBindFieldMissingError(cb.source.GetName(), cb.longhand)
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
