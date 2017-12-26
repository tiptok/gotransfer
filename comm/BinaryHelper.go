package comm

import (
	"encoding/binary"
	"encoding/hex"
)

var BinaryHelper binaryHelper

type binaryHelper struct{}

/*
	bytes to int16
*/
func (binaryHelper) ToInt16(value []byte, startIndex int32) int16 {
	return int16(binary.BigEndian.Uint16(value[startIndex:]))
}

/*
	bytes to int32
*/
func (binaryHelper) ToInt32(value []byte, startIndex int32) int32 {
	return int32(binary.BigEndian.Uint32(value[startIndex:]))
}

/*
	bytes to int64
*/
func (binaryHelper) ToInt64(value []byte, startIndex int32) int64 {
	return int64(binary.BigEndian.Uint64(value[startIndex:]))
}

/*
	bytes to ASCII String 解码
	31 32 33 34 35 36
	startindex 1  length 3
	234
*/
func (binaryHelper) ToASCIIString(value []byte, startIndex int32, length int32) string {
	return string(value[startIndex : startIndex+length])
}
func (binaryHelper) GetASCIIString(value string) []byte {
	return []byte(value)
}

/*
	bytes to BCD String 解码
	31 32 33 34 35 36
	startindex 1  length 3
	313233
*/
func (binaryHelper) ToBCDString(value []byte, startIndex int32, length int32) string {
	return hex.EncodeToString(value[startIndex : startIndex+length])
}

func (binaryHelper) GetBCDString(value string) ([]byte, error) {
	return hex.DecodeString(value)
}

/*
	返回固定长度切片
	31 32 33 34 35 36
	startindex 1  length 3
	32 33 34
*/
func (binaryHelper) CloneRange(value []byte, startIndex int32, length int32) []byte {
	return value[startIndex : startIndex+length]
}
