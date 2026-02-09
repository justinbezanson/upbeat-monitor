# Project Commands

This file contains a list of important commands for the project.

## Running Tests

To run the test suite, use the following command:

```bash
docker compose exec api go test ./...
```

## Database Migrations

This project uses [Atlas](https://atlasgo.io/) for database migrations.

### Creating a New Migration

To create a new migration file, use the `atlas migrate diff` command. This will compare the current state of the database with the schema definition file and generate a new migration file with the necessary SQL statements.

Replace `<migration-name>` with a descriptive name for your migration.

```bash
atlas migrate diff --env local <migration-name>
```

### Applying Migrations

To apply all pending migrations, use the `atlas migrate apply` command.

```bash
atlas migrate apply --env local
```

## Docker

### Starting the Docker Containers

To start the Docker containers in detached mode, use the following command:

```bash
docker compose up -d
```

### Stopping the Docker Containers

To stop the Docker containers, use the following command:

```bash
docker compose down
```

### Rebuilding the Docker Containers

To rebuild the Docker containers, use the following command:

```bash
docker compose up --build -d
```
