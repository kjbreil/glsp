package location

import (
	"fmt"
	protocol "github.com/kjbreil/glsp/protocol_3_16"
)

type Range struct {
	Start, End Point
}
type Point struct {
	Line, Column int
}

func (p Point) IsZero() bool {
	return p.Line == 0 && p.Column == 0
}

func (p Point) NewLine() Point {
	return Point{Line: p.Line + 1, Column: 0}
}

func (p Point) NewColumn() Point {
	return Point{Line: p.Line, Column: p.Column + 1}
}

func (p Point) Equals(point Point) bool {
	return p.Line == point.Line && p.Column == point.Column
}

func (p Point) After(start Point) bool {
	return p.Line > start.Line || (p.Line == start.Line && p.Column >= start.Column)
}

func (p Point) Before(end Point) bool {
	return p.Line < end.Line || (p.Line == end.Line && p.Column < end.Column)
}

func (p Point) String() string {
	return fmt.Sprintf("Line: %d Column: %d", p.Line, p.Column)
}

func ZeroPoint() Point {
	return Point{Line: 0, Column: 0}
}

func FullRange() Range {
	return Range{Start: NegativePoint(), End: MaxPoint()}
}
func MaxRange() Range {
	return Range{Start: MaxPoint(), End: MaxPoint()}
}

func NegativePoint() Point {
	return Point{Line: -1, Column: -1}
}

func MaxPoint() Point {
	return Point{Line: 1<<30 - 1, Column: 1<<30 - 1}
}

func (r *Range) Correct() {
	if r.Start.Line < 0 {
		r.Start.Line = 0
	}
	if r.Start.Column < 0 {
		r.Start.Column = 0
	}
	if r.End.Line < 0 {
		r.End.Line = 0
	}
	if r.End.Column < 0 {
		r.End.Column = 0
	}
}

func (r *Range) Minus(i int) *Range {
	var nr Range
	nr.End = r.End
	nr.Start.Line = r.Start.Line
	nr.Start.Column = r.Start.Column - i
	return &nr
}

func (r *Range) Plus(i int) *Range {
	var nr Range
	nr.Start = r.Start
	nr.End.Line = r.End.Line
	nr.End.Column = r.End.Column + i
	return &nr
}

func (r *Range) Contains(op Point) bool {
	if r.Start.Line == op.Line && r.Start.Line == r.End.Line {
		return r.Start.Column <= op.Column && r.End.Column >= op.Column
	}
	if r.Start.Line == op.Line {
		return r.Start.Column <= op.Column
	}

	if r.End.Line == op.Line {
		return r.End.Column >= op.Column
	}

	return r.Start.Line <= op.Line && r.End.Line >= op.Line
}

func (r *Range) Intersects(r2 *Range) bool {
	// if the r start is greater than r2 start and less than r2 end
	if r.Start.After(r2.Start) && r.Start.Before(r2.End) {
		return true
	}
	// if the r2 start is greater than r start and less than r end
	if r2.Start.After(r.Start) && r2.Start.Before(r.End) {
		return true
	}
	// if the r start is equal to r2 start and less than r2 end
	if r.Start.Equals(r2.Start) {
		return true
	}
	if r.End.Equals(r2.End) {
		return true
	}

	return false
}

func (r *Range) ProtocolRange() protocol.Range {
	return protocol.Range{
		Start: protocol.Position{
			Line:      protocol.UInteger(r.Start.Line),
			Character: protocol.UInteger(r.Start.Column),
		},
		End: protocol.Position{
			Line:      protocol.UInteger(r.End.Line),
			Character: protocol.UInteger(r.End.Column),
		},
	}
}

func (r *Range) LessThan(location *Range) bool {
	return r.Start.Line < location.Start.Line ||
		(r.Start.Line == location.Start.Line && r.Start.Column < location.Start.Column)
}

func (r *Range) Within(ol *Range) bool {
	// if start and end line are same check within start and end column
	if r.Start.Line == ol.Start.Line && r.End.Line == ol.End.Line {
		return r.Start.Column >= ol.Start.Column && r.End.Column < ol.End.Column
	}

	// if within the lines return true
	if r.Start.Line > ol.Start.Line && r.End.Line < ol.End.Line {
		return true
	}
	// if on same start line check start position
	if r.Start.Line == ol.Start.Line && r.Start.Column >= ol.Start.Column {
		return true
	}
	// if on same end check end position
	if r.End.Line == ol.End.Line && r.End.Column < ol.End.Column {
		return true
	}

	return false
}

func (r *Range) Equals(or *Range) bool {
	return r.Start.Equals(or.Start) && r.End.Equals(or.End)
}

func (r *Range) Invalid() bool {
	if r.End.Before(r.Start) {
		return true
	}
	return false
}

func (r Range) MakeSingle() *Range {
	return r.Plus(-1).Minus(-1)

}

func (r *Range) String() string {
	return fmt.Sprintf("Range Start:L%d:C%d End:L%d:C%d)", r.Start.Line, r.Start.Column, r.End.Line, r.End.Column)
}
