package util

import (
	"log"
	"testing"
)

func TestHex(t *testing.T) {

	bs, err := FromHex("0000000000000000000000003926f326e328be2616f7b05f00bcecac652c8f15")

	if err != nil {
		panic(err)
	}

	log.Println(BytesToHexString(bs[12:]))

}

func TestHexToInt(t *testing.T) {

	//log.Println(strconv.ParseInt("0x000000000000000000000000000000000000000000000000000000001dcd6500", 16, 64))
	log.Println(HexToInt("0000000000000000000000000000000000000000000004cb433cf96ff45c0000"))
}
