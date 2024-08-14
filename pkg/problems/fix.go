package problems

import (
	protocol "github.com/kjbreil/glsp/protocol_3_16"
	"github.com/kjbreil/loc-macro/pkg/editreader"
)

type Fix struct {
	Title string
	Range *editreader.Range
	Fixed string
}

func (f *Fix) TextEdit() *protocol.TextEdit {
	return &protocol.TextEdit{
		Range:   f.Range.ProtocolRange(),
		NewText: f.Fixed,
	}
}
