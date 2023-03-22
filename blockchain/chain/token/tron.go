package token

import (
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"google.golang.org/grpc"
	"log"
)

func BalanceOf() {
	trc20Contract := "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t" // USDT
	address := "TMuA6YqfCeX8EhbfYEg5y7S4DqzSJireY9"

	conn := client.NewGrpcClient("grpc.trongrid.io:50051")
	err := conn.Start(grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	balance, err := conn.TRC20ContractBalance(address, trc20Contract)

	log.Println(balance.String())
}
