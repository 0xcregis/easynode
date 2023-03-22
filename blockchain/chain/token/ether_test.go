package token

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math"
	"math/big"
	"testing"
)

func Init() *Token {

	//https://rpc.ankr.com/eth
	client, err := ethclient.Dial("https://eth-mainnet.g.alchemy.com/v2/RzxBjjh_c4y0LVHZ7GNm8zoXEZR3HYop")
	if err != nil {
		log.Fatal(err)
		return nil
	}

	// Golem (GNT) Address
	tokenAddress := common.HexToAddress("0xdac17f958d2ee523a2206206994597c13d831ec7")
	instance, err := NewToken(tokenAddress, client)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return instance
}

func TestTokenCaller_BalanceOf(t *testing.T) {

	instance := Init()
	if instance == nil {
		panic("instance is nil")
	}
	address := common.HexToAddress("0x11cc6083a3f2608e7c6b2862185e4874ea7b2b56")
	bal, err := instance.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		log.Panic(err)
	}
	log.Println(bal.String())

	name, err := instance.Name(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	symbol, err := instance.Symbol(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	decimals, err := instance.Decimals(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("name: %s\n", name)         // "name: Golem Network"
	fmt.Printf("symbol: %s\n", symbol)     // "symbol: GNT"
	fmt.Printf("decimals: %v\n", decimals) // "decimals: 18"

	fmt.Printf("wei: %s\n", bal) // "wei: 74605500647408739782407023"

	fbal := new(big.Float)
	fbal.SetString(bal.String())
	value := new(big.Float).Quo(fbal, big.NewFloat(math.Pow10(int(decimals))))

	fmt.Printf("balance: %f", value) // "balance: 74605500.647409"

}
