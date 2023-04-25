-- 创建数据库：不同公链需要创建不同的数据库
CREATE DATABASE ether;

-- 创建交易表
CREATE TABLE IF NOT EXISTS ether.tx
(
    id UInt64,--时间戳
    hash String,
    tx_time String,
    tx_status String,
    block_number String,
    from_addr String,
    to_addr String,
    value String,
    fee String,
    gas_price String,
    max_fee_per_gas String,
    gas String,
    gas_used String,
    base_fee_per_gas String,
    max_priority_fee_per_gas String,
    input_data String,
    block_hash String,
    tx_type String,
    transaction_index String
) ENGINE = ReplicatedReplacingMergeTree('/clickhouse/tables/{shard}/tx', '{replica}',id)
    ORDER BY hash;

CREATE TABLE ether.queue_tx (
    txData String
) ENGINE = Kafka('kafka:9092', 'ether_tx', 'g_ether_tx_4', 'JSONAsString') settings kafka_thread_per_consumer = 0, kafka_num_consumers = 1;


CREATE MATERIALIZED VIEW ether.consumer_tx TO ether.tx
    AS SELECT JSONExtractUInt(txData,'id') as id,
              JSONExtractString(txData,'hash') as hash,
              JSONExtractString(txData,'txTime') as tx_time,
              JSONExtractString(txData,'txStatus') as tx_status,
              JSONExtractString(txData,'blockNumber') as block_number,
              JSONExtractString(txData,'from') as from_addr,
              JSONExtractString(txData,'to') as to_addr,
              JSONExtractString(txData,'value') as value,
     JSONExtractString(txData,'fee') as fee,
     JSONExtractString(txData,'gasPrice') as gas_price,
     JSONExtractString(txData,'maxFeePerGas') as max_fee_per_gas,
     JSONExtractString(txData,'gas') as gas,
     JSONExtractString(txData,'gasUsed') as gas_used,
     JSONExtractString(txData,'baseFeePerGas') as base_fee_per_gas,
     JSONExtractString(txData,'maxPriorityFeePerGas') as max_priority_fee_per_gas,
     JSONExtractString(txData,'input') as input_data,
     JSONExtractString(txData,'blockHash') as block_hash,
     JSONExtractString(txData,'transactionIndex') as transaction_index,
     JSONExtractString(txData,'type') as tx_type
       FROM ether.queue_tx;



/**创建区块表*/
CREATE TABLE IF NOT EXISTS ether.block
(id UInt64,--时间戳
 hash String,
 block_time String,
 block_status String,
 block_number String,
 parent_hash String,
 block_reward String,
 fee_recipient String,
 total_difficulty String,
 block_size String,
 gas_limit String,
 gas_used String,
 base_fee_per_gas String,
 extra_data String,
 state_root String,
 transactions_root String,
 receipts_root String,
 miner String,
 transactions String,
 nonce String
) ENGINE = ReplicatedReplacingMergeTree('/clickhouse/tables/{shard}/block', '{replica}',id)
    ORDER BY hash;

CREATE TABLE ether.queue_block (
    blockData String
) ENGINE = Kafka('kafka:9092', 'ether_block', 'g_ether_block_7', 'JSONAsString') settings kafka_thread_per_consumer = 0, kafka_num_consumers = 1;


CREATE MATERIALIZED VIEW ether.consumer_block TO ether.block
    AS SELECT JSONExtractUInt(blockData,'id') as id,
              JSONExtractString(blockData,'hash') as hash,
              JSONExtractString(blockData,'timestamp') as block_time,
              JSONExtractString(blockData,'blockStatus') as block_status,
              JSONExtractString(blockData,'number') as block_number,
              JSONExtractString(blockData,'parentHash') as parent_hash,
              JSONExtractString(blockData,'blockReward') as block_reward,
              JSONExtractString(blockData,'feeRecipient') as fee_recipient,
              JSONExtractString(blockData,'totalDifficulty') as total_difficulty,
              JSONExtractString(blockData,'size') as block_size,
              JSONExtractString(blockData,'gasLimit') as gas_limit,
              JSONExtractString(blockData,'gasUsed') as gas_used,
              JSONExtractString(blockData,'baseFeePerGas') as base_fee_per_gas,
              JSONExtractString(blockData,'extraData') as extra_data,
              JSONExtractString(blockData,'stateRoot') as state_root,
              JSONExtractString(blockData,'transactions') as transactions,
              JSONExtractString(blockData,'transactionsRoot') as transactions_root,
              JSONExtractString(blockData,'receiptsRoot') as receipts_root,
              JSONExtractString(blockData,'miner') as miner,
              JSONExtractString(blockData,'nonce') as nonce
       FROM ether.queue_block;




/**创建区块表*/
CREATE TABLE IF NOT EXISTS ether.receipt
(id UInt64,--时间戳
 block_hash String,
 logs_bloom String,
 contract_address String,
 transaction_index String,
 tx_type String,
 transaction_hash String,
 gas_used String,
 block_number String,
 cumulative_gas_used String,
 from_addr String,
 to_addr String,
 effective_gas_price String,
 logs String,
 create_time String,
 status String
) ENGINE = ReplicatedReplacingMergeTree('/clickhouse/tables/{shard}/receipt', '{replica}',id)
    ORDER BY transaction_hash;

CREATE TABLE ether.queue_receipt (
    receiptData String
) ENGINE = Kafka('kafka:9092', 'ether_receipt', 'g_ether_receipt_4', 'JSONAsString') settings kafka_thread_per_consumer = 0, kafka_num_consumers = 1;


CREATE MATERIALIZED VIEW ether.consumer_receipt TO ether.receipt
    AS SELECT JSONExtractUInt(receiptData,'id') as id,
              JSONExtractString(receiptData,'blockHash') as block_hash,
              JSONExtractString(receiptData,'logsBloom') as logs_bloom,
              JSONExtractString(receiptData,'contractAddress') as contract_address,
              JSONExtractString(receiptData,'transactionIndex') as transaction_index,
              JSONExtractString(receiptData,'type') as tx_type,
              JSONExtractString(receiptData,'transactionHash') as transaction_hash,
              JSONExtractString(receiptData,'gasUsed') as gas_used,
              JSONExtractString(receiptData,'blockNumber') as block_number,
              JSONExtractString(receiptData,'cumulativeGasUsed') as cumulative_gas_used,
              JSONExtractString(receiptData,'from') as from_addr,
              JSONExtractString(receiptData,'to') as to_addr,
              JSONExtractString(receiptData,'effectiveGasPrice') as effective_gas_price,
              JSONExtractString(receiptData,'logs') as logs,
              JSONExtractString(receiptData,'createTime') as create_time,
              JSONExtractString(receiptData,'status') as status
       FROM ether.queue_receipt;



