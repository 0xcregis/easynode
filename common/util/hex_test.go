package util

import (
	"log"
	"testing"
)

func TestHex(t *testing.T) {
	//0x0000000000000000000000006b75d8af000000e20b7a7ddf000ba900b4009a80
	//bs, err := FromHex("0x0000000000000000000000006b75d8af000000e20b7a7ddf000ba900b4009a80")
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	//log.Println(BytesToHexString(bs[12:]))

	log.Println(Hex2Address("000000000000000000000000e96051e8da2f0cc02c372252dfacfdb129b5d4d6"))
}

func TestHexToInt(t *testing.T) {
	//log.Println(strconv.ParseInt("0x000000000000000000000000000000000000000000000000000000001dcd6500", 16, 64))
	log.Println(HexToInt("0000000000000000000000000000000000000000000004cb433cf96ff45c0000"))
}
