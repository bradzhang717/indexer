
CREATE DATABASE  `tap_indexer`  DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci ;
USE `tap_indexer`;


DROP TABLE IF EXISTS `address_txs`;
CREATE TABLE `address_txs` (
   `id` bigint unsigned NOT NULL AUTO_INCREMENT,
   `chain` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'chain name',
   `event` tinyint(1) NOT NULL,
   `protocol` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_bin NOT NULL COMMENT 'protocol name',
   `operate` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_bin NOT NULL COMMENT 'operate',
   `tx_hash` varbinary(128) DEFAULT NULL,
   `address` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'from address',
   `amount` decimal(38,18) NOT NULL COMMENT 'amount',
   `tick` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'inscription name',
   `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
   `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   `related_address` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'related address',
   PRIMARY KEY (`id`),
   KEY `idx_tx_hash` (`tx_hash`(12)),
   KEY `idx_address` (`address`(12)),
   KEY `idx_chain_protocol_tick` (`chain`,`protocol`,`operate`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

DROP TABLE IF EXISTS `balance_txn`;
CREATE TABLE `balance_txn` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `chain` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `protocol` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_bin NOT NULL,
  `event` tinyint(1) NOT NULL,
  `address` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `tick` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_bin NOT NULL,
  `amount` decimal(38,18) NOT NULL,
  `available` decimal(38,18) NOT NULL COMMENT 'available',
  `balance` decimal(38,18) NOT NULL,
  `tx_hash` varbinary(128) DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_address` (`address`(12)),
  KEY `idx_tx_hash` (`tx_hash`(12)),
  KEY `idx_chain_protocol_tick` (`chain`,`protocol`,`tick`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


DROP TABLE IF EXISTS `balances`;
CREATE TABLE `balances` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `sid` int unsigned NOT NULL COMMENT 'sid',
  `chain` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'chain name',
  `protocol` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_bin NOT NULL COMMENT 'protocol name',
  `address` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'address',
  `tick` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_bin NOT NULL COMMENT 'inscription code',
  `available` decimal(38,18) NOT NULL COMMENT 'available',
  `balance` decimal(38,18) NOT NULL COMMENT 'balance',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `address` (`address`,`chain`,`protocol`,`tick`),
  UNIQUE KEY `uqx_chain_sid` (`chain`,`sid`),
  KEY `idx_chain_protocol_tick` (`chain`,`protocol`,`tick`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


DROP TABLE IF EXISTS `block`;
CREATE TABLE `block` (
  `chain` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `block_hash` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `block_number` bigint NOT NULL,
  `block_time` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `chain_id` bigint NOT NULL DEFAULT '0' COMMENT 'chain id',
  PRIMARY KEY (`chain`) USING BTREE,
  UNIQUE KEY `uqx_chain` (`chain`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;



DROP TABLE IF EXISTS `chain_info`;
CREATE TABLE `chain_info` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `chain_id` int unsigned NOT NULL COMMENT 'chain id',
  `chain` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'inner chain name',
  `outer_chain` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'outer chain name',
  `name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'name',
  `logo` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'logo url',
  `network_id` int unsigned NOT NULL COMMENT 'network id',
  `ext` varchar(4098) COLLATE utf8mb4_general_ci NOT NULL COMMENT 'ext',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uqx_chain_id_chain_name` (`chain_id`,`chain`,`name`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


INSERT INTO `chain_info` VALUES (1,0,'btc','btc','BTC','https://s3.indexs.io/chain/icon/btc.png',0,'','2024-03-21 02:00:02','2024-03-21 02:00:02'),
                                (2,1,'eth','eth','Ethereum','https://s3.indexs.io/chain/icon/eth.png',1,'','2024-03-21 02:00:02','2024-03-21 02:00:02'),
                                (3,43114,'avalanche','avax','Avalanche','https://s3.indexs.io/chain/icon/avalanche.png',43114,'','2024-03-21 02:00:02','2024-03-21 02:00:02'),
                                (4,42161,'arbitrum','ETH','Arbitrum One','https://s3.indexs.io/chain/icon/arbitrum.png',42161,'','2024-03-21 02:00:02','2024-03-21 02:00:02'),
                                (5,56,'bsc','BSC','BNB Smart Chain Mainnet','https://s3.indexs.io/chain/icon/bsc.png',56,'','2024-03-21 02:00:02','2024-03-21 02:00:02'),
                                (6,250,'fantom','FTM','Fantom Opera','https://s3.indexs.io/chain/icon/fantom.png',250,'','2024-03-21 02:00:02','2024-03-21 02:00:02'),
                                (7,137,'polygon','Polygon','Polygon Mainnet','https://s3.indexs.io/chain/icon/polygon.png',137,'','2024-03-21 02:00:02','2024-03-21 02:00:02');


DROP TABLE IF EXISTS `chain_stats_hour`;
CREATE TABLE `chain_stats_hour` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `chain` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'chain name',
  `date_hour` int unsigned NOT NULL COMMENT 'date_hour',
  `address_count` int unsigned NOT NULL COMMENT 'address_count',
  `address_last_id` bigint unsigned NOT NULL COMMENT 'address_last_id',
  `inscriptions_count` int unsigned NOT NULL COMMENT 'inscriptions_count',
  `balance_sum` decimal(38,18) NOT NULL COMMENT 'balance_sum',
  `balance_last_id` bigint unsigned NOT NULL COMMENT 'balance_last_id',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uqx_chain_date_hour` (`chain`,`date_hour`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


DROP TABLE IF EXISTS `inscriptions`;
CREATE TABLE `inscriptions` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `sid` int unsigned NOT NULL,
  `chain` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `protocol` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_bin NOT NULL,
  `tick` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_bin NOT NULL,
  `name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `limit_per_mint` decimal(38,18) NOT NULL,
  `deploy_by` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `total_supply` decimal(38,18) NOT NULL,
  `decimals` tinyint unsigned NOT NULL,
  `deploy_hash` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `deploy_time` timestamp NOT NULL,
  `transfer_type` tinyint(1) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `inscription_id` varchar(256) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'inscription id',
  `inscription_number` bigint NOT NULL DEFAULT '0' COMMENT 'inscription number',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_chain_protocol_name` (`chain`,`protocol`,`tick`),
  UNIQUE KEY `uq_chain_sid` (`chain`,`sid`),
  KEY `idx_inscription_number` (`inscription_number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;



DROP TABLE IF EXISTS `inscriptions_stats`;
CREATE TABLE `inscriptions_stats` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `sid` int unsigned NOT NULL,
  `chain` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `protocol` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_bin NOT NULL,
  `tick` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_bin NOT NULL,
  `minted` decimal(38,18) unsigned NOT NULL DEFAULT '0.000000000000000000',
  `mint_completed_time` timestamp NULL DEFAULT NULL,
  `mint_first_block` bigint unsigned NOT NULL,
  `mint_last_block` bigint unsigned NOT NULL,
  `last_sn` int unsigned NOT NULL,
  `holders` int unsigned NOT NULL,
  `tx_cnt` bigint unsigned NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `inscription_id` varchar(256) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'inscription id',
  `inscription_number` bigint NOT NULL DEFAULT '0' COMMENT 'inscription number',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_chain_protocol_name` (`chain`,`protocol`,`tick`),
  UNIQUE KEY `uq_chain_sid` (`chain`,`sid`),
  KEY `idx_inscription_number` (`inscription_number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;



DROP TABLE IF EXISTS `txs`;
CREATE TABLE `txs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `chain` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'chain name',
  `protocol` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `block_height` bigint unsigned NOT NULL COMMENT 'block height',
  `position_in_block` bigint unsigned NOT NULL COMMENT 'Position in Block',
  `block_time` timestamp NOT NULL COMMENT 'block time',
  `tx_hash` varbinary(128) DEFAULT NULL,
  `from` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'from address',
  `to` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'to address',
  `op` varchar(38) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'to address',
  `tick` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'to address',
  `amt` decimal(38,18) NOT NULL COMMENT 'to address',
  `gas` bigint NOT NULL COMMENT 'gas, spend fee',
  `gas_price` bigint NOT NULL COMMENT 'gas price',
  `status` tinyint(1) NOT NULL COMMENT 'tx status',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `content` varchar(2048) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'inscription content',
  PRIMARY KEY (`id`,`block_time`),
  KEY `idx_tx_hash_chain` (`tx_hash`(12),`chain`(4)),
  KEY `idx_chain_protocol_tick` (`chain`,`protocol`,`tick`),
  KEY `idx_chain_block_height` (`chain`,`block_height`)
) ENGINE=InnoDB AUTO_INCREMENT=397752991 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
PARTITION BY RANGE(UNIX_TIMESTAMP(block_time)) (
    PARTITION p202301 VALUES LESS THAN (UNIX_TIMESTAMP('2023-01-01 00:00:00')),  -- 2023.1
    PARTITION p202302 VALUES LESS THAN (UNIX_TIMESTAMP('2023-02-01 00:00:00')),  -- 2023.2
    PARTITION p202303 VALUES LESS THAN (UNIX_TIMESTAMP('2023-03-01 00:00:00')),  -- 2023.3
    PARTITION p202304 VALUES LESS THAN (UNIX_TIMESTAMP('2023-04-01 00:00:00')),  -- 2023.4
    PARTITION p202305 VALUES LESS THAN (UNIX_TIMESTAMP('2023-05-01 00:00:00')),  -- 2023.5
    PARTITION p202306 VALUES LESS THAN (UNIX_TIMESTAMP('2023-06-01 00:00:00')),  -- 2023.6
    PARTITION p202307 VALUES LESS THAN (UNIX_TIMESTAMP('2023-07-01 00:00:00')),  -- 2023.7
    PARTITION p202308 VALUES LESS THAN (UNIX_TIMESTAMP('2023-08-01 00:00:00')),  -- 2023.8
    PARTITION p202309 VALUES LESS THAN (UNIX_TIMESTAMP('2023-09-01 00:00:00')),  -- 2023.9
    PARTITION p202310 VALUES LESS THAN (UNIX_TIMESTAMP('2023-10-01 00:00:00')),  -- 2023.10
    PARTITION p202311 VALUES LESS THAN (UNIX_TIMESTAMP('2023-11-01 00:00:00')),  -- 2023.11
    PARTITION p202312 VALUES LESS THAN (UNIX_TIMESTAMP('2023-12-01 00:00:00')),  -- 2023.12
    PARTITION p202401 VALUES LESS THAN (UNIX_TIMESTAMP('2024-01-01 00:00:00')),  -- 2024.1
    PARTITION p202402 VALUES LESS THAN (UNIX_TIMESTAMP('2024-02-01 00:00:00')),  -- 2024.2
    PARTITION p202403 VALUES LESS THAN (UNIX_TIMESTAMP('2024-03-01 00:00:00')),  -- 2024.3
    PARTITION p202404 VALUES LESS THAN (UNIX_TIMESTAMP('2024-04-01 00:00:00')),  -- 2024.4
    PARTITION p202405 VALUES LESS THAN (UNIX_TIMESTAMP('2024-05-01 00:00:00')),  -- 2024.5
    PARTITION p202406 VALUES LESS THAN (UNIX_TIMESTAMP('2024-06-01 00:00:00')),  -- 2024.6
    PARTITION p202407 VALUES LESS THAN (UNIX_TIMESTAMP('2024-07-01 00:00:00')),  -- 2024.7
    PARTITION p202408 VALUES LESS THAN (UNIX_TIMESTAMP('2024-08-01 00:00:00')),  -- 2024.8
    PARTITION p202409 VALUES LESS THAN (UNIX_TIMESTAMP('2024-09-01 00:00:00')),  -- 2024.9
    PARTITION p202410 VALUES LESS THAN (UNIX_TIMESTAMP('2024-10-01 00:00:00')),  -- 2024.10
    PARTITION p202411 VALUES LESS THAN (UNIX_TIMESTAMP('2024-11-01 00:00:00')),  -- 2024.11
    PARTITION p202412 VALUES LESS THAN (UNIX_TIMESTAMP('2024-12-01 00:00:00')),  -- 2024.12
    PARTITION p202501 VALUES LESS THAN (UNIX_TIMESTAMP('2025-01-01 00:00:00')),  -- 2025.1
    PARTITION p202502 VALUES LESS THAN (UNIX_TIMESTAMP('2025-02-01 00:00:00')),  -- 2025.2
    PARTITION p202503 VALUES LESS THAN (UNIX_TIMESTAMP('2025-03-01 00:00:00')),  -- 2025.3
    PARTITION p202504 VALUES LESS THAN (UNIX_TIMESTAMP('2025-04-01 00:00:00'))   -- 2025.4
)



DROP TABLE IF EXISTS `utxos`;
CREATE TABLE `utxos` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `chain` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `protocol` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_bin NOT NULL,
  `address` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `tick` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_bin NOT NULL,
  `amount` decimal(38,18) NOT NULL,
  `root_hash` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `tx_hash` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `status` tinyint(1) NOT NULL COMMENT 'tx status',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `inscription_id` varchar(256) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'inscription id',
  `inscription_number` bigint NOT NULL DEFAULT '0' COMMENT 'inscription number',
  PRIMARY KEY (`id`),
  KEY `idx_address` (`address`),
  KEY `idx_inscription_number` (`inscription_number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
