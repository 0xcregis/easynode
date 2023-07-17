task是easynode系统的基础和核心服务，是其他服务的运行的必要条件。
该服务负责任务产生、任务分发、系统监控等功能。

## Prerequisites

- go version >=1.20
- collect 服务已完成部署

## Building the source

(以linux系统为例)

- mkdir easynode & cd easynode
- git clone https://github.com/0xcregis/easynode.git
- cd easynode/cmd/task
- CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o easynode_task app.go
  (mac下编译linux程序为例，其他交叉编译的命令请自行搜索)

- ./easynode_task -config ./config.json

## config.json 详解

``````
{
  "TaskKafka": { //task kafka 配置
    "Host": "192.168.2.20",
    "Port": 9092
  },
  "AutoCreateBlockTask": true, //是否自动产生区块任务
  "Chains": [// 公链配置 和 AutoCreateBlockTask 配合使用，当AutoCreateBlockTask=true 时，必需
    {
      "NodeHost": "https://eth-mainnet.g.alchemy.com/v2",//公链节点地址
      "NodeKey": "RzxBjjh_c4y0LVHZ7GNm8zoXEZR3HYop", //公链节点 token 非必需
      "BlockChainName": "eth",
      "BlockChainCode": 200,
      "BlockMin": 17113233, //最小区块高度
      "BlockMax": 0 //最大区块高度，如果是0:则 时时获取公链的最新高度
    },
    {
      "NodeHost": "https://api.trongrid.io",
      "NodeKey": "244f918d-56b5-4a16-9665-9637598b1223",
      "BlockChainName": "tron",
      "BlockChainCode": 205,
      "BlockMin": 50563524,
      "BlockMax": 0
    }
  ]
}

``````

## usages

启动服务后，等待分配任务并执行
