
集成 collect,blockchain,task,task_api,store等服务到一个application, 该app拥有easynode全部功能模块
 
## Prerequisites
- go version >=1.20

## Building the source

(以linux系统为例)
- mkdir easynode & cd easynode
- git clone https://github.com/0xcregis/easynode.git
- cd easynode/cmd/easynode
- CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o easynode app.go
  (mac下编译linux程序为例，其他交叉编译的命令请自行搜索)
- ./easynode

## config.json 详解

- blockchain_config.json
  
  请参考easynode/cmd/blockchain/config.json


- collect_config.json

  请参考easynode/cmd/collect/config.json


- store_config.json
 
  请参考easynode/cmd/store/config.json


- task_config.json

  请参考easynode/cmd/task/config.json


- taskapi_config.json

  请参考easynode/cmd/taskapi/config.json

## run command

 ./easynode [OPTIONS]

``````
-collect_config  string  config path of collect server

-task_config  string  config path of task server

-blockchain_config  string  config path of blockchain server

-taskapi_config  string  config path of taskapi server

-taskapi_config  string  config path of taskapi server

-store_config  string  config path of store server

``````

## usages

- blockchain 服务的使用
  
  请参考 easynode/cmd/blockchain/README.md


- collect 服务的使用

  请参考 easynode/cmd/collect/README.md


- task 服务的使用

  请参考 easynode/cmd/task/README.md


- taskapi 服务的使用

  请参考 easynode/cmd/taskapi/README.md


- store 服务的使用

  请参考 easynode/cmd/store/README.md