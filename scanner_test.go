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
