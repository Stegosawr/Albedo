package static

import "errors"

var (
	// ErrURLParseFailed defines URL parse failed error.
	ErrURLParseFailed = errors.New("URL parse failed")
	// ErrCurrencyConversionFailed defines a currency conversion error.
	ErrCurrencyConversionFailed = errors.New("currency conversion failed")
	// ErrParsingMessageFailed defines a message parsing error.
	ErrParsingMessageFailed = errors.New("message parsing failed")
)
