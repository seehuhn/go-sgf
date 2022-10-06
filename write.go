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
	buf := bufio.NewWriter(w)

	for _, g := range c {
		g.write(buf)
	}

	_, _ = buf.WriteRune('\n')
	return buf.Flush()
}

func (g *Tree) write(buf *bufio.Writer) {
	_, _ = buf.WriteRune('(')
	for {
		g.Properties.write(buf)
		if len(g.Children) != 1 {
			break
		}
		g = g.Children[0]
	}
	for _, c := range g.Children {
		_, _ = buf.WriteRune('\n')
		c.write(buf)
	}
	_, _ = buf.WriteRune(')')
}

func (n Properties) write(buf *bufio.Writer) {
	_, _ = buf.WriteRune(';')
	keys := maps.Keys(n)
	sort.Strings(keys)
	for j, key := range keys {
		if j > 0 {
			_, _ = buf.WriteRune('\n')
		}
		_, _ = buf.WriteString(key)
		for _, value := range n[key] {
			_, _ = buf.WriteRune('[')
			_, _ = buf.WriteString(value)
			_, _ = buf.WriteRune(']')
		}
	}
}
