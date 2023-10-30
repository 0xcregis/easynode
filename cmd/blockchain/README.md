### Overview

This service forwards the request initiated by the user to the corresponding node of the blockchain and provides the
optimal node during this period. It is mainly used in scenarios such as submitting transactions, checking balances, and
querying mining fees.

### Prerequisites

- go version>=1.18

### Building the source

(Take linux system as an example)

- mkdir easynode & cd easynode
- git clone https://github.com/0xcregis/easynode.git
- cd easynode/cmd/blockchain
- CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o easynode_chain app.go
  (Compiling a Linux program under Mac is an example. Please search for other cross-compilation commands by yourself)

- ./easynode_chain -config ./config.json

### config.json explanation

``````
{
   "RootPath": "/api/chain", //api root directory
   "Port": 9002, //port
   "Kafka": { //Backup broadcast transaction Kafka
     "Host": "192.168.2.9",//host
     "Port": 9092, //port
     "Topic": "backup_tx", //kafka.topic
     "Partition": 0 //kafka.partition,default:0
   },
   "Cluster": {//Blockchain node cluster
     "200": [{ //ether node configuration
       "NodeUrl": "https://eth-mainnet.g.alchemy.com/v2", //Node address
       "NodeToken": "************************" //Node requires token
     }],
     "205": [{ //tron node configuration
       "NodeUrl": "https://api.trongrid.io",//Node address
       "NodeToken": "**********************"//Node requires token
     }]
   }
}

``````

### usages

- http protocol

``````
//Send transaction 
curl -X POST \
   http://127.0.0.1:9002/api/chain/tx/sendRawTransaction \
   -H 'Content-Type: application/json' \
   -H 'Postman-Token: 84336133-4ab2-47c4-b600-7b56bfdd79d9' \
   -H 'cache-control: no-cache' \
   -d '{
      "chain":200,
      "signed_tx":"0xe5cB067E90D5Cd1F8052B83562Ae670bA4A211a8" // tx,which was signed
      }'

//Query balance by address
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

// Query balance of erc20 
curl -X POST \
   http://127.0.0.1:9002/api/chain/account/tokenBalance \
   -H 'Content-Type: application/json' \
   -H 'Postman-Token: 7b49c120-1bfe-4223-9a7a-e1d2620f3099' \
   -H 'cache-control: no-cache' \
   -d '{
      "chain":205, //chaincode
      "address":"TUtAk64jJqdf1pY3SiHeooVikP2SFWXjZ6", // account
      "contract": "TUtAk64jJqdf1pY3SiHeooVikP2SFWXjZ6" //address of contract
      }'

//The nonce value of the Ether chain
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

//The latest block height of the blockchain
curl -X POST \
   http://127.0.0.1:9002/api/chain/block/latest \
   -H 'Content-Type: application/json' \
   -H 'Postman-Token: 8d1d5eab-482a-4f0b-b31e-98ab59daf924' \
   -H 'cache-control: no-cache'
   -d '{
      "chain":205
      }'


//common function with http-rpc protocol
curl -X POST \
   http://127.0.0.1:9002/api/chain/jsonrpc\
   -H 'Content-Type: application/json' \
   -H 'Postman-Token: 3d6a170b-bd88-42b5-83d1-d3e5bec7babf' \
   -H 'cache-control: no-cache' \
   -d '{
      "chain":205, //chainCode
      "data":{   //ether method
               "id": 1,
               "jsonrpc": "2.0",
               "params": [
                  "0xe5cB067E90D5Cd1F8052B83562Ae670bA4A211a8",
                   "latest"
               ],
               "method": "eth_getBalance"
             }
      }'

//Query block by blockHash
curl -X POST \
   http://127.0.0.1:9002/api/chain/block/hash \
   -H 'Content-Type: application/json' \
   -H 'Postman-Token: 49c02db3-51a2-4a23-8fbd-f78c50cbd8c7' \
   -H 'cache-control: no-cache' \
   -d '{
      "chain":205,
      "hash":"0000000002f2f66e7256eaffa627b521c380f7dcc4d354bf6c7a5ed8e0c4ea72"
      }'

//Query blocks by blockHeight
curl -X POST \
   http://127.0.0.1:9002/api/chain/block/number \
   -H 'Content-Type: application/json' \
   -H 'Postman-Token: 3f33d2b1-9a01-4a48-9f9a-c49334ac2903' \
   -H 'cache-control: no-cache' \
   -d '{
     "chain":205,
      "number":"49477110"
      }'

//Query transaction by txHash
curl -X POST \
   http://127.0.0.1:9002/api/chain/tx/hash \
   -H 'Content-Type: application/json' \
   -H 'Postman-Token: 72b4bc45-c235-4361-b88b-2ccffd42a384' \
   -H 'cache-control: no-cache' \
   -d '{
     "chain":205,
      "hash":"89afac2142e025a13987ed183444ec90e9dcb8028bc7bc0757a21c654aa78b31"
      }'

//Query receipt by txHash
curl -X POST \
   http://127.0.0.1:9002/api/chain/receipts/hash \
   -H 'Content-Type: application/json' \
   -H 'Postman-Token: 9e33035b-fef4-4dfd-a7f3-fd0d2f3de318' \
   -H 'cache-control: no-cache' \
   -d '{
     "chain":205,
      "hash":"89afac2142e025a13987ed183444ec90e9dcb8028bc7bc0757a21c654aa78b31"
      }'

//query nft.tokenUri 
curl -X POST \
   http://127.0.0.1:9002/api/chain/nft/tokenUri \
   -H 'Content-Type: application/json' \
   -H 'Postman-Token: 40192de5-9473-4179-a142-202ea405e368' \
   -H 'cache-control: no-cache' \
   -d '{
     "chain": 200,
     "contract": "0x0000000000664ceffed39244a8312bd895470803", //contract address
     "tokenId":"439034", //id
     "eip":721 //EIP
    }'

//query nft.balanceOf 
curl -X POST \
  http://127.0.0.1:9002/api/chain/nft/balanceOf \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: 7de26f12-1598-458e-b58e-d1a2182cdc92' \
  -H 'cache-control: no-cache' \
  -d '{
    "chain": 200,
    "contract": "0x0000000000664ceffed39244a8312bd895470803",
    "address":"0x99f49B6783f6E1e6D6A9b16e291BbB9D164e54FF",
    "tokenId":"439034",
    "eip":721
    }'

//query nft.owner
curl -X POST \
  http://127.0.0.1:9002/api/chain/nft/owner \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: d7f0aae4-5542-4f1a-88b0-c642f7eb92ec' \
  -H 'cache-control: no-cache' \
  -d '{
    "chain": 200,
    "contract": "0x0000000000664ceffed39244a8312bd895470803",
    "tokenId":"439034",
    "eip":721
    }'

//query nft.totalSupply
curl -X POST \
  http://127.0.0.1:9002/api/chain/nft/totalSupply \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: 81cba267-66cd-4e8e-8278-fa122b1d5032' \
  -H 'cache-control: no-cache' \
  -d '{
    "chain": 200,
    "contract": "0x0000000000664ceffed39244a8312bd895470803",
    "eip":721
    }'

``````