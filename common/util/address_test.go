package util

import (
	"log"
	"testing"
)

func TestBase58ToAddress(t *testing.T) {
	a, err := Base58ToAddress("TVpznkJcNNRMBA1MBsBd9n54Y2sTpoBx1n")
	if err != nil {
		panic(err)
	}
	log.Println(a.Hex())
}

func TestHexToAddress(t *testing.T) {
	a := HexToAddress("41d9a7ab26d45627dbd45856c550550cc9c4cf26c9")
	log.Println(a.Base58())
}
