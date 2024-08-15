package editreader

import (
	"fmt"
	"math"
	"strings"
	"testing"
)

func TestFile_Replace(t *testing.T) {
	type args struct {
		n string
		r Range
	}
	tests := []struct {
		name     string
		args     args
		original string
		want     string
		setCol   int
		setLine  int
	}{
		{
			name:     "replace single character",
			original: "a",
			args: args{
				n: "b",
				r: Range{
					Start: Point{Line: 0, Column: 0},
					End:   Point{Line: 0, Column: 1},
				},
			},
			want: "b",
		},
		{
			name:     "replace single character",
			original: "ac",
			args: args{
				n: "b",
				r: Range{
					Start: Point{Line: 0, Column: 1},
					End:   Point{Line: 0, Column: 1},
				},
			},
			want: "abc",
		},
		{
			name:     "just move column",
			original: "abcdefghijklmnopqrstuvwxyz",
			setCol:   15,
			args: args{
				n: "0",
				r: Range{
					Start: Point{Line: 0, Column: 11},
					End:   Point{Line: 0, Column: 12},
				},
			},
			want: "abcdefghijk0mnopqrstuvwxyz",
		},
		// start is 100 end is 102
		{
			name:     "replace single character",
			original: "abc\nd1f",
			args: args{
				n: "e",
				r: Range{
					Start: Point{Line: 1, Column: 1},
					End:   Point{Line: 1, Column: 2},
				},
			},
			want: "abc\ndef",
		},
		{
			name:     "replace all with empty",
			original: "abcdefghijklmnopqrstuvwxyz\nabcdefghijklmnopqrstuvwxyz",
			args: args{
				n: "",
				r: Range{
					Start: Point{Line: -1, Column: -1},
					End:   Point{Line: math.MaxInt, Column: math.MaxInt},
				},
			},
			want: "",
		},
		{
			name:     "add to nothing",
			original: "",
			args: args{
				n: "abc",
				r: Range{
					Start: Point{Line: 0, Column: 0},
					End:   Point{Line: 0, Column: 0},
				},
			},
			want: "abc",
		},
		{
			name:     "add newline",
			original: "text",
			args: args{
				n: "\n",
				r: Range{
					Start: Point{Line: 0, Column: 4},
					End:   Point{Line: 0, Column: 4},
				},
			},
			want: "text\n",
		},
		{
			name:     "add newline and text",
			original: "text",
			args: args{
				n: "\nabc",
				r: Range{
					Start: Point{Line: 0, Column: 4},
					End:   Point{Line: 0, Column: 4},
				},
			},
			want: "text\nabc",
		},
		{
			name:     "add newline and text",
			original: "beforeafter",
			args: args{
				n: "not\nbut",
				r: Range{
					Start: Point{Line: 0, Column: 6},
					End:   Point{Line: 0, Column: 6},
				},
			},
			want: "beforenot\nbutafter",
		},
		{
			name:     "add newline and text",
			original: "abc\n",
			args: args{
				n: "d",
				r: Range{
					Start: Point{Line: 1, Column: 0},
					End:   Point{Line: 1, Column: 0},
				},
			},
			want: "abc\nd",
		},
		{
			name:     "add newline and text",
			original: "@EXEC(FCT=710);",
			args: args{
				n: "\n@",
				r: Range{
					Start: Point{Line: 0, Column: 15},
					End:   Point{Line: 0, Column: 15},
				},
			},
			want: "@EXEC(FCT=710);\n@",
		},
		{
			name:     "add newline and text",
			original: "@EXEC(FCT=710);\n",
			args: args{
				n: "@",
				r: Range{
					Start: Point{Line: 1, Column: 0},
					End:   Point{Line: 1, Column: 0},
				},
			},
			setLine: 0,
			setCol:  14,
			want:    "@EXEC(FCT=710);\n@",
		},
		{
			name:     "add newline and text",
			original: "@EXEC(FCT=710);\n@",
			args: args{
				n: "E",
				r: Range{
					Start: Point{Line: 1, Column: 1},
					End:   Point{Line: 1, Column: 1},
				},
			},
			setLine: 1,
			want:    "@EXEC(FCT=710);\n@E",
		},
		{
			name:     "Live Test",
			original: "@FMT(2HTML,)",
			args: args{
				n: "\n",
				r: Range{
					Start: Point{
						Line:   0,
						Column: 12,
					},
					End: Point{
						Line:   0,
						Column: 12,
					},
				},
			},
			want: "@FMT(2HTML,)\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := New(strings.NewReader(tt.original))
			if err != nil {
				t.Error(err)
			}
			if tt.setCol != 0 {
				f.edit.gotoCol(tt.setCol)
			}
			if tt.setLine != 0 {
				f.edit.gotoLine(tt.setLine)
			}
			f.Replace(tt.args.n, &tt.args.r)

			if got := f.String(); got != tt.want {
				t.Errorf("File.Replace() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestFile_Replace_check_location(t *testing.T) {
	type args struct {
		n string
		r Range
	}
	tests := []struct {
		name          string
		args          args
		original      string
		want          string
		lastCharPoint Point
	}{
		{
			name:     "replace all with empty",
			original: "abcdefghijklmnopqrstuvwxyz\nabcdefghijklmnopqrstuvwxyz",
			args: args{
				n: "",
				r: Range{
					Start: Point{Line: -1, Column: -1},
					End:   Point{Line: math.MaxInt, Column: math.MaxInt},
				},
			},
			lastCharPoint: Point{0, 0},
			want:          "",
		},
		{
			name:     "add to nothing",
			original: "",
			args: args{
				n: "abc",
				r: Range{
					Start: Point{Line: 0, Column: 0},
					End:   Point{Line: 0, Column: 0},
				},
			},
			lastCharPoint: Point{0, 3},
			want:          "abc",
		},
		{
			name:     "add newline",
			original: "text",
			args: args{
				n: "\n",
				r: Range{
					Start: Point{Line: 0, Column: 4},
					End:   Point{Line: 0, Column: 4},
				},
			},
			lastCharPoint: Point{1, 0},
			want:          "text\n",
		},
		{
			name:     "add newline and text",
			original: "text",
			args: args{
				n: "\nabc",
				r: Range{
					Start: Point{Line: 0, Column: 4},
					End:   Point{Line: 0, Column: 4},
				},
			},
			want:          "text\nabc",
			lastCharPoint: Point{1, 3},
		},
		{
			name:     "add newline and text",
			original: "beforeafter",
			args: args{
				n: "not\nbut",
				r: Range{
					Start: Point{Line: 0, Column: 6},
					End:   Point{Line: 0, Column: 6},
				},
			},
			want:          "beforenot\nbutafter",
			lastCharPoint: Point{1, 8},
		},
		{
			name:     "Live Test",
			original: "a\nb",
			args: args{
				n: "*/",
				r: Range{
					Start: Point{
						Line:   1,
						Column: 1,
					},
					End: Point{
						Line:   1,
						Column: 1,
					},
				},
			},
			want:          "a\nb*/",
			lastCharPoint: Point{1, 3},
		},

		{
			name:     "Live Test",
			original: "a\nb*/",
			args: args{
				n: "/*",
				r: Range{
					Start: Point{
						Line:   1,
						Column: 0,
					},
					End: Point{
						Line:   1,
						Column: 0,
					},
				},
			},
			want:          "a\n/*b*/",
			lastCharPoint: Point{1, 5},
		},
		{
			name:     "Live Test",
			original: "a\n/*b*/",
			args: args{
				n: "*/",
				r: Range{
					Start: Point{
						Line:   0,
						Column: 1,
					},
					End: Point{
						Line:   0,
						Column: 1,
					},
				},
			},
			want:          "a*/\n/*b*/",
			lastCharPoint: Point{1, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := New(strings.NewReader(tt.original))
			if err != nil {
				t.Error(err)
			}

			f.Replace(tt.args.n, &tt.args.r)

			if got := f.String(); got != tt.want {
				t.Errorf("File.Replace() = %v, want %v", got, tt.want)
			}

			if !tt.lastCharPoint.Equals(f.Char().last().point) {
				t.Errorf("File.Replace() did not update lastCharPoint correctly. got %v, want %v", f.Char().last().point, tt.lastCharPoint)
			}

		})
	}
}

func TestFile_ReplaceSlice(t *testing.T) {
	type edit struct {
		n string
		r Range
	}
	type tests struct {
		name     string
		original string
		edits    []edit
		want     string
	}

	ts := []tests{
		{
			name: "",
			want: "@EXEC()\n",
			edits: []edit{
				{
					n: "@",
					r: Range{Start: Point{Line: 0, Column: 0}, End: Point{Line: 0, Column: 0}},
				},
				{
					n: "E",
					r: Range{Start: Point{Line: 0, Column: 1}, End: Point{Line: 0, Column: 1}},
				},

				{
					n: "X",
					r: Range{Start: Point{Line: 0, Column: 2}, End: Point{Line: 0, Column: 2}},
				},

				{
					n: "E",
					r: Range{Start: Point{Line: 0, Column: 3}, End: Point{Line: 0, Column: 3}},
				},
				{
					n: "EXEC()",
					r: Range{Start: Point{Line: 0, Column: 1}, End: Point{Line: 0, Column: 4}},
				},
				{
					n: "\n",
					r: Range{Start: Point{Line: 0, Column: 7}, End: Point{Line: 0, Column: 7}},
				},
			},
		},
		{
			name:     "",
			original: "a\nb",
			want:     "a\n/*b*/",
			edits: []edit{
				{
					n: "*/",
					r: Range{Start: Point{Line: 1, Column: 1}, End: Point{Line: 1, Column: 1}},
				},

				{
					n: "/*",
					r: Range{Start: Point{Line: 1, Column: 0}, End: Point{Line: 1, Column: 0}},
				},
			},
		},
		{
			name:     "undo simple",
			original: "a\nbbb",
			want:     "a\nbbb",
			edits: []edit{
				{
					n: "a",
					r: Range{Start: Point{Line: 1, Column: 2}, End: Point{Line: 1, Column: 2}},
				},
				{
					n: "s",
					r: Range{Start: Point{Line: 1, Column: 3}, End: Point{Line: 1, Column: 3}},
				},
				{
					n: "d",
					r: Range{Start: Point{Line: 1, Column: 4}, End: Point{Line: 1, Column: 4}},
				},
				{
					n: "",
					r: Range{Start: Point{Line: 1, Column: 2}, End: Point{Line: 1, Column: 5}},
				},
			},
		},
		{
			name:     "comments add/edits not right",
			original: "/*a*/\n/*b*/",
			want:     "/*a\na*/\n/*b*/",
			edits: []edit{
				// /*a*/\n/*b*/
				{
					n: "",
					r: Range{Start: Point{Line: 0, Column: 3}, End: Point{Line: 0, Column: 5}},
				},
				// /*a\n/*b*/
				{
					n: "",
					r: Range{Start: Point{Line: 0, Column: 0}, End: Point{Line: 0, Column: 2}},
				},
				// a\n/*b*/
				{
					n: "*/",
					r: Range{Start: Point{Line: 0, Column: 1}, End: Point{Line: 0, Column: 1}},
				},
				// a\n/*b*/
				{
					n: "/*",
					r: Range{Start: Point{Line: 0, Column: 0}, End: Point{Line: 0, Column: 0}},
				},
				// /*/*/*a\n/*b*/
				{
					n: "a\n",
					r: Range{Start: Point{Line: 0, Column: 2}, End: Point{Line: 0, Column: 2}},
				},
			},
		},
		{
			name:     "CRLF",
			original: "a\nb",
			want:     "a\nb",
			edits: []edit{
				{
					n: "",
					r: Range{Start: Point{Line: 0, Column: 0}, End: Point{Line: 1, Column: 1}},
				},
				{
					n: "a\r\nb",
					r: Range{Start: Point{Line: 0, Column: 0}, End: Point{Line: 0, Column: 0}},
				},
				{
					n: "",
					r: Range{Start: Point{Line: 0, Column: 0}, End: Point{Line: 1, Column: 1}},
				},
				{
					n: "a\r\nb",
					r: Range{Start: Point{Line: 0, Column: 0}, End: Point{Line: 0, Column: 0}},
				},
			},
		},
	}

	for _, et := range ts {
		t.Run(et.name, func(t *testing.T) {

			f, err := New(strings.NewReader(et.original))
			if err != nil {
				t.Error(err)
			}
			fmt.Printf("OR: %s\n", escape(et.original))
			for i, e := range et.edits {
				f.Replace(e.n, &e.r)
				fmt.Printf("%02d: %s\n", i, escape(f.String()))
				// f.PrintTestArr(e.n, &e.r)

			}
			if got := f.String(); got != et.want {
				t.Errorf("File.ReplaceSlice() \ngot  = %v\nwant = %v", escape(got), escape(et.want))
			}
		})
	}
}
