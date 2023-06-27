store是easynode系统的基础和核心服务. 该服务负责任监控地址管理、用户订阅和取消、数据推送、数据落盘等功能。

## Prerequisites

- go version >=1.20
- collect 服务已完成部署

## Building the source

(以linux系统为例)

- mkdir easynode & cd easynode
- git clone https://github.com/0xcregis/easynode.git
- cd easynode/cmd/store
- CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o easynode_store app.go
  (mac下编译linux程序为例，其他交叉编译的命令请自行搜索)

- ./easynode_store -config ./config.json

## config.json 详解

``````
{
  "RootPath": "/api/store", //根路径
  "Port": 9003, //端口
  "BaseDb": { //基础库的配置
    "Addr": "192.168.2.9",
    "Port": 9000,
    "User": "test",
    "Password": "test",
    "DbName": "base",
    "AddressTable": "address", //地址表
    "TokenTable": "token" //token表
  },
  "Chains": [ //公链配置
   {
  "BlockChain": 200, //公链代码
  "BlockStore": false, //区块是否落盘
  "TxStore": false, //交易是否落盘
  "ReceiptStore": false, //收据是否落盘
  "SubStore": true,//订阅数据持久化
  "Kafka": { //公链数据所在的Kafka配置
    "SubTx": { //订阅
          "Host": "192.168.2.9",
          "Port": 9092,
          "Topic": "ether_sub_tx",
          "Group": "1",
          "StartOffset": 0,
          "Partition": 0
        },
    "Tx": { //交易配置
      "Host": "192.168.2.20",
      "Port": 9092,
      "Topic": "ether_tx2",
      "Group": "1", //group 后缀
      "StartOffset": 0, //read Kafka 开始未知 -1: latest,-2: first ,0:commited offset
      "Partition": 0
    },
    "Block": { //区块配置
      "Host": "192.168.2.20",
      "Port": 9092,
      "Topic": "ether_block2",
      "Group": "1",
      "StartOffset": 0,
      "Partition": 0
    },
    "Receipt": { //收据配置
      "Host": "192.168.2.20",
      "Port": 9092,
      "Topic": "ether_receipt2",
      "Group": "1",
      "StartOffset": 0,
      "Partition": 0
    }
  },
  "ChainDb": { //clickhouse 数据库配置
    "Addr": "192.168.2.11",
    "Port": 9000,
    "User": "test",
    "Password": "test",
    "DbName": "ether",
    "TxTable": "tx", //交易表
    "BlockTable": "block", //区块表
    "ReceiptTable": "receipt", //收据表
    "SubTxTable": "sub_tx" //订阅数据交易表
  }
 }
]
}

``````

## usages

- http协议

``````
//请求生产token,每个用户 一个即可
curl --location --request POST 'localhost:9003/api/store/monitor/token' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "123@gmail.com"
}'

//提交监控地址
curl --location --request POST 'localhost:9003/api/store/monitor/address' \
--header 'Content-Type: application/json' \
--data-raw '{
    "blockChain": 200,//非必需，如果不传默认0，则表示 跨链监控
    "address": "0x28c6c06298d514db089934071355e5743bf21d61",
    "token": "5fe5f231-7051-4caf-9b52-108db92edbb4"
}'

//查询监控地址
curl --location --request POST 'localhost:9003/api/store/monitor/address/get' \
--header 'Content-Type: application/json' \
--data-raw '{
    "token": "5fe5f231-7051-4caf-9b52-108db92edbb4"
}'

//删除监控地址
curl --location --request POST 'localhost:9003/api/store/monitor/address/delete' \
--header 'Content-Type: application/json' \
--data-raw '{
    "blockChain": 200,//非必需，如果不传默认0，则表示 跨链监控
    "address": "0x28c6c06298d514db089934071355e5743bf21d61",
    "token": "5fe5f231-7051-4caf-9b52-108db92edbb4"
}'

``````

- ws

``````

//入参数据结构：
type WsReqMessage struct {
	Id         int64 //客户端请求序列号
	Code       int64 //1:订阅资产转移交易，2:取消订阅资产转移交易
	BlockChain []int64 `json:"blockChain"` //订阅公链的代码
	Params     map[string]string //非必需
}

//返回数据结构：
type WsRespMessage struct {
	Id         int64 //请求的序列号，和请求保持一致
	Code       int64 //命名code，和请求保持一致
	BlockChain []int64 `json:"blockChain"` //订阅公链的代码
	Status     int   //0:成功 1：失败
	Err        string //错误原因
	Params     map[string]string //请求参数，和请求保持一致
	Resp       interface{} //返回的数据
}
//推送数据结构
type WsPushMessage struct {
	Code       int64 //1:tx,2:block,3:receipt //推送数据业务码
	BlockChain int64 `json:"blockChain"` //公链代码
	Data       interface{} //推送的数据
}

``````
- 提交订阅并接受返回

``````
   url: 
   
   ws://localhost:9003/api/store/ws/{token} 
   
   ws://localhost:9003/api/store/ws/{token}?chainCode={chainCode}
  
   入参：
           {
             "id":1001,
             "code":1,
             "blockChain":[200],
             "Params":{}
            }
   
   订阅返回：
   
            {
              "Id": 1001,
              "Code": 1,
              "blockChain": [200],
              "Status": 0,
              "Err": "",
              "Params": {
              },
              "Resp": null
            }
            
   push 返回：
   
   {
      "Code": 1, //消息类型，1:交易消息
      "blockChain": 200, //公链代码
      "Data": { //交易数据
        "id": 1685094437357929000,
        "blockHash": "0x067fbc694c5ca3540ee965b25c286e55d40f3e5e5fd336d1f398868dfc18feec", //区块hash
        "blockNumber": "17284552", //区块高度
        "chainCode": 200,
        "contractTx": [ //合约交易时EVM 事件
          {
            "contract": "0xdac17f958d2ee523a2206206994597c13d831ec7", //合约地址
            "from": "0x54b50187becd0bbcfd52ec5d538433dab044d2a8", //from 地址
            "method": "Transfer", //合约方法
            "to": "0x408be4b8a862c1a372976521401fd77f9a0178d7", //to 地址
            "value": "59.327379" //交易内容
          }
        ],
        "fee": "0.002053016771146819",//交易费
        "from": "0x54b50187becd0bbcfd52ec5d538433dab044d2a8", //from 地址
        "hash": "0x323c08a889ed99d8bfc6c72b1580432f7a13ca7c992fd1bac523e46bfe7ab98f", //交易hash
        "status": "1", //交易状态 1:成功, 0:交易失败
        "to": "0xdac17f958d2ee523a2206206994597c13d831ec7", //to地址
        "txTime": "1684390019", //交易时间
        "txType": 1, //交易类型 1:合约调用，2:普通资产转移
        "value": "0" //交易额 
      }
  }         
            
                        
   
``````


- notes

    - 同一token ，同一时刻仅能有一个 订阅命令，需要订阅其他命令，需要先取消已订阅的命令
    - 同一token ，多次连接时，会自动关闭上一个连接
    - 客户端 必需 实现 ping ,pong命令，长时间未收到客户端发出ping ，则会自动关闭连接
    - 服务端 定时的发送ping 命令，客户端收到时，需要及时返回 pong命令
    - 如果在路径中指定blockChain,则在消息中 blockChain 则失效
    - 【ws://localhost:9003/api/store/ws/{token}?chainCode={chainCode}】 仅在 同一个token 需要多次连接同一个服务节点时使用(订阅不同链的数据)，但仍然推荐 使用一个连接同时订阅多个链的方案


