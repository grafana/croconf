package croconf

import "testing"

func TestJSONBindIntValue(t *testing.T) {

	var json = []byte(`{"k6_vus":6,"pi":3.14,"k6_config":"./config.json","k6_user_agent":"foo"}`)

	var vus = NewJSONSource(json).From("k6_vus")
	// var k6UserAgent = NewJSONSource(json).From("k6_user_agent")
	// var missed = NewJSONSource(json).From("missed")

	withFixedBytesSizeFunc := func(bytesSize int) {
		val, err := vus.BindIntValue()(bytesSize)
		expected := int64(6)
		if err != nil {
			t.Errorf("BindIntValue error: %s", err)
		}
		if val != expected {
			t.Errorf("BindIntValue: got %d, expected %d", val, expected)
		}

		// _, err = missed.BindIntValue()(bytesSize)
		// if err == nil {
		// 	t.Error("BindIntValue: expected field missing error")
		// }
		// if err.Error() != "field missed is missing in config source jsonment variables" {
		// 	t.Error("BindIntValue: unexpected error message:", err)
		// }

		// _, err = k6UserAgent.BindIntValue()(bytesSize)
		// if err == nil {
		// 	t.Error("BindIntValue: expected syntax error")
		// }
		// if err.Error() != "BindIntValue: parsing \"foo\": invalid syntax" {
		// 	t.Errorf("BindIntValue: unexpected error message")
		// }
	}

	intBytesSizes := []int{0, 8, 16, 32, 64}

	for _, byteSize := range intBytesSizes {
		withFixedBytesSizeFunc(byteSize)
	}
}

// func TestJSONBindUintValue(t *testing.T) {

// 	var json = []byte(`{"k6_vus":6,"pi":3.14,"k6_config":"./config.json","k6_user_agent":"foo"}`)

// 	var vus = NewJSONSource(json).From("k6_vus")
// 	var k6UserAgent = NewJSONSource(json).From("k6_user_agent")
// 	var missed = NewJSONSource(json).From("missed")

// 	withFixedBytesSizeFunc := func(bytesSize int) {
// 		val, err := vus.BindUintValue()(bytesSize)
// 		expected := uint64(6)
// 		if err != nil {
// 			t.Errorf("BindUintValueTo error: %s", err)
// 		}
// 		if val != expected {
// 			t.Errorf("BindUintValue: got %d, expected %d", val, expected)
// 		}

// 		_, err = missed.BindUintValue()(bytesSize)
// 		if err == nil {
// 			t.Error("BindUintValue: expected field missing error")
// 		}
// 		if err.Error() != "field missed is missing in config source jsonment variables" {
// 			t.Error("BindUintValue: unexpected error message:", err)
// 		}

// 		_, err = k6UserAgent.BindUintValue()(bytesSize)
// 		if err == nil {
// 			t.Error("BindUintValue: expected syntax error")
// 		}
// 		if err.Error() != "BindUintValue: parsing \"foo\": invalid syntax" {
// 			t.Errorf("BindUintValue: unexpected error message")
// 		}
// 	}

// 	intBytesSizes := []int{0, 8, 16, 32, 64}

// 	for _, byteSize := range intBytesSizes {
// 		withFixedBytesSizeFunc(byteSize)
// 	}
// }

// func TestJSONFloatValue(t *testing.T) {
// 	var json = []byte(`{"k6_vus":6,"pi":3.14,"k6_config":"./config.json","k6_user_agent":"foo"}`)

// 	var vus = NewJSONSource(json).From("k6_vus")
// 	var pi = NewJSONSource(json).From("pi")
// 	var k6UserAgent = NewJSONSource(json).From("k6_user_agent")
// 	var missed = NewJSONSource(json).From("missed")

// 	withFixedBytesSizeFunc := func(bytesSize int) {
// 		val, err := vus.BindFloatValue()(bytesSize)
// 		expected := float64(6)
// 		if err != nil {
// 			t.Errorf("BindFloatValue error: %s", err)
// 		}
// 		if val != expected {
// 			t.Errorf("BindFloatValue: got %f, expected %f", val, expected)
// 		}

// 		val, err = pi.BindFloatValue()(bytesSize)
// 		expected = float64(3.14)
// 		if err != nil {
// 			t.Errorf("BindFloatValue error: %s", err)
// 		}
// 		// val != expected doesn't work
// 		if (val-3.14) > 1e20 && (val-3.14) < -1e20 {
// 			t.Errorf("BindFloatValue: got %f, expected %f", val, expected)
// 		}

// 		_, err = missed.BindFloatValue()(bytesSize)
// 		if err == nil {
// 			t.Error("BindFloatValue: expected field missing error")
// 		}
// 		if err.Error() != "field missed is missing in config source jsonment variables" {
// 			t.Error("BindFloatValue: unexpected error message:", err)
// 		}

// 		_, err = k6UserAgent.BindFloatValue()(bytesSize)
// 		if err == nil {
// 			t.Error("BindFloatValue: expected syntax error")
// 		}
// 		if err.Error() != "BindFloatValue: parsing \"foo\": invalid syntax" {
// 			t.Errorf("BindFloatValue: unexpected error message")
// 		}
// 	}

// 	intBytesSizes := []int{0, 8, 16, 32, 64}

// 	for _, byteSize := range intBytesSizes {
// 		withFixedBytesSizeFunc(byteSize)
// 	}
// }
