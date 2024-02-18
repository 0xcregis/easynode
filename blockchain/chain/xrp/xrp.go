package xrp

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/0xcregis/easynode/blockchain"
	"github.com/tidwall/gjson"
)

type XRP struct {
}

func (e *XRP) GetToken721(host string, token string, contractAddress string, userAddress string) (map[string]interface{}, error) {
	return nil, fmt.Errorf("not implement the method")
}

func (e *XRP) GetToken1155(host string, token string, contractAddress string, userAddress string) (map[string]interface{}, error) {
	return nil, fmt.Errorf("not implement the method")
}

func (e *XRP) Subscribe(host string, token string) (string, error) {
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

func (e *XRP) UnSubscribe(host string, token string) (string, error) {
	return "", fmt.Errorf("not implement the method")
}

func NewChainClient() blockchain.ChainConn {
	return &XRP{}
}

func (e *XRP) SendRequestToChainByHttp(host string, token string, query string) (string, error) {
	return "", fmt.Errorf("not implement the method")
}

func (e *XRP) GetToken20ByHttp(host string, token string, contractAddress string, userAddress string) (map[string]interface{}, error) {
	return nil, fmt.Errorf("not implement the method")
}

func (e *XRP) SendRequestToChain(host string, token string, query string) (string, error) {
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

	if gjson.ParseBytes(body).Get("result.error").Exists() {
		return "", errors.New(string(body))
	}

	return string(body), nil
}

func (e *XRP) GetToken20(host string, key string, contractAddress string, userAddress string) (map[string]interface{}, error) {
	return nil, fmt.Errorf("not implement the method")
}
