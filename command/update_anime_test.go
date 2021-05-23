package command

import (
	"testing"
)

func TestScrapeAnime(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Get Anime Releases",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := scraperAnime()
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestAddZeroTo2DigitNum(t *testing.T) {
	tests := []struct {
		name   string
		numIn  int
		numOut string
	}{
		{
			name:   "Standard",
			numIn:  2,
			numOut: "02",
		}, {
			name:   "Over",
			numIn:  10,
			numOut: "10",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			num := addZeroTo2DigitNum(tt.numIn)
			if num != tt.numOut {
				t.Errorf("Want: %v, Got: %v", tt.numOut, num)
			}
		})
	}
}

func TestRemoveDuplicatesUnordered(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		out  []string
	}{
		{
			name: "Standard",
			in:   []string{"1", "2", "3", "1"},
			out:  []string{"1", "2", "3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slice := removeDuplicatesUnordered(tt.in)
			if len(slice) != len(tt.out) {
				t.Errorf("Want: %v, Got: %v", len(tt.out), len(slice))
			}
		})
	}
}

func TestAppendStringSeq(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		sep  string
		out  string
	}{
		{
			name: "Standard",
			in:   []string{"abc", "def", "ghi"},
			sep:  ",",
			out:  "abc,def,ghi",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str := appendStringSeq(tt.sep, tt.in...)
			if str != tt.out {
				t.Errorf("Want: %v, Got: %v", tt.out, str)
			}
		})
	}
}

func TestWrapInString(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		sep  string
		out  []string
	}{
		{
			name: "Standard",
			in:   []string{"abc", "def", "ghi"},
			sep:  "**",
			out:  []string{"**abc**", "**def**", "**ghi**"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str := wrappStringIn(tt.sep, tt.in...)
			for i := range str {
				if str[i] != tt.out[i] {
					t.Errorf("Want: %v, Got: %v", tt.out[i], str[i])
				}
			}
		})
	}
}
