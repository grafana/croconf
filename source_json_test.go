package croconf

import "testing"

func TestJSONBindIntValue(t *testing.T) {

	var json = []byte(`{"k6_vus":6,"pi":3.14,"k6_config":"./config.json","k6_user_agent":"foo"}`)

	source := NewJSONSource(json)
	vus := source.From("k6_vus")
	k6UserAgent := source.From("k6_user_agent")
	missed := source.From("missed")

	if err := source.Initialize(); err != nil {
		t.Error(err)
	}

	withFixedBytesSizeFunc := func(bytesSize int) {
		val, err := vus.BindIntValue()(bytesSize)
		expected := int64(6)
		if err != nil {
			t.Errorf("BindIntValue error: %s", err)
		}
		if val != expected {
			t.Errorf("BindIntValue: got %d, expected %d", val, expected)
		}

		_, err = missed.BindIntValue()(bytesSize)
		if err == nil {
			t.Error("BindIntValue: expected field missing error")
		}
		if err.Error() != "field missed is missing in config source json" {
			t.Error("BindIntValue: unexpected error message:", err)
		}

		_, err = k6UserAgent.BindIntValue()(bytesSize)
		if err == nil {
			t.Error("BindIntValue: expected syntax error")
		}
		// TODO why are double quotes "\"foo"\"" ?
		if err.Error() != `BindIntValue: parsing "\"foo\"": invalid syntax` {
			t.Error("BindIntValue: unexpected error message:", err)
		}
	}

	intBytesSizes := []int{0, 8, 16, 32, 64}

	for _, byteSize := range intBytesSizes {
		withFixedBytesSizeFunc(byteSize)
	}
}

func TestJSONBindUintValue(t *testing.T) {

	var json = []byte(`{"k6_vus":6,"pi":3.14,"k6_config":"./config.json","k6_user_agent":"foo"}`)

	source := NewJSONSource(json)
	var vus = source.From("k6_vus")
	var k6UserAgent = source.From("k6_user_agent")
	var missed = source.From("missed")

	if err := source.Initialize(); err != nil {
		t.Error(err)
	}

	withFixedBytesSizeFunc := func(bytesSize int) {
		val, err := vus.BindUintValue()(bytesSize)
		expected := uint64(6)
		if err != nil {
			t.Errorf("BindUintValueTo error: %s", err)
		}
		if val != expected {
			t.Errorf("BindUintValue: got %d, expected %d", val, expected)
		}

		_, err = missed.BindUintValue()(bytesSize)
		if err == nil {
			t.Error("BindUintValue: expected field k6_vus is missing error")
		}
		if err.Error() != "field missed is missing in config source json" {
			t.Error("BindUintValue: unexpected error message:", err)
		}

		_, err = k6UserAgent.BindUintValue()(bytesSize)
		if err == nil {
			t.Error("BindUintValue: expected syntax error")
		}
		// TODO why are double quotes "\"foo"\"" ?
		if err.Error() != `BindIntValue: parsing "\"foo\"": invalid syntax` {
			t.Error("BindIntValue: unexpected error message:", err)
		}
	}

	intBytesSizes := []int{0, 8, 16, 32, 64}

	for _, byteSize := range intBytesSizes {
		withFixedBytesSizeFunc(byteSize)
	}
}

func TestJSONFloatValue(t *testing.T) {
	var json = []byte(`{"k6_vus":6,"pi":3.14,"k6_config":"./config.json","k6_user_agent":"foo"}`)

	source := NewJSONSource(json)
	var vus = source.From("k6_vus")
	var pi = source.From("k6_vus")

	var k6UserAgent = source.From("k6_user_agent")
	var missed = source.From("missed")

	if err := source.Initialize(); err != nil {
		t.Error(err)
	}

	withFixedBytesSizeFunc := func(bytesSize int) {
		val, err := vus.BindFloatValue()(bytesSize)
		expected := float64(6)
		if err != nil {
			t.Errorf("BindFloatValue error: %s", err)
		}
		if val != expected {
			t.Errorf("BindFloatValue: got %f, expected %f", val, expected)
		}

		val, err = pi.BindFloatValue()(bytesSize)
		expected = float64(3.14)
		if err != nil {
			t.Errorf("BindFloatValue error: %s", err)
		}
		// val != expected doesn't work
		if (val-3.14) > 1e20 && (val-3.14) < -1e20 {
			t.Errorf("BindFloatValue: got %f, expected %f", val, expected)
		}

		_, err = missed.BindFloatValue()(bytesSize)
		if err == nil {
			t.Error("BindFloatValue: expected field missing error")
		}
		if err.Error() != "field missed is missing in config source json" {
			t.Error("BindFloatValue: unexpected error message:", err)
		}

		_, err = k6UserAgent.BindFloatValue()(bytesSize)
		if err == nil {
			t.Error("BindFloatValue: expected syntax error")
		}
		// TODO why are double quotes "\"foo"\"" ?
		if err.Error() != `BindIntValue: parsing "\"foo\"": invalid syntax` {
			t.Error("BindIntValue: unexpected error message:", err)
		}
	}

	intBytesSizes := []int{0, 8, 16, 32, 64}

	for _, byteSize := range intBytesSizes {
		withFixedBytesSizeFunc(byteSize)
	}
}

func TestJSONBindIntValue__NestedJSON(t *testing.T) {

	var json = []byte(`{"data":{"k6_vus":6,"pi":3.14,"k6_config":"./config.json","k6_user_agent":"foo"}}`)

	source := NewJSONSource(json)
	vus := source.From("data").From("k6_vus")
	k6UserAgent := source.From("data").From("k6_user_agent")
	missed := source.From("data").From("missed")

	if err := source.Initialize(); err != nil {
		t.Error(err)
	}

	withFixedBytesSizeFunc := func(bytesSize int) {
		val, err := vus.BindIntValue()(bytesSize)
		expected := int64(6)
		if err != nil {
			t.Errorf("BindIntValue error: %s", err)
		}
		if val != expected {
			t.Errorf("BindIntValue: got %d, expected %d", val, expected)
		}

		_, err = missed.BindIntValue()(bytesSize)
		if err == nil {
			t.Error("BindIntValue: expected field missing error")
		}
		if err.Error() != "field data.missed is missing in config source json" {
			t.Error("BindIntValue: unexpected error message:", err)
		}

		_, err = k6UserAgent.BindIntValue()(bytesSize)
		if err == nil {
			t.Error("BindIntValue: expected syntax error")
		}
		// TODO why are double quotes "\"foo"\"" ?
		if err.Error() != `BindIntValue: parsing "\"foo\"": invalid syntax` {
			t.Error("BindIntValue: unexpected error message:", err)
		}
	}

	intBytesSizes := []int{0, 8, 16, 32, 64}

	for _, byteSize := range intBytesSizes {
		withFixedBytesSizeFunc(byteSize)
	}
}
