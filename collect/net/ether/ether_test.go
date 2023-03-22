package ether

import (
	"github.com/sunjiangjun/xlog"
	"log"
	"math"
	"testing"
)

func TestEth_GetTransactionByHash(t *testing.T) {
	log.Println(Eth_GetTransactionByHash("https://eth-mainnet.g.alchemy.com/v2", "demo", "0x5917da4788cdc1383215541744beb93fd804c1902e221d2c5555ce99d9bfff42", nil))
}

func TestGetBlockByHash(t *testing.T) {
	x := xlog.NewXLogger()
	log.Println(Eth_GetBlockByHash("https://eth-mainnet.g.alchemy.com/v2", "RzxBjjh_c4y0LVHZ7GNm8zoXEZR3HYop", "0xb49d607f5b80890531e3e1d57798a7573cf8e18048ec0df34e3c81d48115078f", x))
}

func TestEth_GetBlockByNumber(t *testing.T) {
	//log.Println(Eth_GetBlockByNumber("https://eth-mainnet.g.alchemy.com/v2", "demo", "0xF3F088"))

	//blockNumber := "15986824"
	//if !strings.HasPrefix(blockNumber, "0x") {
	//	n, _ := strconv.ParseInt(blockNumber, 10, 64)
	//	blockNumber = fmt.Sprintf("0x%x", n)
	//}
	//log.Println(blockNumber)

	bs := []byte("abc123a")

	log.Println(bs[len(bs)-1:])

	log.Printf("%d", bs[len(bs)-1])
	log.Printf("%c", rune(bs[0]))

	m := math.Mod(float64(1), float64(20))
	log.Println(m)
}
