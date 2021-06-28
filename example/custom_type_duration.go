package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Duration time.Duration

func (d Duration) String() string {
	return time.Duration(d).String()
}

// ParseExtendedDuration is a helper function that allows for string duration
// values containing days.
func ParseExtendedDuration(data string) (result time.Duration, err error) {
	// Assume millisecond values if data is provided with no units
	if t, errp := strconv.ParseFloat(data, 64); errp == nil {
		return time.Duration(t * float64(time.Millisecond)), nil
	}

	dPos := strings.IndexByte(data, 'd')
	if dPos < 0 {
		return time.ParseDuration(data)
	}

	var hours time.Duration
	if dPos+1 < len(data) { // case "12d"
		hours, err = time.ParseDuration(data[dPos+1:])
		if err != nil {
			return
		}
		if hours < 0 {
			return 0, fmt.Errorf("invalid time format '%s'", data[dPos+1:])
		}
	}

	days, err := strconv.ParseInt(data[:dPos], 10, 64)
	if err != nil {
		return
	}
	if days < 0 {
		hours = -hours
	}
	return time.Duration(days)*24*time.Hour + hours, nil
}

// UnmarshalText converts text data to Duration
func (d *Duration) UnmarshalText(data []byte) error {
	v, err := ParseExtendedDuration(string(data))
	if err != nil {
		return err
	}
	*d = Duration(v)
	return nil
}

// UnmarshalJSON converts JSON data to Duration
// and implements croconf.CustomValue interface
func (d *Duration) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == '"' {
		var str string
		if err := json.Unmarshal(data, &str); err != nil {
			return err
		}

		v, err := ParseExtendedDuration(str)
		if err != nil {
			return err
		}

		*d = Duration(v)
	} else if t, errp := strconv.ParseFloat(string(data), 64); errp == nil {
		*d = Duration(t * float64(time.Millisecond))
	} else {
		return fmt.Errorf("'%s' is not a valid duration value", string(data))
	}

	return nil
}

// ParseFromString implements croconf.CustomValue interface
func (d *Duration) ParseFromString(data string) error {
	duration, err := ParseExtendedDuration(data)
	if err != nil {
		return err
	}

	*d = Duration(duration)
	return nil
}
