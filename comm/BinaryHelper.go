package comm

import (
	"strconv"
	"bytes"
	"errors"
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
	int16 to bytes 小端
*/
func(binaryHelper) Int16ToBytes(value int16)[]byte{
	var rsp = make([]byte,2)
	rsp[0] =byte((value>>8) & 0xFF) 
	rsp[1] =byte(value & 0xFF)
	return rsp
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

/*
	转义808字符
	0x7d 0x01 -> 0x7e
	0x7d 0x02 -> 0x7d
*/
func(binaryHelper)Byte808Descape(value []byte,startIndex int, length int)([]byte,error){
	ilength := len(value)
	if (startIndex+length)>ilength{
		return nil,errors.New("长度不足，下标越界")
	}	
	buf :=new(bytes.Buffer)
	/*去头去尾*/
	for i:=startIndex+1;i<ilength-1;i++{
		if value[i]!=0x7D || i==ilength-2{
			buf.WriteByte(value[i])
		}else if value[i+1]==0x02{
			buf.WriteByte(0x7e)
			i++
		}else if value[i+1]==0x01{
			buf.WriteByte(0x7d)
			i++
		}else{
			return nil,errors.New("终端数据包含非法转义字符7D:"+strconv.Itoa(int(value[i+1])))
		}

	}	
	return buf.Bytes(),nil
}
/*
	转义808字符
	0x7e -> 0x7d 0x01 
	0x7d -> 0x7d 0x02
*/
func(binaryHelper)Byte808Enscape(value[]byte,startIndex int, length int)([]byte){
	ilength := len(value)
	if (startIndex+length)>ilength{
		
	}
	buf :=new(bytes.Buffer)
	buf.WriteByte(0x7e)
	for i:=startIndex;i<ilength;i++{
		if value[i]==0x7D{
			buf.WriteByte(0x7d)
			buf.WriteByte(0x01)
		}else if value[i]==0x7e{
			buf.WriteByte(0x7d)
			buf.WriteByte(0x02)
		}else{
			buf.WriteByte(value[i])
		}
	}
	buf.WriteByte(0x7e)	
	return buf.Bytes()
}


/*CRC Check*/
func(binaryHelper)CRCCheck(value []byte) bool{
	bCRC :=byte(0x00)
	for i:=0;i<len(value)-1;i++{
		bCRC ^= value[i]
	}
	if bCRC == value[len(value)-1]{
		return true
	}else{
		return false
	}
}
