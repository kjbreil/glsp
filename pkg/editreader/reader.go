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

// var (
// 	ErrNoOpener = errors.New("no opener for current character")
// )
//
// func (f *File) Inside() (CharRange, error) {
// 	// check if current or next matches an enclosure
// 	if f.read.curr == nil {
// 		return CharRange{}, nil
// 	}
// 	closer, ok := f.Enclosers[f.read.curr.c]
// 	if !ok {
// 		closer, ok = f.Enclosers[f.read.curr.Next().c]
// 		if !ok {
// 			return CharRange{}, ErrNoOpener
// 		}
// 		f.Advance()
// 	}
// 	start := f.read.curr
// 	// look for matching closer
// 	var ru rune
// 	var err error
// 	for {
// 		ru, err = f.ReadRune()
// 		if err != nil {
// 			if err == io.EOF {
// 				break
// 			}
// 			panic(err)
// 		}
// if
// 	}
// 	return CharRange{
// 		Start: start,
// 	}, nil
// }
