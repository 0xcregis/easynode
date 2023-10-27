### Overview

This service is the API service of task, responsible for submitting tasks (including transactions, blocks, receipts and other tasks), querying nodes and other operations.

### Prerequisites

- go version>=1.20

### Building the source

(Take linux system as an example)

- mkdir easynode & cd easynode
- git clone https://github.com/0xcregis/easynode.git
- cd easynode/cmd/taskapi
- CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o easynode_taskapi app.go
  (Compiling a Linux program under Mac is an example. Please search for other cross-compilation commands by yourself)
- ./easynode_taskapi -config ./config.json

### config.json explanation

``````
{
   "RootPath": "/api/task", //Root directory
   "Port": 9001, //Port
   "BlockChain": [ //Support public chain
     200,
     205
   ],
   "TaskKafka": { //Task Kafka
     "Host": "192.168.2.20",
     "Port": 9092
   }
}

``````

### usages

``````
//Submit block task
curl --location --request POST 'http://127.0.0.1:9001/api/task/block' \
--header 'Content-Type: application/json' \
--data-raw '{
     "blockChain": 200,
     "blockHash": "",
     "blockNumber": "16389175"
}'

//Submit a single transaction task
curl --location --request POST 'http://127.0.0.1:9001/api/task/tx' \
--header 'Content-Type: application/json' \
--data-raw '{
     "blockChain": 200,
     "txHash": "0xc0e81699d2728694cc275521daa9b89414a9e4749f7418c8c69f26b090c99f44"
}'

//Submit a single receipt task
curl --location --request POST 'http://127.0.0.1:9001/api/task/receipt' \
--header 'Content-Type: application/json' \
--data-raw '{
     "blockChain": 200,
     "txHash": "0xc0e81699d2728694cc275521daa9b89414a9e4749f7418c8c69f26b090c99f44"
}'

//Submit all transaction tasks under the block
curl --location --request POST 'http://127.0.0.1:9001/api/task/txs' \
--header 'Content-Type: application/json' \
--data-raw '{
     "blockChain": 200,
     "blockHash": "",
     "blockNumber": "16389175"
}'

//Submit all receipt tasks under the block
curl --location --request POST 'http://127.0.0.1:9001/api/task/receipts' \
--header 'Content-Type: application/json' \
--data-raw '{
     "blockChain": 200,
     "blockHash": "",
     "blockNumber": "16389175"
}'

``````