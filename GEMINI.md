# Gemini CLI Project Context

This project is a **monorepo**. It contains multiple services/applications within a single repository, specifically:
- A Go-based backend API (`cmd/api/main.go`, `internal/`)
- A Vue.js-based frontend SPA (`frontend/`)

The project is **Dockerized**. Development and deployment environments leverage Docker containers.

**Backend API (Go):**
- Go is installed and used **within the Docker containers** for development.
- **Production Hosting:** The backend is designed to be hosted on **AWS Lambda** triggered by **API Gateway**.
- **Rate Limiting:** Do **NOT** implement rate limiting in the Go backend. This is handled at the infrastructure level via AWS API Gateway or AWS WAF to optimize for Lambda cold starts and costs.
- SSL termination and other infrastructure-level concerns are also handled by AWS.

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
