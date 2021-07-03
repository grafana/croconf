package flag

import (
	"strings"
)

type Parser struct {
	unaries map[string]struct{}
	slices  map[string]struct{}
}

func NewParser() *Parser {
	return &Parser{
		unaries: make(map[string]struct{}),
		slices:  make(map[string]struct{}),
	}
}

func (p *Parser) RegisterUnary(long, short string) {
	p.unaries[long] = struct{}{}
	if short != "" {
		p.unaries[short] = struct{}{}
	}
}

func (p *Parser) RegisterSlice(long, short string) {
	p.slices[long] = struct{}{}
	if short != "" {
		p.slices[short] = struct{}{}
	}
}

func (p *Parser) Parse(tt []string) (*Set, error) {
	// TODO: refactor to remove this
	args := make([]string, len(tt))
	copy(args, tt)

	fs := Set{
		flags:   make(map[string]string),
		slices:  make(map[string][]string),
		posArgs: make([]string, 0, len(args)),
	}

	// remove the single or double dash
	nohypens := func(s string) string {
		if s[0] != '-' {
			return s
		}
		if s[1] != '-' {
			// single hypen
			return s[1:]
		}
		// double hypen
		return s[2:]
	}

	addflag := func(key, v string) {
		if _, ok := p.slices[key]; !ok {
			fs.flags[key] = v
		} else {
			fs.slices[key] = append(fs.slices[key], v)
		}
	}

	var err error
	for i := 0; i < len(args); i++ {
		arg := args[i]

		var next *string
		if len(args) > i+1 {
			next = &args[i+1]
		}

		switch {
		case strings.HasPrefix(arg, "--"):
			// -- example.go
			if arg == "--" {
				continue
			}

			arg = nohypens(arg)

			// --opt=value
			opt := strings.Split(arg, "=")
			if len(opt) == 2 {
				addflag(opt[0], opt[1])
				continue
			}

			// --bool cmd1
			if _, ok := p.unaries[arg]; ok {
				fs.flags[arg] = "true"
				continue
			}

			// --opt value
			if next != nil {
				addflag(arg, *next)
				args = args[:i+copy(args[i:], args[i+1:])]
			}

		case strings.HasPrefix(arg, "-"):
			arg = nohypens(arg)

			if len(arg) > 1 {
				// -o=value
				opt := strings.Split(arg, "=")
				if len(opt) == 2 {
					addflag(opt[0], opt[1])
					continue
				}

				// -ob
				for i := 0; i < len(arg); i++ {
					if _, ok := p.unaries[string(arg[i])]; ok {
						fs.flags[string(arg[i])] = "true"
						continue
					}
				}
			}

			// -b cmd1
			if _, ok := p.unaries[arg]; ok {
				fs.flags[arg] = "true"
				continue
			}

			// -o value
			if next != nil {
				addflag(arg, *next)
				args = args[:i+copy(args[i:], args[i+1:])]
			}
		default:
			fs.posArgs = append(fs.posArgs, arg)
		}
	}
	return &fs, err
}

type Set struct {
	slices  map[string][]string
	flags   map[string]string
	posArgs []string
}

func (fs Set) Positional(i uint) (string, bool) {
	index := i - 1
	if int(index) > len(fs.posArgs)-1 {
		return "", false
	}
	return fs.posArgs[index], true
}

func (fs Set) Option(long, short string) (string, bool) {
	opt, ok := fs.flags[long]
	if ok {
		return opt, true
	}
	opt, ok = fs.flags[short]
	return opt, ok
}

func (fs Set) Options(long, short string) (opts []string) {
	longs, ok := fs.slices[long]
	if ok {
		opts = append(opts, longs...)
	}
	shorts, ok := fs.slices[short]
	if ok {
		opts = append(opts, shorts...)
	}
	return
}
