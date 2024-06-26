Store is the basic and core service of the easynode system. This service is responsible for monitoring address management, user subscription and cancellation, data push, data placement and other functions.

## Prerequisites

- go version >=1.20
- The collect service has been deployed

## Building the source

(Take linux system as an example)

- mkdir easynode & cd easynode
- git clone https://github.com/0xcregis/easynode.git
- cd easynode/cmd/store
- CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o easynode_store app.go
  (Compiling a Linux program under Mac is an example. Please search for other cross-compilation commands by yourself)

- ./easynode_store -config ./config.json

## config.json detailed explanation

``````
{
   "RootPath": "/api/store", //Root path
   "Port": 9003, //Port
   "BaseDb": { //Basic library configuration
     "Addr": "192.168.2.9",
     "Port": 9000,
     "User": "test",
     "Password": "test",
     "DbName": "base",
     "AddressTable": "address", //address table
     "TokenTable": "token", //token table
     "FilterTable": "sub_filter" //Subscription rule table
   },
   "Chains": [ //Public chain configuration
      {
       "BlockChain": 200, //Public chain code
       "BlockStore": false, //Whether the block is placed on the market
       "TxStore": false, //Whether the transaction is placed
       "ReceiptStore": false, //Whether the receipt is placed
       "SubStore": true,//Subscription data persistence
       "BackupTxStore": true,//Broadcast transaction persistence
       "Kafka": { //Kafka configuration where the public chain data is located
           "BackupTx": { //Broadcast transaction
                   "Host": "192.168.2.9",
                   "Port": 9092,
                   "Topic": "backup_tx",
                   "Group": "2",
                   "StartOffset": 0,
                   "Partition": 0
                 },
           "SubTx": { //Subscribe
                 "Host": "192.168.2.9",
                 "Port": 9092,
                 "Topic": "ether_sub_tx",
                 "Group": "1",
                 "StartOffset": 0,
                 "Partition": 0
               },
           "Tx": { //Transaction configuration
               "Host": "192.168.2.20",
               "Port": 9092,
               "Topic": "ether_tx2",
               "Group": "1", //group suffix
               "StartOffset": 0, //read Kafka start unknown -1: latest,-2: first,0:commited offset
               "Partition": 0
             },
           "Block": { //Block configuration
               "Host": "192.168.2.20",
               "Port": 9092,
               "Topic": "ether_block2",
               "Group": "1",
               "StartOffset": 0,
               "Partition": 0
             },
           "Receipt": { //receipt configuration
               "Host": "192.168.2.20",
               "Port": 9092,
               "Topic": "ether_receipt2",
               "Group": "1",
               "StartOffset": 0,
               "Partition": 0
             }
     },
     "ChainDb": { //clickhouse database configuration
       "Addr": "192.168.2.11",
       "Port": 9000,
       "User": "test",
       "Password": "test",
       "DbName": "ether",
       "TxTable": "tx", //Transaction table
       "BlockTable": "block", //block table
       "ReceiptTable": "receipt", //receipt table
       "SubTxTable": "sub_tx", //Subscribe data transaction table
       "BackupTxTable": "backup_tx" //Broadcast transaction
     },
     "Redis": { //redis configuration
       "Addr": "192.168.2.9",
       "Port": 6379,
       "DB": 0
     }
    }
  ]

}

``````

## usages

- http protocol

``````
//Request token, one for each user
curl --location --request POST 'localhost:9003/api/store/monitor/token' \
--header 'Content-Type: application/json' \
--data-raw '{
     "email": "123@gmail.com"
}'

//Submit monitoring address
curl --location --request POST 'localhost:9003/api/store/monitor/address' \
--header 'Content-Type: application/json' \
--data-raw '{
     "blockChain": 200, // optional, if not passed the default 0, it means cross-chain monitoring
     "address": "0x28c6c06298d514db089934071355e5743bf21d61",
     "token": "5fe5f231-7051-4caf-9b52-108db92edbb4"
}'

//Query monitoring address
curl --location --request POST 'localhost:9003/api/store/monitor/address/get' \
--header 'Content-Type: application/json' \
--data-raw '{
     "token": "5fe5f231-7051-4caf-9b52-108db92edbb4"
}'

//Delete monitoring address
curl --location --request POST 'localhost:9003/api/store/monitor/address/delete' \
--header 'Content-Type: application/json' \
--data-raw '{
     "blockChain": 200, // optional, if not passed the default 0, it means cross-chain monitoring
     "address": "0x28c6c06298d514db089934071355e5743bf21d61",
     "token": "5fe5f231-7051-4caf-9b52-108db92edbb4"
}'

//Submit subscription rules
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

//Query subscription rules
curl --location --request POST 'localhost:9003/api/store/filter/get' \
--header 'User-Agent: apifox/1.0.0 (https://www.apifox.cn)' \
--header 'Content-Type: application/json' \
--data-raw '{
     "token": "afba013c-0072-4592-b8cd-304fa456f76e",
     "blockChain": 0,
     "txCode": ""
}'

//Delete subscription rules

curl --location --request POST 'localhost:9003/api/store/filter/delete' \
--header 'User-Agent: apifox/1.0.0 (https://www.apifox.cn)' \
--header 'Content-Type: application/json' \
--data-raw '{
     "id": 1692001339482287000
}'

 OR 
 
curl --location --request POST 'localhost:9003/api/store/filter/delete' \
--header 'User-Agent: apifox/1.0.0 (https://www.apifox.cn)' \
--header 'Content-Type: application/json' \
--data-raw '{
     "token": "afba013c-0072-4592-b8cd-304fa456f76e",
     "blockChain": 205
}'


``````


- ws protocol

``````
   url: 
   
   ws://localhost:9003/api/store/ws/{token}  or ws://localhost:9003/api/store/ws/{token}?serialId={serialId}
            
   receive：
   
   {
      "code": 1, //message type
      "blockChain": 200, //chain code
      "data": { 
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

- code: messageType

  1: Asset transfer transaction, 3: Pledged assets 5: Unlock assets 7: Withdraw 9: Agent resources 11: Recycle resources (cancel agent) 13: Activate account

- txs: transaction list

  All events on the chain, including main currency transactions, contract transactions, etc. 
  
   - main currency transactions
    
      ``````
       {
            "contractAddress":"",
            "from":"0x2c2ab61d2506308c0017f26c36e81e5b22942d57",
            "method":"Transfer",
            "to":"0xd9ec62e6927082ad28b73fb5d4b5e9d571e00768",
            "value":"0.209636442559786101",
            "tokenType":-1,
            "index":0
        }
     
     ``````
     
    - contract transactions

      ``````
      {
            "contractAddress":"0xd9ec62e6927082ad28b73fb5d4b5e9d571e00768",
            "from":"0x0000000000000000000000000000000000000000",
            "method":"Transfer",
            "to":"0x2c2ab61d2506308c0017f26c36e81e5b22942d57",
            "value":"19.56",// value of token
            "tokenType":20,
            "index":9
        }
      ``````
      
      ``````
      {
            "contractAddress":"0xd9ec62e6927082ad28b73fb5d4b5e9d571e00768",
            "from":"0x0000000000000000000000000000000000000000",
            "method":"Transfer",
            "to":"0x2c2ab61d2506308c0017f26c36e81e5b22942d57",
            "value":"1315",// the value is tokenId
            "tokenType":721,
            "index":9
        }
      ``````
      `````` 
      {
        "contractAddress": "0x9a06ef3a841316a9e2c1c93b9c21a7342abe484f",
        "from": "0x1e0049783f008a0085193e00003d00cd54003c71",
        "method": "Transfer",
        "to": "0xf698cb85014702653c6c17d55cfa380618e8f098",
        "value": "7009,1",//the value include 2 element，1st is tokenId, 2nd is token value
        "tokenType": 1155,
        "index": 1
      }
      
      ``````      


- data.txs.contractAddress

  The contract address where the event occurred, If it is the main transaction, it may be empty, otherwise it is the contract address

- data.txs.from

  a account which is from

- data.txs.to 

   a account which is to

- data.txs.method

   event type,it includes Transfer etc.

- data.txs.value

  what was transferred and how much was transferred，it has different values, affected by tokenType
  It is the transaction amount, if tokenType=-1 or 20. it is tokenId ,if tokenType=721. it may be tokenId and amount,if  tokenType=1155

- data.txs.tokenType

  int32,contract type ,include -1,20,721,1155

- data.txs.index

  transaction index in txs

- data.chainCode
   
  chain code ,visit [chainCode](https://github.com/0xcregis/easynode/blob/hotfix/release-0.5.1/README.md)
- data.from,data.to
 
  transaction sender and transaction receiver

- data.blockHash

  The block containing this transaction 

- data.blockNumber
  
  The block containing this transaction

- data.hash

  transaction hash

- data.status

  transaction status ,  1:success, 0:failure

- data.txTime

  Transaction execution time on chain

- notes

    - When the same token is connected multiple times, the previous connection will be automatically closed.
    - The client is required to implement the ping command. If a ping is not received from the client for a long time, the connection will be automatically closed.
    - ws://localhost:9003/api/store/ws/{token}?serialId={serialId}


