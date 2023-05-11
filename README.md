
## Overview
该系统使访问各种公链更简易、更稳定，使用用户专注于自己的业务。包含以下子服务

- blockchain: 直接访问公链，选择最优节点。对外提供 http、ws 类型协议与其交互
- collect: 任务(包括：交易任务、区块任务、收据任务)的具体执行者，对公链返回的数据验证、解析后，
  发送到指定的Kafka上
- task: 时时根据公链最新高度，产生区块任务
- taskapi: 接收用户自定义的任务
- store: 接收用户提交的监控地址、接收用户的订阅、主动推送符合条件的交易 和 数据落盘

## Getting Started

### 1. Prerequisites

#### Hardware Requirements

- CPU With 4+ Cores
- 16GB+ RAM
- 200GB of free storage,Recommended that high performance SSD with at least 512GB free space

#### Software Requirements

- go version>=1.20  [install](https://golang.google.cn/doc/install)
- git : Install the latest version of git if it is not already installed. [install](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
- cURL :Install the latest version of cURL if it is not already installed. [install](https://everything.curl.dev/get)
- docker : Install the latest version of Docker if it is not already installed. [install](https://docs.docker.com/desktop/install/linux-install/)
- docker-compose: Install the latest version of Docker-compose if it is not already installed. [install](https://docs.docker.com/compose/install/)

### 2. Download The Source

``````
 mkdir easynode && cd easynode
 
 git clone https://github.com/0xcregis/easynode.git
 
 cd easynode
``````

### 3. Init Config & Database

- ./scripts :数据库初始化脚本,主要修改数据库名称即可（database）默认:ether2
- ./config :系统启动需要的配置文件,以下文件和字段常需要修改
    - blockchain_config.json,collect_config.json :主要修改各个公链的节点地址
    - store_config.json ： DbName 字段
    - task_config.json ：BlockMin、BlockMax等字段

1. clickhouse的工具 [Dbeaver](https://dbeaver.io/download/)
2. config 中每个配置文件 [详细说明](https://github.com/0xcregis/easynode/blob/main/cmd/easynode/README.md)


### 4. Start Dependent Environment

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



### 5. Build & Run Application
 
为提高easynode适用范围，采用组件化的设计思想，因此我们为easynode提供如下三种运行和部署方式。

#### 安装包模式
  - 下载配置文件
    [下载](https://github.com/0xcregis/easynode/releases)
  - 下载安装包 
    [下载](https://github.com/0xcregis/easynode/releases)
  - 运行程序
 
  ``````
  ./easynode -collect ./config/collect_config.json -task ./config/task_config.json -blockchain ./config/blockchain_config.json -taskapi ./config/taskapi_config.json -store ./config/store_config.json
  ``````
#### docker 模式
- build image

``````
   #创建image
   docker build -f Dockerfile -t easynode:1.0 . 
   
   #查看image
   docker images |grep easynode
   
``````
- run easynode

``````
    docker run --name easynode -p 9001:9001 -p 9002:9002 -p 9003:9003 --network easynode_easynode_net -v /root/easy_node/easynode/config/:/app/config/ -v /root/app/log/:/app/log/ -v /root/app/data:/app/data/ -d easynode:1.0  
``````

notes:

1. network easynode_net : 需要和 docker-compose-single-base.yml 中保持一致

2.  -v 文件挂载 : 容器的路径不可变，宿主路径改成本机可用的绝对路径

3. ./config的目录结构如下，每个配置文件的具体配置 [详见](https://github.com/0xcregis/easynode/blob/main/cmd/easynode/README.md)且文件名称不可变

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

### 6. Check

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


## Usage

1. 添加监控地址

 ``````
    # 生成token
    curl --location --request POST 'http://localhost:9003/api/store/monitor/token' \
    --header 'Content-Type: application/json' \
    --data-raw '{
        "blockChain": 200
    }'
    
    #提交监控地址
    curl --location --request POST 'http://localhost:9003/api/store/monitor/address' \
    --header 'Content-Type: application/json' \
    --data-raw '{
        "blockChain": 200,
        "address": "0x28c6c06298d514db089934071355e5743bf21d61",
        "token": "5fe5f231-7051-4caf-9b52-108db92edbb4",
        "txType": 12
    }'
    
 ``````

2. 提交订阅并接受返回

``````
   url: ws://localhost:9003/api/store/ws/{token}
   
   入参：
           {
             "id":1001,
             "code":1,
             "blockChain":[200],
             "Params":{}
            }
   
   返回：
   
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
            
                        
   
``````
[更多](https://github.com/0xcregis/easynode/blob/main/cmd/easynode/README.md)
