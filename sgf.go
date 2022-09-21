// seehuhn.de/go/sgf - read and write Smart Game Format (SGF) files
// Copyright (C) 2022  Jochen Voss <voss@seehuhn.de>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package sgf

import (
	"fmt"
	"unicode"
)

// A Collection is a list of games.
type Collection []*GameTree

// A GameTree represents a single game.
type GameTree struct {
	Nodes    []Node
	Children []*GameTree
}

// MainVariation returns the main variation of the game tree.
func (g *GameTree) MainVariation() []Node {
	var res []Node
	for {
		res = append(res, g.Nodes...)
		if len(g.Children) == 0 {
			break
		}
		g = g.Children[0]
	}
	return res
}

// A Node in the game tree is represented by a map from property names to
// property values.
type Node map[string][]string

func (n Node) getSingle(name string) (string, error) {
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
func (n Node) GetNumber(name string) (int, error) {
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
func (n Node) GetNumberDefault(name string, defaultValue int) (int, error) {
	val, err := n.GetNumber(name)
	if _, ok := err.(*missingError); ok {
		return defaultValue, nil
	}
	return val, err
}

// GetSimpleText returns the value of the property with the given name as a
// simple text.  If the property is missing or has more than one value,
// an error is returned.
func (n Node) GetSimpleText(name string) (string, error) {
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
func (n Node) GetSimpleTextDefault(name string, defaultValue string) (string, error) {
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

type sgfError struct {
	msg string
}

func (e *sgfError) Error() string {
	return e.msg
}

func newErrorf(format string, args ...interface{}) error {
	return &sgfError{fmt.Sprintf(format, args...)}
}
