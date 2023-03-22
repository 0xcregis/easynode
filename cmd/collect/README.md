## easynode_collect

easynode是对区块链节点进一步抽象和简化，封装了普通区块链节点的技术难点，是面向业务方的节点。使业务方即使不了解区块链概念和技术的情况下，基于现有技术栈和经验，仍然可以设计自己基于区块链的业务。具有一下特征：
- 支持多种公链：目前已支持 ether、tron等公链
- 功能可配制化：通过配置文件可以开启或关闭部分功能
- 动态横向扩展：在不重启系统情况下，可以启动新节点，分担当前系统的压力
- 冗余备份：可以配置多个区块链节点，提供系统稳定性，防止部分节点的异常导致业务系统异常
- 自带负载均衡：在不依赖第三方组件的情况下，通过2层均衡设计，解决了单节点性能瓶颈问题
- 过程可视化：整个系统执行过程都是可视化的，方便监控和维护
- 同步数据方式多样：除了通过配置文件指定自动同步数据外，用户还可以通过HTTP协议自定义需要同步的数据
- 系统健壮性：系统自带异常重试、错误检查等功能，解决因网络异常等客观原因导致数据丢失的问题

easynode_collect是easynode系统的基础和核心服务，是其他服务的运行的必要条件。该服务负责同步主网区块数据的服务，根据用户配置的规则，自动同步主网数据到本地，easynode_collect服务支持横向扩展和负载均衡。


## Prerequisites

- go version: >=1.16
- clickhouse 部署和配置
   [详见](https://github.com/uduncloud/easynode_collect/wiki/clickhouse-%E9%85%8D%E7%BD%AE%E5%92%8C%E9%83%A8%E7%BD%B2)
- kafka 部署和配置
   [详见](https://github.com/uduncloud/easynode_collect/wiki/kafka-%E9%85%8D%E7%BD%AE%E5%92%8C%E9%83%A8%E7%BD%B2)
- etcd  部署和配置
   [详见](https://github.com/uduncloud/easynode_collect/wiki/etcd-%E9%85%8D%E7%BD%AE%E5%92%8C%E9%83%A8%E7%BD%B2)
- mysql 部署和配置
   [详见](https://github.com/uduncloud/easynode_collect/wiki/mysql-%E9%85%8D%E7%BD%AE%E5%92%8C%E9%83%A8%E7%BD%B2)

## Building the source
(以linux系统为例)
- mkdir easynode & cd easynode
- git clone https://github.com/uduncloud/easynode_collect.git
- cd easynode_collect
- CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o easynode_collect app.go
(mac下编译linux程序为例，其他交叉编译的命令请自行搜索)

- ./easynode_collect -config ./config.json

## config.json 详解
(不要copy,使用代码库中配置文件，作为模版)
``````
{
  "NodeTaskDb": { //待执行任务表
    "Addr": "192.168.2.11",//mysql 地址（案例内网地址，请勿直接使用）
    "Port": 3306, //端口
    "User": "root",//用户
    "Password": "456789aqet&*#", //密码
    "DbName": "easy_node", //数据库名称
    "Table": "node_task" //表名
  },
  "NodeSourceDb": {//待分配任务表
    "Addr": "192.168.2.11",
    "Port": 3306,
    "User": "root",
    "Password": "123456789",
    "DbName": "easy_node",
    "Table": "node_source"
  },
  "NodeInfoDb": { //部署节点信息表
    "Addr": "192.168.2.11",
    "Port": 3306,
    "User": "root",
    "Password": "123456789",
    "DbName": "easy_node",
    "Table": "node_info"
  },
  "Log": { //日志配置
    "Path": "./log/collect", //日志文件相对路径
    "Delay": 2 //日志清理延迟时间，单位天
  },
  "Chains": [ //公链配置，可以配置多个，目前暂支持 ether、tron
    {
      "Etcd": { //etcd 服务器配置
        "Host": "http://192.168.2.20:2379", //etcd服务器地址
        "Key": "eth", //存储键的前缀，要确保不同链的唯一性
        "Ttl": 10 //有效期 单位分钟
      },
      "Kafka": { //kafka 配置
        "Host": "192.168.2.20", //地址
        "Port": 9092 //端口
      },
      "BlockChainName": "eth",//公链名称1
      "BlockChainCode": 200, //公链的代码
      "NodeWeight":11, //该节点权重
      "PullReceipt": false,//是否自动产生收据任务，默认false，建议 false
      "PullTx": false, //是否自动产生交易任务，默认false，建议 false
      "BlockTask": {//区块任务配置,缺省该配置 则表示 不执行该类型任务
        "FromCluster": [//区块链节点配置，支持多个
          {
            "Host": "https://eth-mainnet.g.alchemy.com/v2",//三方服务
            "Key": "*************************" //三方服务key
          }
        ],
        "Kafka": {//区块任务对应的Kafka
          "Topic": "ether_block",//topic
          "Partition": 0 //partition
        }
      },
      "TxTask": { //交易任务配置 ,缺省该配置 则表示 不执行该类型任务
        "FromCluster": [ //区块链节点配置，支持多个
          {
            "Host": "https://eth-mainnet.g.alchemy.com/v2",
            "Key": "***********************"
          }
        ],
        "Kafka": {//交易任务对应的Kafka
          "Topic": "ether_tx",
          "Partition": 0
        }
      },
      "ReceiptTask": {//收据任务配置,缺省该配置 则表示 不执行该类型任务
        "FromCluster": [//区块链节点配置，支持多个
          {
            "Host": "https://eth-mainnet.g.alchemy.com/v2",
            "Key": "*****************************"
          }
        ],
        "Kafka": {//收据任务对应的Kafka
          "Topic": "ether_receipt",
          "Partition": 0
        }
      }
    }
  ]
}

``````

## usages

启动服务后，等待分配任务并执行

``````
./easynode_collect -config ./config.json
``````

