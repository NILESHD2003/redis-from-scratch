package core

import (
	"reflect"
	"testing"
)

func TestDecodeRESP(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    interface{}
		wantErr error
	}{
		{
			name:    "Empty input",
			input:   []byte(""),
			want:    nil,
			wantErr: ErrEmptyData,
		},
		{
			name:    "Invalid RESP Type",
			input:   []byte("%Invalid\r\n"),
			want:    nil,
			wantErr: ErrInvalidRESPType,
		},
		{
			name:    "Partial input with valid RESP but partial frame",
			input:   []byte("+OK\r"),
			want:    "",
			wantErr: ErrCRLFNotFound,
		},
		{
			name:    "Valid simple string",
			input:   []byte("+OK\r\n"),
			want:    "OK",
			wantErr: nil,
		},
		{
			name:    "Valid Error string",
			input:   []byte("-ERR something went wrong\r\n"),
			want:    "ERR something went wrong",
			wantErr: nil,
		},
		{
			name:    "Valid Integer",
			input:   []byte(":10\r\n"),
			want:    int64(10),
			wantErr: nil,
		},
		{
			name:    "Valid Bulk string",
			input:   []byte("$2\r\nOK\r\n"),
			want:    "OK",
			wantErr: nil,
		},
		{
			name:    "Valid Array",
			input:   []byte("*3\r\n$3\r\nPUT\r\n$1\r\nK\r\n$1\r\nV\r\n"),
			want:    []interface{}{"PUT", "K", "V"},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeRESP(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecodeRESP() got = %v, want = %v", got, tt.want)
			}

			if err != tt.wantErr {
				t.Errorf("DecodeRESP() err = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestDecodeSimpleString(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    string
		wantErr error
	}{
		{
			name:    "Valid simple string",
			input:   []byte("+OK\r\n"),
			want:    "OK",
			wantErr: nil,
		},
		{
			name:    "Missing CRLF",
			input:   []byte("+OK"),
			want:    "",
			wantErr: ErrCRLFNotFound,
		},
		{
			name:    "Empty simple string",
			input:   []byte("+\r\n"),
			want:    "",
			wantErr: nil,
		},
		{
			name:    "Simple string with spaces",
			input:   []byte("+hello world\r\n"),
			want:    "hello world",
			wantErr: nil,
		},
		{
			name:    "Simple string with special characters",
			input:   []byte("+!@#$%^&*()\r\n"),
			want:    "!@#$%^&*()",
			wantErr: nil,
		},
		{
			name:    "String with only CR",
			input:   []byte("+OK\r"),
			want:    "",
			wantErr: ErrCRLFNotFound,
		},
		{
			name:    "String with only LF",
			input:   []byte("+OK\n"),
			want:    "",
			wantErr: ErrCRLFNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := decodeSimpleString(tt.input)
			if got != tt.want || err != tt.wantErr {
				t.Errorf("decodeSimpleString() = got %q, want %q, gotErr %v, wantErr %v", got, tt.want, err, tt.wantErr)
			}
		})
	}
}

func TestDecodeError(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    string
		wantErr error
	}{
		{
			name:    "Valid error string",
			input:   []byte("-ERR something went wrong\r\n"),
			want:    "ERR something went wrong",
			wantErr: nil,
		},
		{
			name:    "Missing CRLF",
			input:   []byte("-ERR something went wrong"),
			want:    "",
			wantErr: ErrCRLFNotFound,
		},
		{
			name:    "Empty error string",
			input:   []byte("-\r\n"),
			want:    "",
			wantErr: nil,
		},
		{
			name:    "Error with special characters",
			input:   []byte("-ERR !@#$%^&*()\r\n"),
			want:    "ERR !@#$%^&*()",
			wantErr: nil,
		},
		{
			name:    "Error with spaces",
			input:   []byte("-ERR something went wrong with   spaces\r\n"),
			want:    "ERR something went wrong with   spaces",
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := decodeError(tt.input)
			if got != tt.want || err != tt.wantErr {
				t.Errorf("decodeError() = got %q, want %q, gotErr %v, wantErr %v", got, tt.want, err, tt.wantErr)
			}
		})
	}
}

func TestDecodeInteger64(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    int64
		wantErr error
	}{
		{
			name:    "Valid integer",
			input:   []byte(":10\r\n"),
			want:    10,
			wantErr: nil,
		},
		{
			name:    "Missing CRLF",
			input:   []byte(":10"),
			want:    0,
			wantErr: ErrCRLFNotFound,
		},
		{
			name:    "Zero integer",
			input:   []byte(":0\r\n"),
			want:    0,
			wantErr: nil,
		},
		{
			name:    "Single digit integer",
			input:   []byte(":1\r\n"),
			want:    1,
			wantErr: nil,
		},
		{
			name:    "Large integer",
			input:   []byte(":123456789\r\n"),
			want:    123456789,
			wantErr: nil,
		},
		{
			name:    "Max integer",
			input:   []byte(":9223372036854775807\r\n"),
			want:    9223372036854775807,
			wantErr: nil,
		},
		{
			name:    "Min integer",
			input:   []byte(":-9223372036854775808\r\n"),
			want:    -9223372036854775808,
			wantErr: nil,
		},
		{
			name:    "Non numeric characters",
			input:   []byte(":10abc\r\n"),
			want:    0,
			wantErr: ErrInvalidInteger,
		},
		{
			name:    "Negative integer",
			input:   []byte(":-10\r\n"),
			want:    -10,
			wantErr: nil,
		},
		{
			name:    "Integer with leading zeros",
			input:   []byte(":00010\r\n"),
			want:    10,
			wantErr: nil,
		},
		{
			name:    "Empty integer",
			input:   []byte(":\r\n"),
			want:    0,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := decodeInteger64(tt.input)
			if got != tt.want || err != tt.wantErr {
				t.Errorf("decodeInteger64() = got %d, want %d, gotErr %v, wantErr %v", got, tt.want, err, tt.wantErr)
			}
		})
	}
}

func TestDecodeBulkString(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    string
		wantErr error
	}{
		{
			name:    "Valid bulk string",
			input:   []byte("$2\r\nOK\r\n"),
			want:    "OK",
			wantErr: nil,
		},
		{
			name:    "Valid bulk string with special characters",
			input:   []byte("$13\r\nHello, World!\r\n"),
			want:    "Hello, World!",
			wantErr: nil,
		},
		{
			name:    "Empty bulk string",
			input:   []byte("$0\r\n\r\n"),
			want:    "",
			wantErr: nil,
		},
		{
			name:    "Bulk string with spaces",
			input:   []byte("$11\r\nHello World\r\n"),
			want:    "Hello World",
			wantErr: nil,
		},
		{
			name:    "Null bulk string",
			input:   []byte("$-1\r\n"),
			want:    "",
			wantErr: nil,
		},
		{
			name:    "Single character bulk string",
			input:   []byte("$1\r\nA\r\n"),
			want:    "A",
			wantErr: nil,
		},
		{
			name:    "Bulk string with CRLF as content",
			input:   []byte("$2\r\n\r\n\r\n"),
			want:    "\r\n",
			wantErr: nil,
		},
		{
			name:    "Bulk string with missing CRLF",
			input:   []byte("$2\r\nOK"),
			want:    "",
			wantErr: ErrCRLFNotFound,
		},
		{
			name:    "Bulk string with invalid length",
			input:   []byte("$-2\r\nOK\r\n"),
			want:    "",
			wantErr: ErrInvalidLength,
		},
		{
			name:    "Bulk string with non-numeric length",
			input:   []byte("$abc\r\nOK\r\n"),
			want:    "",
			wantErr: ErrInvalidLength,
		},
		{
			name:    "Bulk string with payload shorter than specified length",
			input:   []byte("$5\r\nOK\r\n"),
			want:    "",
			wantErr: ErrCRLFNotFound,
		},
		{
			name:    "Bulk string with extra data after content",
			input:   []byte("$2\r\nOK\r\nEXTRA"),
			want:    "OK",
			wantErr: nil,
		},
		{
			name:    "Bulk string with payload longer than specified length",
			input:   []byte("$2\r\nOKOK\r\n"),
			want:    "",
			wantErr: ErrCRLFNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := decodeBulkString(tt.input)
			if got != tt.want || err != tt.wantErr {
				t.Errorf("decodeBulkString() = got %q, want %q, gotErr %v, wantErr %v", got, tt.want, err, tt.wantErr)
			}
		})
	}
}
