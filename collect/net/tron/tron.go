package tron

import (
	"errors"
	"fmt"
	"github.com/sunjiangjun/xlog"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func Eth_GetBlockByHash(host string, token string, blockHash string, log *xlog.XLog) (string, error) {

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
        true
    ]
}
`
	query = fmt.Sprintf(query, blockHash)
	return send(host, token, query)
}

// Eth_GetBlockByNumber number: hex value of block number
func Eth_GetBlockByNumber(host string, token string, number string, log *xlog.XLog) (string, error) {
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
        true
    ]
}
`
	query = fmt.Sprintf(query, number)
	return send(host, token, query)
}

// Eth_GetTransactionByHash number: hex value of block number
func Eth_GetTransactionByHash(host string, token string, hash string, log *xlog.XLog) (string, error) {

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
	return send(host, token, query)
}

// Eth_GetTransactionReceiptByHash number: hex value of block number
func Eth_GetTransactionReceiptByHash(host string, token string, hash string, log *xlog.XLog) (string, error) {

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
	return send(host, token, query)
}

// Eth_GetBlockReceiptByBlockHash number: hex value of block number
func Eth_GetBlockReceiptByBlockHash(host string, token string, hash string, log *xlog.XLog) (string, error) {

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
	return send(host, token, query)
}

// Eth_GetBlockReceiptByBlockNumber number: hex value of block number
func Eth_GetBlockReceiptByBlockNumber(host string, token string, number string, log *xlog.XLog) (string, error) {

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
	return send(host, token, query)
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
