package problems

import (
	"fmt"
	"github.com/kjbreil/glsp/internal/helpers"
	protocol "github.com/kjbreil/glsp/protocol_3_16"
	"github.com/kjbreil/loc-macro/pkg/editreader"
	"github.com/pkg/errors"
)

type Problems struct {
	p []Problem
}

var (
	InfoPossibleMacroFound = errors.New("possible macro found")
)

func (p *Problems) AddErr(err error, loc editreader.CharRange, fix *Fix) {
	p.p = append(p.p, Problem{
		Level:     ProblemLevelError,
		err:       err,
		Location:  loc.Range(),
		charRange: loc,
		Fix:       fix,
	})
}
func (p *Problems) AddWarning(err error, loc editreader.CharRange, fix *Fix) {
	p.p = append(p.p, Problem{
		Level:     ProblemLevelWarning,
		err:       err,
		charRange: loc,
		Location:  loc.Range(),
		Fix:       fix,
	})
}

func (p *Problems) AddInfo(err error, loc editreader.CharRange) {
	p.p = append(p.p, Problem{
		Level:     ProblemLevelInfo,
		err:       err,
		charRange: loc,
		Location:  loc.Range(),
	})
}

func (p *Problems) AddPossible(name string, loc *editreader.Range) {
	p.p = append(p.p, Problem{
		Level:    ProblemLevelInfo,
		err:      fmt.Errorf("%w: %s", InfoPossibleMacroFound, name),
		Location: loc,
	})
}

func (p *Problems) Errors(name string) []error {
	var errs []error
	for _, pr := range p.p {
		if pr.Level == ProblemLevelError {
			errs = append(errs, errors.Wrap(pr.err, name))
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errs
}

func (p *Problems) Add(ap *Problems) {
	if ap == nil {
		return
	}
	p.p = append(p.p, ap.p...)
}

func (p *Problems) Slice() []Problem {
	return p.p
}

func (p *Problems) Append(problems *Problems) {
	if problems == nil {
		return
	}
	p.p = append(p.p, problems.p...)
}

func (p *Problems) Intersects(loc *editreader.Range) *Problems {
	np := New()
	for _, ip := range p.p {
		if ip.Location.Intersects(loc) {
			np.p = append(np.p, ip)
		}
	}
	return np
}

func (p *Problems) Len() int {
	return len(p.p)
}

func (p *Problems) ProtocolDiagnostics(maxLevel ProblemLevel) []protocol.Diagnostic {
	diagnostics := []protocol.Diagnostic{}
	for _, err := range p.Slice() {
		if err.Level <= maxLevel {
			diagnostics = append(diagnostics, protocol.Diagnostic{
				Range:    err.Location.ProtocolRange(),
				Severity: problemLevelToSeverity(err.Level),
				Message:  err.Error().Error(),
			})
		}
	}

	return diagnostics
}

func problemLevelToSeverity(level ProblemLevel) *protocol.DiagnosticSeverity {
	switch level {
	case ProblemLevelError:
		return helpers.Ptr(protocol.DiagnosticSeverityError)
	case ProblemLevelWarning:
		return helpers.Ptr(protocol.DiagnosticSeverityWarning)
	case ProblemLevelInfo:
		return helpers.Ptr(protocol.DiagnosticSeverityInformation)
	case ProblemLevelHint:
		return helpers.Ptr(protocol.DiagnosticSeverityHint)
	default:
		return nil
	}
}

func NewErr(err error, loc *editreader.Range) *Problems {
	return &Problems{
		p: []Problem{
			{
				Level:    ProblemLevelError,
				err:      err,
				Location: loc,
			},
		},
	}
}

func New() *Problems {
	return &Problems{}
}
