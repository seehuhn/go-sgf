package sgf

import "testing"

func TestBoardSize(t *testing.T) {
	type testCase struct {
		in     Properties
		width  int
		height int
	}
	cases := []testCase{
		{Properties{}, 19, 19},
		{Properties{"SZ": []string{"18"}}, 18, 18},
		{Properties{"SZ": []string{"17:16"}}, 17, 16},
		{Properties{"SZ": []string{"0"}}, 0, 0},
	}
	for _, test := range cases {
		tree := &Tree{Properties: test.in}
		sz, err := tree.GetBoardSize()
		if err != nil && test.width != 0 {
			t.Errorf("GetBoardSize() returned error %v", err)
		} else if err == nil && test.width == 0 {
			t.Errorf("GetBoardSize() returned %v, expected error", sz)
		} else if sz.Width != test.width || sz.Height != test.height {
			t.Errorf("GetBoardSize() returned %v, want %v", sz, BoardSize{test.width, test.height})
		}
	}
}

func TestMove1(t *testing.T) {
	mm := []Move{
		{X: -1, Y: -1},
	}
	for x := int8(0); x < 19; x++ {
		for y := int8(0); y < 19; y++ {
			mm = append(mm, Move{X: x, Y: y})
		}
	}

	b := BoardSize{19, 19}
	for _, m := range mm {
		m2 := b.DecodeMove(b.EncodeMove(m))
		if m != m2 {
			t.Errorf("move (%d,%d) != (%d,%d)", m.X, m.Y, m2.X, m2.Y)
		}
	}
}

func TestMove2(t *testing.T) {
	mm := []string{
		"", "aa", "ss", "as", "sa", "ij",
	}
	b := BoardSize{19, 19}
	for _, s := range mm {
		m := b.DecodeMove(s)
		s2 := b.EncodeMove(m)
		if s != s2 {
			t.Errorf("move %q != %q", s, s2)
		}
	}
}
