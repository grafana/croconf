package flag

import (
	"reflect"
	"testing"
)

func TestSet_Option(t *testing.T) {
	fs := Set{
		flags: map[string]string{
			"foo": "bar",
		},
	}
	s, ok := fs.Option("foo", "f")
	if s != "bar" || !ok {
		t.Errorf("got unexpected value from lookup: %v", s)
	}
}

func TestSet_Positional(t *testing.T) {
	fs := Set{
		posArgs: []string{"run", "foo"},
	}
	s, ok := fs.Positional(1)
	if s != "run" || !ok {
		t.Errorf("got unexpected value from lookup: %v", s)
	}

	s, ok = fs.Positional(3)
	if s != "" || ok {
		t.Errorf("got unexpected value from lookup: %v", s)
	}
}

func TestParse(t *testing.T) {
	t.Run("LongOption", func(t *testing.T) {
		//{in: "--test=false", exp: []string{}},
		//{in: "--test false", exp: []string{}},

		args := []string{"run", "--user=u1", "--", "hello.go"}

		p := NewParser()
		fs, err := p.Parse(args)
		if err != nil {
			t.Error(err)
		}
		u, ok := fs.flags["user"]
		if !ok {
			t.Errorf("not found expected key")
		}
		if u != "u1" {
			t.Errorf("unexpected user value %s", u)
		}
		if !reflect.DeepEqual(fs.posArgs, []string{"run", "hello.go"}) {
			t.Errorf("unexpected positional arguments")
		}
	})

	t.Run("LongOptionWithSpace", func(t *testing.T) {
		args := []string{"run", "--user", "u1", "--", "hello.go"}

		p := NewParser()
		fs, err := p.Parse(args)
		if err != nil {
			t.Error(err)
		}
		u, ok := fs.flags["user"]
		if !ok {
			t.Errorf("not found expected key")
		}
		if u != "u1" {
			t.Errorf("unexpected user value %s", u)
		}
		if !reflect.DeepEqual(fs.posArgs, []string{"run", "hello.go"}) {
			t.Errorf("unexpected positional arguments")
		}
	})

	t.Run("LongOptionWithUnary", func(t *testing.T) {
		args := []string{"--bool", "foo"}

		p := NewParser()
		p.RegisterUnary("bool", "")
		fs, err := p.Parse(args)
		if err != nil {
			t.Error(err)
		}
		u, ok := fs.flags["bool"]
		if !ok {
			t.Errorf("not found expected key")
		}
		if u != "true" {
			t.Errorf("unexpected user value %s", u)
		}
		if !reflect.DeepEqual(fs.posArgs, []string{"foo"}) {
			t.Errorf("unexpected positional arguments")
		}

		// TODO: test with duplicates
	})

	t.Run("LongOptionSlice", func(t *testing.T) {
		args := []string{"--nums", "0", "alpha", "--nums=42"}

		p := NewParser()
		p.RegisterSlice("nums", "n")
		fs, err := p.Parse(args)
		if err != nil {
			t.Error(err)
		}
		nums, ok := fs.slices["nums"]
		if !ok {
			t.Errorf("not found slice key")
		}
		if !reflect.DeepEqual(nums, []string{"0", "42"}) {
			t.Errorf("unexpected slice")
		}
		if !reflect.DeepEqual(fs.posArgs, []string{"alpha"}) {
			t.Errorf("unexpected positional arguments")
		}
	})

	t.Run("ShortOption", func(t *testing.T) {
		tests := []struct {
			in  []string
			exp map[string]string
		}{
			{
				in: []string{"-r", "cmd1"},
				exp: map[string]string{
					"r": "true",
				},
			},
			{
				in: []string{"-rbc"},
				exp: map[string]string{
					"r": "true",
					"b": "true",
					"c": "true",
				},
			},
			{
				in: []string{"-o", "file.json"},
				exp: map[string]string{
					"o": "file.json",
				},
			},
			{
				in: []string{"-o=file.json"},
				exp: map[string]string{
					"o": "file.json",
				},
			},
		}

		for _, tt := range tests {
			p := NewParser()
			p.RegisterUnary("roo", "r")
			p.RegisterUnary("coo", "c")
			p.RegisterUnary("boo", "b")
			fs, err := p.Parse(tt.in)
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(fs.flags, tt.exp) {
				t.Errorf("unexpected flag %v %+v", tt.in, fs.flags)
			}

		}
	})
}
