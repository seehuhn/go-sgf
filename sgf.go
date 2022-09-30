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
)

// A Collection is a slice of game trees.
type Collection []*Tree

// A Tree represents a node in a game tree.
type Tree struct {
	Properties
	Children []*Tree
}

// IsLinear checks whether the game tree is linear, i.e. whether all
// nodes have at most one child.
func (t *Tree) IsLinear() bool {
	for {
		switch len(t.Children) {
		case 0:
			return true
		case 1:
			t = t.Children[0]
		default:
			return false
		}
	}
}

// MainVariation returns the main variation of the game tree.
// This is the sequence of nodes starting at the root node and
// following the first child of each node.
func (t *Tree) MainVariation() []Properties {
	var res []Properties
	for {
		res = append(res, t.Properties)
		if len(t.Children) == 0 {
			break
		}
		t = t.Children[0]
	}
	return res
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
