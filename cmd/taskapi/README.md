### 概述

该服务是easynode_task的API服务，负责提交任务（包括：交易，区块，收据等各种任务）、查询节点等操作

### Prerequisites

- go version>=1.18

### Building the source

(以linux系统为例)
- mkdir easynode & cd easynode
- git https://github.com/uduncloud/easynode_taskapi.git
- cd easynode_taskapi
- CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o easynode_taskapi app.go
  (mac下编译linux程序为例，其他交叉编译的命令请自行搜索)
  
  

### config.json 解释

``````
{
  "RootPath": "/api/task",//根目录
  "Port": 9001, //端口
  "BlockChain": [200,205],//公链代码
  "TaskDb": { //任务数据库配置
    "Addr": "192.168.2.11",//地址
    "Port": 3306,//端口
    "User": "root",//用户名
    "Password": "123456789",//密码
    "DbName": "easy_node" //database
  }
}

``````  

### usages

``````
//获取活跃节点接口
curl --location --request GET 'http://127.0.0.1:9001/api/task/node' \

//提交区块任务
curl --location --request POST 'http://127.0.0.1:9001/api/task/block' \
--header 'Content-Type: application/json' \
--data-raw '{
    "blockChain": 200,
    "blockHash": "",
    "blockNumber": "16389175"
}'

//提交单个交易任务
curl --location --request POST 'http://127.0.0.1:9001/api/task/tx' \
--header 'Content-Type: application/json' \
--data-raw '{
    "blockChain": 200,
    "txHash": "0xc0e81699d2728694cc275521daa9b89414a9e4749f7418c8c69f26b090c99f44"
}'

//提交耽搁收据任务
curl --location --request POST 'http://127.0.0.1:9001/api/task/receipt' \
--header 'Content-Type: application/json' \
--data-raw '{
    "blockChain": 200,
    "txHash": "0xc0e81699d2728694cc275521daa9b89414a9e4749f7418c8c69f26b090c99f44"
}'

//提交区块下所有交易任务
curl --location --request POST 'http://127.0.0.1:9001/api/task/txs' \
--header 'Content-Type: application/json' \
--data-raw '{
    "blockChain": 200,
    "blockHash": "",
    "blockNumber": "16389175"
}'

//提交区块下所有收据任务
curl --location --request POST 'http://127.0.0.1:9001/api/task/txs' \
--header 'Content-Type: application/json' \
--data-raw '{
    "blockChain": 200,
    "blockHash": "",
    "blockNumber": "16389175"
}'

``````
