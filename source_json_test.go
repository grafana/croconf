package croconf

import (
	"math"
	"testing"
)

func TestJSONBindIntValue(t *testing.T) {
	t.Parallel()
	json := []byte(`{"k6_vus":6,"pi":3.14,"k6_config":"./config.json","k6_user_agent":"foo"}`)

	source := NewJSONSource(json)
	vus := source.From("k6_vus")
	k6UserAgent := source.From("k6_user_agent")
	missed := source.From("missed")

	if err := source.Initialize(); err != nil {
		t.Fatalf("received an unexpected init error %s", err)
	}

	var val int64
	err := vus.BindIntValueTo(&val)()
	if err != nil {
		t.Errorf("BindIntValue error: %s", err)
	}
	if expected := int64(6); val != expected {
		t.Errorf("BindIntValue: got %d, expected %d", val, expected)
	}

	err = missed.BindIntValueTo(&val)()
	if err == nil {
		t.Error("BindIntValue: expected field missing error")
	}
	if err.Error() != "field missed is missing in config source json" {
		t.Error("BindIntValue: unexpected error message:", err)
	}

	err = k6UserAgent.BindIntValueTo(&val)()
	if err == nil {
		t.Error("BindIntValue: expected syntax error")
	}
	// TODO why are double quotes "\"foo"\"" ?
	if err.Error() != `BindIntValue: parsing "\"foo\"": invalid syntax` {
		t.Error("BindIntValue: unexpected error message:", err)
	}
}

func TestJSONBindUintValue(t *testing.T) {
	t.Parallel()
	json := []byte(`{"k6_vus":6,"pi":3.14,"k6_config":"./config.json","k6_user_agent":"foo"}`)

	source := NewJSONSource(json)
	vus := source.From("k6_vus")
	k6UserAgent := source.From("k6_user_agent")
	missed := source.From("missed")

	if err := source.Initialize(); err != nil {
		t.Fatalf("received an unexpected init error %s", err)
	}

	var val uint64
	err := vus.BindUintValueTo(&val)()
	if err != nil {
		t.Errorf("BindUintValueTo error: %s", err)
	}
	if expected := uint64(6); val != expected {
		t.Errorf("BindUintValue: got %d, expected %d", val, expected)
	}

	err = missed.BindUintValueTo(&val)()
	if err == nil {
		t.Error("BindUintValue: expected field k6_vus is missing error")
	}
	if err.Error() != "field missed is missing in config source json" {
		t.Error("BindUintValue: unexpected error message:", err)
	}

	err = k6UserAgent.BindUintValueTo(&val)()
	if err == nil {
		t.Error("BindUintValue: expected syntax error")
	}
	// TODO why are double quotes "\"foo"\"" ?
	if err.Error() != `BindIntValue: parsing "\"foo\"": invalid syntax` {
		t.Error("BindIntValue: unexpected error message:", err)
	}
}

func TestJSONFloatValue(t *testing.T) {
	t.Parallel()
	json := []byte(`{"k6_vus":6,"pi":3.14,"k6_config":"./config.json","k6_user_agent":"foo"}`)

	source := NewJSONSource(json)
	vus := source.From("k6_vus")
	pi := source.From("k6_vus")

	k6UserAgent := source.From("k6_user_agent")
	missed := source.From("missed")

	if err := source.Initialize(); err != nil {
		t.Fatalf("received an unexpected init error %s", err)
	}

	var val float64
	err := vus.BindFloatValueTo(&val)()
	expected := float64(6)
	if err != nil {
		t.Errorf("BindFloatValue error: %s", err)
	}
	if val != expected {
		t.Errorf("BindFloatValue: got %f, expected %f", val, expected)
	}

	err = pi.BindFloatValueTo(&val)()
	expected = float64(3.14)
	if err != nil {
		t.Errorf("BindFloatValue error: %s", err)
	}
	if math.Abs(val-expected) > 1e20 { // val != expected doesn't work for floats
		t.Errorf("BindFloatValue: got %f, expected %f", val, expected)
	}

	err = missed.BindFloatValueTo(&val)()
	if err == nil {
		t.Error("BindFloatValue: expected field missing error")
	}
	if err.Error() != "field missed is missing in config source json" {
		t.Error("BindFloatValue: unexpected error message:", err)
	}

	err = k6UserAgent.BindFloatValueTo(&val)()
	if err == nil {
		t.Error("BindFloatValue: expected syntax error")
	}
	// TODO why are double quotes "\"foo"\"" ?
	if err.Error() != `BindIntValue: parsing "\"foo\"": invalid syntax` {
		t.Error("BindIntValue: unexpected error message:", err)
	}
}

func TestJSONBindIntValue__NestedJSON(t *testing.T) {
	t.Parallel()
	json := []byte(`{"data":{"k6_vus":6,"pi":3.14,"k6_config":"./config.json","k6_user_agent":"foo"}}`)

	source := NewJSONSource(json)
	vus := source.From("data").From("k6_vus")
	k6UserAgent := source.From("data").From("k6_user_agent")
	missed := source.From("data").From("missed")

	if err := source.Initialize(); err != nil {
		t.Fatalf("received an unexpected init error %s", err)
	}

	var val int64
	err := vus.BindIntValueTo(&val)()
	if err != nil {
		t.Errorf("BindIntValue error: %s", err)
	}
	if expected := int64(6); val != expected {
		t.Errorf("BindIntValue: got %d, expected %d", val, expected)
	}

	err = missed.BindIntValueTo(&val)()
	if err == nil {
		t.Error("BindIntValue: expected field missing error")
	}
	if err.Error() != "field data.missed is missing in config source json" {
		t.Error("BindIntValue: unexpected error message:", err)
	}

	err = k6UserAgent.BindIntValueTo(&val)()
	if err == nil {
		t.Error("BindIntValue: expected syntax error")
	}
	// TODO why are double quotes "\"foo"\"" ?
	if err.Error() != `BindIntValue: parsing "\"foo\"": invalid syntax` {
		t.Error("BindIntValue: unexpected error message:", err)
	}
}
