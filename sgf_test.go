package sgf

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

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
