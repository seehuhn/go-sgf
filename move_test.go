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
