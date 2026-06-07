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

// Helper to find the length of the data stream for bulk strings and arrays.
func readLength(data []byte) (int, int, error) {
	length, pos, error := decodeInteger64(data)

	return int(length), pos, error
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

// Read the integer value from the data stream and return it as an int64,
// along with the number of bytes read(delta) and any error encountered.
func decodeInteger64(data []byte) (int64, int, error) {
	pos := 1

	var value int64 = 0

	relativePos := findCRLFIndex(data[pos:])
	if relativePos == -1 {
		return 0, 0, ErrCRLFNotFound
	}

	crlfPos := pos + relativePos

	for ; pos < crlfPos; pos++ {
		value = value*10 + int64(data[pos]-'0')
	}

	return value, pos + 2, nil
}

// Read the bulk string from the data stream and return it as a string,
// along with the number of bytes read(delta) and any error encountered.
func decodeBulkString(data []byte) (string, int, error) {
	length, pos, err := readLength(data)
	if err != nil {
		return "", 0, err
	}

	if length < 0 {
		return "", pos, nil
	}

	if pos+int(length)+2 > len(data) {
		return "", 0, ErrCRLFNotFound
	}

	return string(data[pos : pos+int(length)]), pos + int(length) + 2, nil
}

// Read the array from the data stream and return it as a slice of interfaces,
// along with the number of bytes read(delta) and any error encountered.
func decodeArray(data []byte) ([]interface{}, int, error) {
	length, pos, err := readLength(data)
	if err != nil {
		return nil, 0, err
	}

	if length < 0 {
		return nil, pos, nil
	}

	var elements []interface{} = make([]interface{}, length)

	for i := range elements {
		elem, delta, err := PartialDecodeRESP(data[pos:])
		if err != nil {
			return nil, 0, err
		}

		elements[i] = elem
		pos += delta
	}
	return elements, pos, nil
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
