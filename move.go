package sgf

import (
	"strconv"
	"strings"
)

// Move represents a move in a game of Go.
// The coordinates are 0-based, with (0,0) being the top left corner.
// If the move is a pass, X and Y are -1.
type Move struct {
	X int8 // column, 0-based, from left to right
	Y int8 // row, 0-based, from bottom to top
}

type BoardSize struct {
	Width  int
	Height int
}

func (t *Tree) GetBoardSize() (BoardSize, error) {
	val, err := t.getSingle("SZ")
	if _, ok := err.(*missingError); ok {
		return BoardSize{19, 19}, nil
	} else if err != nil {
		return BoardSize{}, err
	}

	wh := strings.Split(val, ":")
	switch len(wh) {
	case 1:
		sz, err := strconv.Atoi(val)
		if err != nil || sz < 1 || sz > 52 {
			return BoardSize{}, newErrorf("property SZ has invalid value %q", val)
		}
		return BoardSize{sz, sz}, nil
	case 2:
		w, err := strconv.Atoi(wh[0])
		if err != nil || w < 1 || w > 52 {
			return BoardSize{}, newErrorf("property SZ has invalid value %q", val)
		}
		h, err := strconv.Atoi(wh[1])
		if err != nil || h < 1 || h > 52 {
			return BoardSize{}, newErrorf("property SZ has invalid value %q", val)
		}
		return BoardSize{w, h}, nil
	default:
		return BoardSize{}, newErrorf("property SZ has invalid value %q", val)
	}
}

func (sz BoardSize) String() string {
	return strconv.Itoa(sz.Width) + "x" + strconv.Itoa(sz.Height)
}

// in SGF, "aa" is the top left corner
func (sz BoardSize) decodeSgfMove(sgfMove string) Move {
	if sgfMove == "" || (sz.Width <= 19 && sz.Height <= 19 && sgfMove == "tt") {
		return Move{-1, -1}
	}
	if len(sgfMove) != 2 {
		panic("invalid move " + sgfMove)
	}
	x := int(sgfMove[0] - 'a')
	y := (sz.Height - 1) - int(sgfMove[1]-'a')
	if x < 0 || x >= sz.Width || y < 0 || y >= sz.Height {
		panic("invalid move " + sgfMove)
	}
	return Move{int8(x), int8(y)}
}

func (sz BoardSize) encodeSgfMove(move Move) string {
	if move.X < 0 || move.Y < 0 {
		return ""
	}
	return string([]byte{'a' + byte(move.X), 'a' + byte(sz.Height-1) - byte(move.Y)})
}
