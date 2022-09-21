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

// Simplify simplifies all game trees in the collection.
func (c Collection) Simplify() {
	for i, g := range c {
		c[i] = g.Simplify()
	}
}

// Simplify returns a simplified deep copy of the game tree.
func (g *GameTree) Simplify() *GameTree {
	var nodes []Node
	for {
		for _, origNode := range g.Nodes {
			n := Node{}
			for key, val := range origNode {
				if len(val) > 0 {
					n[key] = val
				}
			}
			if len(n) == 0 {
				continue
			}
			nodes = append(nodes, n)
		}
		if len(g.Children) != 1 {
			break
		}
		g = g.Children[0]
	}

	children := make([]*GameTree, len(g.Children))
	for i, child := range g.Children {
		children[i] = child.Simplify()
	}
	return &GameTree{
		Nodes:    nodes,
		Children: children,
	}
}
