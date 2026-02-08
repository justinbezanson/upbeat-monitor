-- Create "users" table
CREATE TABLE `users` (
  `id` text NOT NULL,
  `email` text NOT NULL,
  `password_hash` text NOT NULL,
  `tier` text NOT NULL DEFAULT 'free',
  `created_at` datetime NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  PRIMARY KEY (`id`)
);
-- Create index "idx_users_email" to table: "users"
CREATE UNIQUE INDEX `idx_users_email` ON `users` (`email`);
-- Create "monitors" table
CREATE TABLE `monitors` (
  `id` text NOT NULL,
  `user_id` text NOT NULL,
  `url` text NOT NULL,
  `interval_seconds` integer NOT NULL DEFAULT 60,
  `next_check_at` datetime NULL,
  `last_checked_at` datetime NULL,
  `status` text NOT NULL DEFAULT 'pending',
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_monitors_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
);
-- Create index "idx_monitors_next_check" to table: "monitors"
CREATE INDEX `idx_monitors_next_check` ON `monitors` (`next_check_at`);
-- Create "checks" table
CREATE TABLE `checks` (
  `id` text NOT NULL,
  `monitor_id` text NOT NULL,
  `status` text NOT NULL,
  `latency_ms` integer NOT NULL,
  `created_at` datetime NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_checks_monitor` FOREIGN KEY (`monitor_id`) REFERENCES `monitors` (`id`) ON DELETE CASCADE
);
-- Create index "idx_checks_monitor_created" to table: "checks"
CREATE INDEX `idx_checks_monitor_created` ON `checks` (`monitor_id`, `created_at`);
