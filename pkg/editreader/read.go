package editreader

import (
	"strings"
	"unicode"
)

func readString(s string) *Char {
	c := &Char{
		c: -1,
	}
	head := c

	for _, ru := range s {
		// assign current char's rune
		c.c = ru
		// make a new char
		nc := Char{
			p: c,
			c: -1,
		}
		// if the character is a newline check for carriage return and skip it
		if c.c == '\n' {
			if c.p != nil && c.p.c == '\r' {
				// handle CRLF for character
				if c.p.p == nil {
					head = c
				} else {
					c.p.p.n = c
					c.p = c.p.p
				}
			}
		}
		// assign the n to the n new one
		c.n = &nc
		c = &nc
	}

	return head
}

func removeLoopers(s string) string {
	rtn := strings.Builder{}
	for _, ru := range s {
		if unicode.IsLetter(ru) || unicode.IsDigit(ru) || ru == ' ' {
			rtn.WriteRune(ru)
		}
	}
	return rtn.String()
}
