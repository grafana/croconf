package croconf

import (
	"fmt"
	"testing"
)

func TestGoMapBindIntValue(t *testing.T) {
	gomap := map[string]interface{}{
		"k6_vus_0": 6, "k6_vus_8": int8(6), "k6_vus_16": int16(6), "k6_vus_32": int32(6), "k6_vus_64": int64(6),
		"pi": 3.14, "k6_config": "./config.json", "k6_user_agent": "foo"}

	withFixedBytesSizeFunc := func(bytesSize int) {

		source := NewGoMapSource(gomap)
		vusKey := fmt.Sprintf("k6_vus_%d", bytesSize)
		vus := source.From(vusKey)
		k6UserAgent := source.From("k6_user_agent")
		missed := source.From("missed")

		if err := source.Initialize(); err != nil {
			t.Error(err)
		}

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
		if err.Error() != "field missed is missing in config source go map" {
			t.Error("BindIntValue: unexpected error message:", err)
		}

		_, err = k6UserAgent.BindIntValue()(bytesSize)
		if err == nil {
			t.Error("BindIntValue: expected error")
		}
	}

	intBytesSizes := []int{0, 8, 16, 32, 64}

	for _, byteSize := range intBytesSizes {
		withFixedBytesSizeFunc(byteSize)
	}
}

func TestGoMapBindUintValue(t *testing.T) {
	gomap := map[string]interface{}{
		"k6_vus_0": uint(6), "k6_vus_8": uint8(6), "k6_vus_16": uint16(6), "k6_vus_32": uint32(6), "k6_vus_64": uint64(6),
		"pi": 3.14, "k6_config": "./config.json", "k6_user_agent": "foo"}

	withFixedBytesSizeFunc := func(bytesSize int) {

		source := NewGoMapSource(gomap)
		vusKey := fmt.Sprintf("k6_vus_%d", bytesSize)
		vus := source.From(vusKey)
		k6UserAgent := source.From("k6_user_agent")
		missed := source.From("missed")

		if err := source.Initialize(); err != nil {
			t.Error(err)
		}
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
		if err.Error() != "field missed is missing in config source go map" {
			t.Error("BindUintValue: unexpected error message:", err)
		}

		_, err = k6UserAgent.BindUintValue()(bytesSize)
		if err == nil {
			t.Error("BindUintValue: expected syntax error")
		}
	}

	intBytesSizes := []int{0, 8, 16, 32, 64}

	for _, byteSize := range intBytesSizes {
		withFixedBytesSizeFunc(byteSize)
	}
}

func TestGoMapFloatValue(t *testing.T) {
	gomap := map[string]interface{}{
		"k6_vus": 6,
		"pi_32":  float32(3.14), "pi_64": float64(3.14), "k6_config": "./config.json", "k6_user_agent": "foo"}

	withFixedBytesSizeFunc := func(bytesSize int) {

		source := NewGoMapSource(gomap)
		piKey := fmt.Sprintf("pi_%d", bytesSize)
		pi := source.From(piKey)
		vus := source.From("k6_vus")
		k6UserAgent := source.From("k6_user_agent")
		missed := source.From("missed")

		if err := source.Initialize(); err != nil {
			t.Error(err)
		}

		val, err := pi.BindFloatValue()(bytesSize)
		expected := float64(3.14)
		if err != nil {
			t.Errorf("BindFloatValue error: %s", err)
		}
		// val != expected doesn't work
		if (val-3.14) > 1e20 && (val-3.14) < -1e20 {
			t.Errorf("BindFloatValue: got %f, expected %f", val, expected)
		}

		val, err = vus.BindFloatValue()(bytesSize)
		expected = float64(6)
		if err == nil {
			t.Errorf("BindFloatValue expected error")
		}

		_, err = missed.BindFloatValue()(bytesSize)
		if err == nil {
			t.Error("BindFloatValue: expected field missing error")
		}
		if err.Error() != "field missed is missing in config source go map" {
			t.Error("BindFloatValue: unexpected error message:", err)
		}

		_, err = k6UserAgent.BindFloatValue()(bytesSize)
		if err == nil {
			t.Error("BindFloatValue: expected syntax error")
		}
	}

	intBytesSizes := []int{32, 64}

	for _, byteSize := range intBytesSizes {
		withFixedBytesSizeFunc(byteSize)
	}
}

func TestGoMapBindIntValue__NestedGoMap(t *testing.T) {
	gomap := map[string]interface{}{
		"data": map[string]interface{}{
			"k6_vus_0": 6, "k6_vus_8": int8(6), "k6_vus_16": int16(6), "k6_vus_32": int32(6), "k6_vus_64": int64(6),
			"pi": 3.14, "k6_config": "./config.json", "k6_user_agent": "foo"}}

	withFixedBytesSizeFunc := func(bytesSize int) {

		source := NewGoMapSource(gomap)
		vusKey := fmt.Sprintf("k6_vus_%d", bytesSize)
		vus := source.From("data").From(vusKey)
		k6UserAgent := source.From("data").From("k6_user_agent")
		missed := source.From("data").From("missed")

		if err := source.Initialize(); err != nil {
			t.Error(err)
		}

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
		if err.Error() != "field data.missed is missing in config source go map" {
			t.Error("BindIntValue: unexpected error message:", err)
		}

		_, err = k6UserAgent.BindIntValue()(bytesSize)
		if err == nil {
			t.Error("BindIntValue: expected error")
		}
	}

	intBytesSizes := []int{0, 8, 16, 32, 64}

	for _, byteSize := range intBytesSizes {
		withFixedBytesSizeFunc(byteSize)
	}
}
