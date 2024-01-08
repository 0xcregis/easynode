package blockchain

type WsReqMessage struct {
	Id     int64
	Code   int64
	Params map[string]string
}

type WsRespMessage struct {
	Id     int64
	Code   int64
	Status int //0:成功 1：失败
	Err    string
	Params map[string]string
	Resp   interface{}
}

/**
  {
            "action": {
                "from": "0x4e68ccd3e89f51c3074ca5072bbac773960dfa36",
                "callType": "staticcall",
                "gas": "0x26e23",
                "input": "0x70a082310000000000000000000000004e68ccd3e89f51c3074ca5072bbac773960dfa36",
                "to": "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
                "value": "0x0"
            },
            "blockHash": "0x6b3f8639a57d744dcbbbe0d84bca25d7ea776282d6f53f329db3cc9408416115",
            "blockNumber": 18938364,
            "result": {
                "gasUsed": "0x9e6",
                "output": "0x00000000000000000000000000000000000000000000036e79c74a878ae9c2cd"
            },
            "subtraces": 0,
            "traceAddress": [
                1,
                1
            ],
            "transactionHash": "0x999cabe1fcca80148290827a8c655734531615cb22d30faa222ec7a67928587b",
            "transactionPosition": 201,
            "type": "call"
        }
*/

type InterTx struct {
	Result struct {
		Output  string `json:"output" gorm:"column:output"`
		GasUsed string `json:"gasUsed" gorm:"column:gasUsed"`
	} `json:"-" gorm:"column:result"`
	BlockHash           string `json:"blockHash" gorm:"column:blockHash"`
	TransactionPosition int    `json:"-" gorm:"column:transactionPosition"`
	BlockNumber         int    `json:"blockNumber" gorm:"column:blockNumber"`
	TraceAddress        []int  `json:"-" gorm:"column:traceAddress"`
	Action              struct {
		Input    string `json:"-" gorm:"column:input"`
		Gas      string `json:"gas" gorm:"column:gas"`
		From     string `json:"from" gorm:"column:from"`
		To       string `json:"to" gorm:"column:to"`
		Value    string `json:"value" gorm:"column:value"`
		CallType string `json:"callType" gorm:"column:callType"`
	} `json:"action" gorm:"column:action"`
	Type            string `json:"-" gorm:"column:type"`
	Subtraces       int    `json:"-" gorm:"column:subtraces"`
	TransactionHash string `json:"transactionHash" gorm:"column:transactionHash"`
}
