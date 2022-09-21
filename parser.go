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

import "fmt"

type parser struct {
	tokens  <-chan *token
	backlog []*token
}

func parse(s string) (Collection, error) {
	tokens := make(chan *token)
	scanner := &scanner{
		input:  s,
		tokens: tokens,
	}
	go scanner.run()

	p := &parser{
		tokens: tokens,
	}
	c, err := p.parseCollection()
	if err != nil {
		// drain the lexer
		for range tokens {
		}
		return nil, err
	}

	return c, nil
}

func (p *parser) parseCollection() (Collection, error) {
	var c Collection

gameLoop:
	for {
		t := p.peek()

		switch t.typ {
		case tokenEOF:
			break gameLoop
		default:
			g, err := p.parseGameTree()
			if err != nil {
				return nil, err
			}
			c = append(c, g)
		}
	}
	return c, nil
}

func (p *parser) parseGameTree() (*GameTree, error) {
	err := p.require(tokenParenOpen, "GameTree")
	if err != nil {
		return nil, err
	}

	g := &GameTree{}

	for {
		n, err := p.parseNode()
		if err != nil {
			return nil, err
		}
		g.Nodes = append(g.Nodes, n)
		if p.peek().typ != tokenSemicolon {
			break
		}
	}

childLoop:
	for {
		t := p.peek()
		switch t.typ {
		case tokenParenClose:
			break childLoop
		default:
			child, err := p.parseGameTree()
			if err != nil {
				return nil, err
			}
			g.Children = append(g.Children, child)
		}
	}

	err = p.require(tokenParenClose, "closing round bracket")
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (p *parser) parseNode() (Node, error) {
	err := p.require(tokenSemicolon, "Node")
	if err != nil {
		return nil, err
	}

	n := make(Node)
	for {
		t := p.next()
		if t.typ != tokenPropIdent {
			p.backup(t)
			break
		}
		key := t.val

		var values []string
		for {
			t = p.next()
			if t.typ != tokenPropValue {
				p.backup(t)
				break
			}
			values = append(values, t.val)
		}
		if len(values) == 0 {
			return nil, makeError(p.peek(), "property %q has no values", key)
		}

		n[key] = values
	}
	return n, nil
}

func (p *parser) next() *token {
	if len(p.backlog) > 0 {
		n := len(p.backlog) - 1
		t := p.backlog[n]
		p.backlog = p.backlog[:n]
		return t
	}
	return <-p.tokens
}

func (p *parser) backup(t *token) {
	p.backlog = append(p.backlog, t)
}

func (p *parser) peek() *token {
	t := p.next()
	p.backup(t)
	return t
}

func (p *parser) require(tp tokenType, desc string) error {
	t := p.next()
	if t.typ != tp {
		return makeError(t, "expected %s, got %q", desc, t)
	}
	return nil
}

type parseError struct {
	next *token
	msg  string
}

func (err *parseError) Error() string {
	return fmt.Sprintf("line %d, column %d: %s",
		err.next.line+1, err.next.col+1, err.msg)
}

func makeError(next *token, format string, a ...interface{}) error {
	msg := fmt.Sprintf(format, a...)
	return &parseError{next: next, msg: msg}
}
