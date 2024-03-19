
Use
tap_indexer;


ALTER TABLE inscriptions ADD inscription_id varchar(256) NOT NULL DEFAULT ""  COMMENT "inscription id";
ALTER TABLE inscriptions ADD inscription_number BIGINT NOT NULL DEFAULT 0  COMMENT "inscription number";


ALTER TABLE inscriptions_stats ADD inscription_id varchar(256) NOT NULL DEFAULT ""  COMMENT "inscription id";
ALTER TABLE inscriptions_stats ADD inscription_number BIGINT NOT NULL DEFAULT 0  COMMENT "inscription number";

