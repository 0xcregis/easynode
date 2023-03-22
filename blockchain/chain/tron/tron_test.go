package tron

import (
	"log"
	"testing"
)

func TestEth_GetToken(t *testing.T) {
	c := NewChainClient()
	log.Println(c.GetTokenBalance("grpc.trongrid.io:50051", "", "TLa2f6VPqDgRE67v1736s7bJ8Ray5wYjU7", "TMuA6YqfCeX8EhbfYEg5y7S4DqzSJireY9"))
}

func TestGetTokenByHttp(t *testing.T) {
	c := NewChainClient()
	log.Println(c.GetTokenBalanceByHttp("https://api.trongrid.io/wallet/triggerconstantcontract", "244f918d-56b5-4a16-9665-9637598b1223", "THb4CqiFdwNHsWsQCs4JhzwjMWys4aqCbF", "TYZJDb3TW2iTQsmgBtfNmzAi1chivK1Zf5"))
}
