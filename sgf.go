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
