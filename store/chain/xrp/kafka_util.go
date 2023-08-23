package xrp

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/0xcregis/easynode/store"
	"github.com/tidwall/gjson"
)

func GetBlockFromKafka(value []byte) (*store.Block, error) {
	/**
	{
	  "accepted": "true",
	  "account_hash": "910FCF5C2CA222286C6D1D3AB40198D4552B639129A0EE516A45E8B4770EA00E",
	  "close_flags": "0",
	  "close_time": "746087402",
	  "close_time_human": "2023-Aug-23 06:30:02.000000000 UTC",
	  "close_time_resolution": "10",
	  "closed": "true",
	  "hash": "224CE7DA0F1D88ED35559E5D1904B20043DD125AFA21CED62FA66780BDC1B491",
	  "ledger_hash": "224CE7DA0F1D88ED35559E5D1904B20043DD125AFA21CED62FA66780BDC1B491",
	  "ledger_index": "82032099",
	  "parent_close_time": "746087401",
	  "parent_hash": "9D4985A3CBD7CACD42C0149C440304C3AFBB1BF3157CA4063D94F6BB8B077614",
	  "seqNum": "82032099",
	  "totalCoins": "99988480317227197",
	  "total_coins": "99988480317227197",
	  "transaction_hash": "AFA1C7E4CD2956BDCF495F597B7EEE0D96457399FA59DFFA190FF20088EE81DE"
	}
	*/
	var block store.Block
	r := gjson.ParseBytes(value)
	hash := r.Get("ledger_hash").String()
	number := r.Get("ledger_index").String()

	parentHash := r.Get("parent_hash").String()
	//coinAddr := bookRoot.Get("Miner").String()
	blockTime := r.Get("close_time").String()
	//baseFee := bookRoot.Get("ParentBaseFee").String()
	txHash := r.Get("transaction_hash").String()
	statusHash := r.Get("account_hash").String()
	status := r.Get("closed").Bool()

	block.Id = uint64(time.Now().UnixNano())
	block.BlockHash = hash
	block.Root = statusHash
	block.TxRoot = txHash
	if status {
		block.BlockStatus = "1"
	} else {
		block.BlockStatus = "0"
	}

	block.BlockNumber = fmt.Sprintf("%v", number)
	block.BlockTime = blockTime
	block.ParentHash = parentHash
	//block.Coinbase = coinAddr
	//block.BaseFee = baseFee
	return &block, nil
}

func GetTxFromKafka(value []byte) (*store.Tx, error) {
	/**
	{
	  "block": "{\"accepted\":\"true\",\"account_hash\":\"2BD14C2BE88F3ED179007E0EFA752BC36127308D824354C55504E31AC9562C21\",\"close_flags\":\"0\",\"close_time\":\"746087000\",\"close_time_human\":\"2023-Aug-23 06:23:20.000000000 UTC\",\"close_time_resolution\":\"10\",\"closed\":\"true\",\"hash\":\"FFF79B99B28E446652B741A682FA40D4C1E485BBE78EA5D3868223C60E6C93A6\",\"ledger_hash\":\"FFF79B99B28E446652B741A682FA40D4C1E485BBE78EA5D3868223C60E6C93A6\",\"ledger_index\":\"82031992\",\"parent_close_time\":\"746086992\",\"parent_hash\":\"8F058D0BA35084C7989834CA1AF3AA195AF47A850FAEB4148C4704C20D887BD0\",\"seqNum\":\"82031992\",\"totalCoins\":\"99988480325556972\",\"total_coins\":\"99988480325556972\",\"transaction_hash\":\"CFD845FDF68FAB6553ADA7F398301C7E347980B00CEC98C5CB011F84F50E34AA\"}",
	  "blockHash": "FFF79B99B28E446652B741A682FA40D4C1E485BBE78EA5D3868223C60E6C93A6",
	  "blockNumber": "82031992",
	  "hash": "EB9F83623CE7CE8802B6D1D12FB55133437005AAB2E4E09BA61C2B1EAE5751B7",
	  "tx": "{\"Account\":\"rJn2zAPdFA193sixJwuFixRkYDUtx3apQh\",\"Amount\":\"2299401567509\",\"Destination\":\"rMvCasZ9cohYrSZRNYPTZfoaaSUQMfgQ8G\",\"Fee\":\"10000\",\"Sequence\":105320,\"SigningPubKey\":\"037DFC657B370962A83DA8035E05B240D2A921AEF5FD59B904D71E3D5B0E3C9D9C\",\"TransactionType\":\"Payment\",\"TxnSignature\":\"3044022075005F0CAFB68AE7FE135AEA8F2496AC4FF1B1EA9DC230D35E3545FDB623EFE3022031B0696D319E3D665E3DF14EC41A2841F60EECE1E9A54817A17E208C11A837B4\",\"hash\":\"EB9F83623CE7CE8802B6D1D12FB55133437005AAB2E4E09BA61C2B1EAE5751B7\",\"metaData\":{\"AffectedNodes\":[{\"ModifiedNode\":{\"FinalFields\":{\"Account\":\"rMvCasZ9cohYrSZRNYPTZfoaaSUQMfgQ8G\",\"Balance\":\"32155814798353\",\"Flags\":0,\"OwnerCount\":1,\"Sequence\":66878203},\"LedgerEntryType\":\"AccountRoot\",\"LedgerIndex\":\"3BFAEEAD933C62413BD9A92BBC11993F2B0FCA28E2CD8FEE69EB8C22EE848E36\",\"PreviousFields\":{\"Balance\":\"29856413230844\"},\"PreviousTxnID\":\"A0BCDF69185C351155D58AC7C4E5C19061637078107EBC33990FE65F12503FD7\",\"PreviousTxnLgrSeq\":82031976}},{\"ModifiedNode\":{\"FinalFields\":{\"Account\":\"rJn2zAPdFA193sixJwuFixRkYDUtx3apQh\",\"Balance\":\"1780244485\",\"Flags\":131072,\"OwnerCount\":1,\"Sequence\":105321},\"LedgerEntryType\":\"AccountRoot\",\"LedgerIndex\":\"C19B36F6B6F2EEC9F4E2AF875E533596503F4541DBA570F06B26904FDBBE9C52\",\"PreviousFields\":{\"Balance\":\"2301181821994\",\"Sequence\":105320},\"PreviousTxnID\":\"C72F3FFDE1E59E8DE391A4E850E9031AE8D33156A7A2EBD1E3E355A364E527B7\",\"PreviousTxnLgrSeq\":82031989}}],\"TransactionIndex\":19,\"TransactionResult\":\"tesSUCCESS\",\"delivered_amount\":\"2299401567509\"}}"
	}
	*/
	var tx store.Tx
	root := gjson.ParseBytes(value)

	tx.TxHash = root.Get("hash").String()

	if root.Get("blockNumber").Exists() {
		tx.BlockNumber = root.Get("blockNumber").String()
	}

	if root.Get("blockHash").Exists() {
		tx.BlockHash = root.Get("blockHash").String()
	}

	txBody := root.Get("tx").String()
	if len(txBody) < 5 {
		return nil, errors.New("tx data is error")
	}
	txRoot := gjson.Parse(txBody)

	var status, txValue, transactionIndex string
	if txRoot.Get("meta").Exists() {
		status = txRoot.Get("meta.TransactionResult").String()
		txValue = txRoot.Get("meta.delivered_amount").String()
		transactionIndex = txRoot.Get("meta.TransactionIndex").String()
	}

	if txRoot.Get("metaData").Exists() {
		status = txRoot.Get("metaData.TransactionResult").String()
		txValue = txRoot.Get("metaData.delivered_amount").String()
		transactionIndex = txRoot.Get("meta.TransactionIndex").String()
	}
	//status := txRoot.Get("meta.TransactionResult").String()
	txTime := txRoot.Get("date").Uint()
	//limit := txRoot.Get("GasLimit").String()
	txType := txRoot.Get("TransactionType").String()
	//txValue := txRoot.Get("meta.delivered_amount").String()
	from := txRoot.Get("Account").String()
	to := txRoot.Get("Destination").String()
	blockNumber := txRoot.Get("ledger_index").String()
	input := txRoot.String()
	//txValue := txRoot.Get("Value").String()
	//gasPremium := txRoot.Get("GasPremium").String()
	//gasFeeCap:=txRoot.Get("GasFeeCap").String()
	fee := txRoot.Get("Fee").String()

	tx.Id = uint64(time.Now().UnixNano())
	tx.Value = txValue
	tx.BlockNumber = blockNumber
	tx.TransactionIndex = transactionIndex
	//tx.BlockHash = blockId
	//tx.TxHash = hash
	tx.TxStatus = status
	tx.TxTime = fmt.Sprintf("%v", txTime)
	tx.ToAddr = to
	tx.FromAddr = from
	tx.Type = txType
	tx.InputData = input
	tx.Fee = fee
	//tx.GasLimit = limit
	//tx.PriorityFee = gasPremium

	blockBody := root.Get("block").String()
	if len(blockBody) > 5 {
		blockRoot := gjson.Parse(blockBody)
		tx.BlockNumber = blockRoot.Get("ledger_index").String()
		tx.BlockHash = blockRoot.Get("ledger_hash").String()
		tx.TxTime = blockRoot.Get("close_time").String()
	}

	return &tx, nil
}
func ParseTx(value []byte, transferTopic string, blockchain int64) (*store.SubTx, error) {
	var tx store.SubTx
	root := gjson.ParseBytes(value)

	tx.BlockChain = uint64(blockchain)
	tx.Hash = root.Get("hash").String()

	if root.Get("blockNumber").Exists() {
		tx.BlockNumber = root.Get("blockNumber").String()
	}

	if root.Get("blockHash").Exists() {
		tx.BlockHash = root.Get("blockHash").String()
	}

	txBody := root.Get("tx").String()
	if len(txBody) < 5 {
		return nil, errors.New("tx data is error")
	}
	txRoot := gjson.Parse(txBody)

	var status, txValue string
	if txRoot.Get("meta").Exists() {
		status = txRoot.Get("meta.TransactionResult").String()
		txValue = txRoot.Get("meta.delivered_amount").String()
		//transactionIndex = txRoot.Get("meta.TransactionIndex").String()
	}

	if txRoot.Get("metaData").Exists() {
		status = txRoot.Get("metaData.TransactionResult").String()
		txValue = txRoot.Get("metaData.delivered_amount").String()
		//transactionIndex = txRoot.Get("meta.TransactionIndex").String()
	}
	//status := txRoot.Get("meta.TransactionResult").String()
	txTime := txRoot.Get("date").Uint()
	//limit := txRoot.Get("GasLimit").String()
	txType := txRoot.Get("TransactionType").String()
	//txValue := txRoot.Get("meta.delivered_amount").String()
	from := txRoot.Get("Account").String()
	to := txRoot.Get("Destination").String()
	blockNumber := txRoot.Get("ledger_index").String()
	input := txRoot.String()
	//txValue := txRoot.Get("Value").String()
	//gasPremium := txRoot.Get("GasPremium").String()
	//gasFeeCap:=txRoot.Get("GasFeeCap").String()
	fee := txRoot.Get("Fee").String()

	tx.Id = uint64(time.Now().UnixNano())
	tx.Value = txValue
	tx.BlockNumber = blockNumber
	//tx.TransactionIndex = transactionIndex
	//tx.BlockHash = blockId
	//tx.TxHash = hash
	if status == "tesSUCCESS" {
		tx.Status = 1
	} else {
		tx.Status = 0
	}

	tx.TxTime = fmt.Sprintf("%v", txTime)
	tx.To = to
	tx.From = from
	if txType == "Payment" {
		tx.TxType = 2
	}
	//tx.TxType = txType
	tx.Input = input
	tx.Fee = fee
	//tx.GasLimit = limit
	//tx.PriorityFee = gasPremium

	blockBody := root.Get("block").String()
	if len(blockBody) > 5 {
		blockRoot := gjson.Parse(blockBody)
		tx.BlockNumber = blockRoot.Get("ledger_index").String()
		tx.BlockHash = blockRoot.Get("ledger_hash").String()
		tx.TxTime = blockRoot.Get("close_time").String()
	}
	return &tx, nil
}

func GetReceiptFromKafka(value []byte) (*store.Receipt, error) {
	/**

	{
	    "account":"rf4VZ6LYKnPY1uaEiZF8qzD8vrW91vvjbb",
	    "date":745986661,
	    "hash":"89CEEBAF42BE602DCFDF0F89B7B9111A4E09CF4D55EC0BFEF5C439246E31C8B5",
	    "ledgerIndex":82005629,
	    "status":"success",
	    "transactionIndex":18,
	    "transactionResult":"tecINSUFFICIENT_RESERVE"
	}

	*/
	var receipt store.Receipt

	root := gjson.ParseBytes(value)

	receipt.Id = uint64(time.Now().UnixNano())
	receipt.TransactionHash = root.Get("hash").String()
	receipt.BlockNumber = root.Get("ledgerIndex").String()
	receipt.TransactionIndex = root.Get("transactionIndex").String()
	receipt.From = root.Get("account").String()
	receipt.Status = root.Get("transactionResult").String()
	return &receipt, nil
}

func GetTxType(body []byte) (uint64, error) {
	root := gjson.ParseBytes(body)
	txBody := root.Get("tx").String()
	txType := gjson.Parse(txBody).Get("TransactionType").String()
	if txType == "Payment" {
		//普通资产转移
		return 2, nil
	}
	return 0, fmt.Errorf("unknow for blockchain")
}

func GetCoreAddr(addr string) string {
	if strings.HasPrefix(addr, "0x") {
		return strings.Replace(addr, "0x", "", 1) //去丢0x
	}
	return addr
}

func CheckAddress(txValue []byte, addrList map[string]*store.MonitorAddress, transferTopic string) bool {
	if len(addrList) < 1 || len(txValue) < 1 {
		return false
	}

	txAddressList := make(map[string]int64, 10)
	root := gjson.ParseBytes(txValue)

	txBody := root.Get("tx").String()
	txRoot := gjson.Parse(txBody)

	fromAddr := txRoot.Get("Account").String()
	txAddressList[GetCoreAddr(fromAddr)] = 1

	if txRoot.Get("Destination").Exists() {
		toAddr := txRoot.Get("Destination").String()
		txAddressList[GetCoreAddr(toAddr)] = 1
	}

	has := false
	for k := range txAddressList {
		//monitorAddr := getCoreAddr(v)
		if _, ok := addrList[k]; ok {
			has = true
			break
		}
	}
	return has
}
