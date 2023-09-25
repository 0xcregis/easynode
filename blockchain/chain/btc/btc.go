package btc

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/0xcregis/easynode/blockchain"
	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/tidwall/gjson"
)

type Btc struct {
	lru *expirable.LRU[string, string]
}

func (e *Btc) GetToken721(host string, key string, contractAddress string, userAddress string) (map[string]interface{}, error) {
	mp := make(map[string]interface{}, 2)
	return mp, nil
}

func (e *Btc) GetToken1155(host string, key string, contractAddress string, userAddress string) (map[string]interface{}, error) {
	mp := make(map[string]interface{}, 2)
	return mp, nil
}

func (e *Btc) Subscribe(host string, token string) (string, error) {
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

func (e *Btc) UnSubscribe(host string, token string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func NewChainClient() blockchain.ChainConn {
	cache := expirable.NewLRU[string, string](100, nil, time.Minute*5)
	return &Btc{
		lru: cache,
	}
}

func (e *Btc) SendRequestToChainByHttp(host string, token string, query string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (e *Btc) GetToken20ByHttp(host string, token string, contractAddress string, userAddress string) (map[string]interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (e *Btc) SendRequestToChain(host string, token string, query string) (string, error) {
	if len(token) > 1 {
		host = fmt.Sprintf("%v/%v", host, token)
	}

	cacheOK := gjson.Parse(query).Get("cache").Bool()

	//return from cache
	var keyCache string
	if e.lru != nil && cacheOK {
		hash := md5.Sum([]byte(query))
		keyCache = hex.EncodeToString(hash[:])
		valueCache, ok := e.lru.Get(keyCache)
		if ok {
			return valueCache, nil
		}
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

	r := gjson.ParseBytes(body)
	if r.Get("error").Exists() {
		if r.Get("error").String() != "" {
			return "", errors.New(r.Get("error").String())
		}
	}

	value := string(body)
	if e.lru != nil && cacheOK {
		e.lru.Add(keyCache, value)
	}
	return value, nil
}

func (e *Btc) GetToken20(host string, key string, contractAddress string, userAddress string) (map[string]interface{}, error) {
	return nil, nil
}
