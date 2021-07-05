package croconf

import (
	"fmt"
	"strconv"
)

func checkIntBitsize(val int64, bitSize int) error {
	// See MinInt and MaxInt values in https://golang.org/pkg/math/#pkg-constants
	min, max := int64(-1<<(bitSize-1)), int64(1<<(bitSize-1)-1)
	if val < min || val > max {
		return fmt.Errorf("invalid value %d, it must be between %d and %d", val, min, max)
	}
	return nil
}

func checkUintBitsize(val uint64, bitSize int) error {
	// See MaxUint values in https://golang.org/pkg/math/#pkg-constants
	if max := uint64(1<<bitSize - 1); val > max {
		return fmt.Errorf("invalid value %d, it must be between 0 and %d", val, max)
	}
	return nil
}

func parseInt(s string) (int64, *BindValueError) {
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, NewBindValueError("parseInt", s, err)
	}
	return val, nil
}

func parseUint(s string) (uint64, *BindValueError) {
	val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, NewBindValueError("parseUint", s, err)
	}
	return val, nil
}

func parseFloat(s string) (float64, *BindValueError) {
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, NewBindValueError("parseFloat", s, err)
	}
	return val, nil
}
