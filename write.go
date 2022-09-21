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
	"bufio"
	"io"
	"sort"

	"golang.org/x/exp/maps"
)

func (c Collection) Write(w io.Writer) error {
	wb := bufio.NewWriter(w)

	for _, g := range c {
		g.write(wb)
	}

	wb.WriteRune('\n')
	return wb.Flush()
}

func (g *GameTree) write(w *bufio.Writer) {
	w.WriteRune('(')
	for _, n := range g.Nodes {
		n.write(w)
	}
	for _, c := range g.Children {
		w.WriteRune('\n')
		c.write(w)
	}
	w.WriteRune(')')
}

func (n Node) write(w *bufio.Writer) {
	w.WriteRune(';')
	keys := maps.Keys(n)
	sort.Strings(keys)
	for j, key := range keys {
		if j > 0 {
			w.WriteRune('\n')
		}
		w.WriteString(key)
		for _, value := range n[key] {
			w.WriteRune('[')
			w.WriteString(value)
			w.WriteRune(']')
		}
	}
}
