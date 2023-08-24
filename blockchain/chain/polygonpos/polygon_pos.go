package polygonpos

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/0xcregis/easynode/blockchain"
	"github.com/0xcregis/easynode/blockchain/chain/token"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tidwall/gjson"
)

type PolygonPos struct {
}

func (e *PolygonPos) Subscribe(host string, token string) (string, error) {
	if len(token) > 1 {
		host = fmt.Sprintf("%v/%v", host, token)
	}

	if strings.HasPrefix(host, "ws") {
		host = strings.ReplaceAll(host, "http", "ws")
	}
	if strings.HasPrefix(host, "wss") {
		host = strings.ReplaceAll(host, "http", "wss")
	}
	return host, nil
}

func (e *PolygonPos) UnSubscribe(host string, token string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func NewChainClient() blockchain.ChainConn {
	return &PolygonPos{}
}

func (e *PolygonPos) SendRequestToChainByHttp(host string, token string, query string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (e *PolygonPos) GetTokenBalanceByHttp(host string, token string, contractAddress string, userAddress string) (map[string]interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (e *PolygonPos) SendRequestToChain(host string, token string, query string) (string, error) {
	if len(token) > 1 {
		host = fmt.Sprintf("%v/%v", host, token)
	}
	payload := strings.NewReader(query)

	req, err := http.NewRequest("POST", host, payload)
	if err != nil {
		return "", err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("cache-control", "no-cache")
	//req.Header.Add("Postman-Token", "181e4572-a9db-453a-b7d4-17974f785de0")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	if gjson.ParseBytes(body).Get("error").Exists() {
		return "", errors.New(string(body))
	}

	return string(body), nil
}

func (e *PolygonPos) GetTokenBalance(host string, key string, contractAddress string, userAddress string) (map[string]interface{}, error) {

	if len(key) > 1 {
		host = fmt.Sprintf("%v/%v", host, key)
	}

	client, err := ethclient.Dial(host)
	if err != nil {
		return nil, err
	}

	tokenAddress := common.HexToAddress(contractAddress)
	instance, err := token.NewToken(tokenAddress, client)
	if err != nil {
		return nil, err
	}

	address := common.HexToAddress(userAddress)
	bal, err := instance.BalanceOf(&bind.CallOpts{Pending: false}, address)
	if err != nil {
		log.Printf("contract:%v,from:%v,err=%v", contractAddress, userAddress, err.Error())
		return nil, err
	}
	mp := make(map[string]interface{}, 2)
	mp["balance"] = bal.String()

	decimals, err := instance.Decimals(&bind.CallOpts{Pending: false, From: address})
	if err != nil {
		log.Printf("contract:%v,from:%v,err=%v", contractAddress, userAddress, err.Error())
		return nil, err
	} else {
		mp["decimals"] = decimals
	}

	//name, err := instance.Name(&bind.CallOpts{})
	//if err != nil {
	//	log.Println("err=", err)
	//} else {
	//	mp["name"] = name
	//}
	//
	//symbol, err := instance.Symbol(&bind.CallOpts{})
	//if err != nil {
	//	log.Println("err=", err)
	//} else {
	//	mp["symbol"] = symbol
	//}

	return mp, nil
}
