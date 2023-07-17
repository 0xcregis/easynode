### 概述

本服务转发用户端发起请求到区块链相应的节点并在此期间提供最优节点，主要应用在提交交易、余额查阅，查询矿工费等场景

### 架构简述

1.网关由gin实现，负责接受客户端发起的请求，

2.收到请求并通过验证后，将请求转发到service中

3.判断区块链类型且选择最佳的节点

4.把请求数据 发送到最佳节点，且等待响应

5.把请求的最终结果，告知客户端

### 限制

- 目前仅支持 ether、 tron 等2种公链

### Prerequisites

- go version>=1.18

### Building the source

(以linux系统为例)

- mkdir easynode & cd easynode
- git clone https://github.com/0xcregis/easynode.git
- cd easynode/cmd/blockchain
- CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o easynode_chain app.go
  (mac下编译linux程序为例，其他交叉编译的命令请自行搜索)

- ./easynode_chain -config ./config.json

### config.json 解释

``````
{
  "RootPath": "/api/chain", //api根目录
  "Port": 9002, //端口
  "BlockChain": [200,205], //支持的公链代码
  "Cluster": {//区块链节点集群
    "200": [{ //ether 节点配置
      "NodeUrl": "https://eth-mainnet.g.alchemy.com/v2", //节点地址
      "NodeToken": "************************" //节点需要token
    }],
    "205": [{ //tron 节点配置
      "NodeUrl": "https://api.trongrid.io",//节点地址
      "NodeToken": "********************"//节点需要token
    }]
  }
}

``````

### usages

- http 协议

``````
//发送交易接口
curl -X POST \
  http://127.0.0.1:9002/api/chain/tx/sendRawTransaction \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: 84336133-4ab2-47c4-b600-7b56bfdd79d9' \
  -H 'cache-control: no-cache' \
  -d '{
    "chain":200,
	"signed_tx":"0xe5cB067E90D5Cd1F8052B83562Ae670bA4A211a8"
}'

//查询地址余额接口
curl -X POST \
  http://127.0.0.1:9002/api/chain/account/balance \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: 3eac86c7-c5aa-4dd3-93d3-a3b23678512d' \
  -H 'cache-control: no-cache' \
  -d '{
    "chain":200,
	"address":"0xe5cB067E90D5Cd1F8052B83562Ae670bA4A211a8",
	"tag":"latest"
}'

curl -X POST \
  http://127.0.0.1:9002/api/chain/account/balance \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: 7b49c120-1bfe-4223-9a7a-e1d2620f3099' \
  -H 'cache-control: no-cache' \
  -d '{
     "chain":205,
	"address":"TUtAk64jJqdf1pY3SiHeooVikP2SFWXjZ6"
}'

//Ether 链的 nonce值
curl -X POST \
  http://127.0.0.1:9002/api/chain/account/nonce \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: 8fe90095-369c-453c-8d59-6fa18b0a83a0' \
  -H 'cache-control: no-cache' \
  -d '{
    "chain":200,
	"address":"0xe5cB067E90D5Cd1F8052B83562Ae670bA4A211a8",
	"tag":"latest"
}'

// 区块链的最新区块高度
curl -X POST \
  http://127.0.0.1:9002/api/chain/block/latest \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: 8d1d5eab-482a-4f0b-b31e-98ab59daf924' \
  -H 'cache-control: no-cache'
  -d '{
     "chain":205
}'


// 代币余额查询
curl -X POST \
  http://127.0.0.1:9002/api/chain/account/tokenBalance \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: 380f20a5-b023-42bd-961c-ac5f84c580e1' \
  -H 'cache-control: no-cache' \
  -d '{
    "chain":205,
	"address":"TMuA6YqfCeX8EhbfYEg5y7S4DqzSJireY9",
	"codeHash":"TLa2f6VPqDgRE67v1736s7bJ8Ray5wYjU7",
}'


//区块链通用接口，基于区块链标准，以http-rpc协议实现
curl -X POST \
  http://127.0.0.1:9002/api/chain/jsonrpc \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: 3d6a170b-bd88-42b5-83d1-d3e5bec7babf' \
  -H 'cache-control: no-cache' \
  -d '{
     "chain":205,
     "data":{
          "id": 1,
          "jsonrpc": "2.0",
          "params": [
             "0xe5cB067E90D5Cd1F8052B83562Ae670bA4A211a8",
              "latest"
          ],
          "method": "eth_getBalance"
          }
  }'

//根据hash 查询区块
curl -X POST \
  http://127.0.0.1:9002/api/chain/block/hash \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: 49c02db3-51a2-4a23-8fbd-f78c50cbd8c7' \
  -H 'cache-control: no-cache' \
  -d '{
    "chain":205,
	"hash":"0000000002f2f66e7256eaffa627b521c380f7dcc4d354bf6c7a5ed8e0c4ea72"
}'

//根据高度 查询区块
curl -X POST \
  http://127.0.0.1:9002/api/chain/block/number \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: 3f33d2b1-9a01-4a48-9f9a-c49334ac2903' \
  -H 'cache-control: no-cache' \
  -d '{
    "chain":205,
	"number":"49477110"
}'

//查询交易
curl -X POST \
  http://127.0.0.1:9002/api/chain/tx/hash \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: 72b4bc45-c235-4361-b88b-2ccffd42a384' \
  -H 'cache-control: no-cache' \
  -d '{
    "chain":205,
	"hash":"89afac2142e025a13987ed183444ec90e9dcb8028bc7bc0757a21c654aa78b31"
}'

//查询收据
curl -X POST \
  http://127.0.0.1:9002/api/chain/receipts/hash \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: 9e33035b-fef4-4dfd-a7f3-fd0d2f3de318' \
  -H 'cache-control: no-cache' \
  -d '{
    "chain":205,
	"hash":"89afac2142e025a13987ed183444ec90e9dcb8028bc7bc0757a21c654aa78b31"
}'

``````

- ws 协议

``````
入参结构体：
type WsReqMessage struct {
	Id     int64
	Code   int64
	Params map[string]string
}

输出结构体：
type WsRespMessage struct {
	Id     int64
	Code   int64
	Status int //0:成功 1：失败
	Err    string
	Params map[string]string
	Resp   interface{}
}

1. code=1 //交易查询
   params:
   blockChain int //公链编号
   txHash string //交易hash
   
2. code=2 //区块查询
   params:
   blockChain int //公链编号
   blockHash string //区块hash  
   
3. code=3 //区块查询
   params:
   blockChain int //公链编号
   number string //区块高度 
   
4. code=4 //最新区块
   params:
   blockChain int //公链编号
            

5. code=8 //广播交易
   params:
   blockChain int //公链编号
   signed string //签名后交易
   
6. code=9 //jsonRpc 请求
   params:
   blockChain int //公链编号
   jsonRpc string //jsonRpc  
   
7. code=10 //主币余额查询
   params:
   blockChain int //公链编号
   address string //账户
 
 8. code=11 //代币余额查询
   params:
   blockChain int //公链编号
   address string //账户  
   contract string //代币合约  
   
 9. code=12 //nonce
   params:
   blockChain int //公链编号
   address string //账户  
     
``````
