package token

import (
	"context"
	address2 "github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/common"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"google.golang.org/grpc"
	"log"
	"strconv"
	"testing"
)

func TestCallContract(t *testing.T) {
	a := "TAUN6FwrnwwmaEqYcckffC7wYmbaS6cBiX"
	c1, _ := address2.Base58ToAddress(a)

	log.Printf("hex=%s \n", c1.Hex())

	c2 := c1.Hex()[4:]
	log.Printf("del 0x :hex=%s", c2)

	m := "0000000000000000000000000000000000000000000000000000000000000000"

	log.Println(string(m[:len(m)-len(c2)]) + string(c2))

	log.Println(strconv.ParseInt("0000000000000000000000000000000000000000000000000000000000d0bfbd", 16, 64))

	//a := address2.HexToAddress("410583A68A3BCD86C25AB1BEE482BAC04A216B0261")
	//log.Println(a.String())
	//
	//c := address2.HexToAddress("41a614f803b6fd780986a42c78ec9c7f77e6ded13c")
	//log.Println(c.String())
}

func TestTRC20_Balance(t *testing.T) {
	trc20Contract := "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t" // USDT
	address := "TAUN6FwrnwwmaEqYcckffC7wYmbaS6cBiX"

	conn := client.NewGrpcClient("18.196.99.16:50051")
	conn.SetAPIKey("244f918d-56b5-4a16-9665-9637598b1223")
	err := conn.Start(grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	balance, err := conn.TRC20ContractBalance(address, trc20Contract)

	log.Println(balance.String())

	log.Println(conn.TRC20GetDecimals(trc20Contract))
}

func TestBalance(t *testing.T) {

	conn := client.NewGrpcClient("161.117.83.38:50051")
	conn.SetAPIKey("244f918d-56b5-4a16-9665-9637598b1223")
	err := conn.Start(grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	//addr, _ := address2.Base58ToAddress("TUtAk64jJqdf1pY3SiHeooVikP2SFWXjZ6")
	addr, _ := common.DecodeCheck("TUtAk64jJqdf1pY3SiHeooVikP2SFWXjZ6")
	a := new(core.AccountBalanceRequest)
	a.AccountIdentifier = &core.AccountIdentifier{Address: addr}
	a.BlockIdentifier = &core.BlockBalanceTrace_BlockIdentifier{Hash: []byte("0000000002f1acf091fe4eb7d35794d04e810e0422d49153477a7a6528be20dc"), Number: 49392880}
	r, err := conn.Client.GetAccountBalance(ctx, a)
	if err != nil {
		panic(err)
	}

	log.Println(r.Balance)
}
