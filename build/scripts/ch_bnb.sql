-- 创建数据库：不同公链需要创建不同的数据库
CREATE
DATABASE IF NOT EXISTS bnb;

-- 创建交易表
CREATE TABLE IF NOT EXISTS bnb.tx
(
    id
    UInt64,--时间戳
    hash
    String,
    tx_time
    String,
    tx_status
    String,
    block_number
    String,
    from_addr
    String,
    to_addr
    String,
    value
    String,
    fee
    String,
    gas_price
    String,
    max_fee_per_gas
    String,
    gas
    String,
    gas_used
    String,
    base_fee_per_gas
    String,
    max_priority_fee_per_gas
    String,
    input_data
    String,
    block_hash
    String,
    tx_type
    String,
    transaction_index
    String
) ENGINE = ReplacingMergeTree
(
    id
) ORDER BY hash;


/**创建区块表*/
CREATE TABLE IF NOT EXISTS bnb.block
(
    id
    UInt64,--时间戳
    hash
    String,
    block_time
    String,
    block_status
    String,
    block_number
    String,
    parent_hash
    String,
    block_reward
    String,
    fee_recipient
    String,
    total_difficulty
    String,
    block_size
    String,
    gas_limit
    String,
    gas_used
    String,
    base_fee_per_gas
    String,
    extra_data
    String,
    state_root
    String,
    transactions_root
    String,
    receipts_root
    String,
    miner
    String,
    transactions
    String,
    nonce
    String
) ENGINE = ReplacingMergeTree
(
    id
) ORDER BY hash;


/**创建区块表*/
CREATE TABLE IF NOT EXISTS bnb.receipt
(
    id
    UInt64,--时间戳
    block_hash
    String,
    logs_bloom
    String,
    contract_address
    String,
    transaction_index
    String,
    tx_type
    String,
    transaction_hash
    String,
    gas_used
    String,
    block_number
    String,
    cumulative_gas_used
    String,
    from_addr
    String,
    to_addr
    String,
    effective_gas_price
    String,
    logs
    String,
    create_time
    String,
    status
    String
) ENGINE = ReplacingMergeTree
(
    id
) ORDER BY transaction_hash;

CREATE TABLE IF NOT EXISTS bnb.sub_tx
(
    id
    UInt64,--时间戳
    block_chain
    UInt64,
    hash
    String,
    tx_time
    String,
    tx_status
    UInt64,
    block_number
    String,
    block_hash
    String,
    from_addr
    String,
    to_addr
    String,
    value
    String,
    fee
    String,
    fee_detail
    String,
    input_data
    String,
    tx_type
    UInt64,
    contract_tx
    String
) ENGINE = ReplacingMergeTree
(
    id
) ORDER BY hash;

CREATE TABLE IF NOT EXISTS bnb.backup_tx
(
    id
    UInt64,--时间戳
    block_chain
    UInt64,
    extra
    String,
    signed
    String,
    tx_status
    UInt8,
    from_addr
    String,
    to_addr
    String,
    response
    String
) ENGINE = ReplacingMergeTree
(
    id
) ORDER BY id;