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
  "Nodes": [
    {
      "BlockChain": 200,
      "NodeUrl": "https://eth-mainnet.g.alchemy.com/v2",
      "NodeToken": "RzxBjjh_c4y0LVHZ7GNm8zoXEZR3HYop"
    },
    {
      "BlockChain": 205,
      "NodeUrl": "https://api.trongrid.io",
      "NodeToken": "244f918d-56b5-4a16-9665-9637598b1223"
    }
  ]
}

``````

### usages

- http protocol

``````
//Send transaction 
curl -X POST \
   http://127.0.0.1:9002/api/chain/origin/tx/sendRawTransaction \
   -H 'Content-Type: application/json' \
   -H 'Postman-Token: 84336133-4ab2-47c4-b600-7b56bfdd79d9' \
   -H 'cache-control: no-cache' \
   -d '{
      "chain":200,
      "signed_tx":"0xe5cB067E90D5Cd1F8052B83562Ae670bA4A211a8" // tx,which was signed
      }'

    response:
     {
      "code": 0,//0: success ,1:failure
      "data": "{origin blockchain data for the method}"
     }
    
//Query balance by address
curl -X POST \
   http://127.0.0.1:9002/api/chain/origin/account/balance \
   -H 'Content-Type: application/json' \
   -H 'Postman-Token: 3eac86c7-c5aa-4dd3-93d3-a3b23678512d' \
   -H 'cache-control: no-cache' \
   -d '{
      "chain":200,
      "address":"0xe5cB067E90D5Cd1F8052B83562Ae670bA4A211a8",
      "tag":"latest"
      }'

    response:
     {
      "code": 0,//0: success ,1:failure
      "data": "{origin blockchain data for the method}"
     }
   
// Query balance of erc20 
curl -X POST \
   http://127.0.0.1:9002/api/chain/origin/account/tokenBalance \
   -H 'Content-Type: application/json' \
   -H 'Postman-Token: 7b49c120-1bfe-4223-9a7a-e1d2620f3099' \
   -H 'cache-control: no-cache' \
   -d '{
      "chain":205, //chaincode
      "address":"TUtAk64jJqdf1pY3SiHeooVikP2SFWXjZ6", // account
      "contract": "TUtAk64jJqdf1pY3SiHeooVikP2SFWXjZ6" //address of contract
      }'
      
        response:
       {
        "code": 0,//0: success ,1:failure
        "data": "{origin blockchain data for the method}"
       }

//The nonce value of the Ether chain
curl -X POST \
   http://127.0.0.1:9002/api/chain/origin/account/nonce \
   -H 'Content-Type: application/json' \
   -H 'Postman-Token: 8fe90095-369c-453c-8d59-6fa18b0a83a0' \
   -H 'cache-control: no-cache' \
   -d '{
     "chain":200,
      "address":"0xe5cB067E90D5Cd1F8052B83562Ae670bA4A211a8",
      "tag":"latest"
      }'
    response:
     {
      "code": 0,//0: success ,1:failure
      "data": "{origin blockchain data for the method}"
     }      

//The latest block height of the blockchain
curl -X POST \
   http://127.0.0.1:9002/api/chain/origin/block/latest \
   -H 'Content-Type: application/json' \
   -H 'Postman-Token: 8d1d5eab-482a-4f0b-b31e-98ab59daf924' \
   -H 'cache-control: no-cache'
   -d '{
      "chain":205
      }'
    response:
     {
      "code": 0,//0: success ,1:failure
      "data": "{origin blockchain data for the method}"
     }

//common function with http-rpc protocol
curl -X POST \
   http://127.0.0.1:9002/api/chain/origin/jsonrpc\
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

    response:
     {
      "code": 0,//0: success ,1:failure
      "data": "{origin blockchain data for the method}"
     }
     
//Query block by blockHash
curl -X POST \
   http://127.0.0.1:9002/api/chain/origin/block/hash \
   -H 'Content-Type: application/json' \
   -H 'Postman-Token: 49c02db3-51a2-4a23-8fbd-f78c50cbd8c7' \
   -H 'cache-control: no-cache' \
   -d '{
      "chain":205,
      "hash":"0000000002f2f66e7256eaffa627b521c380f7dcc4d354bf6c7a5ed8e0c4ea72"
      }'

    response:
     {
      "code": 0,//0: success ,1:failure
      "data": "{origin blockchain data for the method}"
     }
     
//Query blocks by blockHeight
curl -X POST \
   http://127.0.0.1:9002/api/chain/origin/block/number \
   -H 'Content-Type: application/json' \
   -H 'Postman-Token: 3f33d2b1-9a01-4a48-9f9a-c49334ac2903' \
   -H 'cache-control: no-cache' \
   -d '{
     "chain":205,
      "number":"49477110"
      }'
      
    response:
     {
      "code": 0,//0: success ,1:failure
      "data": "{origin blockchain data for the method}"
     }
     
//Query transaction by txHash
curl -X POST \
   http://127.0.0.1:9002/api/chain/origin/tx/hash \
   -H 'Content-Type: application/json' \
   -H 'Postman-Token: 72b4bc45-c235-4361-b88b-2ccffd42a384' \
   -H 'cache-control: no-cache' \
   -d '{
     "chain":205,
      "hash":"89afac2142e025a13987ed183444ec90e9dcb8028bc7bc0757a21c654aa78b31"
      }'

    response:
     {
      "code": 0,//0: success ,1:failure
      "data": "{origin blockchain data for the method}"
     }
     
//Query receipt by txHash
curl -X POST \
   http://127.0.0.1:9002/api/chain/origin/receipts/hash \
   -H 'Content-Type: application/json' \
   -H 'Postman-Token: 9e33035b-fef4-4dfd-a7f3-fd0d2f3de318' \
   -H 'cache-control: no-cache' \
   -d '{
     "chain":205,
      "hash":"89afac2142e025a13987ed183444ec90e9dcb8028bc7bc0757a21c654aa78b31"
      }'

    response:
     {
      "code": 0,//0: success ,1:failure
      "data": "{origin blockchain data for the method}"
     }
     
//query nft.tokenUri 
curl -X POST \
   http://127.0.0.1:9002/api/chain/origin/nft/tokenUri \
   -H 'Content-Type: application/json' \
   -H 'Postman-Token: 40192de5-9473-4179-a142-202ea405e368' \
   -H 'cache-control: no-cache' \
   -d '{
     "chain": 200,
     "contract": "0x0000000000664ceffed39244a8312bd895470803", //contract address
     "tokenId":"439034", //id
     "eip":721 //EIP
    }'

    response:
     {
      "code": 0,//0: success ,1:failure
      "data": "{origin blockchain data for the method}"
     }
     
//query nft.balanceOf 
curl -X POST \
  http://127.0.0.1:9002/api/chain/origin/nft/balanceOf \
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

    response:
     {
      "code": 0,//0: success ,1:failure
      "data": "{origin blockchain data for the method}"
     }
     
//query nft.owner
curl -X POST \
  http://127.0.0.1:9002/api/chain/origin/nft/owner \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: d7f0aae4-5542-4f1a-88b0-c642f7eb92ec' \
  -H 'cache-control: no-cache' \
  -d '{
    "chain": 200,
    "contract": "0x0000000000664ceffed39244a8312bd895470803",
    "tokenId":"439034",
    "eip":721
    }'

    response:
     {
      "code": 0,//0: success ,1:failure
      "data": "{origin blockchain data for the method}"
     }
     
//query nft.totalSupply
curl -X POST \
  http://127.0.0.1:9002/api/chain/origin/nft/totalSupply \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: 81cba267-66cd-4e8e-8278-fa122b1d5032' \
  -H 'cache-control: no-cache' \
  -d '{
    "chain": 200,
    "contract": "0x0000000000664ceffed39244a8312bd895470803",
    "eip":721
    }'

    response:
     {
      "code": 0,//0: success ,1:failure
      "data": "{origin blockchain data for the method}"
     }
     
// query GasPrice
curl -X POST \
  http://127.0.0.1:9002/api/chain/origin/gas/price \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: e65c8170-65f9-482f-92d6-49ed01194444' \
  -H 'cache-control: no-cache' \
  -d '{
    "chain": 200
    }'

    response:
     {
      "code": 0,//0: success ,1:failure
      "data": "{origin blockchain data for the method}"
     }
     
//query EstimateGas
curl -X POST \
  http://127.0.0.1:9002/api/chain/origin/gas/estimateGas \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: 9a8fb36b-0ed3-47c5-b82f-eecd8631c2be' \
  -H 'cache-control: no-cache' \
  -d '{
    "chain": 202,
    "from":"0xd46e8dd67c5d32be8058bb8eb970870f07244567",
    "to":"0xd46e8dd67c5d32be8058bb8eb970870f07244567",
    "data":"0x" // signed tx data
    }'

    response:
     {
      "code": 0,//0: success ,1:failure
      "data": "{origin blockchain data for the method}"
     }

// there are many method of easynode space that response data is returned uniformly, ignoring chain differences.

// query balance of main in easynode space
curl -X POST \
  http://127.0.0.1:9002/api/chain/easynode/account/balance \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: 1762686f-e680-4c06-b18e-d33b92ab390f' \
  -H 'cache-control: no-cache' \
  -d '{
	"chain":205,
	"address":"TUtAk64jJqdf1pY3SiHeooVikP2SFWXjZ6"
}'

response:
{
    "code": 0,//0: success ,1: fail
    "data": {
        "balance": "2" //balance of the address
    }
}

//query block by hash  in easynode space
curl -X POST \
  http://127.0.0.1:9002/api/chain/easynode/block/hash \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: ee71a75f-f3af-4f1e-9ff3-fbc9fa4ed8e2' \
  -H 'cache-control: no-cache' \
  -d '{
	"chain":205,
	"hash":"0000000002f2f66e7256eaffa627b521c380f7dcc4d354bf6c7a5ed8e0c4ea72"
}'
response:
{
    "code": 0,
    "data": {
        "blockHash": "0000000002f2f66e7256eaffa627b521c380f7dcc4d354bf6c7a5ed8e0c4ea72",
        "blockNumber": "49477230",
        "timestamp": "1679043555000"
    }
}

//query tx by hash  in easynode space
curl -X POST \
  http://127.0.0.1:9002/api/chain/easynode/tx/hash \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: f90e9adf-e82c-48e2-9e75-84035fe67680' \
  -H 'cache-control: no-cache' \
  -d '{
	"chain":200,
	"hash":"0x4553316c698c8c89aa979a5dff71eb531d31284b36132e2e02ba8348114286d7"
}'

response:
{
    "code": 0,
    "data": {
        "blockHash": "0xeafd01a251d7d861190b4ef0175a5905bc4941dedab2ee939e58c60271909473",
        "blockNumber": "18626183",
        "from": "0x8186b214a917fb4922eb984fb80cfafa30ee8810",
        "status": 1, //1: success ,0:fail
        "to": "0x35cab8e0d48f40fd9a4ddaf252e63e2b8d4755f5",
        "txHash": "0x4553316c698c8c89aa979a5dff71eb531d31284b36132e2e02ba8348114286d7"
    }
}

//query block by number in easynode space
curl -X POST \
  http://127.0.0.1:9002/api/chain/easynode/block/number \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: 80523f3e-ec91-45e5-b05e-df156dc6a088' \
  -H 'cache-control: no-cache' \
  -d '{
	"chain":200,
	"number":"18718806"
}'

response:
{
    "code": 0,
    "data": {
        "blockHash": "0xae321a54986dfaccf45032d05270d8a6212c0e8642225c57b133a109f3e06b5a",
        "blockNumber": "18718806",
        "timestamp": "1701763115000"
    }
}

//query token balance in easynode space
curl -X POST \
  http://127.0.0.1:9002/api/chain/easynode/account/tokenBalance \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: 3fdef499-681a-4c21-8dc3-3ca0e9eb1a11' \
  -H 'cache-control: no-cache' \
  -d '{
	"chain":200,
	"address":"0xdac17f958d2ee523a2206206994597c13d831ec7",
	"contract":"0xdac17f958d2ee523a2206206994597c13d831ec7",
	"abi":""
}'

response:
{
    "code": 0,
    "data": {
        "balance": "1235900525719" //balance of the address
    }
}

//query nonce in easynode space
curl -X POST \
  http://127.0.0.1:9002/api/chain/easynode/account/nonce \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: d6c68d0b-a29f-4321-b008-b27910b4bfe1' \
  -H 'cache-control: no-cache' \
  -d '{
	"chain":200,
	"address":"0xe5cB067E90D5Cd1F8052B83562Ae670bA4A211a8",
	"tag":"latest"
}'

response:
{
    "code": 0,
    "data": "49" //nonce of the address
}

// query block in easynode space
curl -X POST \
  http://127.0.0.1:9002/api/chain/easynode/block/latest \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: 709ff247-53e9-4c90-bccc-3e1c3a6f241a' \
  -H 'cache-control: no-cache' \
  -d '{
    "chain": 200
}'

response:
{
    "code": 0,
    "data": "18719165" //latest block number
}

//query price to easynode space on current latest block
curl -X POST \
  http://127.0.0.1:9002/api/chain/easynode/gas/price \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: 17e6be68-ef2e-4d18-aba8-aa5ccde36e09' \
  -H 'cache-control: no-cache' \
  -d '{
    "chain": 200
}'

response:
{
    "code": 0,
    "data": "45133954243" //the value is the smallest unit of different chains
}

//estimateGas fee in easynode space 
curl -X POST \
  http://127.0.0.1:9002/api/chain/easynode/gas/estimateGas \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: 8cc2e412-ac93-4d6e-9fa7-6001c3376733' \
  -H 'cache-control: no-cache' \
  -d '{
    "chain": 200,
    "from":"0xd46e8dd67c5d32be8058bb8eb970870f07244567",
    "to":"0xd46e8dd67c5d32be8058bb8eb970870f07244567",
    "data":"0x" // signed tx 
}'

response:
{
    "code": 0,
    "data": "21000" //estimated cost to complete the transaction
}

``````