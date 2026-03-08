# Gemini CLI Project Context

This project is a **monorepo**. It contains multiple services/applications within a single repository, specifically:
- A Go-based backend API (`cmd/api/main.go`, `internal/`)
- A Vue.js-based frontend SPA (`frontend/`)

The project is **Dockerized**. Development and deployment environments leverage Docker containers.

**Backend API (Go):**
- Go is installed and used **within the Docker containers** for development.
- **Production Hosting:** The backend is designed to be hosted on **AWS Lambda** triggered by **API Gateway**.
- Infrastructure-level concerns like rate limiting and SSL termination should be handled via AWS WAF or API Gateway configuration in production.

**Frontend Styling:**
Tailwind CSS v4 is used for styling. Do not use `tailwind.config.js/cjs` files as v4 is configured via CSS imports and the `@theme` directive in `src/style.css`.

**Database Schema:**
The canonical source of truth for the database schema is the `internal/repository/schema.hcl` file. All database structure changes should be made there.

**Key Directories:**
- `cmd/api/`: Go backend application entry point.
- `internal/`: Internal Go packages (handlers, repository, etc.).
- `frontend/`: Vue.js single-page application.
- `Dockerfile`, `Dockerfile.backend`, `Dockerfile.frontend`: Docker build definitions.
- `docker-compose.yml`: Defines multi-service Docker application.

**Frontend Testing:**
Browser tests for the frontend can be run using `npm run test:browser` within the `frontend/` directory.
