# Database Migrations (Atlas + Turso)

We use [Atlas](https://atlasgo.io) to manage the database schema for Turso (libSQL/SQLite).

## 1. Setup

### Install Atlas CLI
```bash
curl -sSf https://atlasgo.sh | sh
```

### Directory Structure
*   `internal/repository/schema.hcl`: The **Source of Truth**. Edit this file to change the DB structure.
*   `migrations/`: Contains the generated SQL files. Do not edit these manually.
*   `atlas.hcl`: Configuration file.

## 2. Workflow

### A. Making Schema Changes
1.  Modify `internal/repository/schema.hcl`.
2.  Run the following command to generate a new migration file:

```bash
# This compares schema.hcl against the migration directory
atlas migrate diff --env local
```

3.  Review the new file created in `migrations/`.

### B. Applying to Local Dev
To apply the pending migrations to your local SQLite file (`local.db`):

```bash
atlas migrate apply --env local
```

### C. Applying to Turso (Production)
To apply migrations to the remote Turso database:

1.  Get your Turso connection URL:
    ```bash
    turso db show upbeat-prod --url
    turso db tokens create upbeat-prod
    ```
2.  Construct the URL: `libsql://[URL]?authToken=[TOKEN]`
3.  Run apply:

```bash
TURSO_DATABASE_URL="libsql://upbeat-prod-org.turso.io?authToken=ey..." \
atlas migrate apply --env turso
```

## 3. Troubleshooting
*   **Checksum Mismatch:** If you manually edited a migration file, Atlas will complain. Use `atlas migrate hash` to re-calculate checksums.