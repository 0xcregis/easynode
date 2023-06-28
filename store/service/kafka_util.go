package service

import (
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/tidwall/gjson"
	"github.com/uduncloud/easynode/collect/service"
	"github.com/uduncloud/easynode/common/util"
	"math/big"
	"strconv"
	"strings"
	"time"
)

func ParseTx(blockchain int64, msg *kafka.Message) (*SubTx, error) {
	if blockchain == 200 {
		return ParseTxForEther(msg.Value)
	}
	if blockchain == 205 {
		return ParseTxForTron(msg.Value)
	}
	return nil, nil
}

func GetTxType(blockchain int64, msg *kafka.Message) (uint64, error) {
	if blockchain == 200 {
		return GetTxTypeForEther(msg.Value)
	}
	if blockchain == 205 {
		return GetTxTypeForTron(msg.Value)
	}
	return 0, nil
}

func GetTxTypeForEther(body []byte) (uint64, error) {
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

func ParseTxForEther(body []byte) (*SubTx, error) {
	var r SubTx
	root := gjson.ParseBytes(body)
	r.BlockChain = 200
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
	r.Value = div(v, 18)

	tp, err := GetTxTypeForEther(body)
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
		r.Fee = div(fmt.Sprintf("%v", fee), 18)
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

	logs := receiptRoot.Get("logs").Array()
	contractTx := make([]*ContractTx, 0, 5)
	for _, v := range logs {
		contract := v.Get("address").String()
		data := v.Get("data").String()
		r := gjson.Parse(data)
		if r.IsObject() {
			//mp := make(map[string]string, 2)
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
				data = div(bigFee, decimals)
			} else {
				return nil, errors.New("tx.log.contract is error")
			}
		} else {
			continue
		}

		tps := v.Get("topics").Array()
		var from, to string
		if len(tps) >= 3 && tps[0].String() == EthTopic {
			//method = tps[0].String()
			from = tps[1].String()
			to = tps[2].String()
			var m ContractTx
			m.Contract = contract
			m.Value = data
			m.From, _ = util.Hex2Address(from)
			m.To, _ = util.Hex2Address(to)
			m.Method = "Transfer"
			contractTx = append(contractTx, &m)
		}

	}
	r.ContractTx = contractTx
	return &r, nil
}

func GetTxTypeForTron(body []byte) (uint64, error) {

	root := gjson.ParseBytes(body)
	txBody := root.Get("tx").String()
	if len(txBody) < 5 {
		return 0, errors.New("tx is error")
	}

	txRoot := gjson.Parse(txBody)

	var tx uint64
	txType := txRoot.Get("raw_data.contract.0.type").String()
	if txType == "TransferContract" {
		tx = 2
	} else if txType == "TriggerSmartContract" {
		tx = 1
	} else if txType == "DelegateResourceContract" {
		tx = 3
	} else if txType == "UnDelegateResourceContract" {
		tx = 4
	} else if txType == "AccountCreateContract" {
		tx = 5
	} else if txType == "FreezeBalanceV2Contract" {
		tx = 6
	} else if txType == "UnfreezeBalanceV2Contract" {
		tx = 7
	} else if txType == "WithdrawExpireUnfreezeContract" {
		tx = 8
	} else {
		return 0, errors.New("undefined type of tx")
	}
	return tx, nil
}

func ParseTxForTron(body []byte) (*SubTx, error) {

	root := gjson.ParseBytes(body)

	blockId := root.Get("blockId").String()
	txBody := root.Get("tx").String()
	if len(txBody) < 5 {
		return nil, errors.New("tx is error")
	}
	//r := make(map[string]interface{}, 10)
	var r SubTx
	r.BlockChain = 205
	r.Id = uint64(time.Now().UnixNano())
	txRoot := gjson.Parse(txBody)
	status := txRoot.Get("ret.0.contractRet").String()
	if status == "SUCCESS" { //交易成功
		r.Status = 1
	} else {
		//交易失败
		r.Status = 0
	}

	hash := txRoot.Get("txID").String()
	r.Hash = hash
	//tron block_hash 比较特殊，是hash 部分，暂时不返回了
	//blockHash := txRoot.Get("raw_data.ref_block_hash").String()
	r.BlockHash = blockId
	txTime := txRoot.Get("raw_data.timestamp").String()
	r.TxTime = txTime

	tp, err := GetTxTypeForTron(body)
	if err != nil {
		return nil, err
	}
	r.TxType = tp

	v := txRoot.Get("raw_data.contract.0.parameter.value")
	from := v.Get("owner_address").String()
	r.From = util.HexToAddress(from).Base58()
	var to string
	//DelegateResourceContract
	if v.Get("receiver_address").Exists() {
		to = v.Get("receiver_address").String()
	}

	//TransferContract 转帐
	if v.Get("to_address").Exists() {
		to = v.Get("to_address").String()
	}

	//TriggerSmartContract
	if v.Get("contract_address").Exists() {
		to = v.Get("contract_address").String()
	}

	//AccountCreateContract
	if v.Get("account_address").Exists() {
		to = v.Get("account_address").String()
	}

	r.To = util.HexToAddress(to).Base58()

	var input string
	if v.Get("data").Exists() {
		input = v.Get("data").String()
	}
	r.Input = input

	if r.TxType == 2 { //普通交易
		r.Value = div(v.Get("amount").String(), 6)
	} else if r.TxType == 1 { //合约调用
		r.Value = "0"
	} else { //其他
		r.Value = v.String()
	}

	if !root.Get("receipt").Exists() { //收据不存在的交易，则放弃
		return nil, errors.New("receipt is error")
	}

	receiptBody := root.Get("receipt").String()
	if len(receiptBody) > 5 {
		receiptRoot := gjson.Parse(receiptBody)
		fee := receiptRoot.Get("fee").String()
		r.Fee = div(fee, 6)
		gasFee := receiptRoot.Get("receipt").Map()
		delete(gasFee, "result")
		r.FeeDetail = gasFee
		number := receiptRoot.Get("blockNumber").String()
		r.BlockNumber = number

		logs := receiptRoot.Get("log").Array()
		contractTx := make([]*ContractTx, 0, 5)
		for _, v := range logs {
			contract := v.Get("address").String()
			//value := v.Get("data").String()
			tps := v.Get("topics").Array()
			data := v.Get("data").String()
			r := gjson.Parse(data)
			if r.IsObject() {

				//mp := make(map[string]string, 2)
				contractDecimals := r.Get("contractDecimals").String()
				if len(contractDecimals) < 1 {
					return nil, errors.New("tx.log.contract is error")
				}
				decimals, err := strconv.Atoi(contractDecimals)
				if err != nil {
					return nil, errors.New("tx.log.contract is error")
				}
				//bigDecimals := math.Pow10(decimals)

				//mp["contractDecimals"] = contractDecimals

				fee := r.Get("data").String()
				bigFee, err := util.HexToInt(fee)
				if err == nil {
					//fmt.Sprintf("%.5f", new(big.Float).Quo(new(big.Float).SetFloat64(float64(bigFee)), new(big.Float).SetFloat64(bigDecimals)))
					data = div(bigFee, decimals)
				} else {
					return nil, errors.New("tx.log.contract is error")
				}
				//bs, _ := json.Marshal(mp)
				//data = string(bs)
			} else {
				//data, _ = util.HexToInt(data)
				continue
			}

			var from, to string
			if len(tps) >= 3 && tps[0].String() == TronTopic {
				//method = tps[0].String()
				from = tps[1].String()
				to = tps[2].String()
				var m ContractTx
				m.Contract = util.HexToAddress(fmt.Sprintf("41%v", contract)).Base58()
				m.Value = data

				from, _ := util.Hex2Address2(from)
				m.From = util.HexToAddress(from).Base58()

				to, _ := util.Hex2Address2(to)
				m.To = util.HexToAddress(to).Base58()
				m.Method = "Transfer"
				contractTx = append(contractTx, &m)
			}

		}
		r.ContractTx = contractTx
	}

	return &r, nil
}

func div(str string, pos int) string {

	if str == "" || str == "0" {
		return "0"
	}

	if pos == 0 {
		return str
	}

	r := make([]string, 0, 10)
	for true {
		if len(str) <= pos {
			str = "0" + str
		} else {
			break
		}
	}

	list := []byte(str)
	l := len(list)
	p := 0
	for i := l - 1; i >= 0; i-- {
		s := fmt.Sprintf("%c", list[l-1-i])
		r = append(r, s)

		if l-1 > 0 && (l-1-p == pos) {
			r = append(r, ".")
		}

		p++
	}

	result := fmt.Sprintf("%s", strings.Join(r, ""))

	for strings.HasSuffix(result, "0") || strings.HasSuffix(result, ".") {
		if strings.HasSuffix(result, "0") {
			result = strings.TrimSuffix(result, "0")
		}

		if strings.HasSuffix(result, ".") {
			result = strings.TrimSuffix(result, ".")
		}
	}
	return result
}

func CheckAddress(blockChain int64, msg *kafka.Message, list []*MonitorAddress) bool {
	if blockChain == 200 {
		return CheckAddressForEther(msg, list)
	} else if blockChain == 205 {
		return CheckAddressForTron(msg, list)
	} else {
		return false
	}

}

func CheckAddressForEther(msg *kafka.Message, list []*MonitorAddress) bool {
	root := gjson.ParseBytes(msg.Value)
	fromAddr := root.Get("from").String()
	toAddr := root.Get("to").String()
	has := false
	for _, v := range list {
		//已经判断出该交易 符合要求了，不需要在检查其他地址了
		if has {
			break
		}

		// 普通交易且 地址包含订阅地址
		if strings.HasPrefix(strings.ToLower(fromAddr), strings.ToLower(v.Address)) || strings.HasPrefix(strings.ToLower(toAddr), strings.ToLower(v.Address)) {
			has = true
			break
		}

		//合约交易
		monitorAddr := strings.TrimLeft(v.Address, "0x") //去丢0x
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
				//Transfer()
				if len(topics) >= 3 && topics[0].String() == service.EthTopic {
					if strings.HasSuffix(strings.ToLower(topics[1].String()), strings.ToLower(monitorAddr)) || strings.HasSuffix(strings.ToLower(topics[2].String()), strings.ToLower(monitorAddr)) {
						has = true
						break
					}

				}

			}

		}
	}
	return has
}

func CheckAddressForTron(msg *kafka.Message, list []*MonitorAddress) bool {
	root := gjson.ParseBytes(msg.Value)
	tx := root.Get("tx").String()
	txRoot := gjson.Parse(tx)
	contracts := txRoot.Get("raw_data.contract").Array()
	if len(contracts) < 1 {
		return false
	}
	r := contracts[0]
	txType := r.Get("type").String()

	var fromAddr, toAddr string
	var logs []gjson.Result
	var internalTransactions []gjson.Result

	fromAddr = r.Get("parameter.value.owner_address").String()

	if txType == "TransferContract" {
		toAddr = r.Get("parameter.value.to_address").String()
		//r.Get("parameter.value.amount").String()
	}

	//DelegateResourceContract,UnDelegateResourceContract
	if r.Get("parameter.value.receiver_address").Exists() {
		toAddr = r.Get("parameter.value.receiver_address").String()
	}

	//TriggerSmartContract
	if r.Get("parameter.value.contract_address").Exists() {
		toAddr = r.Get("parameter.value.contract_address").String()
	}

	//AccountCreateContract
	if r.Get("parameter.value.account_address").Exists() {
		toAddr = r.Get("parameter.value.account_address").String()
	}

	if txType == "TriggerSmartContract" {
		receipt := root.Get("receipt").String()
		receiptRoot := gjson.Parse(receipt)
		if receiptRoot.Get("receipt.result").String() != "SUCCESS" {
			return false
		}
		logs = receiptRoot.Get("log").Array()
		internalTransactions = receiptRoot.Get("internal_transactions").Array()
	}

	has := false
	for _, v := range list {
		//已经判断出该交易 符合要求了，不需要在检查其他地址了
		if has {
			break
		}

		var monitorAddr string
		monitorAddr = v.Address

		if strings.HasPrefix(v.Address, "0x") {
			monitorAddr = strings.TrimLeft(v.Address, "0x") //去丢0x
		}

		if strings.HasPrefix(v.Address, "41") {
			monitorAddr = strings.TrimLeft(v.Address, "41") //去丢41
		}

		if strings.HasPrefix(v.Address, "0x41") {
			monitorAddr = strings.TrimLeft(v.Address, "0x41") //去丢41
		}

		// 普通交易且 地址包含订阅地址
		if strings.HasSuffix(fromAddr, monitorAddr) || strings.HasSuffix(toAddr, monitorAddr) {
			has = true
			break
		}

		if txType == "TriggerSmartContract" {
			//合约交易 合约调用下的TRC20
			if len(logs) > 0 {
				for _, v := range logs {

					//过滤没有合约信息的交易，出现这种情况原因：1. 合约获取失败会重试 2:非20合约
					data := v.Get("data").String()
					if !gjson.Parse(data).IsObject() {
						continue
					}

					topics := v.Get("topics").Array()
					//Transfer()
					if len(topics) >= 3 && topics[0].String() == service.TronTopic {
						if strings.HasSuffix(topics[1].String(), monitorAddr) || strings.HasSuffix(topics[2].String(), monitorAddr) {
							has = true
							break
						}
					}
				}
			}

			//合约调用下的内部交易TRX转帐和TRC10转账：
			if len(internalTransactions) > 0 {
				for _, v := range internalTransactions {
					fromAddr = v.Get("caller_address").String()
					toAddr = v.Get("transferTo_address").String()
					if strings.HasSuffix(fromAddr, monitorAddr) || strings.HasSuffix(toAddr, monitorAddr) {
						has = true
						break
					}
				}
			}
		}
	}
	return has
}
