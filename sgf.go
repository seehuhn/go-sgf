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

import "io"

// A Collection is a list of games.
type Collection []*GameTree

// A GameTree represents a single game.
type GameTree struct {
	Nodes    []Node
	Children []*GameTree
}

// A Node in the game tree is represented by a map from property names to
// property values.
type Node map[string][]string

// Read reads a collection of games from r.
func Read(r io.Reader) (Collection, error) {
	body, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return parse(string(body))
}
