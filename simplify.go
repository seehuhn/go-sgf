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
