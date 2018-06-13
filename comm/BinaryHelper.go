package comm

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

var BinaryHelper binaryHelper

type binaryHelper struct{}

/*
	bytes to int16
*/
func (binaryHelper) ToInt16(value []byte, startIndex int32) int16 {
	return int16(binary.BigEndian.Uint16(value[startIndex:]))
}
func (binaryHelper) ToUInt16(value []byte, startIndex int32) uint16 {
	return binary.BigEndian.Uint16(value[startIndex:])
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
func (binaryHelper) Int16ToBytes(value int16) []byte {
	var rsp = make([]byte, 2)
	rsp[0] = byte((value >> 8) & 0xFF)
	rsp[1] = byte(value & 0xFF)
	return rsp
}
func (binaryHelper) UInt16ToBytes(value uint16) []byte {
	var rsp = make([]byte, 2)
	rsp[0] = byte((value >> 8) & 0xFF)
	rsp[1] = byte(value & 0xFF)
	return rsp
}

/*
	int to bytes 小端
*/
func (binaryHelper) Int32ToBytes(value int) []byte {
	var rsp = make([]byte, 4)
	rsp[0] = byte((value >> 24) & 0xFF)
	rsp[1] = byte(value >> 16 & 0xFF)
	rsp[2] = byte((value >> 8) & 0xFF)
	rsp[3] = byte(value & 0xFF)
	return rsp
}
func (binaryHelper) UInt32ToBytes(value uint) []byte {
	var rsp = make([]byte, 4)
	rsp[0] = byte((value >> 24) & 0xFF)
	rsp[1] = byte(value >> 16 & 0xFF)
	rsp[2] = byte((value >> 8) & 0xFF)
	rsp[3] = byte(value & 0xFF)
	return rsp
}

/*
	int to bytes 小端
*/
func (binaryHelper) Int64ToBytes(value int64) []byte {
	var rsp = make([]byte, 4)
	rsp[0] = byte((value >> 56) & 0xFF)
	rsp[1] = byte(value >> 48 & 0xFF)
	rsp[2] = byte((value >> 40) & 0xFF)
	rsp[3] = byte(value >> 32 & 0xFF)
	rsp[4] = byte((value >> 24) & 0xFF)
	rsp[5] = byte(value >> 16 & 0xFF)
	rsp[6] = byte((value >> 8) & 0xFF)
	rsp[7] = byte(value & 0xFF)
	return rsp
}
func (binaryHelper) UInt64ToBytes(value uint64) []byte {
	var rsp = make([]byte, 4)
	rsp[0] = byte((value >> 56) & 0xFF)
	rsp[1] = byte(value >> 48 & 0xFF)
	rsp[2] = byte((value >> 40) & 0xFF)
	rsp[3] = byte(value >> 32 & 0xFF)
	rsp[4] = byte((value >> 24) & 0xFF)
	rsp[5] = byte(value >> 16 & 0xFF)
	rsp[6] = byte((value >> 8) & 0xFF)
	rsp[7] = byte(value & 0xFF)
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
func (binaryHelper) GetASCIIStringWL(value string, length int) []byte {
	data := []byte(value)
	rsp := make([]byte, length)
	if len(data) < length {
		copy(rsp, data)
		return rsp
	}
	return data[:length]
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
func (binaryHelper) Byte808Descape(value []byte, startIndex int, length int) ([]byte, error) {
	ilength := len(value)
	if (startIndex + length) > ilength {
		return nil, errors.New("长度不足，下标越界")
	}
	buf := new(bytes.Buffer)
	/*去头去尾*/
	for i := startIndex + 1; i < ilength-1; i++ {
		if value[i] != 0x7D || i == ilength-2 {
			buf.WriteByte(value[i])
		} else if value[i+1] == 0x02 {
			buf.WriteByte(0x7e)
			i++
		} else if value[i+1] == 0x01 {
			buf.WriteByte(0x7d)
			i++
		} else {
			return nil, errors.New("终端数据包含非法转义字符7D:" + strconv.Itoa(int(value[i+1])))
		}

	}
	return buf.Bytes(), nil
}

/*
	转义808字符
	0x7e -> 0x7d 0x01
	0x7d -> 0x7d 0x02
*/
func (binaryHelper) Byte808Enscape(value []byte, startIndex int, length int) []byte {
	ilength := len(value)
	if (startIndex + length) > ilength {

	}
	buf := new(bytes.Buffer)
	buf.WriteByte(0x7e)
	for i := startIndex; i < ilength; i++ {
		if value[i] == 0x7D {
			buf.WriteByte(0x7d)
			buf.WriteByte(0x01)
		} else if value[i] == 0x7e {
			buf.WriteByte(0x7d)
			buf.WriteByte(0x02)
		} else {
			buf.WriteByte(value[i])
		}
	}
	buf.WriteByte(0x7e)
	return buf.Bytes()
}

/*CRC Check*/
func (binaryHelper) CRCCheck(value []byte) bool {
	bCRC := byte(0x00)
	for i := 0; i < len(value)-1; i++ {
		bCRC ^= value[i]
	}
	if bCRC == value[len(value)-1] {
		return true
	} else {
		return false
	}
}

/*CRC16 Check*/
func (binaryHelper) CRC16Check(data []byte) int16 {
	var crc_reg int16 = -1
	var current int16 = 0
	for i := 0; i < len(data); i++ {
		current = int16(data[i]) << 8
		for j := 0; j < 8; j++ {
			if (crc_reg ^ current) < 0 {
				crc_reg = ((crc_reg << 1) ^ 0x1021)
			} else {
				crc_reg <<= 1
			}
			current <<= 1
		}
	}
	return int16(crc_reg)
}

//生成GUID
func readMachineId() []byte {
	var sum [3]byte
	id := sum[:]
	hostname, err1 := os.Hostname()
	if err1 != nil {
		_, err2 := io.ReadFull(rand.Reader, id)
		if err2 != nil {
			panic(fmt.Errorf("cannot get hostname: %v; %v", err1, err2))
		}
		return id
	}
	hw := md5.New()
	hw.Write([]byte(hostname))
	copy(id, hw.Sum(nil))
	fmt.Println("readMachineId:" + string(id))
	return id
}

var objectIdCounter uint32 = 0
var machineId = readMachineId()

// NewObjectId returns a new unique ObjectId.
// 4byte 时间，
// 3byte 机器ID
// 2byte pid
// 3byte 自增ID
//长度24
func (binaryHelper) NewObjectId() string {
	var b [12]byte
	// Timestamp, 4 bytes, big endian
	binary.BigEndian.PutUint32(b[:], uint32(time.Now().Unix()))
	// Machine, first 3 bytes of md5(hostname)
	b[4] = machineId[0]
	b[5] = machineId[1]
	b[6] = machineId[2]
	// Pid, 2 bytes, specs don't specify endianness, but we use big endian.
	pid := os.Getpid()
	b[7] = byte(pid >> 8)
	b[8] = byte(pid)
	// Increment, 3 bytes, big endian
	i := atomic.AddUint32(&objectIdCounter, 1)
	b[9] = byte(i >> 16)
	b[10] = byte(i >> 8)
	b[11] = byte(i)
	return BinaryHelper.ToBCDString(b[:], 0, int32(len(b)))
}

//生成32位md5字串
func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

//生成Guid字串
func (binaryHelper) UniqueId() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	guid := GetMd5String(base64.URLEncoding.EncodeToString(b))
	return fmt.Sprintf("%s-%s-%s-%s-%s", guid[0:8], guid[8:12], guid[12:16], guid[16:20], guid[20:])
}
