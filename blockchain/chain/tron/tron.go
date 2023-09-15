package tron

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/0xcregis/easynode/blockchain"
	"github.com/0xcregis/easynode/common/util"
	"github.com/tidwall/gjson"
)

type Tron struct {
}

func (t *Tron) GetToken721(host string, token string, contractAddress string, userAddress string) (map[string]interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (t *Tron) GetToken1155(host string, token string, contractAddress string, userAddress string) (map[string]interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (t *Tron) Subscribe(host string, token string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (t *Tron) UnSubscribe(host string, token string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func NewChainClient() blockchain.ChainConn {
	return &Tron{}
}

func (t *Tron) SendRequestToChain(host string, token string, query string) (string, error) {
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
	body, err := io.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	if gjson.ParseBytes(body).Get("error").Exists() {
		return "", errors.New(string(body))
	}

	return string(body), nil
}

func (t *Tron) SendRequestToChainByHttp(host string, token string, query string) (string, error) {
	query = strings.Replace(query, "\t", "", -1)
	query = strings.Replace(query, "\n", "", -1)
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
	body, err := io.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	if gjson.ParseBytes(body).Get("Error").Exists() {
		return "", errors.New(string(body))
	}

	return string(body), nil
}

func (t *Tron) GetToken20ByHttp(host string, token string, contractAddress string, userAddress string) (map[string]interface{}, error) {
	mp := make(map[string]interface{}, 2)

	if !strings.HasPrefix(userAddress, "0x41") && !strings.HasPrefix(userAddress, "41") {
		userAddr, err := util.Base58ToAddress(userAddress)
		if err != nil {
			return nil, err
		}
		userAddress = userAddr.Hex()
	}

	if !strings.HasPrefix(contractAddress, "0x41") && !strings.HasPrefix(contractAddress, "41") {
		contractAddr, err := util.Base58ToAddress(contractAddress)
		if err != nil {
			return nil, err
		}
		contractAddress = contractAddr.Hex()
	}

	balance, err := t.GetTokenBalanceByHttp2(host, token, contractAddress, userAddress)
	if err != nil {
		return nil, err
	}
	mp["balance"] = balance

	decimal, err := t.GetTokenDecimalsByHttp(host, token, contractAddress, userAddress)
	if err != nil {
		return nil, err
	}
	mp["decimals"] = decimal

	if (len(balance) == 0 || balance == "0") && (len(decimal) == 0 || decimal == "0") {
		return nil, errors.New("contract is error")
	}

	return mp, nil
}

func (t *Tron) GetTokenDecimalsByHttp(host string, token string, contractAddress string, userAddress string) (string, error) {
	//var query string
	query := `
			{
			  "owner_address": "%v",
			  "contract_address": "%v",
			  "function_selector": "decimals()"
			}
			`
	query = fmt.Sprintf(query, userAddress, contractAddress)
	resp, err := t.SendRequestToChainByHttp(host, token, query)
	if err != nil {
		return "", err
	}
	//log.Println(resp)
	r := gjson.Parse(resp).Get("constant_result")
	if r.Exists() {
		decimals, _ := strconv.ParseInt(r.Array()[0].String(), 16, 64)
		return fmt.Sprintf("%v", decimals), nil
	}

	return "", errors.New("no data")
}

func (t *Tron) GetTokenBalanceByHttp2(host string, token string, contractAddress string, userAddress string) (string, error) {

	var query string
	query = `
			{
			  "owner_address": "%v",
			  "contract_address": "%v",
			  "function_selector": "balanceOf(address)",
			  "parameter": "%v"
			}
			`
	var c2 string
	if strings.HasPrefix(userAddress, "41") {
		c2 = userAddress[2:]
	} else if strings.HasPrefix(userAddress, "0x41") {
		c2 = userAddress[4:]
	}

	m := "0000000000000000000000000000000000000000000000000000000000000000"
	params := m[:len(m)-len(c2)] + c2
	query = fmt.Sprintf(query, userAddress, contractAddress, params)
	resp, err := t.SendRequestToChainByHttp(host, token, query)
	if err != nil {
		return "", err
	}
	//log.Println(resp)
	r := gjson.Parse(resp).Get("constant_result")
	if r.Exists() {
		balance, _ := strconv.ParseInt(r.Array()[0].String(), 16, 64)
		return fmt.Sprintf("%v", balance), nil
	}

	return "", errors.New("no data")
}

func (t *Tron) GetToken20(host string, key string, contractAddress string, userAddress string) (map[string]interface{}, error) {
	return nil, nil
}
