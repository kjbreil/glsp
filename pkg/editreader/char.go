package editreader

import (
	"fmt"
	"github.com/kjbreil/glsp/pkg/location"
	"strings"
	"unicode"
)

var utf8CheckChars map[rune]rune

func init() {
	utf8CheckChars = map[rune]rune{
		174: 194,
		169: 194,
	}
}

type Char struct {
	p     *Char
	c     rune
	point location.Point
	n     *Char
}

func (c *Char) Debug() string {
	return fmt.Sprintf("Char %c, line: %d, column: %d", c.c, c.point.Line, c.point.Column)
}

// String returns a string representation of the linked list of characters starting from the current character.
// It iterates through the linked list, appending each character to a strings.Builder, and then returns the resulting string.
// If the current character is empty (nil or -1), the function immediately returns an empty string.
func (c *Char) String() string {
	sb := strings.Builder{}
	for ; !c.IsEmpty(); c = c.n {
		sb.WriteRune(c.c)
	}
	return sb.String()
}

func CharRangeSingle(c *Char) CharRange {
	return CharRange{Start: c, End: c}
}

// Lines returns the number of newlines in the linked list of characters starting from the current character.
// It recursively traverses the linked list, incrementing a counter each time it encounters a newline character.
// If the current character is nil, the function immediately returns 0.
func (c *Char) Lines() int {
	if c == nil {
		return 0
	}
	if c.newLine() {
		return 1 + c.n.Lines()
	}
	return c.n.Lines()
}

func (c *Char) IsEmpty() bool {
	return c == nil || c.c == -1
}

// Next returns the next character in the linked list.
// If the next character exists, it is returned.
// If the next character does not exist, the current character is returned.
func (c *Char) Next() *Char {
	if c.n != nil {
		return c.n
	}
	return c
}

// Previous returns the previous character in the linked list.
// If the previous character exists, it is returned.
// If the previous character does not exist, the current character is returned.
func (c *Char) Previous() *Char {
	if c.p != nil {
		c = c.p
	}
	return c
}

// ContainsSpecialChars checks if the linked list of characters starting from the current character contains any special characters.
// Special characters in this context are '®' and '©'.
//
// The function iterates through the linked list, checking each character.
// If it encounters a special character, it immediately returns true.
// If it reaches the end of the linked list without finding any special characters, it returns false.
//
// Parameters:
// - c: A pointer to the current character in the linked list.
//
// Returns:
// - bool: True if the linked list contains a special character, false otherwise.
func (c *Char) ContainsSpecialChars() bool {
	for ; c.n != nil; c = c.n {
		switch c.c {
		case '®', '©':
			return true
		}
	}
	return false
}

func (c *Char) IsAlphaNumeric() bool {
	return unicode.IsLetter(c.c) || unicode.IsDigit(c.c)
}

// NextNewline returns the next newline character in the linked list of characters starting from the current character.
// If a newline character is found, it is returned.
// If no newline character is found, the last character in the linked list is returned.
//
// Parameters:
// - c: A pointer to the current character in the linked list.
//
// Returns:
// - *Char: A pointer to the next newline character in the linked list, or the last character if no newline is found.
func (c *Char) NextNewline() *Char {
	next := c.until('\n')
	if next == nil {
		return c.last()
	}
	return next
}

// equals compares the current character and its subsequent characters with another character and its subsequent characters.
// It returns true if all characters in both linked lists are equal, and false otherwise.
//
// Parameters:
// - od: A pointer to the first character in the other linked list to compare with.
//
// Returns:
// - bool: True if all characters in both linked lists are equal, false otherwise.
func (c *Char) equals(od *Char) bool {
	if c.c == -1 && c.n != nil {
		c = c.n
	}
	if c.c != od.c {
		return false
	}
	if c.nextIsEmpty() && od.nextIsEmpty() {
		return true
	}
	if c.nextIsEmpty() || od.nextIsEmpty() {
		return false
	}
	return c.n.equals(od.n)
}

func (c *Char) newLine() bool {
	return c.c == '\n'
}

func (c *Char) nextIsEmpty() bool {
	return c.n == nil || c.n.c == -1
}

func (c *Char) mostLikely1252() bool {
	for {
		if c == nil || c.c == 0 {
			return false
		}

		if c.c == 65533 {
			return true
		}

		if pr, ok := utf8CheckChars[c.c]; ok {
			if c.previousEquals(pr) {
				return true
			}
		}

		return c.n.mostLikely1252()
	}
}

func (c *Char) mostLikelyUTF8() bool {

	for {
		if c == nil || c.c == 0 {
			return false
		}

		if pr, ok := utf8CheckChars[c.c]; ok {
			if c.previousEquals(pr) {
				return true
			}
		}

		return c.n.mostLikely1252()
	}
}

func (c *Char) previousEquals(ru rune) bool {
	if c.p == nil {
		return false
	}
	if c.p.c == ru {
		return true
	}
	return false
}

func (c *Char) last() *Char {
	for c.n != nil {
		c = c.n
	}
	return c
}

func (c *Char) delete() {

	if c.n != nil {
		c.n.p = c.p
	}
	if c.p != nil {
		c.p.n = c.n
	}
}

func (c *Char) setNext(rc *Char) {
	c.n = rc
	if rc != nil {
		rc.p = c
	}
}

func (c *Char) insert(nc *Char) {
	if nc == nil {
		return
	}

	// c's next previous needs to be set to rc.last()
	c.n.p = nc.last()
	// rc's lasts next needs to be set to current c
	nc.last().n = c.n
	// current next needs to be set to rc
	c.n = nc
	// rc previous needs to be set to current
	nc.p = c
}

func (c *Char) Rune() rune {

	return c.c
}

func (c *Char) Point() location.Point {
	return c.point
}

func (c *Char) Escaped() string {
	return escapeRune(c.c)
}

func (c *Char) setLoc() {

	for ; c.n != nil; c = c.n {

		if c.p != nil {
			if c.p.newLine() {
				c.point = c.p.point.NewLine()
			} else {
				c.point = c.p.point.NewColumn()
			}
		} else {
			c.point = location.Point{
				Line:   0,
				Column: -1,
			}
		}
	}

	if c.p != nil {
		if c.p.newLine() {
			c.point = c.p.point.NewLine()
		} else {
			c.point = c.p.point.NewColumn()
		}
	} else {
		c.point = location.Point{
			Line:   0,
			Column: -1,
		}
	}

	if c.n != nil {
		c.n.setLoc()
	}
}

func (c *Char) first() *Char {
	for c.p != nil {
		c = c.p
	}
	return c
}

func (c *Char) until(ru rune) *Char {
	for {
		if c == nil {
			return nil
		}
		if c.c == ru {
			return c
		}
		c = c.n
	}
}
