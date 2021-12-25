package currency

import "testing"

func TestConvCurrency(t *testing.T) {
	tests := []struct {
		Name    string
		Content string
		Source  string
		Target  string
		Want    string
	}{
		{
			Name:    "EUR to USD",
			Content: "62.50 EUR",
			Source:  "EUR",
			Target:  "💵",
			Want:    "70.25 USD",
		}, {
			Name:    "USD to EUR",
			Content: "70.25 USD",
			Source:  "USD",
			Target:  "💶",
			Want:    "62.50 EUR",
		}, {
			Name:    "EUR to JPY",
			Content: "62.50 EUR",
			Source:  "EUR",
			Target:  "💴",
			Want:    "7992 JPY",
		},
	}
	exchangeRates = map[string]float64{
		"EUR": 1.12407,
		"USD": 1,
		"JPY": 0.00879076,
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			newCurr := convertCurr(tt.Content, tt.Source, tt.Target)
			if newCurr != tt.Want {
				t.Errorf("Want: %v, Got: %v", tt.Want, newCurr)
			}
		})
	}
}
