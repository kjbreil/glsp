package editreader

import (
	"bufio"
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"os"
	"strings"
	"testing"
)

func TestNew_fromString(t *testing.T) {

	tests := []struct {
		name    string
		input   string
		want    *File
		wantErr bool
	}{
		{
			name:  "single character",
			input: "a",
			want: &File{
				head: &Char{
					c: 'a',
					n: nil,
				},
			},
		},
		{
			name:  "space",
			input: " ",
			want: &File{
				head: &Char{
					c: ' ',
					n: nil,
				},
			},
		},
		{
			name:  "newline test",
			input: "a\n",
			want: &File{
				head: &Char{
					c: 'a',
					n: &Char{
						c: '\n',
						n: nil,
					},
				},
			},
		},
		{
			name:  "CRLF First",
			input: "\r\n",
			want: &File{
				head: &Char{
					c: '\n',
					n: nil,
				},
			},
		},
		{
			name:  "CRLF",
			input: "a\r\n",
			want: &File{
				head: &Char{
					c: 'a',
					n: &Char{
						c: '\n',
						n: nil,
					},
				},
			},
		},
		{
			name:  "Newline without -1 @ end",
			input: "a\n",
			want: &File{
				head: &Char{
					c: 'a',
					n: &Char{
						c: '\n',
						n: nil,
					},
				},
			},
		},
		{
			name:  "newline test",
			input: "a\nb",
			want: &File{
				head: &Char{
					c: 'a',
					n: &Char{
						c: '\n',
						n: &Char{
							c: 'b',
							n: nil,
						},
					},
				},
			},
		},
		{
			name:  "registered",
			input: "®",
			want: &File{
				head: &Char{
					c: '®',
					n: nil,
				},
			},
		},
		{
			name:  "registered",
			input: "©",
			want: &File{
				head: &Char{
					c: '©',
					n: nil,
				},
			},
		},
		{
			name:  "UTF-8 Read as Windows-1252",
			input: "Â®",
			want: &File{
				head: &Char{
					c: '®',
					n: nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			got, err := New(r)

			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Equals(tt.want) {
				t.Errorf("New() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew_fromString_1252(t *testing.T) {

	tests := []struct {
		name    string
		input   string
		want    *File
		wantErr bool
	}{
		{
			name:  "registered",
			input: "®",
			want: &File{
				head: &Char{
					c: '®',
					n: nil,
				},
			},
		},
		{
			name:  "registered",
			input: "©",
			want: &File{
				head: &Char{
					c: '©',
					n: nil,
				},
			},
		},
	}
	for _, tt := range tests {

		encoder := charmap.Windows1252.NewEncoder()

		t.Run(tt.name, func(t *testing.T) {
			s, err := encoder.String(tt.input)
			if err != nil {
				t.Fatal(err)
			}

			r := strings.NewReader(s)
			got, err := New(r)
			fmt.Println(got.Encoding)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Equals(tt.want) {
				t.Errorf("New() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Benchmark_SeekLargeFile(b *testing.B) {

	//b.StopTimer()
	f, err := os.Open("testdata/romeo_juliet.txt")
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = New(r)
		if err != nil {
			b.Fatal(err)
		}
	}
	file, err := New(r)
	if err != nil {
		b.Fatal(err)
	}
	start := file.Head()
	end := file.Tail()

	//b.StartTimer()

	for n := 0; n < b.N; n++ {
		file.GoTo(start)
		file.GoTo(end)
	}

}
