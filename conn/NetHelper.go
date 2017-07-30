package conn

import (
	"bytes"
	"log"
	"strconv"
)

func MyRecover() {
	if err := recover(); err != nil {
		log.Println("On Recover", err)
		//fmt.Println(err)
	}
}

func ToHex(d []byte) string {
	buffer := new(bytes.Buffer)
	for _, b := range d {
		buffer.WriteString(strconv.FormatInt(int64(b&0xff), 16))
	}
	return buffer.String()
}
