package tron

import (
	"errors"
	"fmt"
	"github.com/sunjiangjun/xlog"
	"github.com/tidwall/gjson"
	chainService "github.com/uduncloud/easynode/blockchain/service"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func Eth_GetCode(client chainService.API, chainCode int64, address string) (string, error) {
	query := `{
				"id": 1,
				"jsonrpc": "2.0",
				"method": "eth_getCode",
				"params": [
					"%v",
					"latest"
				]
			}`
	query = fmt.Sprintf(query, address)
	return client.SendJsonRpc(chainCode, query)
}

func Eth_GetAddressType(client chainService.API, chainCode int64, address string) (string, error) {
	query := `{
				"id": 1,
				"jsonrpc": "2.0",
				"method": "eth_getCode",
				"params": [
					"%v",
					"latest"
				]
			}`
	query = fmt.Sprintf(query, address)
	resp, err := client.SendJsonRpc(chainCode, query)
	if err != nil {
		return "", err
	}

	code := gjson.Parse(resp).Get("result").String()
	if len(code) > 5 {
		//合约地址
		return "0x12", nil
	} else {
		//外部地址
		return "0x11", nil
	}
}

func Eth_GetBlockByHash(client chainService.API, blockHash string, log *xlog.XLog,flag bool) (string, error) {

	start := time.Now()
	defer func() {
		log.Printf("Eth_GetBlockByHash | duration=%v", time.Now().Sub(start))
	}()

	//host = fmt.Sprintf("%v/%v", host, "jsonrpc")
	query := `
{
    "id": 1,
    "jsonrpc": "2.0",
    "method": "eth_getBlockByHash",
    "params": [
        "%v",
        %v
    ]
}
`
	query = fmt.Sprintf(query, blockHash,flag)
	return client.SendJsonRpc(205, query)
	//return send(host, token, query)
}

// Eth_GetBlockByNumber number: hex value of block number
func Eth_GetBlockByNumber(client chainService.API, number string, log *xlog.XLog,flag bool) (string, error) {
	start := time.Now()
	defer func() {
		log.Printf("Eth_GetBlockByNumber | duration=%v", time.Now().Sub(start))
	}()

	//host = fmt.Sprintf("%v/%v", host, "jsonrpc")
	query := `
{
    "id": 1,
    "jsonrpc": "2.0",
    "method": "eth_getBlockByNumber",
    "params": [
        "%v",
        %v
    ]
}
`
	query = fmt.Sprintf(query, number,flag)
	return client.SendJsonRpc(205, query)
	//return send(host, token, query)
}

// Eth_GetTransactionByHash number: hex value of block number
func Eth_GetTransactionByHash(client chainService.API, hash string, log *xlog.XLog) (string, error) {

	start := time.Now()
	defer func() {
		log.Printf("Eth_GetTransactionByHash | duration=%v", time.Now().Sub(start))
	}()

	//host = fmt.Sprintf("%v/%v", host, "jsonrpc")
	query := `
{
    "id": 1,
    "jsonrpc": "2.0",
    "method": "eth_getTransactionByHash",
    "params": [
        "%v"
    ]
}
`

	query = fmt.Sprintf(query, hash)

	return client.SendJsonRpc(205, query)
	//return send(host, token, query)
}

// Eth_GetTransactionReceiptByHash number: hex value of block number
func Eth_GetTransactionReceiptByHash(client chainService.API, hash string, log *xlog.XLog) (string, error) {

	start := time.Now()
	defer func() {
		log.Printf("Eth_GetTransactionReceiptByHash | duration=%v", time.Now().Sub(start))
	}()

	//host = fmt.Sprintf("%v/%v", host, "jsonrpc")
	query := `
{
    "id": 1,
    "jsonrpc": "2.0",
    "method": "eth_getTransactionReceipt",
    "params": [
        "%v"
    ]
}
`
	query = fmt.Sprintf(query, hash)
	return client.SendJsonRpc(205, query)
	//return send(host, token, query)
}

// Eth_GetBlockReceiptByBlockHash number: hex value of block number
func Eth_GetBlockReceiptByBlockHash(client chainService.API, hash string, log *xlog.XLog) (string, error) {

	start := time.Now()
	defer func() {
		log.Printf("Eth_GetBlockReceiptByBlockHash | duration=%v", time.Now().Sub(start))
	}()

	//host = fmt.Sprintf("%v/%v", host, "jsonrpc")
	query := `
{
    "id": 1,
    "jsonrpc": "2.0",
    "method": "eth_getBlockReceipts",
    "params": [
        "%v"
    ]
}
`
	query = fmt.Sprintf(query, hash)
	return client.SendJsonRpc(205, query)
	//return send(host, token, query)
}

// Eth_GetBlockReceiptByBlockNumber number: hex value of block number
func Eth_GetBlockReceiptByBlockNumber(client chainService.API, number string, log *xlog.XLog) (string, error) {

	start := time.Now()
	defer func() {
		log.Printf("Eth_GetBlockReceiptByBlockNumber | duration=%v", time.Now().Sub(start))
	}()

	//host = fmt.Sprintf("%v/%v", host, "jsonrpc")
	query := `
{
    "id": 1,
    "jsonrpc": "2.0",
    "method": "eth_getBlockReceipts",
    "params": [
        "%v"
    ]
}
`

	query = fmt.Sprintf(query, number)
	return client.SendJsonRpc(205, query)
	//return send(host, token, query)
}

func send(host, token string, query string) (string, error) {
	payload := strings.NewReader(query)

	req, err := http.NewRequest("POST", host, payload)
	if err != nil {
		return "", err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("TRON_PRO_API_KEY", token)

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
