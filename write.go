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
