// Copyright (c) 2012-2017 The ANTLR Project. All rights reserved.
// Use of this file is governed by the BSD 3-clause license that
// can be found in the LICENSE.txt file in the project root.

package antlr

import "strings"

type TokenSource interface {
	NextToken() Token
	Skip()
	More()
	GetLine() int
	GetCharPositionInLine() int
	GetInputStream() CharStream
	GetSourceName() string
	setTokenFactory(factory TokenFactory)
	GetTokenFactory() TokenFactory
}

type ListTokenSource struct {
	tokens  []Token
	eof     Token
	idx     int
	Name    string
	Factory TokenFactory
}

func NewListTokenSource(tokens []Token, name string) *ListTokenSource {
	return &ListTokenSource{
		tokens:  tokens,
		Name:    name,
		Factory: CommonTokenFactoryDEFAULT,
	}
}

func (s *ListTokenSource) makeEOFToken() {
	if s.eof != nil {
		return
	}

	start := -1
	stop := -1

	if len(s.tokens) > 0 {
		prevStop := s.tokens[len(s.tokens)-1].GetStop()
		if prevStop != -1 {
			start = prevStop + 1
		}
	}

	if start > -1 {
		stop = start - 1
	}

	pair := &TokenSourceCharStreamPair{s, s.GetInputStream()}

	s.eof = s.Factory.Create(pair, TokenEOF, "EOF", TokenDefaultChannel, start, stop, s.GetLine(), s.GetCharPositionInLine())
}

func (s *ListTokenSource) NextToken() Token {
	if s.idx >= len(s.tokens) {
		s.makeEOFToken()
		return s.eof
	}

	t := s.tokens[s.idx]
	if s.idx == len(s.tokens) && t.GetTokenType() == TokenEOF {
		s.eof = t
	}

	s.idx++
	return t
}

func (s *ListTokenSource) GetInputStream() CharStream {
	if s.idx < len(s.tokens) {
		return s.tokens[s.idx].GetInputStream()
	}

	if s.eof != nil {
		return s.eof.GetInputStream()
	}

	if len(s.tokens) > 0 {
		return s.tokens[len(s.tokens)-1].GetInputStream()
	}

	return nil
}

func (s *ListTokenSource) GetLine() int {
	if s.idx < len(s.tokens) {
		return s.tokens[s.idx].GetLine()
	}

	if s.eof != nil {
		return s.eof.GetLine()
	}

	if len(s.tokens) == 0 {
		return -1
	}

	last := s.tokens[len(s.tokens)-1]
	numLines := last.GetLine() + strings.Count(last.GetText(), "\n")

	return numLines
}

func (s *ListTokenSource) GetCharPositionInLine() int {
	if s.idx < len(s.tokens) {
		return s.tokens[s.idx].GetColumn()
	} else if s.eof != nil {
		return s.eof.GetColumn()
	} else if len(s.tokens) == 0 {
		return 0
	}

	last := s.tokens[len(s.tokens)-1]
	text := last.GetText()

	lastNewLine := strings.LastIndex(text, "\n")
	if lastNewLine == -1 {
		return len(text) - lastNewLine - 1
	}

	return last.GetColumn() + last.GetStop() - last.GetStart() + 1
}

func (s *ListTokenSource) Skip() {}

func (s *ListTokenSource) More() {}

func (s *ListTokenSource) GetSourceName() string {
	return s.Name
}

func (s *ListTokenSource) setTokenFactory(f TokenFactory) {
	s.Factory = f
}

func (s *ListTokenSource) GetTokenFactory() TokenFactory {
	return s.Factory
}
