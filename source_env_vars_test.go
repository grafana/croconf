package croconf

import (
	"math"
	"testing"
)

func TestEnvVarsBindIntValue(t *testing.T) {
	t.Parallel()
	environ := []string{"K6_VUS=6", "PI=3.14", "K6_CONFIG=./config.json", "K6_USER_AGENT=foo"}

	source := NewSourceFromEnv(environ)
	vus := source.From("K6_VUS")
	k6UserAgent := source.From("K6_USER_AGENT")
	missed := source.From("MISSED")

	if err := source.Initialize(); err != nil {
		t.Fatalf("received an unexpected init error %s", err)
	}

	var val int64
	err := vus.BindIntValueTo(&val)()
	if err != nil {
		t.Errorf("BindIntValueTo error: %s", err)
	}
	if expected := int64(6); val != expected {
		t.Errorf("BindIntValue: got %d, expected %d", val, expected)
	}

	err = missed.BindIntValueTo(&val)()
	if err == nil {
		t.Error("BindIntValue: expected field missing error")
	}
	if err.Error() != "field MISSED is missing in config source environment variables" {
		t.Error("BindIntValue: unexpected error message:", err)
	}

	err = k6UserAgent.BindIntValueTo(&val)()
	if err == nil {
		t.Error("BindIntValue: expected syntax error")
	}
	if err.Error() != "BindIntValue: parsing \"foo\": invalid syntax" {
		t.Errorf("BindIntValue: unexpected error message")
	}
}

func TestEnvVarsBindUintValue(t *testing.T) {
	t.Parallel()
	environ := []string{"K6_VUS=6", "PI=3.14", "K6_CONFIG=./config.json", "K6_USER_AGENT=foo"}

	source := NewSourceFromEnv(environ)
	vus := source.From("K6_VUS")
	k6UserAgent := source.From("K6_USER_AGENT")
	missed := source.From("MISSED")

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
		t.Error("BindUintValue: expected field missing error")
	}
	if err.Error() != "field MISSED is missing in config source environment variables" {
		t.Error("BindUintValue: unexpected error message:", err)
	}

	err = k6UserAgent.BindUintValueTo(&val)()
	if err == nil {
		t.Error("BindUintValue: expected syntax error")
	}
	if err.Error() != "BindUintValue: parsing \"foo\": invalid syntax" {
		t.Errorf("BindUintValue: unexpected error message")
	}
}

func TestEnvVarsFloatValue(t *testing.T) {
	t.Parallel()
	environ := []string{"K6_VUS=6", "PI=3.14", "K6_CONFIG=./config.json", "K6_USER_AGENT=foo"}

	source := NewSourceFromEnv(environ)
	vus := source.From("K6_VUS")
	pi := source.From("PI")
	k6UserAgent := source.From("K6_USER_AGENT")
	missed := source.From("MISSED")

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
	if err.Error() != "field MISSED is missing in config source environment variables" {
		t.Error("BindFloatValue: unexpected error message:", err)
	}

	err = k6UserAgent.BindFloatValueTo(&val)()
	if err == nil {
		t.Error("BindFloatValue: expected syntax error")
	}
	if err.Error() != "BindFloatValue: parsing \"foo\": invalid syntax" {
		t.Errorf("BindFloatValue: unexpected error message")
	}
}
