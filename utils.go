package croconf

import "strconv"

func parseInt(s string, base int, bitSize int) (int64, *BindValueError) {
	val, err := strconv.ParseInt(s, base, bitSize)
	if err != nil {
		return 0, NewBindValueError("parseInt", s, err)
	}
	return val, nil
}

func parseUint(s string, base int, bitSize int) (uint64, *BindValueError) {
	val, err := strconv.ParseUint(s, base, bitSize)
	if err != nil {
		return 0, NewBindValueError("parseUint", s, err)
	}
	return val, nil
}

func parseFloat(s string, bitSize int) (float64, *BindValueError) {
	val, err := strconv.ParseFloat(s, bitSize)
	if err != nil {
		return 0, NewBindValueError("parseFloat", s, err)
	}
	return val, nil
}
