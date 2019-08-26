package comm

import (
	"fmt"
	"log"
	"testing"
)

func TestParsePart(t *testing.T) {
	data, _ := BinaryHelper.GetBCDString("5b000000480000000e100100bc614e010001000000000000bc614e34323334353637383132372e302e302e3100000000000000000000000000000000000000000000004671cc385d5b000000480000000e100100bc614e010001000000000000bc614e35323334353637383132372e302e302e3100000000000000000000000000000000000000000000004671cc385d5b000000480000000e100100bc614e010001000000000000bc614e36323334353637383132372e302e302e3100000000000000000000000000000000000000000000004671cc385d5c5b000000480000000e100100bc614e010001000000000000bc614e37323334353637383132372e302e302e3100000000000000000000000000000000000000000000004671cc385d11")
	pack, left, err := ParseHelper.ParsePart(data, 0x5b, 0x5d)
	if err != nil {
		log.Println(err)
	}
	if len(pack) > 0 {
		for i, v := range pack {
			log.Println(i, BinaryHelper.ToBCDString(v, int32(0), int32(len(v))))
		}
	}
	log.Println("package size:", len(pack), " left", len(left))
}

func TestParsePart2(t *testing.T) {
	data, _ := BinaryHelper.GetBCDString("5b0000005a020000ab74120000bc614e0100010000000000c2b34748334e3133000000000000000000000000000112020000002400100807e21203150719843f02302e900024002700005629005a020029000c000300000004baad5d")
	pack, left, err := ParseHelper.ParsePart(data, 0x5b, 0x5d)
	if err != nil {
		log.Println(err)
	}
	if len(pack) > 0 {
		for i, v := range pack {
			log.Println(i, BinaryHelper.ToBCDString(v, int32(0), int32(len(v))))
		}
	}
	log.Println("package size:", len(pack), " left", len(left))
}

func TestCopyBuf(t *testing.T) {
	buf := make([]byte, 3)
	data := []byte{0x01, 0x02, 0x03}
	copy(buf, data)
	fmt.Println(buf)
	fmt.Println(data)
}
