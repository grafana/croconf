package croconf

import "testing"

func TestCLIBinding_BindStringValueTo(t *testing.T) {
	tests := map[string]struct {
		args []string
	}{
		//"shorthand": {
		//args: []string{"-b"},
		//},
		"longname": {
			args: []string{"--bool", "true"},
		},
		//"shorthand-with-value": {
		//args: []string{"-b", "true"},
		//},
		"single-arg": {
			args: []string{"--bool=true"},
		},
		"with-cmd": {
			args: []string{"cmd1", "--bool", "true"},
		},
	}
	for name, tt := range tests {
		var b bool

		src := NewSourceFromCLIFlags(tt.args)
		field := src.FromNameAndShorthand("bool", "b")

		err := src.Initialize()
		if err != nil {
			t.Errorf("test: %s: got unexpected error: %v", name, err)
			return
		}

		bindfn := field.BindBoolValueTo(&b)
		err = bindfn()
		if err != nil {
			t.Errorf("test: %s: got unexpected error: %v", name, err)
			return
		}

		if !b {
			t.Errorf("expected a true boolean value")
		}
	}
}
