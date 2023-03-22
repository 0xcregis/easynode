
easynode_task是easynode系统的基础和核心服务，是其他服务的运行的必要条件。
该服务负责任务产生、任务分发、系统监控等功能。

## Prerequisites
- go version >=1.16
- easynode_collect 服务已完成部署

## Building the source

(以linux系统为例)
- mkdir easynode & cd easynode
- git clone https://github.com/uduncloud/easynode_task.git
- cd easynode_task
- CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o easynode_task app.go
  (mac下编译linux程序为例，其他交叉编译的命令请自行搜索)

- ./easynode_task -config ./config.json

## config.json 详解

``````
{
  "NodeTaskDb": {//待执行任务表
    "Addr": "192.168.2.11",//地址
    "Port": 3306,//端口
    "User": "root",//用户名
    "Password": "123456789",//密码
    "DbName": "easy_node",//数据库
    "Table": "node_task" //表名
  },
  "NodeSourceDb": {//待分配的任务表
    "Addr": "192.168.2.11",
    "Port": 3306,
    "User": "root",
    "Password": "123456789",
    "DbName": "easy_node",
    "Table": "node_source"
  },
  "NodeInfoDb": {//节点表
    "Addr": "192.168.2.11",
    "Port": 3306,
    "User": "root",
    "Password": "123456789",
    "DbName": "easy_node",
    "Table": "node_info"
  },
  "NodeErrorDb": {//数据缺失表
    "Addr": "192.168.2.11",
    "Port": 3306,
    "User": "root",
    "Password": "123456789",
    "DbName": "easy_node",
    "Table": "node_error"
  },
  "Chains": [ //公链配置
    {
      "NodeHost": "https://eth-mainnet.g.alchemy.com/v2",//三方区块链节点的地址
      "NodeKey": "**********************",//三方区块链节点key
      "BlockChainName": "eth",//公链名称
      "BlockChainCode": 200,//公链代码
      "BlockMin": 16103500,//区块最低高度
      "BlockMax":16125881//区块最高高度，如果是0，则表示时时获取公链最新高度
    },
    {
      "NodeHost": "https://api.trongrid.io",
      "NodeKey": "************************",
      "BlockChainName": "tron",
      "BlockChainCode": 205,
      "BlockMin": 47153472,
      "BlockMax":0//区块最高高度，如果是0，则表示时时获取公链最新高度
    }
  ]
}

``````

## usages

启动服务后，等待分配任务并执行
