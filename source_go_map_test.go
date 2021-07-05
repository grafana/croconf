package croconf

import (
	"testing"
)

func TestGoMapBindIntValue(t *testing.T) {
	t.Parallel()
	gomap := map[string]interface{}{
		"k6_vus": 6, "pi": 3.14, "k6_config": "./config.json", "k6_user_agent": "foo"}

	source := NewGoMapSource(gomap)
	vus := source.From("k6_vus")
	k6UserAgent := source.From("k6_user_agent")
	missed := source.From("missed")

	if err := source.Initialize(); err != nil {
		t.Fatalf("received an unexpected init error %s", err)
	}

	var val int64
	valBinding := vus.BindIntValueTo(&val)
	err := valBinding.Apply()
	if err != nil {
		t.Errorf("BindIntValueTo error: %s", err)
	}
	if expected := int64(6); val != expected {
		t.Errorf("BindIntValueTo: got %d, expected %d", val, expected)
	}

	valBinding = missed.BindIntValueTo(&val)
	err = valBinding.Apply()
	if err == nil {
		t.Error("BindIntValueTo: expected field missing error")
	}
	if err.Error() != "field missed is missing in config source go map" {
		t.Error("BindIntValueTo: unexpected error message:", err)
	}

	valBinding = k6UserAgent.BindIntValueTo(&val)
	err = valBinding.Apply()
	if err == nil {
		t.Error("BindIntValueTo: expected syntax error")
	}
	if err.Error() != `BindIntValueTo: parsing "foo": casting any int* type failed` {
		t.Error("BindIntValueTo: unexpected error message:", err)
	}
}

func TestGoMapBindUintValue(t *testing.T) {
	t.Parallel()
	gomap := map[string]interface{}{
		"k6_vus": uint(6), "pi": 3.14, "k6_config": "./config.json", "k6_user_agent": "foo"}

	source := NewGoMapSource(gomap)
	vus := source.From("k6_vus")
	k6UserAgent := source.From("k6_user_agent")
	missed := source.From("missed")

	if err := source.Initialize(); err != nil {
		t.Fatalf("received an unexpected init error %s", err)
	}

	var val uint64
	valBinding := vus.BindUintValueTo(&val)
	err := valBinding.Apply()
	if err != nil {
		t.Errorf("BindUintValueTo error: %s", err)
	}
	if expected := uint64(6); val != expected {
		t.Errorf("BindUintValueTo: got %d, expected %d", val, expected)
	}

	valBinding = missed.BindUintValueTo(&val)
	err = valBinding.Apply()
	if err == nil {
		t.Error("BindUintValueTo: expected field missing error")
	}
	if err.Error() != "field missed is missing in config source go map" {
		t.Error("BindUintValueTo: unexpected error message:", err)
	}

	valBinding = k6UserAgent.BindUintValueTo(&val)
	err = valBinding.Apply()
	if err == nil {
		t.Error("BindUintValueTo: expected syntax error")
	}
	if err.Error() != `BindUintValueTo: parsing "foo": casting any uint* type failed` {
		t.Error("BindUintValueTo: unexpected error message:", err)
	}
}

func TestGoMapBindFloatValue(t *testing.T) {
	t.Parallel()
	gomap := map[string]interface{}{
		"k6_vus": 6, "pi": 3.14, "k6_config": "./config.json", "k6_user_agent": "foo"}

	source := NewGoMapSource(gomap)
	pi := source.From("pi")
	k6UserAgent := source.From("k6_user_agent")
	missed := source.From("missed")

	if err := source.Initialize(); err != nil {
		t.Fatalf("received an unexpected init error %s", err)
	}

	var val float64
	valBinding := pi.BindFloatValueTo(&val)
	err := valBinding.Apply()
	if err != nil {
		t.Errorf("BindFloatValueTo error: %s", err)
	}
	if expected := float64(3.14); val != expected {
		t.Errorf("BindFloatValueTo: got %f, expected %f", val, expected)
	}

	valBinding = missed.BindFloatValueTo(&val)
	err = valBinding.Apply()
	if err == nil {
		t.Error("BindFloatValueTo: expected field missing error")
	}
	if err.Error() != "field missed is missing in config source go map" {
		t.Error("BindFloatValueTo: unexpected error message:", err)
	}

	valBinding = k6UserAgent.BindFloatValueTo(&val)
	err = valBinding.Apply()
	if err == nil {
		t.Error("BindFloatValueTo: expected syntax error")
	}
	if err.Error() != `BindFloatValueTo: parsing "foo": casting any float* type failed` {
		t.Error("BindFloatValueTo: unexpected error message:", err)
	}
}

type person struct {
	name string
}

func (p *person) UnmarshalText(txt []byte) error {
	p.name = string(txt)
	return nil
}

func TestGoMapBindTextBasedValueTo(t *testing.T) {
	t.Parallel()

	gomap := map[string]interface{}{"name": "Alice"}

	source := NewGoMapSource(gomap)
	name := source.From("name")

	if err := source.Initialize(); err != nil {
		t.Fatalf("received an unexpected init error %s", err)
	}

	var val person
	valBinding := name.BindTextBasedValueTo(&val)
	err := valBinding.Apply()
	if err != nil {
		t.Errorf("BindTextBasedValueTo error: %s", err)
	}
	if expected := "Alice"; val.name != expected {
		t.Errorf("BindTextBasedValueTo: got %s, expected %s", val, expected)
	}

}
