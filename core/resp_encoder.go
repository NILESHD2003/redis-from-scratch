package core

import "strconv"

func EncodeRESP(value RESPValue) ([]byte, error) {
	return partialEncodeRESP(value)
}

func partialEncodeRESP(value RESPValue) ([]byte, error) {

	switch v := value.(type) {
	case SimpleString:
		return encodeSimpleString(v), nil
	case BulkString:
		return encodeBulkString(v), nil
	case ErrorString:
		return encodeError(v), nil
	case Integer:
		return encodeInteger(v), nil
	case Array:
		return encodeArray(v)
	case NullBulkString:
		return encodeNullBulkString(), nil
	case NullArray:
		return encodeNullArray(), nil
	default:
		return nil, ErrUnsupportedRESPValue
	}
}

func encodeSimpleString(value SimpleString) []byte {
	return []byte(string(SimpleStringPrefix) + string(value) + "\r\n")
}

func encodeBulkString(value BulkString) []byte {
	return []byte(string(BulkStringPrefix) + strconv.Itoa(len(value)) + "\r\n" + string(value) + "\r\n")
}

func encodeError(value ErrorString) []byte {
	return []byte(string(ErrorPrefix) + string(value) + "\r\n")
}

func encodeInteger(value Integer) []byte {
	return []byte(string(IntegerPrefix) + strconv.FormatInt(int64(value), 10) + "\r\n")
}

func encodeArray(value Array) ([]byte, error) {
	result := []byte(string(ArrayPrefix) + strconv.Itoa(len(value)) + "\r\n")

	for _, v := range value {
		encodedValue, err := partialEncodeRESP(v)
		if err != nil {
			return nil, err
		}
		result = append(result, encodedValue...)
	}

	return result, nil
}

func encodeNullBulkString() []byte {
	return []byte("$-1\r\n")
}

func encodeNullArray() []byte {
	return []byte("*-1\r\n")
}
