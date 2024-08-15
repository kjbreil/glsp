package editreader

import (
	"golang.org/x/text/encoding/charmap"
	"testing"
)

func Test_readString_1252(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  *Char
	}{
		{
			name:  "single character",
			input: "a",
			want: &Char{
				c: 'a',
				n: nil,
			},
		},
		{
			name:  "registered",
			input: "®",
			want: &Char{
				c: '®',
				n: nil,
			},
		},
		{
			name:  "registered",
			input: "©",
			want: &Char{
				c: '©',
				n: nil,
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

			got := readString(s)
			if !got.equals(tt.want) {
				t.Errorf("readString() = %v, want %v", got, tt.want)
			}
		})
	}
}
