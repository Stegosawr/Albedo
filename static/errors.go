package static

import "errors"

var (
	// ErrCurrencyConversionFailed defines a currency conversion error.
	ErrCurrencyConversionFailed = errors.New("currency conversion failed")
	// ErrParsingMessageFailed defines a message parsing error.
	ErrParsingMessageFailed = errors.New("message parsing failed")
)
