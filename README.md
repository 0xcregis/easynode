# Easynode
[![Go](https://github.com/0xcregis/easynode/actions/workflows/go.yml/badge.svg)](https://github.com/0xcregis/easynode/actions/workflows/go.yml)

> English | [中文](https://github.com/0xcregis/easynode/blob/main/README_zh.md)

Multi-chain one-stop service platform, through API, can quickly connect to various on-chain services, such as NFT, Defi, transfer, query, etc.
Allow all developers to use on-chain services without deploying nodes or knowing professional knowledge, including the following core services:

- blockchain: directly access the public chain, select the optimal node, and provide http and ws type protocols to interact with it
- collect: the specific executor of the task, to filter, verify and analyze the data on the acquisition chain
- task: Generate block tasks according to the latest height of the public chain
- tasksapi: receive user-defined tasks
- store: Receive monitoring addresses submitted by users, receive user subscriptions, actively push eligible transactions, and place data on the market


## Key Features of Easynode

- You can get the transactions you care about, and you can also get the original transactions of the public chain
- Configure multiple nodes for each public chain, intelligently select the best node
- Provide http, ws and other protocols to adapt to more scenarios
- Components can be used independently or in combination
- Functions are configurable and can be configured freely according to business needs
- Support various public chains
- Data traceable and replayable
- Possess exception monitoring and task retry capabilities
- Visualization of task execution process
- Submit custom tasks
- Set transaction filter conditions
- user subscription
- Active push
- Support subscription for various business scenarios: asset transaction, token transfer, pledge, activation, etc.
- Historical data backup


## Blockchain for table

  | blockchain  | chaincode |  desc   |
  |:------------|:---------:|:-------:|
  | ethereum    |    200    |   ETH   |
  | tron        |    205    |   TRX   |
  | polygon-pos |    201    | POLYGON |
  | bitcoin     |    300    |   BTC   |
  | filecoin    |    301    |   FIL   |
  | ripple      |    310    |   XRP   |
  | bnb         |    202    |   BNB   |

## Getting Started

### 1. Prerequisites

#### Hardware Requirements

- CPU With 4+ Cores
- 16GB+ RAM
- 200GB of free storage,Recommended that high performance SSD with at least 512GB free space

#### Software Requirements

- go version>=1.20  [install](https://golang.google.cn/doc/install)
- git : Install the latest version of git if it is not already
  installed. [install](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
- cURL :Install the latest version of cURL if it is not already installed. [install](https://everything.curl.dev/get)
- docker : Install the latest version of Docker if it is not already
  installed. [install](https://docs.docker.com/desktop/install/linux-install/)
- docker-compose: Install the latest version of Docker-compose if it is not already
  installed. [install](https://docs.docker.com/compose/install/)

### 2. Download The Source

``````
 mkdir easynode && cd easynode
 
 git clone https://github.com/0xcregis/easynode.git
 
 cd easynode
``````

### 3. Init Config & Database

- ./scripts : Database initialization script, mainly modify the database name (database) default: ether2
- ./config : The configuration file required for system startup, the following files and fields often need to be modified
  - blockchain_config.json, collect_config.json: mainly modify the node address of each public chain
  - store_config.json : DbName field
  - task_config.json : BlockMin, BlockMax and other fields

1. Clickhouse tool [Dbeaver](https://dbeaver.io/download/)
2. Each configuration file in config folder [details](https://github.com/0xcregis/easynode/blob/main/cmd/easynode/README.md)

### 4. Start Dependent Environment

- Initialize the environment for easynode operation

``````
   docker-compose -f docker-compose-single-base.yml up -d
``````

*Quick deployment, execute the following command*

``````
   docker-compose -f docker-compose-single-base-app.yml up -d
``````

*Shortcut mode deployment, step 5 can be skipped, but it is often applicable to simple business requirements*

- Other commands that may be used

``````
   #view
   docker-compose -f docker-compose-single-base.yml ps
     
   #del
   docker-compose -f docker-compose-single-base.yml down  
   docker-compose -f docker-compose-single-base.yml down -v 
   
   #rebuild
   docker-compose -f docker-compose-single-base-app.yml build easynode
  
``````

notes:

- The default of clickhouse: account: test, password: test

### 5. Build and Run easynode

In order to improve the scope of application of easynode, we adopt a componentized design idea, so we provide the following three operation and deployment methods for easynode.

#### Binary package mode

- Download configuration files
  [Download](https://github.com/0xcregis/easynode/releases)
- Download the installation package
  [Download](https://github.com/0xcregis/easynode/releases)
- Add hosts item

  ``````
  192.168.2.9 easykafka
  ``````

  *Among them, 192.168.2.9: Kafka server IP, replace it with your own server IP, easykafka: Kafka service name needs to be consistent with the environment script, and the general situation remains unchanged*

- Run easynode

  ``````
  ./easynode -collect ./config/collect_config.json -task ./config/task_config.json -blockchain ./config/blockchain_config.json -taskapi ./config/taskapi_config.json -store ./config/store_config.json
  ``````

#### Docker mode

- Build image

``````
   #create image
   docker build -f Dockerfile -t easynode:1.0 . 
   
   #scan image
   docker images |grep easynode
   
``````

- Run easynode

``````
    docker run --name easynode -p 9001:9001 -p 9002:9002 -p 9003:9003 --network easynode_easynode_net -v /root/easy_node/easynode/config/:/app/config/ -v /root/app/log/:/app/log/ -v /root/app/data:/app/data/ -d easynode:1.0  
``````

notes:

1. network easynode_net : need to be consistent with docker-compose-single-base.yml

2. -v file mount: the path of the container is immutable, and the path of the host is changed to the absolute path available on the machine

3. The directory structure of ./config is as follows, the specific configuration of each configuration file [see](https://github.com/0xcregis/easynode/blob/main/cmd/easynode/README.md) and the file name is immutable

  ``````
       ./config
          ./blockchaiin_config.json
          ./collect_config.json
          ./store_config.json
          ./task_config.json
          ./taskapi_config.json

   ``````

4. If *step 4* executes docker-compose-single-base-app.yml, you can skip *step 5*

#### Docker-compose cluster mode

- Networks setting

  Determine the network name in step 4, use the following command to view, and modify the networks.default.name field in docker-compose-cluster-easynode.yml

  ``````
  docker network ls|grep easynode_net
  ``````

- Manager service

  According to the needs of specific scenarios, add or delete related services [learn ](https://github.com/0xcregis/easynode/wiki/Overall-Design-For-Easynode)

- Run cluster

  ``````
  docker-compose -f docker-compose-cluster-easynode.yml up -d
  ``````

  The use and function of each service, [details](https://github.com/0xcregis/easynode/blob/main/cmd/easynode/README.md)

### 6. Check

- Check whether the operating environment is started

``````
     docker-compose -f docker-compose-single-base.yml ps
``````

- Check Kafka data

``````
     # review Kafka 
     docker exec -it 25032fc8414e kafka-topics.sh --list --bootstrap-server easykafka:9092
    
     #review kakfa data
     docker exec -it 25032fc8414e kafka-console-consumer.sh --group g1 --topic ether_tx   --bootstrap-server easykafka:9092
``````

- Check easynode

 ``````
    #review app container
    
    docker ps |grep easynode
    
    # review app log
    
    docker logs -n 10 24a81a2a8e89
 ``````

notes:

- Port 9001 of easynode [Instructions](https://github.com/0xcregis/easynode/blob/main/cmd/taskapi/README.md)
- Port 9002 of easynode [Instructions](https://github.com/0xcregis/easynode/blob/main/cmd/blockchain/README.md)
- Port 9003 of easynode [Instructions](https://github.com/0xcregis/easynode/blob/main/cmd/store/README.md)

### 7. Usage

1. Add monitoring address

 ``````
    # create token
    curl --location --request POST 'http://localhost:9003/api/store/monitor/token' \
    --header 'Content-Type: application/json' \
    --data-raw '{
      "email": "123@gmail.com"
    }'
    
    #submit monitoring address
    curl --location --request POST 'http://localhost:9003/api/store/monitor/address' \
    --header 'Content-Type: application/json' \
    --data-raw '{
        "blockChain": 200,//Not required, if not passed the default 0, it means cross-chain monitoring
        "address": "0x28c6c06298d514db089934071355e5743bf21d61",
        "token": "5fe5f231-7051-4caf-9b52-108db92edbb4"
    }'
    
    #submit subscription rules
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

2. Submit subscription and accept return

``````
   url: ws://localhost:9003/api/store/ws/{token}
               
   return of push：
   {
     "code": 1, //message type, 1: transaction message
     "blockChain": 200, //public chain code
     "data": { //transaction data
       "id": 1685094437357929000,
       "blockHash": "0x067fbc694c5ca3540ee965b25c286e55d40f3e5e5fd336d1f398868dfc18feec", //block hash
       "blockNumber": "17284552", //block height
       "chainCode": 200,
       "contractTx": [ //EVM event during contract transaction
         {
           "contract": "0xdac17f958d2ee523a2206206994597c13d831ec7", //contract address
           "from": "0x54b50187becd0bbcfd52ec5d538433dab044d2a8", //from address
           "method": "Transfer", //contract method
           "to": "0x408be4b8a862c1a372976521401fd77f9a0178d7", //to address
           "value": "59.327379" //transaction content
         }
       ],
       "fee": "0.002053016771146819",//transaction fee
       "from": "0x54b50187becd0bbcfd52ec5d538433dab044d2a8", //from address
       "hash": "0x323c08a889ed99d8bfc6c72b1580432f7a13ca7c992fd1bac523e46bfe7ab98f", //transaction hash
       "status": "1", //transaction status 1: successful, 0: transaction failed
       "to": "0xdac17f958d2ee523a2206206994597c13d831ec7", //to address
       "txTime": "1684390019", //transaction time
       "txType": 1, //transaction type 1: contract call, 2: normal asset transfer
       "value": "0" //transaction amount
     }
  }        
                             
``````

[more](https://github.com/0xcregis/easynode/blob/main/cmd/easynode/README.md)
