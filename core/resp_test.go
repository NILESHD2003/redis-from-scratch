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
