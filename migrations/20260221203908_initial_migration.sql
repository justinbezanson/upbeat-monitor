-- Create "accounts" table
CREATE TABLE `accounts` (
  `id` text NOT NULL,
  `name` text NULL,
  `tier` text NULL DEFAULT 'free',
  `created_at` datetime NULL,
  PRIMARY KEY (`id`)
);
-- Create "users" table
CREATE TABLE `users` (
  `id` text NOT NULL,
  `account_id` text NULL,
  `email` text NOT NULL,
  `password_hash` text NULL,
  `role` text NULL DEFAULT 'owner',
  `created_at` datetime NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_users_account` FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`)
);
-- Create index "idx_users_email" to table: "users"
CREATE UNIQUE INDEX `idx_users_email` ON `users` (`email`);
-- Create "oauth_clients" table
CREATE TABLE `oauth_clients` (
  `client_id` text NOT NULL,
  `account_id` text NULL,
  `client_secret_hash` text NULL,
  `name` text NULL,
  `created_at` datetime NULL,
  PRIMARY KEY (`client_id`),
  CONSTRAINT `fk_oauth_clients_account` FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`)
);
-- Create "monitors" table
CREATE TABLE `monitors` (
  `id` text NOT NULL,
  `account_id` text NULL,
  `friendly_name` text NULL,
  `type` text NULL,
  `url` text NULL,
  `interval_seconds` integer NULL,
  `status` text NULL,
  `last_checked_at` datetime NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_monitors_account` FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`)
);
-- Create "checks" table
CREATE TABLE `checks` (
  `id` text NOT NULL,
  `monitor_id` text NULL,
  `status_code` integer NULL,
  `latency_ms` integer NULL,
  `success` bool NULL,
  `created_at` datetime NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_checks_monitor` FOREIGN KEY (`monitor_id`) REFERENCES `monitors` (`id`)
);
