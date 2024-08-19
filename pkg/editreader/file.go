package editreader

import (
	"bufio"
	"fmt"
	"github.com/kjbreil/glsp/pkg/location"
	"golang.org/x/text/encoding/charmap"
	"io"
	"strings"
	"sync"
)

type File struct {
	head *Char
	path string

	edit     tracker
	read     tracker
	m        sync.Mutex
	Encoding Encoding
}

type Encoding int

const (
	Unknown Encoding = iota
	Windows1252
	UTF8
)

func New(r io.Reader) (*File, error) {
	c := &Char{
		c: -1,
	}

	f := &File{
		head: c,
		m:    sync.Mutex{},
	}

	// Read each line and create a new char for each rune
	// use bufio to buffer the reading input
	decoder := charmap.Windows1252.NewDecoder()

	br := bufio.NewReader(r)
	var err error
	for {
		line, readErr := br.ReadString('\n')
		if readErr != nil {
			if readErr != io.EOF {
				return nil, err
			}
		}
		if f.Encoding == Windows1252 {
			line, err = decoder.String(line)
			if err != nil {
				return nil, err
			}
		}

		nc := readString(line)
		if nc == nil {
			break
		}
		if nc.mostLikelyUTF8() {
			f.Encoding = UTF8
			encoder := charmap.Windows1252.NewEncoder()
			newLine, err := encoder.String(line)
			if err == nil {
				nc = readString(newLine)
			}
		}

		if nc.mostLikely1252() {
			f.Encoding = Windows1252
			newLine, err := decoder.String(line)
			if err == nil {
				nc = readString(newLine)
			}
		}

		c.n = nc
		nc.p = c
		c = nc

		if readErr == io.EOF {
			break
		}
		// go to the second to last CHAR and then delete the last element since its not the end of the file
		c = c.last().p
		c.n.delete()
	}

	f.edit.curr = f.head
	f.read.curr = f.head
	if f.head.c != -1 || f.head.last().c != -1 {
		return nil, fmt.Errorf("first and last characters are not -1")
	}
	f.Reset()

	f.head.setLoc()

	return f, nil
}

// Equals confirms the data is the same in two files
func (f *File) Equals(of *File) bool {
	f.m.Lock()
	defer f.m.Unlock()
	of.m.Lock()
	defer of.m.Unlock()
	return f.head.equals(of.head)
}

func (f *File) String() string {
	f.m.Lock()
	defer f.m.Unlock()
	var sb strings.Builder
	c := f.head
	if c.IsEmpty() && !c.nextIsEmpty() {
		c = c.n
	}
	for c != nil && c.c != -1 {
		sb.WriteRune(c.c)
		c = c.n
	}
	return sb.String()
}

// func (f *File) setCol(c int) {
//	f.col = c
// }

func (f *File) Reset() {
	f.m.Lock()
	defer f.m.Unlock()
	f.read.curr = f.head
	f.read.reset()
}

func (f *File) SetPath(path string) {
	f.path = path
}

func (f *File) PrintTest(text string, editRange *location.Range) {
	curr := f.String()
	// 		{
	//			name:     "replace single character",
	//			original: "a",
	//			args: args{
	//				n: "b",
	//				r: Range{
	//					Start: Point{Line: 0, Column: 0},
	//					End:   Point{Line: 0, Column: 1},
	//				},
	//			},
	//			want: "b",
	//		},

	sb := strings.Builder{}
	sb.WriteString("\n")
	sb.WriteString("{\n")
	sb.WriteString("\tname:     \"Live Test\",\n")
	sb.WriteString("\toriginal: \"" + curr + "\",\n")
	sb.WriteString("\targs: args{\n")
	sb.WriteString("\t\tn: \"" + text + "\",\n")
	sb.WriteString("\t\tr: Range{\n")
	sb.WriteString("\t\t\tStart: Point{\n")
	sb.WriteString("\t\t\t\tLine: " + fmt.Sprint(editRange.Start.Line) + ",\n")
	sb.WriteString("\t\t\t\tCol: " + fmt.Sprint(editRange.Start.Column) + ",\n")
	sb.WriteString("\t\t\t},\n")
	sb.WriteString("\t\t\tEnd: Point{\n")
	sb.WriteString("\t\t\t\tLine: " + fmt.Sprint(editRange.End.Line) + ",\n")
	sb.WriteString("\t\t\t\tCol: " + fmt.Sprint(editRange.End.Column) + ",\n")
	sb.WriteString("\t\t\t},\n")
	sb.WriteString("\t\t},\n")
	sb.WriteString("\t},\n")
	sb.WriteString("\twant: \"" + text + "\",\n")
	sb.WriteString("},\n")

	fmt.Println(sb.String())
}

func escapeRune(c rune) string {
	switch c {
	case '\n':
		return "\\n"
	case '\r':
		return "\\r"
	case '\t':
		return "\\t"
	case '\v':
		return "\\v"
	case '\f':
		return "\\f"
	case '\a':
		return "\\a"
	case '\b':
		return "\\b"
	case '\033':
		return "\\e"
	default:
		if c < 32 || c > 126 {
			return fmt.Sprintf("\\x%02x", c)
		}
		return string(c)
	}
}

func escape(o string) string {
	sb := strings.Builder{}
	for _, c := range o {
		sb.WriteString(escapeRune(c))
	}
	return sb.String()
}

func (f *File) PrintTestArr(text string, editRange *location.Range) {
	curr := f.String()

	// 				{
	//					n: "*/",
	//					r: Range{Start: Point{Line: 1, Column: 1}, End: Point{Line: 1, Column: 1}},
	//				},

	sb := strings.Builder{}
	sb.WriteString("// ")
	sb.WriteString(escape(curr))
	sb.WriteString("\n{\n")
	sb.WriteString("\tn: \"")
	sb.WriteString(escape(text))
	sb.WriteString("\",\n\tr: Range{Start: Point{Line: ")
	sb.WriteString(fmt.Sprint(editRange.Start.Line))
	sb.WriteString(",Column: ")
	sb.WriteString(fmt.Sprint(editRange.Start.Column))
	sb.WriteString("}, End: Point{Line: ")
	sb.WriteString(fmt.Sprint(editRange.End.Line))
	sb.WriteString(", Column: ")
	sb.WriteString(fmt.Sprint(editRange.End.Column))
	sb.WriteString("}},\n},")

	fmt.Println(sb.String())
}

func (f *File) Head() *Char {
	f.m.Lock()
	defer f.m.Unlock()
	return f.head
}

func (f *File) Tail() *Char {
	f.m.Lock()
	defer f.m.Unlock()
	c := f.head
	if c == nil {
		return nil
	}
	for c.n != nil {
		c = c.n
	}
	return c
}

func (f *File) GoTo(c *Char) {
	f.read.GoTo(c)
}
