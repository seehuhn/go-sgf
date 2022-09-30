package sgf

import (
	"strings"
	"testing"
)

func TestRead(t *testing.T) {
	in := `(;FF[4]CA[UTF-8]GM[1]SZ[9])`
	c, err := Read(strings.NewReader(in))
	if err != nil {
		t.Error(err)
	}
	if len(c) != 1 {
		t.Errorf("expected 1 game, got %d", len(c))
	}
	tree := c[0]
	if len(tree.Children) != 0 {
		t.Errorf("expected 0 children, got %d", len(tree.Children))
	}
	if len(tree.Properties) != 4 {
		t.Errorf("expected 4 properties, got %d", len(tree.Properties))
	}
	if len(tree.Properties["SZ"]) != 1 || tree.Properties["SZ"][0] != "9" {
		t.Errorf("expected SZ[9], got %v", tree.Properties["SZ"])
	}
}
