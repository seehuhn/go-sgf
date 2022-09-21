package sgf

import (
	"testing"
	"unicode"
)

func TestEOF(t *testing.T) {
	if unicode.IsSpace(eof) {
		t.Error("eof cannot be a space character")
	}
}

func TestScanner(t *testing.T) {
	s := scanner{
		input: "a\n12\n",
	}

	cases := []struct {
		r         rune
		lineNo    int
		lineStart int
	}{
		{'a', 0, 0},
		{'\n', 0, 0},
		{'1', 1, 2},
		{'2', 1, 2},
		{'\n', 1, 2},
		{eof, 2, 5},
	}
	for _, test := range cases {
		r := s.next()
		if r != test.r {
			t.Errorf("expected %q, got %q", test.r, r)
		}
		if s.lineNo != test.lineNo {
			t.Errorf("expected lineNo %d, got %d", test.lineNo, s.lineNo)
		}
		if s.lineStart != test.lineStart {
			t.Errorf("expected lineStart %d, got %d", test.lineStart, s.lineStart)
		}
	}
}
