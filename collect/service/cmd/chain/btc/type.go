package btc

/**

  {
      "txid": "b9618d971a472d59c219c3b8071291e502686d7d21c8e7e044e63ebb5b198580",
      "vout": 1,
      "scriptSig": {
          "asm": "",
          "hex": ""
      },
      "txinwitness": [
          "304502210095d6aa9adb36b654f712c5f8f15fde292b703b0f6d9611a089d8d6559a176e0c022008eae9ec8ace176c676b1a9ba324dc762e370471e4825b4130de79a6f0afdf4101",
          "023e86a944b9cff6bea3e61d3594820cacbfb122b817934f981256b2bc4e6b8504"
      ],
      "prevout": {
          "generated": false,
          "height": 808535,
          "value": 1.35084678,
          "scriptPubKey": {
              "asm": "0 3877cc298fc5b8a591e34982a91183bb35f0c242",
              "desc": "addr(bc1q8pmuc2v0cku2ty0rfxp2jyvrhv6lpsjzq9y6s8)#44qaxdmw",
              "hex": "00143877cc298fc5b8a591e34982a91183bb35f0c242",
              "address": "bc1q8pmuc2v0cku2ty0rfxp2jyvrhv6lpsjzq9y6s8",
              "type": "witness_v0_keyhash"
          }
      },
      "sequence": 4294967295
  }

*/

type BVin struct {
	Sequence  int64 `json:"sequence" gorm:"column:sequence"`
	ScriptSig struct {
		Asm string `json:"asm" gorm:"column:asm"`
		Hex string `json:"hex" gorm:"column:hex"`
	} `json:"scriptSig" gorm:"column:scriptSig"`
	Prevout     *Prevout `json:"prevout" gorm:"column:prevout"`
	Txid        string   `json:"txid" gorm:"column:txid"`
	Txinwitness []string `json:"txinwitness" gorm:"column:txinwitness"`
	Vout        int64    `json:"vout" gorm:"column:vout"`
}

type Prevout struct {
	ScriptPubKey struct {
		Address string `json:"address" gorm:"column:address"`
		Asm     string `json:"asm" gorm:"column:asm"`
		Hex     string `json:"hex" gorm:"column:hex"`
		Type    string `json:"type" gorm:"column:type"`
		Desc    string `json:"desc" gorm:"column:desc"`
	} `json:"scriptPubKey" gorm:"column:scriptPubKey"`
	Generated bool   `json:"generated" gorm:"column:generated"`
	Value     string `json:"value" gorm:"column:value"`
	Height    int    `json:"height" gorm:"column:height"`
	Bestblock string `json:"bestblock"`
}
