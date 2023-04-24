# easynode

### overview
该系统使访问各种公链更简易、更稳定，使用用户专注与自己的业务。包含一些子服务

 - blockchain: 直接访问公链，选择最优节点。对外提供 http、ws 类型协议与其交互
 - collect: 任务(包括：交易任务、区块任务、收据任务)的具体执行者，对公链返回的数据验证、解析后，
 发送到指定的Kafka上
 - task: 时时根据公链最新高度，产生区块任务
 - taskapi: 接收用户自定义的任务
 - store: 接收用户提交的监控地址、接收用户的订阅、主动推送符合条件的交易 和 数据落盘

### install & deploy

 1. Dependent Environment
   
   run script 

   ``````
   docker-compose -f docker-compose-single-ch.yml up -d
   ``````

 2. Init Database

   run script (./build/ch.sql) in command line which is clickhouse client 

 3. Deploy & run application

   - 



 