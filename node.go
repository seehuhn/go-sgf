package sgf

import (
	"fmt"
	"strconv"
	"unicode"
)

// Properties of a game tree node are given as a map from property names to
// property values.
type Properties map[string][]string

func (n Properties) getSingle(name string) (string, error) {
	vals, ok := n[name]
	if !ok {
		return "", &missingError{name}
	}
	if len(vals) != 1 {
		return "", newErrorf("property %q has %d values, expected 1", name, len(vals))
	}
	return vals[0], nil
}

// GetNumber returns the value of the property with the given name as a
// number.  If the property is missing, has more than one value, or
// the value is not an integer, an error is returned.
func (n Properties) GetNumber(name string) (int, error) {
	str, err := n.getSingle(name)
	if err != nil {
		return 0, err
	}

	s := 1
	if str[0] == '-' {
		s = -1
		str = str[1:]
	}
	abs := 0
	for _, c := range str {
		if c < '0' || c > '9' || abs > 1<<31/10 {
			return 0, newErrorf("property %q has invalid value %q", name, str)
		}
		abs = 10*abs + int(c-'0')
	}

	return s * abs, nil
}

// GetNumberDefault returns the value of the property with the given name as a
// number.  If the property is missing, the defaultValue is returned. It the
// property has more than one value, or the value is not an integer, an error
// is returned.
func (n Properties) GetNumberDefault(name string, defaultValue int) (int, error) {
	val, err := n.GetNumber(name)
	if _, ok := err.(*missingError); ok {
		return defaultValue, nil
	}
	return val, err
}

// GetReal returns the value of the property with the given name as a floating
// point number.  If the property is missing, has more than one value, or the
// value is not a valid number, an error is returned.
func (n Properties) GetReal(name string) (float64, error) {
	str, err := n.getSingle(name)
	if err != nil {
		return 0, err
	}

	x, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, newErrorf("property %q has invalid value %q", name, str)
	}

	return x, nil
}

// GetRealDefault returns the value of the property with the given name as a
// floating point number.  If the property is missing, the defaultValue is
// returned. It the property has more than one value, or the value is not a
// valid number, an error is returned.
func (n Properties) GetRealDefault(name string, defaultValue float64) (float64, error) {
	val, err := n.GetReal(name)
	if _, ok := err.(*missingError); ok {
		return defaultValue, nil
	}
	return val, err
}

// GetSimpleText returns the value of the property with the given name as a
// simple text.  If the property is missing or has more than one value,
// an error is returned.
func (n Properties) GetSimpleText(name string) (string, error) {
	s, err := n.getSingle(name)
	if err != nil {
		return "", err
	}

	res := make([]rune, 0, len(s))
	spaceSeen := false
	escSeen := false
	var nlIgnore rune
	for _, r := range s {
		if escSeen {
			escSeen = false
			if r == '\n' {
				nlIgnore = '\r'
				continue
			} else if r == '\r' {
				nlIgnore = '\n'
				continue
			}
		} else if r == '\\' {
			escSeen = true
			continue
		}

		skip := r == nlIgnore
		nlIgnore = 0
		if skip {
			continue
		}

		if unicode.IsSpace(r) {
			if !spaceSeen {
				res = append(res, ' ')
				spaceSeen = true
			}
			continue
		}
		spaceSeen = false

		res = append(res, r)
	}
	return string(res), nil
}

// GetSimpleTextDefault returns the value of the property with the given name
// as a simple text.  If the property is missing, the defaultValue is returned.
// If the property has more than one value, an error is returned.
func (n Properties) GetSimpleTextDefault(name string, defaultValue string) (string, error) {
	val, err := n.GetSimpleText(name)
	if _, ok := err.(*missingError); ok {
		return defaultValue, nil
	}
	return val, err
}

type missingError struct {
	name string
}

func (e *missingError) Error() string {
	return fmt.Sprintf("missing property %q", e.name)
}
