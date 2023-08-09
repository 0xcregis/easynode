package filecoin

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
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
	/**
	{
	  "block": "{\"Miner\":\"f01740934\",\"Ticket\":{\"VRFProof\":\"gqF930YfBc98OQqJ58cRVG0P9xKAjWSgZz8clGNKwhQOZ/fptd/kbfe+7F3XrDT+Fc1rRYJxZSkXG6/+77y2sKZnd+mTISj2C1gV0eIsL8blAZnitbOVcrFRbiOUELX8\"},\"ElectionProof\":{\"WinCount\":1,\"VRFProof\":\"pGth0oztjSrhgJBdNWD2i3N7cG1DXznXZ3/i2Ioeng6uB6Za4rQxLrFjh98tennYFDnyv6xH2Be8YmsqFXOhyfR0KX3pA5EnO0ovquP2iWgRgQ8AZzW/NkIaKwqyVFTV\"},\"BeaconEntries\":[{\"Round\":3190391,\"Data\":\"gzL6qDFWfBXluwwpcpTBe8G5/lcqZGpxqpTgxJHGV/njb3Tvu4aB9hv6O31atYLEA5PRxI7fVqFbxw2LM4D+hIAQjS1v5V0/5YVT4rE90yRNm7cKr4vM861cScWZlN9H\"}],\"WinPoStProof\":[{\"PoStProof\":3,\"ProofBytes\":\"rp8v9iFqv5vYAqS0ybR3aBbvS9Ozs6UbehGqXqBc56hz++eeCiTxJJUfmRaj6RRqhtdl6/mB08HoIN4lt7lGjIOdM5cjgfqIPj+9ena+GPcyLkVvnJzNlyna/5Vdh6sbFtPU2lr3J7w5tysRO3vQw+1tN02XNs/bFjaCWSyIgFtGySaKEBgR1eV7LYd7wT65i+rmHoKlp0SQono9KpjInobCrnyYAA6JPX8wUG1IQrtz+1oA51f3XtL9d55TmC9m\"}],\"Parents\":[{\"/\":\"bafy2bzaced3ln6uq7t5vvvfrhot3m2nvbe4wapogyiwkymygvuc46wvuuqogk\"}],\"ParentWeight\":\"73069933935\",\"Height\":3094546,\"ParentStateRoot\":{\"/\":\"bafy2bzaceayytpshy4onplwqx6bbah2li26dq5qwyxs3i5p5a5ppmlmdprbg2\"},\"ParentMessageReceipts\":{\"/\":\"bafy2bzacecsks4vtbvmamvzhja447ami2qe3medzqh7q5i2cftj5j2s62cgyq\"},\"Messages\":{\"/\":\"bafy2bzacecgsziwhiesugjsrz4qkuguinw3ebk6tsa7is7hlj2nj3aizv6a72\"},\"BLSAggregate\":{\"Type\":2,\"Data\":\"jChb0//B0UO6YvhQ9CVWgnmQ65KCwuO8ftOVWzv37v6cqjPE11bxwF8pCZ+ktEBiEY//ONjM3O+4cAbaDRHj+0cNDt3cGB0lNja3xGGEZDHHBrI0B2FPpuWUGA8hlPU1\"},\"Timestamp\":1691142780,\"BlockSig\":{\"Type\":2,\"Data\":\"tuC8b+ocjh2YEK4CllU4R6mOZTzzcE41Teo+wDxNYzE2aHmWdC2kxdU6Y9Gkl4ksD3rwRcINbp1MgnGAoScNSwPzFOqRMjuY747HXDyB9Uma3VUt32gUMTUm8Noy+TgN\"},\"ForkSignaling\":0,\"ParentBaseFee\":\"234397907\"}",
	  "blockHash": "bafy2bzacebxm5vvowzqiuyrcwvgimlyaoyxd2qnxo4uxtro3tetww5sebp4vi",
	  "number": "3094546"
	}
	*/
	var block store.Block
	r := gjson.ParseBytes(value)
	hash := r.Get("blockHash").String()
	number := r.Get("number").String()
	blockBody := r.Get("block").String()
	bookRoot := gjson.Parse(blockBody)
	parentHash := bookRoot.Get("Parents").String()
	coinAddr := bookRoot.Get("Miner").String()
	blockTime := bookRoot.Get("Timestamp").String()
	baseFee := bookRoot.Get("ParentBaseFee").String()

	block.Id = uint64(time.Now().UnixNano())
	block.BlockHash = hash
	block.BlockNumber = fmt.Sprintf("%v", number)
	block.BlockTime = blockTime
	block.ParentHash = parentHash
	block.Coinbase = coinAddr
	block.BaseFee = baseFee
	return &block, nil
}

func GetTxFromKafka(value []byte) (*store.Tx, error) {
	/**
	{
	   m["blockNumber"] = blockNumber(可选)
	   m["blockHash"] = blockHash （可选）
	  [可选] "block": "{\"Miner\":\"f01264319\",\"Ticket\":{\"VRFProof\":\"iGTzfHwjx9sONc6iLxcyolaGmyf0zFGf5eBk5L/xNxuJbKUaITQKyg0lf5aesd6zDmA/syrPCW0h22hndls9z8hUzoeXKBKUyBQtMaTgk6hGxhrFBYsONj3iYlVodZmE\"},\"ElectionProof\":{\"WinCount\":1,\"VRFProof\":\"uA0+CVOgZhLF6u91WEW+roUp7/05+0JCjajzpysqFweYR8/eid749pEkjYN/RyKSFL6Kit8bARvUZtZhOVjCSBRoqngwdt9LS9BPLOF6kXlKn4bJ6OA98ZMNnD0q7Lml\"},\"BeaconEntries\":[{\"Round\":3190339,\"Data\":\"pJizw+DWXvWXR07I47H9Ckzw6q4JoZafGGazTvDKHEKPTT5ugJLztnNXTwKVS8kFDCenVCuZGHhPpyHwtmbsR3aFfalfVbxWqYKVatf1Wb7J4HdznLTOMelaNnQcaWMK\"}],\"WinPoStProof\":[{\"PoStProof\":3,\"ProofBytes\":\"mCtLcPQnsAiHQMlkSvLXG3lfTGncrxTgXKYHXJ2LygFbO4P7RTY2ANR9H/XexCm5tWd7wSvG9Y4KPfPFzXwSTJgDr8h2Pbda24T/TNSl6ZKX45qpwej+pm9DcRqfDGS9Cjm0AaRDjL+/gkBs3wON2Jj+iTggDBV67a8TaL/cYBfWumbIYVn/fwa6IqadedwCpKzf0UnDINNCSlE00qPn6tjXYAHY64aJobE8OMjlMJLcOxhMQyKW395EfHwpQN/Y\"}],\"Parents\":[{\"/\":\"bafy2bzacecje3zplhpl7wr33u5o2mcm4s35z3qcs3c2xubslonwq4geabaucy\"},{\"/\":\"bafy2bzacecuwzv2attmd7arl4wst4bpxjdsyuvgba657fi3hittfic5yl55o6\"},{\"/\":\"bafy2bzacea2rjdxarxkylgcscmfgncetvrq2usw5vgmcisghkox4vgnwuty4w\"},{\"/\":\"bafy2bzaceammtxnbebkj2rl276z7fo2bua5cxc2ome6s6ysi6dqecggqlmh6g\"}],\"ParentWeight\":\"73068677305\",\"Height\":3094494,\"ParentStateRoot\":{\"/\":\"bafy2bzacebw63prfd2lgpa2p5lfx4e2fsuh2aoli6a5mcnhd5fbymguczjykg\"},\"ParentMessageReceipts\":{\"/\":\"bafy2bzacebwx5ctbjtr74ss3fbm7wctoco226visyiz5iqhetejmsocqkqnsc\"},\"Messages\":{\"/\":\"bafy2bzaceagmdzb4nzo5hwfixj7xdiq5ncgi7r7a4pdjh2v4ejrhfp2jmjdag\"},\"BLSAggregate\":{\"Type\":2,\"Data\":\"ilA1MFVzrUQaq1nKkrwxBHpsa9ABWy4+0k1pcLIYX9a0G1B/lu3ihqnc4vxDcdwgBm0n6zxjOp7b4UgYZ3O9VAAr8LwfHCIlIUEMNOc0e9bSKu3nEI4CSCTixiTd9SM3\"},\"Timestamp\":1691141220,\"BlockSig\":{\"Type\":2,\"Data\":\"jFzmlAZT8CPP58/1B+X98nb7XsdGEfQ2CX522Z8gFFxIje1/MOalG7eBe2KlqRvVAJ6lJpNnv34L+9JeKdFMvXupXvxtcW41xo94QVv5fdToimeZO6jI0+GlE+S4tLGG\"},\"ForkSignaling\":0,\"ParentBaseFee\":\"445758004\"}",
	  "hash": "bafy2bzaceb2jsvg4bvzgf2uy4sqdajibgypgxhstid2gzbu3faz77nt6enfps",
	  "receipt": "{\"id\":1691143029869217000,\"blockHash\":\"0x55445443269a2de8a00f22ef9c64c4de3b18d6e84d7fa7274cb63cd398b11109\",\"logsBloom\":\"0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffeffefffffffffffff7fffffffffffffffbfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffbffffffffffffffffffffffffffffffffffffffffffffffffffffffffdffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffbffffffffffffffffffffffffff\",\"contractAddress\":\"\",\"transactionIndex\":\"0x60\",\"type\":\"0x2\",\"transactionHash\":\"0x749954dc0d7262ea98e4a0302501361e6b9e5340f46c869b2833ffb67e234af9\",\"gasUsed\":\"0x136740f\",\"blockNumber\":\"3094494\",\"cumulativeGasUsed\":\"0x0\",\"from\":\"0xff000000000000000000000000000000000a5194\",\"to\":\"0xff0000000000000000000000000000000007f51b\",\"effectiveGasPrice\":\"0x1d6296e0\",\"logs\":[],\"createTime\":\"2023-08-04\",\"status\":\"0x1\"}",
	  "tx": "{\"Version\":0,\"To\":\"f0521499\",\"From\":\"f3qm3mmhcbxlltk4e5lfjf5oldn5kxkr33avr4nszupfz5wgrqnihkclmeozf2ot7z6umecdqz6sw2jcsulacq\",\"Nonce\":17513,\"Value\":\"0\",\"GasLimit\":25351378,\"GasFeeCap\":\"4877209233\",\"GasPremium\":\"1437816\",\"Method\":5,\"Params\":\"hRWBggBAgYIOWMCK627TXRkVtz4jP9a4ZCdy/SfgeWUshXtJEYYlZPedygMO4eXK037hOx5u5rXclpy3K4Hr+5a4NJfXjcHU5Ptxjd2mvwdLYvcQXzqKqxMaoAk28zjM+iFgzeRxwzS4r0AKAzpLhcc7e+CV7LhavqlxicbFFRaYLWVQ0vbwFUtKva1KqcRF2ClTkVnreVspTpiJsY0DHArzxE29IU3rHKhsTFQXkjqrjoqLlrPxeyW+UP3MDqH5wvLsib5LT0NppVsaAC83xVgghYbRFUdLq5DXXWtgOIMe/+rjwEESSL/TcR9umo+bMQE\u003d\",\"CID\":{\"/\":\"bafy2bzaceb2jsvg4bvzgf2uy4sqdajibgypgxhstid2gzbu3faz77nt6enfps\"}}"
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

	//status := txRoot.Get("ret.0.contractRet").String()
	//hash := txRoot.Get("txID").String()
	//blockHash := txRoot.Get("raw_data.ref_block_hash").String()
	//txTime := txRoot.Get("raw_data.timestamp").Uint()
	limit := txRoot.Get("GasLimit").String()
	//txType := txRoot.Get("raw_data.contract.0.type").String()
	//v := txRoot.Get("raw_data.contract.0.parameter.value")
	from := txRoot.Get("From").String()
	to := txRoot.Get("To").String()
	input := txRoot.Get("Params").String()
	txValue := txRoot.Get("Value").String()
	gasPremium := txRoot.Get("GasPremium").String()
	//gasFeeCap:=txRoot.Get("GasFeeCap").String()

	var txType string
	if len(input) > 1 {
		txType = "1"
	} else {
		txType = "2"
	}

	tx.Id = uint64(time.Now().UnixNano())
	tx.Value = txValue
	//tx.BlockHash = blockId
	//tx.TxHash = hash
	//tx.TxStatus = status
	//tx.TxTime = fmt.Sprintf("%v", txTime)
	tx.ToAddr = to
	tx.FromAddr = from
	tx.Type = txType
	tx.InputData = input
	tx.GasLimit = limit
	tx.PriorityFee = gasPremium
	receiptBody := root.Get("receipt").String()
	if len(receiptBody) > 5 {
		receiptRoot := gjson.Parse(receiptBody)
		if len(tx.BlockHash) < 1 {
			tx.BlockHash = receiptRoot.Get("blockHash").String()
		}
		if len(tx.BlockNumber) < 1 {
			tx.BlockNumber = receiptRoot.Get("blockNumber").String()
		}
		tx.GasUsed = receiptRoot.Get("gasUsed").String()
		tx.TxStatus = receiptRoot.Get("status").String()
	}

	blockBody := root.Get("block").String()
	if len(blockBody) > 5 {
		blockRoot := gjson.Parse(blockBody)
		tx.BaseFee = blockRoot.Get("ParentBaseFee").String()
		tx.TxTime = blockRoot.Get("Timestamp").String()
	}

	return &tx, nil
}

func GetTxType(body []byte) (uint64, error) {
	root := gjson.ParseBytes(body)
	txBody := root.Get("tx").String()
	input := gjson.Parse(txBody).Get("Params").String()
	if len(input) > 5 {
		//合约调用
		return 1, nil
	} else {
		//普通资产转移
		return 2, nil
	}
}
func ParseTx(value []byte, transferTopic string, blockchain int64) (*store.SubTx, error) {
	var tx store.SubTx
	root := gjson.ParseBytes(value)

	tx.BlockChain = uint64(blockchain)
	tx.Id = uint64(time.Now().UnixNano())
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

	//todo
	gasLimit := txRoot.Get("GasLimit").String()
	from := txRoot.Get("From").String()
	to := txRoot.Get("To").String()
	input := txRoot.Get("Params").String()
	txValue := txRoot.Get("Value").String()
	gasPremium := txRoot.Get("GasPremium").String()
	//gasFeeCap:=txRoot.Get("GasFeeCap").String()

	var txType uint64
	if len(input) > 1 {
		txType = 1
	} else {
		txType = 2
	}

	tx.Id = uint64(time.Now().UnixNano())
	tx.Value = util.Div(txValue, 18)
	tx.To = to
	tx.From = from
	tx.TxType = txType
	tx.Input = input

	var gasUsed, baseFee string

	//tx = limit
	receiptBody := root.Get("receipt").String()
	if len(receiptBody) > 5 {
		receiptRoot := gjson.Parse(receiptBody)
		if len(tx.BlockHash) < 1 {
			tx.BlockHash = receiptRoot.Get("blockHash").String()
		}
		if len(tx.BlockNumber) < 1 {
			tx.BlockNumber = receiptRoot.Get("blockNumber").String()
		}
		gasUsed = receiptRoot.Get("gasUsed").String()
		gasUsed, _ = util.HexToInt(gasUsed)
		status := receiptRoot.Get("status").String()
		if status == "0x0" {
			tx.Status = 0
		} else if status == "0x1" {
			tx.Status = 1
		}
	}

	blockBody := root.Get("block").String()
	if len(blockBody) > 5 {
		blockRoot := gjson.Parse(blockBody)
		txTime := blockRoot.Get("Timestamp").String()
		txTime = fmt.Sprintf("%v000", txTime)
		tx.TxTime = txTime
		baseFee = blockRoot.Get("ParentBaseFee").String()
	}

	if len(baseFee) > 0 && len(gasUsed) > 0 && len(gasLimit) > 0 && len(gasPremium) > 0 {
		var ok bool
		var bfee, used, limit, premium *big.Int
		bfee, ok = new(big.Int).SetString(baseFee, 0)
		used, ok = new(big.Int).SetString(gasUsed, 0)
		limit, ok = new(big.Int).SetString(gasLimit, 0)
		premium, ok = new(big.Int).SetString(gasPremium, 0)

		if ok {
			burn := bfee.Mul(bfee, used)
			miner := limit.Mul(limit, premium)
			fee := burn.Add(burn, miner)
			tx.Fee = util.Div(fee.String(), 18)
		}

	}

	return &tx, nil
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

	fromAddr := txRoot.Get("From").String()
	txAddressList[GetCoreAddr(fromAddr)] = 1

	toAddr := txRoot.Get("To").String()
	txAddressList[GetCoreAddr(toAddr)] = 1

	if root.Get("receipt").Exists() {
		receipt := root.Get("receipt").String()
		receiptRoot := gjson.Parse(receipt)
		list := receiptRoot.Get("logs").Array()
		for _, v := range list {
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
