package util

import (
	"log"
	"testing"
)

func TestBase58ToAddress(t *testing.T) {
	a, err := Base58ToAddress("TXsmKpEuW7qWnXzJLGP9eDLvWPR2GRn1FS")
	if err != nil {
		panic(err)
	}
	log.Println(a.Hex())
}

func TestHexToAddress(t *testing.T) {
	a := HexToAddress("0x4142a1e39aefa49290f2b3f9ed688d7cecf86cd6e0")
	log.Println(a.Base58())
}
