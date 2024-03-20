
Use
tap_indexer;


ALTER TABLE inscriptions ADD inscription_id varchar(256) NOT NULL DEFAULT ""  COMMENT "inscription id";
ALTER TABLE inscriptions ADD inscription_number BIGINT NOT NULL DEFAULT 0  COMMENT "inscription number";
ALTER TABLE `inscriptions` ADD INDEX `idx_inscription_number` (`inscription_number`);

ALTER TABLE txs ADD content varchar(2048) NOT NULL DEFAULT ""  COMMENT "inscription content";

ALTER TABLE inscriptions_stats ADD inscription_id varchar(256) NOT NULL DEFAULT ""  COMMENT "inscription id";
ALTER TABLE inscriptions_stats ADD inscription_number BIGINT NOT NULL DEFAULT 0  COMMENT "inscription number";
ALTER TABLE `inscriptions_stats` ADD INDEX `idx_inscription_number` (`inscription_number`);


ALTER TABLE utxos ADD inscription_id varchar(256) NOT NULL DEFAULT ""  COMMENT "inscription id";
ALTER TABLE utxos ADD inscription_number BIGINT NOT NULL DEFAULT 0  COMMENT "inscription number";
ALTER TABLE utxos DROP COLUMN sn;
ALTER TABLE `utxos` ADD INDEX `idx_inscription_number` (`inscription_number`);
