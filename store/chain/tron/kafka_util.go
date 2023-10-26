package tron

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/0xcregis/easynode/common/util"
	"github.com/0xcregis/easynode/store"
	"github.com/tidwall/gjson"
)

func GetReceiptFromKafka(value []byte) (*store.Receipt, error) {
	var receipt store.Receipt
	root := gjson.ParseBytes(value)
	hash := root.Get("id").String()
	fee := root.Get("fee").Uint()
	number := root.Get("blockNumber").Uint()
	blockTime := root.Get("blockTimeStamp").Uint()
	contractAddress := root.Get("contract_address").String()
	status := root.Get("receipt.result").String()

	logs := root.Get("log").String()
	var list *store.Logs
	_ = json.Unmarshal([]byte(logs), &list)
	receipt.Id = uint64(time.Now().UnixNano())
	receipt.TransactionHash = hash
	receipt.BlockNumber = fmt.Sprintf("%v", number)
	receipt.GasUsed = fmt.Sprintf("%v", fee)
	receipt.Status = status
	receipt.ContractAddress = contractAddress
	receipt.CreateTime = fmt.Sprintf("%v", blockTime)
	receipt.Logs = list
	receipt.CumulativeGasUsed = root.Get("receipt").String()

	return &receipt, nil
}

func GetBlockFromKafka(value []byte) (*store.Block, error) {
	var block store.Block
	r := gjson.ParseBytes(value)
	hash := r.Get("blockID").String()
	number := r.Get("block_header.raw_data").Get("number").Int()
	txRoot := r.Get("block_header.raw_data").Get("txTrieRoot").String()
	parentHash := r.Get("block_header.raw_data").Get("parentHash").String()
	coinAddr := r.Get("block_header.raw_data").Get("witness_address").String()
	blockTime := r.Get("block_header.raw_data").Get("timestamp").String()

	block.Id = uint64(time.Now().UnixNano())
	block.BlockHash = hash
	block.BlockNumber = fmt.Sprintf("%v", number)
	block.BlockTime = blockTime
	block.ParentHash = parentHash
	block.TxRoot = txRoot
	block.Coinbase = coinAddr

	array := r.Get("transactions.#.txID").Array()
	list := make([]string, 0, len(array))
	for _, v := range array {
		list = append(list, v.String())
	}
	block.Transactions = list

	return &block, nil
}

func GetTxFromKafka(value []byte) (*store.Tx, error) {
	var tx store.Tx
	blockId := gjson.ParseBytes(value).Get("blockId").String()
	txBody := gjson.ParseBytes(value).Get("tx").String()
	if len(txBody) < 5 {
		return nil, errors.New("tx data is error")
	}
	txRoot := gjson.Parse(txBody)

	status := txRoot.Get("ret.0.contractRet").String()
	hash := txRoot.Get("txID").String()
	//blockHash := txRoot.Get("raw_data.ref_block_hash").String()
	txTime := txRoot.Get("raw_data.timestamp").Uint()
	limit := txRoot.Get("raw_data.fee_limit").Uint()
	txType := txRoot.Get("raw_data.contract.0.type").String()
	v := txRoot.Get("raw_data.contract.0.parameter.value")
	from := v.Get("owner_address").String()
	var to string
	if v.Get("receiver_address").Exists() {
		to = v.Get("receiver_address").String()
	}

	if v.Get("to_address").Exists() {
		to = v.Get("to_address").String()
	}

	if v.Get("contract_address").Exists() {
		to = v.Get("contract_address").String()
	}

	var input string
	if v.Get("data").Exists() {
		input = v.Get("data").String()
	}

	txValue := v.String()
	tx.Id = uint64(time.Now().UnixNano())
	tx.Value = txValue
	tx.BlockHash = blockId
	tx.TxHash = hash
	tx.TxStatus = status
	tx.TxTime = fmt.Sprintf("%v", txTime)
	tx.ToAddr = to
	tx.FromAddr = from
	tx.Type = txType
	tx.InputData = input
	tx.GasLimit = fmt.Sprintf("%v", limit)
	receiptBody := gjson.ParseBytes(value).Get("receipt").String()
	if len(receiptBody) > 5 {
		receiptRoot := gjson.Parse(receiptBody)
		fee := receiptRoot.Get("fee").Uint()
		number := receiptRoot.Get("blockNumber").Uint()
		tx.Fee = fmt.Sprintf("%v", fee)
		tx.BlockNumber = fmt.Sprintf("%v", number)
	}

	return &tx, nil
}

func GetTxType(body []byte) (uint64, error) {
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
func ParseTx(body []byte, transferTopic string, blockchain int64) (*store.SubTx, error) {
	root := gjson.ParseBytes(body)

	blockId := root.Get("blockId").String()
	txBody := root.Get("tx").String()
	if len(txBody) < 5 {
		return nil, errors.New("tx is error")
	}
	//r := make(map[string]interface{}, 10)
	var r store.SubTx
	r.BlockChain = uint64(blockchain)
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

	tp, err := GetTxType(body)
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
		r.Value = util.Div(v.Get("amount").String(), 6)
	} else if r.TxType == 1 { //合约调用
		r.Value = "0"
	} else { //其他
		r.Value = "-1"
	}

	if !root.Get("receipt").Exists() { //收据不存在的交易，则放弃
		return nil, errors.New("receipt is error")
	}

	receiptBody := root.Get("receipt").String()
	if len(receiptBody) > 5 {
		receiptRoot := gjson.Parse(receiptBody)
		fee := receiptRoot.Get("fee").String()
		r.Fee = util.Div(fee, 6)
		gasFee := receiptRoot.Get("receipt").Map()
		delete(gasFee, "result")
		r.FeeDetail = gasFee
		number := receiptRoot.Get("blockNumber").String()
		r.BlockNumber = number
		contractTx := make([]*store.ContractTx, 0, 5)

		if r.TxType != 1 {
			var c store.ContractTx
			c.Contract = ""
			c.Index = 0
			c.Value = r.Value
			c.From = r.From
			c.To = r.To
			c.Method = "Transfer"
			c.EIP = 0
			c.Token = ""
			contractTx = append(contractTx, &c)
		}

		logs := receiptRoot.Get("log").Array()
		l := len(contractTx)
		for index, v := range logs {
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
					data = util.Div(bigFee, decimals)
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
			if len(tps) >= 3 && tps[0].String() == transferTopic {
				//method = tps[0].String()
				from = tps[1].String()
				to = tps[2].String()
				var m store.ContractTx
				m.Contract = util.HexToAddress(fmt.Sprintf("41%v", contract)).Base58()
				m.Value = data
				m.Index = int64(index + l)

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
func GetCoreAddr(addr string) string {
	if strings.HasPrefix(addr, "0x41") {
		return strings.Replace(addr, "0x41", "", 1) //去丢41
	}

	if strings.HasPrefix(addr, "0x") {
		return strings.Replace(addr, "0x", "", 1) //去丢0x
	}

	if strings.HasPrefix(addr, "41") {
		return strings.Replace(addr, "41", "", 1) //去丢41
	}

	return addr
}
func CheckAddress(txValue []byte, addrList map[string]*store.MonitorAddress, transferTopic string) bool {
	if len(addrList) < 1 || len(txValue) < 1 {
		return false
	}
	txAddressList := make(map[string]int64, 10)

	root := gjson.ParseBytes(txValue)
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
	txAddressList[GetCoreAddr(fromAddr)] = 1

	//TransferContract,TransferAssetContract
	if r.Get("parameter.value.to_address").Exists() {
		toAddr = r.Get("parameter.value.to_address").String()
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

	txAddressList[GetCoreAddr(toAddr)] = 1

	if txType == "TriggerSmartContract" {
		receipt := root.Get("receipt").String()
		receiptRoot := gjson.Parse(receipt)
		if receiptRoot.Get("receipt.result").String() != "SUCCESS" {
			return false
		}
		logs = receiptRoot.Get("log").Array()
		internalTransactions = receiptRoot.Get("internal_transactions").Array()

		//合约交易 合约调用下的TRC20
		if len(logs) > 0 {
			for _, v := range logs {
				topics := v.Get("topics").Array()
				//Transfer()
				if len(topics) >= 3 && topics[0].String() == transferTopic {
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

		//合约调用下的内部交易TRX转帐和TRC10转账：
		if len(internalTransactions) > 0 {
			for _, v := range internalTransactions {
				fromAddr = v.Get("caller_address").String()
				toAddr = v.Get("transferTo_address").String()
				txAddressList[GetCoreAddr(fromAddr)] = 1
				txAddressList[GetCoreAddr(toAddr)] = 1
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
