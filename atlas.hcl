# Environment for local development (SQLite file)
env "local" {
  # Source of truth: The HCL file we just created
  src = "file://internal/repository/schema.hcl"

  # Target database: A local SQLite file
  url = "sqlite://local.db?_fk=1"

  # Dev database: Used by Atlas to calculate diffs (In-memory)
  dev = "sqlite://file?mode=memory&_fk=1"

  migration {
    dir = "file://migrations"
  }
}

# Environment for Production (Turso)
env "turso" {
  src = "file://internal/repository/schema.hcl"
  url = getenv("TURSO_DATABASE_URL")
  dev = "sqlite://file?mode=memory&_fk=1"
  migration {
    dir = "file://migrations"
  }
}