# Easynode
[![Go](https://github.com/0xcregis/easynode/actions/workflows/go.yml/badge.svg)](https://github.com/0xcregis/easynode/actions/workflows/go.yml)

> [English](https://github.com/0xcregis/easynode/blob/main/README.md) | 中文

多链一站式服务平台，通过API可快速对接各种链上服务，NFT、Defi、transfer、查询等，
让所有开发者都能使用链上服务，不需要部署节点、不需要了解专业知识,包含以下核心服务：

- blockchain: 直接访问公链，选择最优节点,对外提供 http、ws 类型协议与其交互
- collect: 任务的具体执行者，对获取链上的数据筛选、验证、解析
- task: 根据公链最新高度，产生区块任务
- taskapi: 接收用户自定义的任务
- store: 接收用户提交的监控地址、接收用户的订阅、主动推送符合条件的交易、数据落盘

## 核心能力

- 可获取你关心交易也可获取公链原始交易
- 为每个公链配置多个节点，智能选择最佳节点
- 提供http、ws等协议，适配更多场景
- 组件可独立使用，也可以组合使用
- 功能可配置，根据业务需要，自由配置
- 支持各个公链
- 数据可追溯、可回放
- 具备异常监控，任务重试能力
- 任务执行过程可视化
- 提交自定义任务
- 设定交易过滤条件
- 用户订阅
- 主动推送
- 支持订阅多种业务场景：资产交易，代币转移，质押，激活等
- 历史数据备份

## 公链表

| blockchain  | chaincode |  desc   |
|:------------|:---------:|:-------:|
| ethereum    |    200    |   ETH   |
| tron        |    205    |   TRX   |
| polygon-pos |    201    | POLYGON |
| bitcoin     |    300    |   BTC   |
| filecoin    |    301    |   FIL   |
| ripple      |    310    |   XRP   |
| bnb         |    202    |   BNB   |

## 入门指引

### 1. 必要条件

#### 硬件要求

- CPU With 4+ Cores
- 16GB+ RAM
- 200GB 存储空间,推荐使用SSD且至少512GB空间 

#### 软件要求

- go 版本 需要 v1.20及以上  [install](https://golang.google.cn/doc/install)
- git : 安装最新版本git
   [install](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
- cURL :安装最新版本cURL. [install](https://everything.curl.dev/get)
- docker : 安装最新版本docker
   [install](https://docs.docker.com/desktop/install/linux-install/)
- docker-compose: 安装最新版本docker-compose
   [install](https://docs.docker.com/compose/install/)

### 2. 源码下载

``````
 mkdir easynode && cd easynode
 
 git clone https://github.com/0xcregis/easynode.git
 
 cd easynode
``````

### 3. 初始化配置文件和数据库脚本

- ./build/scripts :数据库初始化脚本,主要修改数据库名称即可（database）默认:ether2
- ./build/config :系统启动需要的配置文件,以下文件和字段常需要修改
    - blockchain_config.json,collect_config.json :主要修改各个公链的节点地址
    - store_config.json ： DbName 字段
    - task_config.json ：BlockMin、BlockMax等字段

1. clickhouse的工具 [Dbeaver](https://dbeaver.io/download/)
2. config 中每个配置文件 [详细说明](https://github.com/0xcregis/easynode/blob/main/cmd/easynode/README.md)

### 4. 运行依赖环境

- 初始化easynode 运行需要的环境

``````
   docker-compose -f docker-compose-single-base.yml up -d
``````

*快捷部署,执行如下命令*

``````
   docker-compose -f docker-compose-single-base-app.yml up -d
``````

*快捷模式部署，则可以跳过步骤5，但其常适用于简单的业务需求*

- 可能用到的其他命令

``````
   #查看命令
   docker-compose -f docker-compose-single-base.yml ps
     
   #删除命令
   docker-compose -f docker-compose-single-base.yml down  
   docker-compose -f docker-compose-single-base.yml down -v 
   
   #重新构建
   docker-compose -f docker-compose-single-base-app.yml build easynode

   #启动命令(启动依赖环境和 应用程序)
   docker-compose -f docker-compose-single-base-app.yml up -d
  
``````

notes:

- clickhouse 的默认：账户：test,密码：test

### 5. 构建和运行easynode

为提高easynode适用范围，采用组件化的设计思想，因此我们为easynode提供如下三种运行和部署方式。

#### 安装包模式

- 下载配置文件
  [下载](https://github.com/0xcregis/easynode/releases)
- 下载安装包
  [下载](https://github.com/0xcregis/easynode/releases)
- 新增hosts项

  ``````
  192.168.2.9 easykafka
  ``````

*其中 192.168.2.9:Kafka服务器IP,换成自己服务器IP，easykafka: Kafka服务名称 需要和 环境脚本保持一致，一般情况不变*

- 运行程序

  ``````
  ./easynode -collect ./config/collect_config.json -task ./config/task_config.json -blockchain ./config/blockchain_config.json -taskapi ./config/taskapi_config.json -store ./config/store_config.json
  ``````

#### docker 模式

- 构建镜像

``````
   #创建image
   docker build -f Dockerfile -t easynode:1.0 . 
   
   #查看image
   docker images |grep easynode
   
``````

- 运行 easynode

``````
    docker run --name easynode -p 9001:9001 -p 9002:9002 -p 9003:9003 --network easynode_easynode_net -v /root/easy_node/easynode/config/:/app/config/ -v /root/app/log/:/app/log/ -v /root/app/data:/app/data/ -d easynode:1.0  
``````

notes:

1. network easynode_net : 需要和 docker-compose-single-base.yml 中保持一致

2. -v 文件挂载 : 容器的路径不可变，宿主路径改成本机可用的绝对路径

3. ./config的目录结构如下，每个配置文件的具体配置 [详见](https://github.com/0xcregis/easynode/blob/main/cmd/easynode/README.md)
且文件名称不可变

  ``````
       ./config
          ./blockchaiin_config.json
          ./collect_config.json
          ./store_config.json
          ./task_config.json
          ./taskapi_config.json

   ``````

4. 如果 *步骤4* 执行 docker-compose-single-base-app.yml 则 可以跳过 *步骤5*

#### docker-compose 集群模式

- networks 设置

  确定步骤4中网络名称，使用如下命令查看，并修改 docker-compose-cluster-easynode.yml 中networks.default.name 字段

  ``````
  docker network ls|grep easynode_net
  ``````

- 配置服务

    根据具体场景需要，增加或删除相关服务 [learn ](https://github.com/0xcregis/easynode/wiki/Overall-Design-For-Easynode)


- 运行集群

    ``````
    docker-compose -f docker-compose-cluster-easynode.yml up -d
    ``````
    每个服务的使用和作用，[详见](https://github.com/0xcregis/easynode/blob/main/cmd/easynode/README.md)

### 6. 检测

- 检查运行环境是否启动

``````
     docker-compose -f docker-compose-single-base.yml ps
``````

- 检查Kafka数据

``````
     # 查看Kafka 
     docker exec -it 25032fc8414e kafka-topics.sh --list --bootstrap-server easykafka:9092
    
     #查看kakfa 数据
     docker exec -it 25032fc8414e kafka-console-consumer.sh --group g1 --topic ether_tx   --bootstrap-server easykafka:9092
``````

- 检查 easynode 是否正常

 ``````
    #查看 app 容器
    
    docker ps |grep easynode
    
    # 查看app 的命令行日志
    
    docker logs -n 10 24a81a2a8e89
 ``````

notes:

- easynode的 9001 端口 [使用说明](https://github.com/0xcregis/easynode/blob/main/cmd/taskapi/README.md)
- easynode的 9002 端口 [使用说明](https://github.com/0xcregis/easynode/blob/main/cmd/blockchain/README.md)
- easynode的 9003 端口 [使用说明](https://github.com/0xcregis/easynode/blob/main/cmd/store/README.md)

### 7.使用

1. 添加监控地址

 ``````
    # 生成token
    curl --location --request POST 'http://localhost:9003/api/store/monitor/token' \
    --header 'Content-Type: application/json' \
    --data-raw '{
      "email": "123@gmail.com"
    }'
    
    #提交监控地址
    curl --location --request POST 'http://localhost:9003/api/store/monitor/address' \
    --header 'Content-Type: application/json' \
    --data-raw '{
        "blockChain": 200,//非必需，如果不传默认0，则表示 跨链监控
        "address": "0x28c6c06298d514db089934071355e5743bf21d61",
        "token": "5fe5f231-7051-4caf-9b52-108db92edbb4"
    }'
    
    //提交订阅规则
    curl --location --request POST 'localhost:9003/api/store/filter/new' \
    --header 'User-Agent: apifox/1.0.0 (https://www.apifox.cn)' \
    --header 'Content-Type: application/json' \
    --data-raw '[
        {
            "token": "afba013c-0072-4592-b8cd-304fa456f76e",
            "blockChain": 205,
            "txCode": "1",
            "params": ""
        }
    
    ]'
    
 ``````

2. 提交订阅并接受返回

``````
   url: ws://localhost:9003/api/store/ws/{token}
            
   接受：
   {
      "code": 1, //消息类型，1:交易消息
      "blockChain": 200, //公链代码
      "data": { //交易数据
          "id":1698395758827420000,
          "chainCode":200,
          "blockHash":"0xbe36cdcfce377f7415bd91be3be10555fc705cd9c48ac077b3de9a1c298c4a36",
          "blockNumber":"18117360",
          "txs":[
              {
                  "contractAddress":"0xd9ec62e6927082ad28b73fb5d4b5e9d571e00768",
                  "from":"0x0000000000000000000000000000000000000000",
                  "method":"Transfer",
                  "to":"0x2c2ab61d2506308c0017f26c36e81e5b22942d57",
                  "value":"1315",
                  "tokenType":721,
                  "index":9
              }
          ],
          "fee":"0.00182692522485181",
          "from":"0x2c2ab61d2506308c0017f26c36e81e5b22942d57",
          "hash":"0x2b7b684d469c365e0f8d9e2bf94bee672878aff4604b7715a48a7f37432f1a21",
          "status":1,
          "to":"0xd9ec62e6927082ad28b73fb5d4b5e9d571e00768",
          "txTime":"1684390019000"
       }
  }        
                       
``````

[更多](https://github.com/0xcregis/easynode/blob/main/cmd/easynode/README.md)
