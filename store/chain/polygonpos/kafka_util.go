package polygonpos

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/0xcregis/easynode/common/util"
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
	root := gjson.ParseBytes(body)
	input := root.Get("input").String()
	if len(input) > 5 {
		//合约调用
		return 1, nil
	} else {
		//普通资产转移
		return 2, nil
	}
}

func ParseTx(body []byte, transferTopic string, blocchain int64) (*store.SubTx, error) {
	var r store.SubTx
	root := gjson.ParseBytes(body)
	r.BlockChain = uint64(blocchain)
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
	value := root.Get("value").String()
	v, err := util.HexToInt(value)
	if err != nil {
		return nil, err
	}
	r.Value = util.Div(v, 18)

	tp, err := GetTxType(body)
	if err != nil {
		return nil, err
	}
	r.TxType = tp

	txTime := root.Get("txTime").String()
	txTime = fmt.Sprintf("%v000", txTime)
	r.TxTime = txTime

	gasPrice := root.Get("gasPrice").String() //单位：wei
	//price, err := util.HexToInt(gasPrice)
	//if err != nil {
	//	return nil, err
	//}

	if !root.Get("receipt").Exists() { //收据不存在的交易，则放弃
		return nil, errors.New("receipt is error")
	}

	receipt := root.Get("receipt").String()
	receiptRoot := gjson.Parse(receipt)

	gasUsed := receiptRoot.Get("gasUsed").String()
	//gas, _ := util.HexToInt2(gasUsed)

	bigPrice, b := new(big.Int).SetString(gasPrice, 0)
	bigGas, b2 := new(big.Int).SetString(gasUsed, 0)

	if b && b2 {
		fee := bigPrice.Mul(bigPrice, bigGas)
		r.Fee = util.Div(fmt.Sprintf("%v", fee), 18)
		r.FeeDetail = map[string]string{"gasPrice": fmt.Sprintf("%v", bigPrice.String()), "gasUsed": fmt.Sprintf("%v", bigGas.String())}
	} else {
		return nil, errors.New("price or gas is wrong")
	}

	status := receiptRoot.Get("status").String() //0x0:失败，0x1:成功
	if status == "0x0" {
		r.Status = 0
	} else if status == "0x1" {
		r.Status = 1
	}

	contractTx := make([]*store.ContractTx, 0, 5)
	if r.TxType != 1 {
		var c store.ContractTx
		c.Contract = ""
		c.Value = r.Value
		c.From = r.From
		c.To = r.To
		c.Index = 0
		c.Method = "Transfer"
		c.EIP = -1
		c.Token = ""
		contractTx = append(contractTx, &c)
	}

	logs := receiptRoot.Get("logs").Array()
	l := len(contractTx)
	for index, v := range logs {
		contract := v.Get("address").String()
		data := v.Get("data").String()
		r := gjson.Parse(data)
		if r.IsObject() {
			//mp := make(map[string]string, 2)

			eip := r.Get("eip").Int()
			token := r.Get("token").String()

			if eip == 20 {
				contractDecimals := r.Get("contractDecimals").String()
				if len(contractDecimals) < 1 {
					return nil, errors.New("tx.log.contract is error")
				}
				decimals, err := strconv.Atoi(contractDecimals)
				if err != nil {
					return nil, errors.New("tx.log.contract is error")
				}

				fee := r.Get("data").String()
				bigFee, err := util.HexToInt(fee)
				if err == nil {
					data = util.Div(bigFee, decimals)
				} else {
					return nil, errors.New("tx.log.contract is error")
				}

				tps := v.Get("topics").Array()
				var from, to string
				if len(tps) == 3 && tps[0].String() == transferTopic {
					//method = tps[0].String()
					from = tps[1].String()
					to = tps[2].String()
					var m store.ContractTx
					m.Contract = contract
					m.Value = data
					m.Index = int64(index + l)
					m.From, _ = util.Hex2Address(from)
					m.To, _ = util.Hex2Address(to)
					m.Method = "Transfer"
					m.EIP = eip
					m.Token = token
					contractTx = append(contractTx, &m)
				}

			}

			if eip == 721 {
				tps := v.Get("topics").Array()
				var from, to string
				if len(tps) == 4 && tps[0].String() == transferTopic {
					//method = tps[0].String()
					from = tps[1].String()
					to = tps[2].String()
					var m store.ContractTx
					m.Contract = contract
					v, _ := util.HexToInt(tps[3].String())
					m.Value = v
					m.Index = int64(index + l)
					m.From, _ = util.Hex2Address(from)
					m.To, _ = util.Hex2Address(to)
					m.Method = "Transfer"
					m.EIP = eip
					m.Token = token
					contractTx = append(contractTx, &m)
				}
			}

		}
		//else {
		// 合约数据不完整，直接放弃了
		//}
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
func CheckAddress(tx []byte, addrList map[string]*store.MonitorAddress, transferTopic string) bool {

	if len(addrList) < 1 || len(tx) < 1 {
		return false
	}
	txAddressList := make(map[string]int64, 10)
	root := gjson.ParseBytes(tx)

	fromAddr := root.Get("from").String()
	txAddressList[GetCoreAddr(fromAddr)] = 1

	toAddr := root.Get("to").String()
	txAddressList[GetCoreAddr(toAddr)] = 1

	if root.Get("receipt").Exists() {
		receipt := root.Get("receipt").String()
		receiptRoot := gjson.Parse(receipt)
		list := receiptRoot.Get("logs").Array()
		for _, v := range list {
			//过滤没有合约信息的交易，出现这种情况原因：1. 合约获取失败会重试 2:非20合约
			data := v.Get("data").String()
			if !gjson.Parse(data).IsObject() {
				continue
			}
			topics := v.Get("topics").Array()
			//Transfer(),erc20
			if len(topics) == 3 && topics[0].String() == transferTopic {
				from, _ := util.Hex2Address(topics[1].String())
				if len(from) > 0 {
					txAddressList[GetCoreAddr(from)] = 1
				}
				to, _ := util.Hex2Address(topics[2].String())
				if len(to) > 0 {
					txAddressList[GetCoreAddr(to)] = 1
				}
			}
			//Transfer(),erc721
			if len(topics) == 4 && topics[0].String() == transferTopic {
				from, _ := util.Hex2Address(topics[1].String())
				if len(from) > 0 {
					txAddressList[GetCoreAddr(from)] = 1
				}
				to, _ := util.Hex2Address(topics[2].String())
				if len(to) > 0 {
					txAddressList[GetCoreAddr(to)] = 1
				}
			}
		}
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
