package core

import "errors"

// Common Errors
var (
	ErrCRLFNotFound         = errors.New("CRLF not found in data stream.")
	ErrInvalidRESPType      = errors.New("Invalid RESP type.")
	ErrEmptyData            = errors.New("Empty data.")
	ErrInvalidInteger       = errors.New("Invalid integer format.")
	ErrInvalidLength        = errors.New("Invalid length for bulk string or array.")
	ErrUnsupportedRESPValue = errors.New("Unsupported RESP Value")
)

func DecodeRESPString(data []byte) ([]string, error) {
	value, err := DecodeRESP(data)

	if err != nil {
		return nil, err
	}

	ts := value.([]interface{})

	tokens := make([]string, len(ts))

	for i := range tokens {
		tokens[i] = ts[i].(string)
	}

	return tokens, nil
}
