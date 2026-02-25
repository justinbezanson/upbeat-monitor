-- Create "password_reset_tokens" table
CREATE TABLE `password_reset_tokens` (
  `id` text NOT NULL,
  `user_id` text NOT NULL,
  `token_hash` text NOT NULL,
  `expires_at` datetime NOT NULL,
  `created_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_password_reset_tokens_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
);
-- Create index "idx_password_reset_tokens_token_hash" to table: "password_reset_tokens"
CREATE UNIQUE INDEX `idx_password_reset_tokens_token_hash` ON `password_reset_tokens` (`token_hash`);
