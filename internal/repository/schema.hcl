schema "main" {}

table "accounts" {
  schema = schema.main
  column "id" {
    type = text
    null = false
  }
  column "name" {
    type = text
    null = true
  }
  column "tier" {
    type    = text
    null    = true
    default = "free"
  }
  column "created_at" {
    type = datetime
    null = true
  }
  primary_key {
    columns = [column.id]
  }
}

table "users" {
  schema = schema.main
  column "id" {
    null = false
    type = text
  }
  column "account_id" {
    null = true
    type = text
  }
  column "email" {
    null = false
    type = text
  }
  column "password_hash" {
    null = true
    type = text
  }
  column "role" {
    null    = true
    type    = text
    default = "owner"
  }
  column "created_at" {
    null = true
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_users_account" {
    columns     = [column.account_id]
    ref_columns = [table.accounts.column.id]
  }
  index "idx_users_email" {
    unique  = true
    columns = [column.email]
  }
}

table "oauth_clients" {
  schema = schema.main
  column "client_id" {
    null = false
    type = text
  }
  column "account_id" {
    null = true
    type = text
  }
  column "client_secret_hash" {
    null = true
    type = text
  }
  column "name" {
    null = true
    type = text
  }
  column "created_at" {
    null = true
    type = datetime
  }
  primary_key {
    columns = [column.client_id]
  }
  foreign_key "fk_oauth_clients_account" {
    columns     = [column.account_id]
    ref_columns = [table.accounts.column.id]
  }
}

table "monitors" {
  schema = schema.main
  column "id" {
    null = false
    type = text
  }
  column "account_id" {
    null = true
    type = text
  }
  column "friendly_name" {
    null = true
    type = text
  }
  column "type" {
    null = true
    type = text
  }
  column "url" {
    null = true
    type = text
  }
  column "interval_seconds" {
    null = true
    type = integer
  }
  column "status" {
    null = true
    type = text
  }
  column "last_checked_at" {
    null = true
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_monitors_account" {
    columns     = [column.account_id]
    ref_columns = [table.accounts.column.id]
  }
}

table "checks" {
  schema = schema.main
  column "id" {
    null = false
    type = text
  }
  column "monitor_id" {
    null = true
    type = text
  }
  column "status_code" {
    null = true
    type = integer
  }
  column "latency_ms" {
    null = true
    type = integer
  }
  column "success" {
    null = true
    type = bool
  }
  column "created_at" {
    null = true
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_checks_monitor" {
    columns     = [column.monitor_id]
    ref_columns = [table.monitors.column.id]
  }
}

table "password_reset_tokens" {
  schema = schema.main
  column "id" {
    null = false
    type = text
  }
  column "user_id" {
    null = false
    type = text
  }
  column "token_hash" {
    null = false
    type = text
  }
  column "expires_at" {
    null = false
    type = datetime
  }
  column "created_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_password_reset_tokens_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
  }
  index "idx_password_reset_tokens_token_hash" {
    unique  = true
    columns = [column.token_hash]
  }
}
