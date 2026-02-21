# Gemini CLI Project Context

This project is a **monorepo**. It contains multiple services/applications within a single repository, specifically:
- A Go-based backend API (`cmd/api/main.go`, `internal/`)
- A Vue.js-based frontend SPA (`frontend/`)

The project is **Dockerized**. Development and deployment environments leverage Docker containers.

**Go Environment:**
Go is installed and used **within the Docker containers** for the backend API. When performing operations that require the Go toolchain (e.g., `go build`, `go test`), these commands should ideally be executed inside the relevant Docker container or by building and running the Docker services.

**Database Schema:**
The canonical source of truth for the database schema is the `internal/repository/schema.hcl` file. All database structure changes should be made there.

**Key Directories:**
- `cmd/api/`: Go backend application entry point.
- `internal/`: Internal Go packages (handlers, repository, etc.).
- `frontend/`: Vue.js single-page application.
- `Dockerfile`, `Dockerfile.backend`, `Dockerfile.frontend`: Docker build definitions.
- `docker-compose.yml`: Defines multi-service Docker application.
