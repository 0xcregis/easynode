package ether

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"strings"

	"github.com/0xcregis/easynode/blockchain"
	"github.com/0xcregis/easynode/blockchain/chain/token"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tidwall/gjson"
)

type Ether struct {
}

func (e *Ether) TokenURI(host string, key string, contractAddress string, tokenId string, eip int64) (string, error) {

	if len(key) > 1 {
		host = fmt.Sprintf("%v/%v", host, key)
	}

	client, err := ethclient.Dial(host)
	if err != nil {
		return "", err
	}

	tokenAddress := common.HexToAddress(contractAddress)
	if eip == 721 {
		instance, err := token.NewToken721(tokenAddress, client)
		if err != nil {
			return "", err
		}

		i, b := new(big.Int).SetString(tokenId, 0)

		if !b {
			return "", fmt.Errorf("tokenId is error,tokenId=%v", tokenId)
		}

		return instance.TokenURI(&bind.CallOpts{Pending: false}, i)

	} else if eip == 1155 {
		instance, err := token.NewToken1155(tokenAddress, client)
		if err != nil {
			return "", err
		}

		i, b := new(big.Int).SetString(tokenId, 0)

		if !b {
			return "", fmt.Errorf("tokenId is error,tokenId=%v", tokenId)
		}

		return instance.Uri(&bind.CallOpts{Pending: false}, i)
	} else {
		return "", fmt.Errorf("unknow eip:%v", eip)
	}
}

func (e *Ether) BalanceOf(host string, key string, contractAddress string, address string, tokenId string, eip int64) (string, error) {
	if len(key) > 1 {
		host = fmt.Sprintf("%v/%v", host, key)
	}

	client, err := ethclient.Dial(host)
	if err != nil {
		return "", err
	}

	tokenAddress := common.HexToAddress(contractAddress)
	if eip == 721 {
		instance, err := token.NewToken721(tokenAddress, client)
		if err != nil {
			return "", err
		}

		count, err := instance.BalanceOf(&bind.CallOpts{Pending: false}, common.HexToAddress(address))
		if err == nil {
			return count.String(), nil
		} else {
			return "", err
		}

	} else if eip == 1155 {
		instance, err := token.NewToken1155(tokenAddress, client)
		if err != nil {
			return "", err
		}

		i, b := new(big.Int).SetString(tokenId, 0)
		if !b {
			return "", fmt.Errorf("tokenId is error,tokenId=%v", tokenId)
		}
		count, err := instance.BalanceOf(&bind.CallOpts{Pending: false}, common.HexToAddress(address), i)
		if err == nil {
			return count.String(), nil
		} else {
			return "", err
		}
	} else {
		return "", fmt.Errorf("unknow eip:%v", eip)
	}
}

func (e *Ether) OwnerOf(host string, key string, contractAddress string, tokenId string, eip int64) (string, error) {
	if len(key) > 1 {
		host = fmt.Sprintf("%v/%v", host, key)
	}

	client, err := ethclient.Dial(host)
	if err != nil {
		return "", err
	}

	tokenAddress := common.HexToAddress(contractAddress)
	if eip == 721 {
		instance, err := token.NewToken721(tokenAddress, client)
		if err != nil {
			return "", err
		}
		i, b := new(big.Int).SetString(tokenId, 0)
		if !b {
			return "", fmt.Errorf("tokenId is error,tokenId=%v", tokenId)
		}
		addr, err := instance.OwnerOf(&bind.CallOpts{Pending: false}, i)
		if err == nil {
			return addr.String(), nil
		} else {
			return "", err
		}

	} else {
		return "", fmt.Errorf("unknow eip:%v", eip)
	}
}

func (e *Ether) TotalSupply(host string, key string, contractAddress string, eip int64) (string, error) {
	if len(key) > 1 {
		host = fmt.Sprintf("%v/%v", host, key)
	}

	client, err := ethclient.Dial(host)
	if err != nil {
		return "", err
	}

	tokenAddress := common.HexToAddress(contractAddress)
	if eip == 721 {
		instance, err := token.NewToken721(tokenAddress, client)
		if err != nil {
			return "", err
		}

		count, err := instance.TotalSupply(&bind.CallOpts{Pending: false})
		if err == nil {
			return count.String(), nil
		} else {
			return "", err
		}

	} else {
		return "", fmt.Errorf("unknow eip:%v", eip)
	}
}

func (e *Ether) GetToken721(host string, key string, contractAddress string, userAddress string) (map[string]interface{}, error) {
	mp := make(map[string]interface{}, 2)
	if len(key) > 1 {
		host = fmt.Sprintf("%v/%v", host, key)
	}

	client, err := ethclient.Dial(host)
	if err != nil {
		return nil, err
	}

	tokenAddress := common.HexToAddress(contractAddress)
	instance, err := token.NewToken721(tokenAddress, client)
	if err != nil {
		return nil, err
	}

	address := common.HexToAddress(userAddress)
	name, err := instance.Name(&bind.CallOpts{Pending: false, From: address})
	if err != nil {
		log.Printf("contract:%v,from:%v,err=%v", contractAddress, userAddress, err.Error())
		return nil, err
	}
	mp["name"] = name

	symbol, err := instance.Symbol(&bind.CallOpts{Pending: false, From: address})
	if err != nil {
		log.Printf("contract:%v,from:%v,err=%v", contractAddress, userAddress, err.Error())
		return nil, err
	} else {
		mp["symbol"] = symbol
	}

	//balance, err := instance.BalanceOf(&bind.CallOpts{Pending: false, From: address}, address)
	//if err == nil {
	//	mp["balance"] = balance.String()
	//}
	return mp, nil
}

func (e *Ether) GetToken1155(host string, key string, contractAddress string, userAddress string) (map[string]interface{}, error) {
	mp := make(map[string]interface{}, 2)
	//if len(key) > 1 {
	//	host = fmt.Sprintf("%v/%v", host, key)
	//}

	//client, err := ethclient.Dial(host)
	//if err != nil {
	//	return nil, err
	//}

	//tokenAddress := common.HexToAddress(contractAddress)
	//instance, err := token.NewToken1155(tokenAddress, client)
	//if err != nil {
	//	return nil, err
	//}
	//
	//address := common.HexToAddress(userAddress)
	//name, err := instance.Uri(&bind.CallOpts{Pending: false, From: address})
	//if err != nil {
	//	log.Printf("contract:%v,from:%v,err=%v", contractAddress, userAddress, err.Error())
	//	return nil, err
	//}
	//mp["name"] = name
	//
	//symbol, err := instance.Symbol(&bind.CallOpts{Pending: false, From: address})
	//if err != nil {
	//	log.Printf("contract:%v,from:%v,err=%v", contractAddress, userAddress, err.Error())
	//	return nil, err
	//} else {
	//	mp["symbol"] = symbol
	//}

	return mp, nil
}

func (e *Ether) Subscribe(host string, token string) (string, error) {
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

func (e *Ether) UnSubscribe(host string, token string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func NewChainClient() blockchain.ChainConn {
	return &Ether{}
}

func NewNFTClient() blockchain.NFT {
	return &Ether{}
}

func (e *Ether) SendRequestToChainByHttp(host string, token string, query string) (string, error) {
	return "", fmt.Errorf("not implement the method")
}

func (e *Ether) GetToken20ByHttp(host string, token string, contractAddress string, userAddress string) (map[string]interface{}, error) {
	return nil, fmt.Errorf("not implement the method")
}

func (e *Ether) SendRequestToChain(host string, token string, query string) (string, error) {
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

func (e *Ether) GetToken20(host string, key string, contractAddress string, userAddress string) (map[string]interface{}, error) {

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

	name, err := instance.Name(&bind.CallOpts{Pending: false, From: address})
	if err != nil {
		log.Println("err=", err)
	} else {
		mp["name"] = name
	}

	symbol, err := instance.Symbol(&bind.CallOpts{Pending: false, From: address})
	if err != nil {
		log.Println("err=", err)
	} else {
		mp["symbol"] = symbol
	}

	return mp, nil
}
