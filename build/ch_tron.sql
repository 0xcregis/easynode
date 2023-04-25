-- 创建数据库：不同公链需要创建不同的数据库
CREATE DATABASE tron;

CREATE TABLE IF NOT EXISTS  tron.address
(

    `token` String,

    `address` String,

    `tx_type` String,

    `block_chain` Int64,

    `id` Int64
) ENGINE = ReplacingMergeTree ORDER BY id  SETTINGS index_granularity = 8192;

-- 创建交易表
CREATE TABLE IF NOT EXISTS tron.tx
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
) ENGINE = ReplacingMergeTree(id) ORDER BY hash;


/**创建区块表*/
CREATE TABLE IF NOT EXISTS tron.block
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
) ENGINE = ReplacingMergeTree(id) ORDER BY hash;


/**创建区块表*/
CREATE TABLE IF NOT EXISTS tron.receipt
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
) ENGINE = ReplacingMergeTree(id) ORDER BY transaction_hash;