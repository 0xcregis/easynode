package btc

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/0xcregis/easynode/store"
	"github.com/tidwall/gjson"
)

func GetReceiptFromKafka(value []byte) (*store.Receipt, error) {
	var receipt store.Receipt
	err := json.Unmarshal(value, &receipt)
	if err != nil {
		return nil, err
	}

	return &receipt, nil
}

func GetBlockFromKafka(value []byte) (*store.Block, error) {
	var block store.Block
	err := json.Unmarshal(value, &block)
	if err != nil {
		return nil, err
	}

	return &block, nil
}

func GetTxFromKafka(value []byte) (*store.Tx, error) {
	var tx store.Tx
	err := json.Unmarshal(value, &tx)
	if err != nil {
		//s.log.Errorf("ReadTxFromKafka|blockChain=%v,error=%v", blockChain, err.Error())
		//continue
		return nil, err
	}
	return &tx, nil
}

func GetTxType(body []byte) (uint64, error) {
	//普通资产转移
	return 2, nil
}
func ParseTx(body []byte, blockchain int64) (*store.SubTx, error) {
	var r store.SubTx
	root := gjson.ParseBytes(body)
	r.BlockChain = uint64(blockchain)
	r.Id = uint64(time.Now().UnixNano())
	blockHash := root.Get("blockHash").String()
	r.BlockHash = blockHash
	blockNumber := root.Get("blockNumber").String()
	r.BlockNumber = blockNumber
	hash := root.Get("hash").String()
	r.Hash = hash
	from := root.Get("from").String()
	r.From = from
	to := root.Get("to").String()
	r.To = to
	input := root.Get("input").String()
	r.Input = input
	r.Value = root.Get("value").String()

	tp, err := GetTxType(body)
	if err != nil {
		return nil, err
	}
	r.TxType = tp

	txTime := root.Get("txTime").String()
	txTime = fmt.Sprintf("%v000", txTime)
	r.TxTime = txTime

	r.Fee = root.Get("fee").String() //单位：wei

	//r.Status = root.Get("status").Int() //0x0:失败，0x1:成功
	//r.ContractTx = contractTx

	contractTx := make([]*store.ContractTx, 0, 5)
	if r.TxType != 1 {
		var c store.ContractTx
		c.Contract = ""
		c.Value = r.Value
		c.From = r.From
		c.To = r.To
		c.Method = "Transfer"
		c.EIP = -1
		c.Index = 0
		c.Token = ""
		contractTx = append(contractTx, &c)
	}
	r.ContractTx = contractTx
	return &r, nil
}
func GetCoreAddr(addr string) string {
	addr = strings.ToLower(addr)
	if strings.HasPrefix(addr, "0x") {
		return strings.Replace(addr, "0x", "", 1) //去丢0x
	}
	return addr
}
func CheckAddress(tx []byte, addrList map[string]*store.MonitorAddress) bool {

	if len(addrList) < 1 || len(tx) < 1 {
		return false
	}
	txAddressList := make(map[string]int64, 10)
	root := gjson.ParseBytes(tx)

	fromAddr := root.Get("from").String()
	fromList := gjson.Parse(fromAddr).Array()
	for _, v := range fromList {
		addr := v.Get("prevout.scriptPubKey.address").String()
		txAddressList[GetCoreAddr(addr)] = 1
	}

	toAddr := root.Get("to").String()
	toList := gjson.Parse(toAddr).Array()
	for _, v := range toList {
		addr := v.Get("scriptPubKey.address").String()
		txAddressList[GetCoreAddr(addr)] = 1
	}

	has := false
	for k := range txAddressList {
		//monitorAddr := getCoreAddrEth(v)
		if _, ok := addrList[k]; ok {
			has = true
			break
		}
	}
	return has
}
