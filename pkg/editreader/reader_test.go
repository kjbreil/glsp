package editreader

import (
	"io"
	"strings"
	"testing"
)

// func TestFile_Read(t *testing.T) {
// 	f, err := New(strings.NewReader("Hello\nWorld\n"))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	all, err := io.ReadAll(f)
// 	if err != nil {
// 		return
// 	}
// 	fmt.Println(string(all))
// }

func TestFile_Read(t *testing.T) {

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "single character",
			input:   "a",
			wantErr: false,
		},
		{
			name:    "multi-line",
			input:   "Hello\nWorld\n",
			wantErr: false,
		},
		{
			name:    ">512 characters",
			input:   "abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvwxyz1234567890",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := New(strings.NewReader(tt.input))
			got, err := io.ReadAll(f)

			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != tt.input {
				t.Errorf("Read() gotN = %v, want %v", string(got), tt.input)
			}
		})
	}
}

func TestFile_ReadLine(t *testing.T) {

	tests := []struct {
		name  string
		input string

		want    string
		wantErr bool
	}{
		{
			name:    "single line newline before EOF",
			input:   "Hello World\n",
			want:    "Hello World\n",
			wantErr: false,
		},
		{
			name:    "single line crlf before EOF",
			input:   "Hello World\r\n",
			want:    "Hello World\n",
			wantErr: false,
		},
		{
			name:    "single line EOF",
			input:   "Hello World",
			want:    "Hello World",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := New(strings.NewReader(tt.input))
			got, err := f.ReadLine()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.String() != tt.want {
				t.Errorf("ReadLine() got = %v, want %v", got.String(), tt.want)
			}
		})
	}
}
