package converter

import (
	"fmt"
	"strconv"
)

// ToInt64 takes an interface value and if it is an integer type,
// converts and returns int64. If it's any other type,
// forces it to a string and attempts to do a strconv.Atoi
// to get an integer out.
func ToInt64(v any) (int64, error) {
	switch i := v.(type) {
	case int:
		return int64(i), nil
	case int8:
		return int64(i), nil
	case int16:
		return int64(i), nil
	case int32:
		return int64(i), nil
	case int64:
		return i, nil
	}

	// Force it to a string and try to convert.
	f, err := strconv.ParseFloat(fmt.Sprintf("%v", v), 64)
	if err != nil {
		return 0, err
	}

	return int64(f), nil
}

// ToFloat64 takes a `v any` value and if it is a float type,
// converts and returns a `float64`. If it's any other type, forces it to a
// string and attempts to get a float out using `strconv.ParseFloat`.
func ToFloat64(v any) (float64, error) {
	switch i := v.(type) {
	case float32:
		return float64(i), nil
	case float64:
		return i, nil
	}

	// Force it to a string and try to convert.
	f, err := strconv.ParseFloat(fmt.Sprintf("%v", v), 64)
	if err != nil {
		return f, err
	}

	return f, nil
}

// ToBool takes an interface value and if it is a bool type,
// returns it. If it's any other type, forces it to a string and attempts
// to parse it as a bool using strconv.ParseBool.
func ToBool(v any) (bool, error) {
	if b, ok := v.(bool); ok {
		return b, nil
	}

	// Force it to a string and try to convert.
	b, err := strconv.ParseBool(fmt.Sprintf("%v", v))
	if err != nil {
		return b, err
	}
	return b, nil
}
