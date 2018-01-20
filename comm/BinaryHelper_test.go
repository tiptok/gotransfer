package comm

import (
	"log"
	"testing"
)

func TestBinary(t *testing.T) {
	log.Println(BinaryHelper.ToInt16([]byte{0x00, 0x02}, 0))
	log.Println(BinaryHelper.ToInt32([]byte{0x00, 0x00, 0x01, 0x00}, 0))
	t.Log(BinaryHelper.ToInt16([]byte{0x00, 0x02}, 0))
}

func TestToASCIIString(t *testing.T) {
	data := []byte{0x31, 0x32, 0x33, 0x34, 0x35, 0x36}
	log.Println(BinaryHelper.ToASCIIString(data, 1, 3))
	log.Println(BinaryHelper.ToASCIIString(data, 2, 3))
	log.Println(BinaryHelper.ToASCIIString(data, 3, 3))

	log.Println(BinaryHelper.ToBCDString(BinaryHelper.GetASCIIString("789"), 0, 3))
	log.Println(BinaryHelper.GetASCIIString("567"))
	log.Println(BinaryHelper.GetASCIIString("123"))
	t.Log("end")
}

func TestByte808Descape(t *testing.T){
	data := []byte{0x7e,0x30,0x7d,0x02,0x08,0x7d,0x01,0x7d,0x02,0x7e} 
	dDsp,err := BinaryHelper.Byte808Descape(data,0,len(data))
	if err!=nil{
		log.Println(err.Error())
	}else{
		log.Println(BinaryHelper.ToBCDString(dDsp,0,int32(len(dDsp))))
	}
	dEsp := BinaryHelper.Byte808Enscape(dDsp,0,len(dDsp))
	log.Println(BinaryHelper.ToBCDString(dEsp,0,int32(len(dEsp))))
	t.Log("end")
}
