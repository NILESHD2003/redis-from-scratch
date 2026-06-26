package core

type RESPValue interface {
	isRESPValue()
}

const (
	SimpleStringPrefix = '+'
	ErrorPrefix        = '-'
	IntegerPrefix      = ':'
	BulkStringPrefix   = '$'
	ArrayPrefix        = '*'
)

type SimpleString string

type BulkString string

type ErrorString string

type Integer int64

type Array []RESPValue

type NullBulkString struct{}

type NullArray struct{}

func (SimpleString) isRESPValue() {}

func (BulkString) isRESPValue() {}

func (ErrorString) isRESPValue() {}

func (Integer) isRESPValue() {}

func (Array) isRESPValue() {}

func (NullBulkString) isRESPValue() {}

func (NullArray) isRESPValue() {}
