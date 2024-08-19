package editreader

import (
	"github.com/kjbreil/glsp/pkg/location"
	protocol "github.com/kjbreil/glsp/protocol_3_16"
	"strings"
	"unicode"
)

type CharRange struct {
	Start *Char
	End   *Char
}

func (r CharRange) String() string {
	if r.Start == nil || r.End == nil {
		return ""
	}
	sb := strings.Builder{}
	curr := r.Start
	sb.WriteRune(curr.c)
	for curr != r.End {
		curr = curr.n
		sb.WriteRune(curr.c)
	}
	return sb.String()

}

func (r CharRange) LowerString() string {
	if r.Start == nil || r.End == nil {
		return ""
	}
	sb := strings.Builder{}
	curr := r.Start
	sb.WriteRune(curr.c)
	var currRune rune
	for curr != r.End {
		curr = curr.n
		currRune = curr.c
		if 'A' <= currRune && currRune <= 'Z' {
			currRune += 'a' - 'A'
		}
		sb.WriteRune(currRune)
	}
	return sb.String()
}

func (r CharRange) ProtocolRange() *protocol.Range {

	return &protocol.Range{
		Start: protocol.Position{
			Line:      protocol.UInteger(r.Start.point.Line),
			Character: protocol.UInteger(r.Start.point.Column),
		},
		End: protocol.Position{
			Line:      protocol.UInteger(r.End.point.Line),
			Character: protocol.UInteger(r.End.point.Column + 1),
		},
	}
}

func (r CharRange) Range() *location.Range {
	if r.Start == nil || r.End == nil {
		return nil
	}
	return &location.Range{
		Start: location.Point{
			Line:   r.Start.point.Line,
			Column: r.Start.point.Column,
		},
		End: location.Point{
			Line:   r.End.point.Line,
			Column: r.End.point.Column + 1,
		},
	}
}

// ShiftLeft returns a new CharRange that "adds" the i amount of characters to the left of the start to the CharRange
func (r CharRange) ShiftLeft(i int) CharRange {
	for ii := 0; ii < i; ii++ {
		if r.Start.p != nil {
			r.Start = r.Start.p
		}
	}
	return r
}

// ShiftRight returns a new CharRange that "adds" the i amount of characters to the End of the CharRange
func (r CharRange) ShiftRight(i int) CharRange {
	for ii := 0; ii < i; ii++ {
		if r.End.n != nil {
			r.End = r.End.n
		}
	}
	return r
}

// Until returns a new CharRange that ends at the first occurrence of the specified rune
func (r CharRange) Until(ru rune) CharRange {
	c := r.Start
	for {
		if c == nil {
			return CharRange{Start: r.Start, End: r.End}
		}
		if c == r.End {
			return CharRange{Start: r.Start, End: r.End}
		}
		if c.n != nil && c.n.c == ru {
			return CharRange{Start: r.Start, End: c}
		}
		c = c.n
	}
}

// At returns a new CharRange that starts and ends at the first occurrence of the specified rune
func (r CharRange) At(ru rune) CharRange {
	c := r.Start
	for {
		if c == nil {
			return r
		}
		if c == r.End {
			return r
		}
		if c.c == ru {
			return CharRange{Start: c, End: c}
		}
		c = c.n
	}
}

// After returns a new CharRange that starts at the first occurrence of the specified rune
func (r CharRange) After(ru rune) CharRange {
	c := r.Start
	start := r.Start
	for {
		if c == nil {
			return CharRange{
				Start: start,
				End:   r.End,
			}
		}
		if c.c == ru {
			start = c
		}
		if c == r.End {
			return CharRange{
				Start: start,
				End:   r.End,
			}
		}
		c = c.n
	}
}

func (r *CharRange) Pad(ru rune) {
	if r.Start == nil || r.End == nil {
		return

	}

	lPad := &Char{
		c: ru,
	}

	rPad := &Char{
		c: ru,
	}

	r.End.insert(rPad)
	r.Start.p.insert(lPad)
	r.Start.p.setLoc()
	r.End = rPad
	r.Start = lPad
}

func (r CharRange) Contains(ch *Char) bool {
	if r.Start == nil || r.End == nil {
		return false
	}

	if r.Start.point.Before(ch.point) && r.End.point.After(ch.point) {
		return true
	}

	return false
}

// Trim returns a new CharRange that trims leading and trailing whitespace
func (r CharRange) Trim() CharRange {
	return r.trimLeft().trimRight()
}

func (r CharRange) trimLeft() CharRange {
	ncr := CharRange{
		End: r.End,
	}
	curr := r.Start
	for {
		if curr == nil {
			break
		}
		if !unicode.IsSpace(curr.c) {
			ncr.Start = curr
			break
		}
		if curr == r.End {
			break
		}
		curr = curr.n
	}
	return ncr
}

func (r CharRange) trimRight() CharRange {
	ncr := CharRange{
		Start: r.Start,
	}
	curr := r.End
	for {
		if curr == nil {
			break
		}
		if !unicode.IsSpace(curr.c) {
			ncr.End = curr
			break
		}
		if curr == r.Start {
			break
		}
		curr = curr.p
	}
	return ncr
}

func (r CharRange) Valid() bool {
	if r.Start == nil || r.End == nil {
		return false
	}
	if r.End.point.Before(r.Start.point) {
		return false
	}

	return true
}

func (r *CharRange) Iter(yield func(c *Char) bool) {
	c := r.Start
	for {
		if c == nil {
			return
		}
		if !yield(c) {
			return
		}
		c = c.n
	}
}

func (r *CharRange) Empty() bool {
	return false
}
