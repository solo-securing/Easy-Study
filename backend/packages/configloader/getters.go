package configloader

import (
	"fmt"
	"time"
)

// Int64 returns the int64 value of a given key path or 0 if the path
// does not exist or if the value is not a valid int64.
func (cl *ConfigLoader) Int64(path string) int64 {
	if v := cl.Get(path); v != nil {
		i, _ := toInt64(v)
		return i
	}
	return 0
}

// MustInt64 returns the int64 value of a given key path or panics
// if the value is not set or set to default value of 0.
func (cl *ConfigLoader) MustInt64(path string) int64 {
	val := cl.Int64(path)
	if val == 0 {
		panic(fmt.Sprintf("invalid value: %s=%v", path, val))
	}
	return val
}

// Int64s returns the []int64 slice value of a given key path or an
// empty []int64 slice if the path does not exist or if the value
// is not a valid int slice.
func (cl *ConfigLoader) Int64s(path string) []int64 {
	o := cl.Get(path)
	if o == nil {
		return []int64{}
	}

	var out []int64
	switch v := o.(type) {
	case []int64:
		return v
	case []int:
		out = make([]int64, 0, len(v))
		for _, vi := range v {
			i, err := toInt64(vi)

			// On error, return as it's not a valid
			// int slice.
			if err != nil {
				return []int64{}
			}
			out = append(out, i)
		}
		return out
	case []any:
		out = make([]int64, 0, len(v))
		for _, vi := range v {
			i, err := toInt64(vi)

			// On error, return as it's not a valid
			// int slice.
			if err != nil {
				return []int64{}
			}
			out = append(out, i)
		}
		return out
	}

	return []int64{}
}

// MustInt64s returns the []int64 slice value of a given key path or panics
// if the value is not set or its default value.
func (cl *ConfigLoader) MustInt64s(path string) []int64 {
	val := cl.Int64s(path)
	if len(val) == 0 {
		panic(fmt.Sprintf("invalid value: %s=%v", path, val))
	}
	return val
}

// Int64Map returns the map[string]int64 value of a given key path
// or an empty map[string]int64 if the path does not exist or if the
// value is not a valid int64 map.
func (cl *ConfigLoader) Int64Map(path string) map[string]int64 {
	var (
		out = map[string]int64{}
		o   = cl.Get(path)
	)
	if o == nil {
		return out
	}

	mp, ok := o.(map[string]any)
	if !ok {
		return out
	}

	out = make(map[string]int64, len(mp))
	for k, v := range mp {
		switch i := v.(type) {
		case int64:
			out[k] = i
		default:
			// Attempt a conversion.
			iv, err := toInt64(i)
			if err != nil {
				return map[string]int64{}
			}
			out[k] = iv
		}
	}
	return out
}

// MustInt64Map returns the map[string]int64 value of a given key path
// or panics if it isn't set or set to default value.
func (cl *ConfigLoader) MustInt64Map(path string) map[string]int64 {
	val := cl.Int64Map(path)
	if len(val) == 0 {
		panic(fmt.Sprintf("invalid value: %s=%v", path, val))
	}
	return val
}

// Int returns the int value of a given key path or 0 if the path
// does not exist or if the value is not a valid int.
func (cl *ConfigLoader) Int(path string) int {
	return int(cl.Int64(path))
}

// MustInt returns the int value of a given key path or panics
// if it isn't set or set to default value of 0.
func (cl *ConfigLoader) MustInt(path string) int {
	val := cl.Int(path)
	if val == 0 {
		panic(fmt.Sprintf("invalid value: %s=%v", path, val))
	}
	return val
}

// Ints returns the []int slice value of a given key path or an
// empty []int slice if the path does not exist or if the value
// is not a valid int slice.
func (cl *ConfigLoader) Ints(path string) []int {
	o := cl.Get(path)
	if o == nil {
		return []int{}
	}

	var out []int
	switch v := o.(type) {
	case []int:
		return v
	case []int64:
		out = make([]int, 0, len(v))
		for _, vi := range v {
			out = append(out, int(vi))
		}
		return out
	case []any:
		out = make([]int, 0, len(v))
		for _, vi := range v {
			i, err := toInt64(vi)

			// On error, return as it's not a valid
			// int slice.
			if err != nil {
				return []int{}
			}
			out = append(out, int(i))
		}
		return out
	}

	return []int{}
}

// MustInts returns the []int slice value of a given key path or panics
// if the value is not set or set to default value.
func (cl *ConfigLoader) MustInts(path string) []int {
	val := cl.Ints(path)
	if len(val) == 0 {
		panic(fmt.Sprintf("invalid value: %s=%v", path, val))
	}
	return val
}

// IntMap returns the map[string]int value of a given key path
// or an empty map[string]int if the path does not exist or if the
// value is not a valid int map.
func (cl *ConfigLoader) IntMap(path string) map[string]int {
	var (
		mp  = cl.Int64Map(path)
		out = make(map[string]int, len(mp))
	)
	for k, v := range mp {
		out[k] = int(v)
	}
	return out
}

// MustIntMap returns the map[string]int value of a given key path or panics
// if the value is not set or set to default value.
func (cl *ConfigLoader) MustIntMap(path string) map[string]int {
	val := cl.IntMap(path)
	if len(val) == 0 {
		panic(fmt.Sprintf("invalid value: %s=%v", path, val))
	}
	return val
}

// Float64 returns the float64 value of a given key path or 0 if the path
// does not exist or if the value is not a valid float64.
func (cl *ConfigLoader) Float64(path string) float64 {
	if v := cl.Get(path); v != nil {
		f, _ := toFloat64(v)
		return f
	}
	return 0
}

// MustFloat64 returns the float64 value of a given key path or panics
// if it isn't set or set to default value 0.
func (cl *ConfigLoader) MustFloat64(path string) float64 {
	val := cl.Float64(path)
	if val == 0 {
		panic(fmt.Sprintf("invalid value: %s=%v", path, val))
	}
	return val
}

// Float64s returns the []float64 slice value of a given key path or an
// empty []float64 slice if the path does not exist or if the value
// is not a valid float64 slice.
func (cl *ConfigLoader) Float64s(path string) []float64 {
	o := cl.Get(path)
	if o == nil {
		return []float64{}
	}

	var out []float64
	switch v := o.(type) {
	case []float64:
		return v
	case []any:
		out = make([]float64, 0, len(v))
		for _, vi := range v {
			i, err := toFloat64(vi)

			// On error, return as it's not a valid
			// int slice.
			if err != nil {
				return []float64{}
			}
			out = append(out, i)
		}
		return out
	}

	return []float64{}
}

// MustFloat64s returns the []Float64 slice value of a given key path or panics
// if the value is not set or set to default value.
func (cl *ConfigLoader) MustFloat64s(path string) []float64 {
	val := cl.Float64s(path)
	if len(val) == 0 {
		panic(fmt.Sprintf("invalid value: %s=%v", path, val))
	}
	return val
}

// Float64Map returns the map[string]float64 value of a given key path
// or an empty map[string]float64 if the path does not exist or if the
// value is not a valid float64 map.
func (cl *ConfigLoader) Float64Map(path string) map[string]float64 {
	var (
		out = map[string]float64{}
		o   = cl.Get(path)
	)
	if o == nil {
		return out
	}

	mp, ok := o.(map[string]any)
	if !ok {
		return out
	}

	out = make(map[string]float64, len(mp))
	for k, v := range mp {
		switch i := v.(type) {
		case float64:
			out[k] = i
		default:
			// Attempt a conversion.
			iv, err := toFloat64(i)
			if err != nil {
				return map[string]float64{}
			}
			out[k] = iv
		}
	}
	return out
}

// MustFloat64Map returns the map[string]float64 value of a given key path or panics
// if the value is not set or set to default value.
func (cl *ConfigLoader) MustFloat64Map(path string) map[string]float64 {
	val := cl.Float64Map(path)
	if len(val) == 0 {
		panic(fmt.Sprintf("invalid value: %s=%v", path, val))
	}
	return val
}

// Duration returns the time.Duration value of a given key path assuming
// that the key contains a valid numeric value.
func (cl *ConfigLoader) Duration(path string) time.Duration {
	// Look for a parsable string representation first.
	if v := cl.Int64(path); v != 0 {
		return time.Duration(v)
	}

	v, _ := time.ParseDuration(cl.String(path))
	return v
}

// MustDuration returns the time.Duration value of a given key path or panics
// if it isn't set or set to default value 0.
func (cl *ConfigLoader) MustDuration(path string) time.Duration {
	val := cl.Duration(path)
	if val == 0 {
		panic(fmt.Sprintf("invalid value: %s=%v", path, val))
	}
	return val
}

// Time attempts to parse the value of a given key path and return time.Time
// representation. If the value is numeric, it is treated as a UNIX timestamp
// and if it's string, a parse is attempted with the given layout.
func (cl *ConfigLoader) Time(path, layout string) time.Time {
	// Unix timestamp?
	v := cl.Int64(path)
	if v != 0 {
		return time.Unix(v, 0)
	}

	// String representation.
	s := cl.String(path)
	if s != "" {
		t, _ := time.Parse(layout, s)
		return t
	}

	return time.Time{}
}

// MustTime attempts to parse the value of a given key path and return time.Time
// representation. If the value is numeric, it is treated as a UNIX timestamp
// and if it's string, a parse is attempted with the given layout. It panics if
// the parsed time is zero.
func (cl *ConfigLoader) MustTime(path, layout string) time.Time {
	val := cl.Time(path, layout)
	if val.IsZero() {
		panic(fmt.Sprintf("invalid value: %s=%v", path, val))
	}
	return val
}

// String returns the string value of a given key path or "" if the path
// does not exist or if the value is not a valid string.
func (cl *ConfigLoader) String(path string) string {
	if v := cl.Get(path); v != nil {
		if i, ok := v.(string); ok {
			return i
		}
		return fmt.Sprintf("%v", v)
	}
	return ""
}

// MustString returns the string value of a given key path
// or panics if it isn't set or set to default value "".
func (cl *ConfigLoader) MustString(path string) string {
	val := cl.String(path)
	if val == "" {
		panic(fmt.Sprintf("invalid value: %s=%v", path, val))
	}
	return val
}

// Strings returns the []string slice value of a given key path or an
// empty []string slice if the path does not exist or if the value
// is not a valid string slice.
func (cl *ConfigLoader) Strings(path string) []string {
	o := cl.Get(path)
	if o == nil {
		return []string{}
	}

	var out []string
	switch v := o.(type) {
	case []any:
		out = make([]string, 0, len(v))
		for _, u := range v {
			if s, ok := u.(string); ok {
				out = append(out, s)
			} else {
				out = append(out, fmt.Sprintf("%v", u))
			}
		}
		return out
	case []string:
		out := make([]string, len(v))
		copy(out, v)
		return out
	}

	return []string{}
}

// MustStrings returns the []string slice value of a given key path or panics
// if the value is not set or set to default value.
func (cl *ConfigLoader) MustStrings(path string) []string {
	val := cl.Strings(path)
	if len(val) == 0 {
		panic(fmt.Sprintf("invalid value: %s=%v", path, val))
	}
	return val
}

// StringMap returns the map[string]string value of a given key path
// or an empty map[string]string if the path does not exist or if the
// value is not a valid string map.
func (cl *ConfigLoader) StringMap(path string) map[string]string {
	var (
		out = map[string]string{}
		o   = cl.Get(path)
	)
	if o == nil {
		return out
	}

	switch mp := o.(type) {
	case map[string]string:
		out = make(map[string]string, len(mp))
		for k, v := range mp {
			out[k] = v
		}
	case map[string]any:
		out = make(map[string]string, len(mp))
		for k, v := range mp {
			switch s := v.(type) {
			case string:
				out[k] = s
			default:
				// There's a non string type. Return.
				return map[string]string{}
			}
		}
	}

	return out
}

// MustStringMap returns the map[string]string value of a given key path or panics
// if the value is not set or set to default value.
func (cl *ConfigLoader) MustStringMap(path string) map[string]string {
	val := cl.StringMap(path)
	if len(val) == 0 {
		panic(fmt.Sprintf("invalid value: %s=%v", path, val))
	}
	return val
}

// StringsMap returns the map[string][]string value of a given key path
// or an empty map[string][]string if the path does not exist or if the
// value is not a valid strings map.
func (cl *ConfigLoader) StringsMap(path string) map[string][]string {
	var (
		out = map[string][]string{}
		o   = cl.Get(path)
	)
	if o == nil {
		return out
	}

	switch mp := o.(type) {
	case map[string][]string:
		out = make(map[string][]string, len(mp))
		for k, v := range mp {
			out[k] = append(out[k], v...)
		}
	case map[string][]any:
		out = make(map[string][]string, len(mp))
		for k, v := range mp {
			for _, v := range v {
				switch sv := v.(type) {
				case string:
					out[k] = append(out[k], sv)
				default:
					return map[string][]string{}
				}
			}
		}
	case map[string]any:
		out = make(map[string][]string, len(mp))
		for k, v := range mp {
			switch s := v.(type) {
			case []string:
				out[k] = append(out[k], s...)
			case []any:
				for _, v := range s {
					switch sv := v.(type) {
					case string:
						out[k] = append(out[k], sv)
					default:
						return map[string][]string{}
					}
				}
			default:
				// There's a non []interface type. Return.
				return map[string][]string{}
			}
		}
	}

	return out
}

// MustStringsMap returns the map[string][]string value of a given key path or panics
// if the value is not set or set to default value.
func (cl *ConfigLoader) MustStringsMap(path string) map[string][]string {
	val := cl.StringsMap(path)
	if len(val) == 0 {
		panic(fmt.Sprintf("invalid value: %s=%v", path, val))
	}
	return val
}

// Bytes returns the []byte value of a given key path or an empty
// []byte slice if the path does not exist or if the value is not a valid string.
func (cl *ConfigLoader) Bytes(path string) []byte {
	return []byte(cl.String(path))
}

// MustBytes returns the []byte value of a given key path or panics
// if the value is not set or set to default value.
func (cl *ConfigLoader) MustBytes(path string) []byte {
	val := cl.Bytes(path)
	if len(val) == 0 {
		panic(fmt.Sprintf("invalid value: %s=%v", path, val))
	}
	return val
}

// Bool returns the bool value of a given key path or false if the path
// does not exist or if the value is not a valid bool representation.
// Accepted string representations of bool are the ones supported by strconv.ParseBool.
func (cl *ConfigLoader) Bool(path string) bool {
	if v := cl.Get(path); v != nil {
		b, _ := toBool(v)
		return b
	}
	return false
}

// Bools returns the []bool slice value of a given key path or an
// empty []bool slice if the path does not exist or if the value
// is not a valid bool slice.
func (cl *ConfigLoader) Bools(path string) []bool {
	o := cl.Get(path)
	if o == nil {
		return []bool{}
	}

	var out []bool
	switch v := o.(type) {
	case []any:
		out = make([]bool, 0, len(v))
		for _, u := range v {
			b, err := toBool(u)
			if err != nil {
				return nil
			}
			out = append(out, b)
		}
		return out
	case []bool:
		return out
	}
	return nil
}

// MustBools returns the []bool value of a given key path or panics
// if the value is not set or set to default value.
func (cl *ConfigLoader) MustBools(path string) []bool {
	val := cl.Bools(path)
	if len(val) == 0 {
		panic(fmt.Sprintf("invalid value: %s=%v", path, val))
	}
	return val
}

// BoolMap returns the map[string]bool value of a given key path
// or an empty map[string]bool if the path does not exist or if the
// value is not a valid bool map.
func (cl *ConfigLoader) BoolMap(path string) map[string]bool {
	var (
		out = map[string]bool{}
		o   = cl.Get(path)
	)
	if o == nil {
		return out
	}

	mp, ok := o.(map[string]any)
	if !ok {
		return out
	}
	out = make(map[string]bool, len(mp))
	for k, v := range mp {
		switch i := v.(type) {
		case bool:
			out[k] = i
		default:
			// Attempt a conversion.
			b, err := toBool(i)
			if err != nil {
				return map[string]bool{}
			}
			out[k] = b
		}
	}

	return out
}

// MustBoolMap returns the map[string]bool value of a given key path or panics
// if the value is not set or set to default value.
func (cl *ConfigLoader) MustBoolMap(path string) map[string]bool {
	val := cl.BoolMap(path)
	if len(val) == 0 {
		panic(fmt.Sprintf("invalid value: %s=%v", path, val))
	}
	return val
}
