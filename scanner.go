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
	"fmt"
	"unicode"
	"unicode/utf8"
)

// stateFn represents the state of the scanner
// as a function that returns the next state.
type stateFn func(*scanner) stateFn

// A scanner breaks a string into tokens.
type scanner struct {
	input  string
	start  int // start position of current token
	pos    int // current position in input
	width  int // width of last rune read from input
	tokens chan<- *token

	eolSeen   bool
	lineStart int
	lineNo    int // 0 based
}

func (s *scanner) run() {
	for state := scanStart; state != nil; {
		state = state(s)
	}
	close(s.tokens)
}

func (s *scanner) next() rune {
	if s.eolSeen {
		s.eolSeen = false
		s.lineStart = s.pos
		s.lineNo++
	}

	if s.pos >= len(s.input) {
		s.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(s.input[s.pos:])
	s.pos += w
	s.width = w

	if r == '\n' {
		s.eolSeen = true
	}

	return r
}

func (s *scanner) backup() {
	s.pos -= s.width
	s.width = 0
}

func (s *scanner) emit(t tokenType) {
	s.tokens <- &token{
		typ:  t,
		val:  s.input[s.start:s.pos],
		line: s.lineNo,
		col:  s.start - s.lineStart,
	}
	s.start = s.pos
}

func (s *scanner) error(msg string) {
	s.tokens <- &token{
		typ:  tokenError,
		val:  msg,
		line: s.lineNo,
		col:  s.start - s.lineStart,
	}
}

func (s *scanner) ignore() {
	s.start = s.pos
}

func (s *scanner) skipWhiteSpace() {
	for {
		r := s.next()
		if !unicode.IsSpace(r) {
			// note that eof (rune 0) doesn't count as white space
			s.backup()
			s.ignore()
			return
		}
	}
}

const eof = rune(0)

func scanStart(s *scanner) stateFn {
	s.skipWhiteSpace()

	r := s.next()
	switch {
	case r == eof:
		s.emit(tokenEOF)
		return nil
	case r == '(':
		s.emit(tokenParenOpen)
	case r == ')':
		s.emit(tokenParenClose)
	case r == '[':
		return scanPropValue
	case r == ';':
		s.emit(tokenSemicolon)
	case r >= 'A' && r <= 'Z':
		s.backup()
		return scanPropIdent
	default:
		s.error(fmt.Sprintf("unexpected character %q", r))
	}

	return scanStart
}

func scanPropIdent(s *scanner) stateFn {
	for {
		r := s.next()
		if r < 'A' || r > 'Z' {
			s.backup()
			s.emit(tokenPropIdent)
			return scanStart
		}
	}
}

func scanPropValue(s *scanner) stateFn {
	s.ignore()
	escaped := false
	for {
		r := s.next()
		if escaped {
			escaped = false
			continue
		} else if r == '\\' {
			escaped = true
			continue
		} else if r == ']' {
			s.backup()
			s.emit(tokenPropValue)
			s.next()
			return scanStart
		} else if r == eof {
			s.error("EOF while scanning PropValue")
			return nil
		}
	}
}

type tokenType int

const (
	tokenError tokenType = iota
	tokenEOF
	tokenParenOpen
	tokenParenClose
	tokenSemicolon
	tokenPropIdent
	tokenPropValue
)

type token struct {
	typ  tokenType
	val  string
	line int // 0 based
	col  int // 0 based
}

func (i token) String() string {
	t := i.typ
	if t == tokenEOF {
		return "EOF"
	}

	v := i.val
	if t != tokenError && len(i.val) > 10 {
		v = fmt.Sprintf("%.10s...", v)
	}
	return v
}
