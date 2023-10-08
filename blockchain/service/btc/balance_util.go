package btc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/tidwall/gjson"
)

/**

{
                   "txid": "f030614a5a1430075624b9cfec8234996b04300f9fefb1e471437d64fd931fc2",
                   "height": "810920",
                   "blockTime": "1696603869",
                   "address": "1PL6qjNjEMRhTLAnHEFJWwnvjjKGWAwFws",
                   "unspentAmount": "1338.97117226",
                   "index": "1"
               }
*/

type UTXO struct {
	Address       string `json:"address" gorm:"column:address"`
	TxId          string `json:"txid" gorm:"column:txid"`
	Index         string `json:"index" gorm:"column:index"`
	BlockTime     string `json:"blockTime" gorm:"column:blockTime"`
	UnspentAmount string `json:"unspentAmount" gorm:"column:unspentAmount"`
	Height        string `json:"height" gorm:"column:height"`
}

func GetBalanceEx(uris []string, address string) ([]*UTXO, error) {
	array := make([][]*UTXO, 0)
	for _, v := range uris {
		if strings.Contains(v, "#") {
			ss := strings.Split(v, "#")
			list, err := GetBalance(ss[0], ss[1], address)
			if err == nil && len(list) > 0 {
				array = append(array, list)
			}
		}
	}

	if len(array) == 0 {
		return nil, nil
	}
	if len(array) == 1 {
		return array[0], nil
	}

	//todo 判断是否存在作恶节点
	return array[1], nil

}

func GetBalance(uri, token string, address string) ([]*UTXO, error) {
	//https://www.oklink.com/api/v5/explorer/address/utxo
	url := "%v?chainShortName=btc&address=%v&page=1&limit=100"
	url = fmt.Sprintf(url, uri, address)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Content-Type", "application/json")
	if len(token) > 1 {
		req.Header.Add("Ok-Access-Key", "a6938ee9-4678-4cd1-90ca-ba13ee472ede")
	}
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	root := gjson.ParseBytes(body)
	if root.Get("code").String() != "0" {
		return nil, fmt.Errorf("%s", body)
	}

	utxo := root.Get("data.0.utxoList").String()

	var r []*UTXO
	_ = json.Unmarshal([]byte(utxo), &r)
	return r, nil
}
