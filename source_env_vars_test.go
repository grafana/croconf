package croconf

import "testing"

var environ = []string{"K6_VUS=6", "PI=3.14", "K6_CONFIG=./config.json", "K6_USER_AGENT=foo"}

var vus = NewSourceFromEnv(environ).From("K6_VUS")
var pi = NewSourceFromEnv(environ).From("PI")
var k6UserAgent = NewSourceFromEnv(environ).From("K6_USER_AGENT")
var missed = NewSourceFromEnv(environ).From("MISSED")

func TestBindIntValue(t *testing.T) {

	withFixedBytesSizeFunc := func(bytesSize int) {
		val, err := vus.BindIntValue()(0)
		expected := int64(6)
		if err != nil {
			t.Errorf("BindIntValueTo error: %s", err)
		}
		if val != expected {
			t.Errorf("BindIntValue: got %d, expected %d", val, expected)
		}

		_, err = k6UserAgent.BindIntValue()(0)
		if err == nil {
			t.Error("BindIntValue: expected syntax error")
		}
		if err.Error() != "BindIntValueTo: parsing \"foo\": invalid syntax" {
			t.Errorf("BindIntValueTo: unexpected error message")
		}

		_, err = missed.BindIntValue()(0)
		if err == nil {
			t.Error("BindIntValueTo: expected field missing error")
		}
		if err.Error() != "BindIntValueTo: binding name MISSED not found in config source" {
			t.Error("BindIntValueTo: unexpected error message", err)
		}
	}

	intBytesSizes := []int{0, 8, 16, 32, 64}

	for _, byteSize := range intBytesSizes {
		withFixedBytesSizeFunc(byteSize)
	}
}
