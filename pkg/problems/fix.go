package problems

import (
	"github.com/kjbreil/glsp/pkg/location"
	protocol "github.com/kjbreil/glsp/protocol_3_16"
)

type Fix struct {
	Title string
	Range *location.Range
	Fixed string
}

func (f *Fix) TextEdit() *protocol.TextEdit {
	return &protocol.TextEdit{
		Range:   f.Range.ProtocolRange(),
		NewText: f.Fixed,
	}
}
