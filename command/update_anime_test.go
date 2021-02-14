package command

import "testing"

func TestScrapeHentai(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Get Hentai Releases",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := scraperHentai()
			if err != nil {
				t.Error(err)
			}
		})
	}
}

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
			name:   "OVer",
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
