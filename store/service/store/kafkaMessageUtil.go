package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/uduncloud/easynode/store/service"
	"time"
)

func GetReceiptFromKafka(value []byte, blockChain int64) (*service.Receipt, error) {
	var receipt service.Receipt
	if blockChain == 200 {
		err := json.Unmarshal(value, &receipt)
		if err != nil {
			//s.log.Errorf("ReadReceiptFromKafka|blockChain=%v,error=%v", blockChain, err.Error())
			//continue
			return nil, err
		}
	} else if blockChain == 205 {

		root := gjson.ParseBytes(value)
		hash := root.Get("id").String()
		fee := root.Get("fee").Uint()
		number := root.Get("blockNumber").Uint()
		blockTime := root.Get("blockTimeStamp").Uint()
		contractAddress := root.Get("contract_address").String()
		status := root.Get("receipt.result").String()

		logs := root.Get("log").String()
		var list *service.Logs
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
	} else {
		return nil, errors.New(fmt.Sprintf("blockchain:%v does not support", blockChain))
	}

	return &receipt, nil
}

func GetBlockFromKafka(value []byte, blockChain int64) (*service.Block, error) {
	var block service.Block
	if blockChain == 200 {
		err := json.Unmarshal(value, &block)
		if err != nil {
			//s.log.Errorf("ReadBlockFromKafka|blockChain=%v,error=%v", blockChain, err.Error())
			//continue
			return nil, err
		}
	} else if blockChain == 205 {
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
	} else {
		return nil, errors.New(fmt.Sprintf("blockchain:%v does not support", blockChain))
	}

	return &block, nil
}

func GetTxFromKafka(value []byte, blockChain int64) (*service.Tx, error) {
	var tx service.Tx
	if blockChain == 200 {
		err := json.Unmarshal(value, &tx)
		if err != nil {
			//s.log.Errorf("ReadTxFromKafka|blockChain=%v,error=%v", blockChain, err.Error())
			//continue
			return nil, err
		}
	} else if blockChain == 205 {

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

	} else {
		return nil, errors.New(fmt.Sprintf("blockchain:%v does not support", blockChain))
	}
	return &tx, nil
}
