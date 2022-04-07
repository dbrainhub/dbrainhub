CREATE DATABASE IF NOT EXISTS `dbrainhub`;

CREATE TABLE IF NOT EXISTS `dbcluster` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(64) NOT NULL,
  `description` TEXT NOT NULL,
  `db_type` VARCHAR(16) NOT NULL,
  `ct` INT NOT NULL,
  `ut` INT NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `unq_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARACTER SET utf8mb4;

CREATE TABLE IF NOT EXISTS `dbcluster_member` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `cluster_id` INT UNSIGNED NOT NULL,
  `hostname` VARCHAR(128) NOT NULL,
  `db_type` VARCHAR(16) NOT NULL,
  `db_version` VARCHAR(32) NOT NULL,
  `role` INT UNSIGNED NOT NULL,
  `ipaddr` VARCHAR(40) NOT NULL,
  `port` INT UNSIGNED NOT NULL,
  `os` VARCHAR(32) NOT NULL,
  `os_version` VARCHAR(32) NOT NULL,
  `host_type` INT UNSIGNED NOT NULL,
  `ct` INT NOT NULL,
  `ut` INT NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `unq_ip_port` (`ipaddr`, `port`),
  INDEX `idx_cluster_id` (`cluster_id`)
) ENGINE=InnoDB DEFAULT CHARACTER SET utf8mb4;

CREATE TABLE IF NOT EXISTS `tag_item` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `item_id` INT UNSIGNED NOT NULL,
  `item_type` VARCHAR(32) NOT NULL,
  `tag` VARCHAR(32) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `unq_item_tag` (`item_id`, `item_type`, `tag`)
) ENGINE=InnoDB DEFAULT CHARACTER SET utf8mb4;
