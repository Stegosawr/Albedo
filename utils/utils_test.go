package utils

import "testing"

func TestUnShortenURL(t *testing.T) {
	tests := []struct {
		Name string
		URL  string
		Want string
	}{
		{
			Name: "Post 1",
			URL:  "https://t.co/TM9aZbnuOX",
			Want: "https://www.amiami.com/eng/detail/?gcode=FIGURE-134930",
		}, {
			Name: "Post 2",
			URL:  "https://t.co/CheQksbjjs",
			Want: "https://www.amiami.com/eng/search/list/?s_keywords=Kizuna%20AI%20Exclusive%20Sale&pagecnt=1&s_st_list_preorder_available=1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			URL, err := UnShortenURL(tt.URL)
			if err != nil {
				t.Error(err)
			}
			if URL != tt.Want {
				t.Errorf("Want: %v, Got: %v", tt.Want, URL)
			}
		})
	}
}
