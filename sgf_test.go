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
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSimpleText(t *testing.T) {
	cases := []struct {
		in  string
		out string
	}{
		{"", ""},
		{"a", "a"},
		{"This is a test.", "This is a test."},
		{"a\nb", "a b"},
		{"a\rb", "a b"},
		{"a\n\rb", "a b"},
		{"a\r\nb", "a b"},
		{"a[\\]b", "a[]b"},
		{"a\\:b", "a:b"},
		{"a\\\nb", "ab"},
		{"a\\\rb", "ab"},
		{"a\\\n\rb", "ab"},
		{"a\\\r\nb", "ab"},
		{"a \n\r\t b", "a b"},
	}
	for i, test := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			n := &Properties{
				"TEST": []string{test.in},
			}
			got, err := n.GetSimpleText("TEST")
			if err != nil {
				t.Error(err)
			}
			if got != test.out {
				t.Errorf("simpleText(%q) = %q, want %q", test.in, got, test.out)
			}
		})
	}
}

func TestExamples(t *testing.T) {
	for _, test := range examples {
		r := strings.NewReader(test)
		_, err := Read(r)
		if err != nil {
			t.Errorf("Read(%q) failed: %v", test, err)
		}
	}

	for _, test := range counterExamples {
		r := strings.NewReader(test)
		_, err := Read(r)
		if err == nil {
			t.Errorf("Read(%q) succeeded, want failure", test)
		}
	}
}

func FuzzSGF(f *testing.F) {
	for _, test := range examples {
		f.Add(test)
	}
	f.Fuzz(func(t *testing.T, a string) {
		r := strings.NewReader(a)
		c1, err := Read(r)
		if err != nil {
			return
		}

		buf := &bytes.Buffer{}
		err = c1.Write(buf)
		if err != nil {
			t.Fatal(err)
		}

		c2, err := Read(buf)
		if err != nil {
			t.Fatal(err)
		}

		if d := cmp.Diff(c1, c2); d != "" {
			t.Errorf("Read(Write(c)) mismatch (-want +got):\n%s", d)
		}
	})
}

var examples = []string{
	`(;FF[4]C[root](;C[a];C[b](;C[c])
(;C[d];C[e]))
(;C[f](;C[g];C[h];C[i])
(;C[j])))`,
	`(;FF[4]GM[1]SZ[19];B[aa];W[bb];B[cc];W[dd];B[ad];W[bd])`,
	`(;FF[4]GM[1]SZ[19];B[aa];W[bb](;B[cc];W[dd];B[ad];W[bd])
(;B[hh];W[hg]))`,
	`(;FF[4]GM[1]SZ[19];B[aa];W[bb](;B[cc]N[Var A];W[dd];B[ad];W[bd])
(;B[hh]N[Var B];W[hg])
(;B[gg]N[Var C];W[gh];B[hh];W[hg];B[kk]))`,
	`(;FF[4]GM[1]SZ[19];B[aa];W[bb](;B[cc];W[dd](;B[ad];W[bd])
(;B[ee];W[ff]))
(;B[hh];W[hg]))`,
	`(;FF[4]GM[1]SZ[19];B[aa];W[bb](;B[cc]N[Var A];W[dd];B[ad];W[bd])
(;B[hh]N[Var B];W[hg])
(;B[gg]N[Var C];W[gh];B[hh]  (;W[hg]N[Var A];B[kk])  (;W[kl]N[Var B])))`,
	`(;)`,
	`(;;;(;;;;)(;;)(;;;(;;)(;)))`,
	`(;(;)(;)(;(;)(;)))`,
	`(;W[tt])`,
}

var counterExamples = []string{
	`()`,
	`(W[tt])`,
	`(;)W[tt]`,
}
