# easynode

### overview
该系统使访问各种公链更简易、更稳定，使用用户专注与自己的业务。包含一些子服务

 - blockchain: 直接访问公链，选择最优节点。对外提供 http、ws 类型协议与其交互
 - collect: 任务(包括：交易任务、区块任务、收据任务)的具体执行者，对公链返回的数据验证、解析后，
 发送到指定的Kafka上
 - task: 时时根据公链最新高度，产生区块任务
 - taskapi: 接收用户自定义的任务
 - store: 接收用户提交的监控地址、接收用户的订阅、主动推送符合条件的交易 和 数据落盘

## Load Source

- mkdir easynode & cd easynode
- git clone https://github.com/0xcregis/easynode.git 或 git clone -b feature/xxx https://github.com/0xcregis/easynode.git
- cd easynode

### install & deploy

 ####  Dependent Environment
   
   执行下面命令，在执行之前确保 docker,docker-compose 都已经安装

   ``````
   //启动命令
   docker-compose -f docker-compose-single-ch.yml up -d
   
   //查看命令
   docker-compose -f docker-compose-single-ch.yml ps
     
   //删除命令
   docker-compose -f docker-compose-single-ch.yml down  

   ``````
   
 运行上述命令后，会启动一些服务和端口
  
   zk: 

       2181:2181
   kafka:

       9092:9092
   kafka_manager:

       9003:9000
   clickhouse:
     
       user:test
       pwd:test
       http.port:8123
       tcp.port:9000
       mysql.port:9004
       postgre.port:9005

   redis:

      6379:6379

 notes:

  1. [docker 安装和使用](https://docs.docker.com/get-docker/)
  2. [docker-compose 安装和使用](https://docs.docker.com/compose/)
  3. docker-compose-single-ch.yml 启动的仅仅是单节点服务，如需要 多kafka 节点、多clickhouse节点不适用该文件
   
 #### Init Database

   使用命令行或工具，关联到上一步已经安装的 clickhouse服务，使用工具执行ch.sql脚本
   (./build/ch_ether.sql) 

   notes：
 
  1. clickhouse的工具 [Dbeaver](https://dbeaver.io/download/)

 #### Deploy & Run Application

   - 创建easynode docker image
  
   ``````
   //创建image
   docker build -f Dockerfile -t easynode:1.0 . 
   
   //查看image 是否生成
   easynode % docker images |grep easynode
   
   ``````
   - 启动 easynode服务
   
   ``````
   docker run --name easynode -p 9001:9001 -p 9002:9002 -p 9003:9003 --network easynode_net -v ./config:/app/config -d easynode:1.0
   ``````

   notes:

   1. network easynode_net : 需要和 docker-compose-single-ch.yml 中保持一致
    
   2. v ./config:/app/config : 把配置文件挂载到镜像中，默认配置文件的名称不可变

   3. ./config的目录结构如下，每个配置文件的具体配置 详见 (/easynode/cmd/easynode/README.md) 

  ``````
       ./config
          ./blockchaiin_config.json
          ./collect_config.json
          ./store_config.json
          ./task_config.json
          ./taskapi_config.json

   ``````
    
   4. 如果配置文件名称需要改变，请 docker run --entrypoint /bin/sh 命令

 - usage
 
请参考 /easynode/cmd/easynode/README.md


 