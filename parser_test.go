package sgf

import (
	"os"
	"testing"
)

func TestParser(t *testing.T) {
	fd, err := os.Open("test.sgf")
	if err != nil {
		t.Fatal(err)
	}
	defer fd.Close()

	c, err := Read(fd)
	if err != nil {
		t.Fatal(err)
	}

	c.Simplify()

	err = c.Write(os.Stdout)
	if err != nil {
		t.Fatal(err)
	}
}
