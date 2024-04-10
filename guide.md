#### easynode 整体架构
 - 整体架构
 easynode 整体上分上下两层；下层主要有 kafka,clickhouse,redis等组成，提供服务间数据流转、数据存储等功能,该层所有组件由dockerfile 一键启动和管理，详细请参考readme文件；
 
    上层根据业务功能，划分了多个独立模块，它们可以独立服务，也可以相互配合完成其它服务


 - 各个模块的作用
  1. task:  负责生产区块任务 并 分发到合理的节点上
 
  2. taskapi： 通过API 添加、查询各种类型的任务，包括：区块任务、交易任务、批量交易任务、收据任务等

 3. blockchin：转发用户请求到最优的fullnode节点，并对返回的数据格式化处理

 4. collect ：根据配置文件的规则，接受、处理各种类型任务；高效、稳定、可靠的同步链上数据；清洗、格式化、转存链上数据到指定位置

 5. store : 接受和管理用户订阅规则；持久化用户指定的数据；推送用户需要的数据；格式化输出数据；计算各种场景矿工费；

#### 交易异常处理机制

1. 确认交易是否发送成功且在链上是否执行成功
  
    通过区块浏览器观察确认

2. 确认与easynode建立websocket连接的token 

    通过使用方确认

3. 与交易相关的监控地址是否提交

    通过查询base库下 sub_filter 、address表 确认用户提交的规则；需要特别注意 用户提交后 需要最长1分钟后才能生效

4. 确认该交易是否正在等待队列中
  
   进入docker/redis
    `````
   docker exec -it 714120a37836 /bin/sh
   redis-cli -r 100
   keys *
   `````
   查看最新区块高度和已经下发的区块高度
   ``````
   //redis 执行命令
   hgetall blockChain
   
   //最新链上区块高度,2510代表chainCode
   "latestNumber_2510" 
     "37730663"
   //最近分发的区块高度,2510代表chainCode
   "recentNumber_2510"
    "37730663"
    ``````
   查看已经执行的区块或交易的任务
   ``````
    //redis 执行命令
   hgetall latestBlock
   
   //已经执行 区块任务最高高度 ,chainCode:198
   "198-LatestBlock-number"
   "43079992"
   
   //已经执行 交易任务最高高度,chainCode:198
   "198-LatestTx-number"
   "43079988"

    ``````
   通过以上 判断该交易 所在的区块是否已经被执行


5. 确认该交易是否在重试队列中

  ``````
  //redis 执行命令
  
  //判断该交易或交易所在区块 是否被重试了： 因合约被重试
   hget errTx_198 xxxxxxxx
  
    //判断该交易或交易所在区块 是否被重试了： 因其他原因被重试
   hget nodeTask_198 xxxxxxxx 

  ``````

6. 日志排查

  通常问题通过以上步骤 就可以得到问题答案； 如果 仍然无法确认问题原因，那么 我们就需要通过日志来排查问题了。
  
 ``````
 //应用启动目录下 通常可以看到如下
 //具体日志路径 可以通过配置文件更改
 [root@amazonlinux easynode]# ls -al ./log/collect/bnb
总用量 2326904
drwxr-xr-x 2 root root        196 4月  10 02:25 .
drwxr-xr-x 7 root root         73 2月  18 18:31 ..
-rw-r--r-- 1 root root  164467894 4月   9 23:59 chain_info_log_202404090000
-rw-r--r-- 1 root root   74372169 4月  10 15:03 chain_info_log_202404100000
-rw-r--r-- 1 root root 1341101746 4月   9 23:59 cmd_log_202404090000
-rw-r--r-- 1 root root  614711707 4月  10 15:03 cmd_log_202404100000
-rw-r--r-- 1 root root   64848955 4月   9 23:59 monitor_log_202404090000
-rw-r--r-- 1 root root   32410541 4月  10 15:03 monitor_log_202404100000
[root@amazonlinux easynode]# 

//日志文件通常按天划分一个新文件

 //查找类似这个文件 ，检查链上数据情况
 chain_info_log_202404100000
 
 //查找类似这个文件 ，检查easynode系统处理逻辑情况
 cmd_log_202404100000
 
 ``````
