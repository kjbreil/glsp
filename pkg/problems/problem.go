package problems

import (
	"errors"
	"fmt"
	"github.com/kjbreil/loc-macro/pkg/editreader"
	"strings"
)

type Problem struct {
	Level     ProblemLevel
	err       error
	Location  *editreader.Range
	filename  string
	Fix       *Fix
	charRange editreader.CharRange
}

func (p *Problem) Error() error {
	return p.err
}

func (p *Problem) ErrorWithLocation() error {
	return fmt.Errorf("%w: %s", p.err, p.Location.String())
}

func (p *Problem) PossibleMacro() (string, string, bool) {
	if errors.Is(p.err, InfoPossibleMacroFound) {
		name := extractMacroName(p.err.Error())
		return name, p.Location.String(), true
	}
	return "", "", false
}

func (p *Problem) CharRange() editreader.CharRange {
	return p.charRange
}

type ProblemLevel int

const (
	ProblemLevelError ProblemLevel = iota
	ProblemLevelWarning
	ProblemLevelInfo
	ProblemLevelHint
	ProblemLevelNone
)

func extractMacroName(str string) string {
	start := 0
	for i, c := range str {
		if c == ':' {
			start = i + 2
			break
		}
	}
	return strings.ToUpper(str[start:])
}
