package editreader

import "github.com/kjbreil/glsp/pkg/location"

// TODO: tracker.curr should never be nil - prevent that

type tracker struct {
	curr *Char
}

// gotoLine returns the Char at the start of a line
func (t *tracker) gotoLine(i int) {
	// go down the lines until we reach the desired line
	for t.curr.point.Line < i {
		if t.curr.n == nil {
			return
		}
		t.forward()
	}
	// go up the lines until we reach the desired line
	for t.curr.point.Line > i {
		if t.curr.p == nil {
			return
		}
		t.reverse()
	}
}

// gotoCol returns the character at the column of the current line
func (t *tracker) gotoCol(col int) {
	// go to the desired column
	for t.curr.point.Column < col {
		if t.curr == nil || t.curr.n == nil {
			return
		}
		t.forward()
	}
	// go to the desired column
	for t.curr.point.Column > col {
		if t.curr == nil || t.curr.p == nil {
			return
		}
		t.reverse()
	}
}

func (t *tracker) forward() {
	if t.curr.n == nil {
		return
	}
	t.curr = t.curr.n
}

func (t *tracker) advance() *Char {
	if t.curr.n == nil {
		return nil
	}
	t.curr = t.curr.n
	return t.curr
}

func (t *tracker) reverse() {
	if t.curr.p == nil {
		return
	}
	t.curr = t.curr.p
}

func (t *tracker) next() *Char {
	if t.curr.n != nil {
		return t.curr.n
	}
	return t.curr
}
func (t *tracker) prev() *Char {
	if t.curr.p != nil {
		return t.curr.p
	}
	return t.curr
}

func (t *tracker) gotoPoint(p location.Point) {
	if t.curr.point.Line != p.Line {
		t.gotoLine(p.Line)
	}
	t.gotoCol(p.Column)
}

func (t *tracker) reset() {
	t.curr = t.curr.first()
}

func (t *tracker) GoTo(c *Char) {
	t.curr = c
	// if c == nil {
	// 	t.curr = t.curr.first()
	// 	return
	// }
	// from the current go forward looking for the c
	// if t.curr.point.Before(c.point) {
	// 	for ; t.curr != c; t.forward() {
	// 		if t.curr.n == nil {
	// 			break
	// 		}
	// 	}
	// }
	//
	// // from the current go backward looking for the c
	// // if we do not find the char then the curr should be at start
	// for ; t.curr != c; t.reverse() {
	// 	if t.curr.p == nil {
	// 		break
	// 	}
	// }

}
