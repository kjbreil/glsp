package editreader

import (
	"io"
	"unicode/utf8"
)

// read implements io.Reader interface
func (f *File) Read(p []byte) (n int, err error) {
	for n < len(p) {
		if f.read.curr == nil {
			f.Reset()
			return
		}
		if f.read.curr.IsEmpty() {
			if f.read.curr.nextIsEmpty() {
				err = io.EOF
				return
			}
			f.read.forward()
			continue
		}
		n += utf8.EncodeRune(p[n:], f.read.curr.c)
		f.read.forward()
	}
	return
}
func (f *File) ReadRune() (rune, int, error) {
	if f.read.curr == nil {
		f.Reset()
	}

	f.read.forward()

	if f.read.curr == nil || (f.read.curr.IsEmpty() && f.read.curr.nextIsEmpty()) {
		return 0, 0, io.EOF
	}
	return f.read.curr.c, 0, nil
}

func (f *File) ReadRuneUpper() (rune, error) {
	c, _, err := f.ReadRune()
	if 'a' <= c && c <= 'z' {
		c -= 'a' - 'A'
	}
	return c, err
}

func (f *File) Unread() {
	f.read.reverse()
}

func (f *File) Advance() *Char {
	if f.read.curr == nil {
		f.Reset()
	}
	f.read.forward()
	return f.read.curr
}

func (f *File) Peek() *Char {
	if f.read.curr.Next() != nil {
		return f.read.curr.Next()
	}
	return f.read.curr
}

func (f *File) Char() *Char {
	return f.read.curr
}

func (f *File) ReadUntilRune(ru rune) (*CharRange, error) {
	f.m.Lock()
	defer f.m.Unlock()
	cr := &CharRange{
		Start: f.read.curr,
		End:   nil,
	}
	for cr.Start.c == -1 {
		cr.Start = f.read.advance()
	}
	for {
		cr.End = f.read.advance()
		if cr.End.c == -1 {
			cr.End = cr.End.p
			return cr, io.EOF
		}

		if cr.End.c == ru {
			break
		}
	}

	return cr, nil
}

// ReadUntilString reads until a given string is found in the file.
// Does not work if the string to find contains a character that is not single width
func (f *File) ReadUntilString(s string) (*CharRange, error) {
	f.m.Lock()
	defer f.m.Unlock()
	cr := &CharRange{
		Start: f.read.curr,
		End:   nil,
	}
	for cr.Start.c == -1 {
		cr.Start = f.read.advance()
	}

	checkIndex := 0
	for {
		cr.End = f.read.advance()
		if cr.End.c == -1 {
			return nil, io.EOF
		}

		if cr.End.Rune() == rune(s[checkIndex]) {
			if checkIndex == len(s)-1 {
				break
			}
			checkIndex++
			continue
		}
		checkIndex = 0
	}

	return cr, nil
}

func (f *File) ReadLine() (*CharRange, error) {

	return f.ReadUntilRune('\n')
}
