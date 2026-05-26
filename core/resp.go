package core

import "errors"

// Common Errors
var (
	ErrCRLFNotFound    = errors.New("CRLF not found in data stream.")
	ErrInvalidRESPType = errors.New("Invalid RESP type.")
	ErrEmptyData       = errors.New("Empty data.")
)

// Helper to find the next CLRF in the data stream
func findCRLFIndex(data []byte) int {
	for i := 0; i < len(data)-1; i++ {
		if data[i] == '\r' && data[i+1] == '\n' {
			return i
		}
	}
	return -1
}

// Read the simple string from the data stream and return it as a string,
// along with the number of bytes read(delta) and any error encountered.
func decodeSimpleString(data []byte) (string, int, error) {
	pos := 1

	relativePos := findCRLFIndex(data[pos:])
	if relativePos == -1 {
		return "", 0, ErrCRLFNotFound
	}

	crlfPos := pos + relativePos

	return string(data[pos:crlfPos]), crlfPos + 2, nil
}

// Read the error message from the data stream and return it as a string,
// along with the number of bytes read(delta) and any error encountered.
func decodeError(data []byte) (string, int, error) {
	return decodeSimpleString(data)
}

func decodeInteger64(data []byte) (int64, int, error) {
	// pending implementation
}

func decodeBulkString(data []byte) (string, int, error) {
	// pending implementation
}

func decodeArray(data []byte) ([]interface{}, int, error) {
	// pending implementation
}

func PartialDecodeRESP(data []byte) (interface{}, int, error) {
	if len(data) == 0 {
		return nil, 0, ErrEmptyData
	}

	switch data[0] {
	case '+':
		return decodeSimpleString(data)
	case '-':
		return decodeError(data)
	case ':':
		return decodeInteger64(data)
	case '$':
		return decodeBulkString(data)
	case '*':
		return decodeArray(data)
	default:
		return nil, 0, ErrInvalidRESPType
	}
}

func DecodeRESP(data []byte) (interface{}, error) {
	if len(data) == 0 {
		return nil, ErrEmptyData
	}

	value, _, err := PartialDecodeRESP(data)

	return value, err
}
