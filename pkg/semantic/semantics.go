package semantic

import (
	"github.com/kjbreil/glsp/pkg/location"
	protocol "github.com/kjbreil/glsp/protocol_3_16"
	"sort"
)

type Semantics struct {
	s []Semantic
}

func (s *Semantics) Sort() {
	sort.Slice(s, func(i, j int) bool {
		return s.s[i].Location.LessThan(s.s[j].Location)
	})
}

func (s *Semantics) TokenMap() [][]Token {

	lines := s.Lines()
	tokenMap := make([][]Token, lines)
	for line := range lines {
		semLine := s.Line(line)
		for column := range semLine.LastColumn() {
			t := semLine.TokenAt(line, column)
			tokenMap[line] = append(tokenMap[line], t)
		}
	}
	return tokenMap
}

func TokenMapToProtocol(tokenMap [][]Token) []protocol.UInteger {
	tokens := []protocol.UInteger{}
	var length int
	lastStart := 0
	currStart := 0
	lineDiff := 0
	currToken := TokenNone
	for _, line := range tokenMap {
		for column, t := range line {
			if t != currToken {
				if currToken != TokenNone {
					tokens = append(tokens, makeTokenSlice(lineDiff, currStart-lastStart, length, int(currToken))...)
					lineDiff = 0
					lastStart = currStart
				}
				length = 0
				currToken = t
				currStart = column
			}
			length++
		}
		lineDiff++
		// if currToken != TokenNone {
		// 	tokens = append(tokens, makeTokenSlice(lineDiff, columnDiff, length, int(currToken))...)
		// 	lineDiff = 0
		// }
		currToken = TokenNone
		lastStart = 0
		length = 0
	}

	// 15 = 2, 20 = 0

	return tokens
}

func makeTokenSlice(lineDiff, columnDiff, length, token int) []protocol.UInteger {
	tokens := make([]protocol.UInteger, 5)
	tokens[0] = protocol.UInteger(lineDiff)
	tokens[1] = protocol.UInteger(columnDiff)
	tokens[2] = protocol.UInteger(length)
	tokens[3] = protocol.UInteger(token)

	return tokens
}

func (s *Semantics) TokenAt(l, c int) Token {
	t := TokenNone

	loc := &location.Range{
		Start: location.Point{
			Line:   l,
			Column: c,
		},
		End: location.Point{
			Line:   l,
			Column: c,
		},
	}
	for _, sem := range s.s {
		if loc.Within(sem.Location) && sem.Token < t {
			t = sem.Token
		}
	}
	return t
}

func (s *Semantics) LastColumn() int {
	var lastColumn int
	for _, sem := range s.s {
		lastColumn = max(lastColumn, sem.Location.End.Column)
	}
	return lastColumn + 1
}

func (s *Semantics) Line(l int) Semantics {
	newS := Semantics{}
	for _, sem := range s.s {
		if l >= sem.Location.Start.Line && l <= sem.Location.End.Line {
			newS.s = append(newS.s, sem)
		}
	}
	return newS
}

func (s *Semantics) Lines() int {
	if len(s.s) == 0 {
		return 0
	}
	var lines int
	for _, sem := range s.s {
		lines = max(lines, sem.Location.End.Line)
	}
	return lines + 1
}

func (s *Semantics) Append(as ...Semantic) {
	s.s = append(s.s, as...)
}

func (s *Semantics) Slice() []Semantic {
	if s == nil {
		return nil
	}
	return s.s
}

func New(semantic ...Semantic) *Semantics {
	return &Semantics{s: semantic}
}
