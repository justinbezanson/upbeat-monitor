schema "main" {}

table "users" {
  schema = schema.main
  column "id" {
    null = false
    type = text
  }
  column "email" {
    null = false
    type = text
  }
  column "password_hash" {
    null = false
    type = text
  }
  column "tier" {
    null = false
    type = text
    default = "free"
  }
  column "created_at" {
    null = false
    type = datetime
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_users_email" {
    unique = true
    columns = [column.email]
  }
}

table "monitors" {
  schema = schema.main
  column "id" {
    null = false
    type = text
  }
  column "user_id" {
    null = false
    type = text
  }
  column "url" {
    null = false
    type = text
  }
  column "interval_seconds" {
    null = false
    type = integer
    default = 60
  }
  column "next_check_at" {
    null = true
    type = datetime
  }
  column "last_checked_at" {
    null = true
    type = datetime
  }
  column "status" {
    null = false
    type = text
    default = "pending"
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_monitors_user" {
    columns = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete = CASCADE
  }
  index "idx_monitors_next_check" {
    columns = [column.next_check_at]
  }
}

table "checks" {
  schema = schema.main
  column "id" {
    null = false
    type = text
  }
  column "monitor_id" {
    null = false
    type = text
  }
  column "status" {
    null = false
    type = text
  }
  column "latency_ms" {
    null = false
    type = integer
  }
  column "created_at" {
    null = false
    type = datetime
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_checks_monitor" {
    columns = [column.monitor_id]
    ref_columns = [table.monitors.column.id]
    on_delete = CASCADE
  }
  index "idx_checks_monitor_created" {
    columns = [column.monitor_id, column.created_at]
  }
}
