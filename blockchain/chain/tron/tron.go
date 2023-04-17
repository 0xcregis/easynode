package tron

import (
	"errors"
	"fmt"
	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/tidwall/gjson"
	"github.com/uduncloud/easynode/blockchain/chain"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Tron struct {
}

func (t *Tron) EthSubscribe(host string, token string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (t *Tron) EthUnSubscribe(host string, token string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func NewChainClient() chain.BlockChain {
	return &Tron{}
}

func (t *Tron) EthSendRequestToChain(host string, token string, query string) (string, error) {

	//host = fmt.Sprintf("%v/%v", host, "jsonrpc")
	payload := strings.NewReader(query)

	req, err := http.NewRequest("POST", host, payload)
	if err != nil {
		return "", err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("cache-control", "no-cache")
	if len(token) > 1 {
		req.Header.Add("TRON_PRO_API_KEY", token)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	if gjson.ParseBytes(body).Get("error").Exists() {
		return "", errors.New(string(body))
	}

	return string(body), nil
}

func (t *Tron) SendRequestToChainByHttp(host string, token string, query string) (string, error) {
	payload := strings.NewReader(query)

	query = strings.Replace(query, "\t", "", -1)
	query = strings.Replace(query, "\n", "", -1)
	req, err := http.NewRequest("POST", host, payload)
	if err != nil {
		return "", err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("cache-control", "no-cache")
	if len(token) > 1 {
		req.Header.Add("TRON_PRO_API_KEY", token)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	if gjson.ParseBytes(body).Get("Error").Exists() {
		return "", errors.New(string(body))
	}

	return string(body), nil
}

func (t *Tron) GetTokenBalanceByHttp(host string, token string, contractAddress string, userAddress string) (map[string]interface{}, error) {

	var query string
	mp := make(map[string]interface{}, 2)

	query = `
			{
			  "owner_address": "%v",
			  "contract_address": "%v",
			  "function_selector": "balanceOf(address)",
			  "parameter": "%v"
			}
			`

	userAddr, err := address.Base58ToAddress(userAddress)
	if err != nil {
		return nil, err
	}

	contractAddr, err := address.Base58ToAddress(contractAddress)
	if err != nil {
		return nil, err
	}

	c2 := userAddr.Hex()[4:]
	m := "0000000000000000000000000000000000000000000000000000000000000000"
	params := m[:len(m)-len(c2)] + c2
	query = fmt.Sprintf(query, userAddr.Hex(), contractAddr.Hex(), params)
	resp, err := t.SendRequestToChainByHttp(host, token, query)
	if err != nil {
		return nil, err
	}
	log.Println(resp)
	r := gjson.Parse(resp).Get("constant_result")
	if r.Exists() {
		balance, _ := strconv.ParseInt(r.Array()[0].String(), 16, 64)
		mp["balance"] = balance
	} else {
		return nil, errors.New(resp)
	}

	return mp, nil
}

func (t *Tron) GetTokenBalance(host string, key string, contractAddress string, userAddress string) (map[string]interface{}, error) {

	//host = "grpc.trongrid.io:50051"
	conn := client.NewGrpcClient(host)
	_ = conn.SetAPIKey(key) // todo 没有发现设置意义
	err := conn.Start(grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	mp := make(map[string]interface{}, 2)
	balance, err := conn.TRC20ContractBalance(userAddress, contractAddress)

	if err != nil {
		log.Println("err=", err)
	} else {
		mp["balance"] = balance.String()
	}

	name, err := conn.TRC20GetName(contractAddress)
	if err != nil {
		log.Println("err=", err)
	} else {
		mp["name"] = name
	}

	symbol, err := conn.TRC20GetSymbol(contractAddress)
	if err != nil {
		log.Println("err=", err)
	} else {
		mp["symbol"] = symbol
	}

	decimals, err := conn.TRC20GetDecimals(contractAddress)
	if err != nil {
		log.Println("err=", err)
	} else {
		mp["decimals"] = decimals
	}
	return mp, nil
}
