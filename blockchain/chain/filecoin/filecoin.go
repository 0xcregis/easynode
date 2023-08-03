package filecoin

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/0xcregis/easynode/blockchain"
	"github.com/tidwall/gjson"
)

type Filecoin struct {
}

func (e *Filecoin) EthSubscribe(host string, token string) (string, error) {
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

func (e *Filecoin) EthUnSubscribe(host string, token string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func NewChainClient() blockchain.BlockChain {
	return &Filecoin{}
}

func (e *Filecoin) SendRequestToChainByHttp(host string, token string, query string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (e *Filecoin) GetTokenBalanceByHttp(host string, token string, contractAddress string, userAddress string) (map[string]interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (e *Filecoin) SendRequestToChain(host string, token string, query string) (string, error) {
	payload := strings.NewReader(query)
	req, err := http.NewRequest("POST", host, payload)
	if err != nil {
		return "", err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("cache-control", "no-cache")
	if len(token) > 1 {
		req.Header.Add("Authorization", "Bearer "+token)
	}

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

func (e *Filecoin) GetTokenBalance(host string, key string, contractAddress string, userAddress string) (map[string]interface{}, error) {
	return nil, nil
}
