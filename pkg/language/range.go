package language

import (
	protocol "github.com/kjbreil/glsp/protocol_3_16"
)

type Range struct {
	Start, End Point
}
type Point struct {
	Line, Column int
}

func ProtocolRange(r *protocol.Range) *Range {

	if r == nil {
		return nil
	}
	return &Range{
		Start: Point{Line: int(r.Start.Line), Column: int(r.Start.Character)},
		End:   Point{Line: int(r.End.Line), Column: int(r.End.Character)},
	}
}
func ProtocolPositionPoint(r protocol.TextDocumentPositionParams) Point {
	return Point{Line: int(r.Position.Line), Column: int(r.Position.Character)}
}
