package ether

import (
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Eth_GetBlockNumber 获取最新区块高度
func Eth_GetBlockNumber(host string, token string) (string, error) {
	if len(token) > 1 {
		host = fmt.Sprintf("%v/%v", host, token)
	}
	query := `
{
    "id": 1,
    "jsonrpc": "2.0",
    "method": "eth_blockNumber"
}
`
	return send(host, query)
}

func send(host string, query string) (string, error) {
	payload := strings.NewReader(query)

	req, err := http.NewRequest("POST", host, payload)
	if err != nil {
		return "", err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("cache-control", "no-cache")
	//req.Header.Add("Postman-Token", "181e4572-a9db-453a-b7d4-17974f785de0")

	client := http.Client{Timeout: 40 * time.Second}
	res, err := client.Do(req)
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
