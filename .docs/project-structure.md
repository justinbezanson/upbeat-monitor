# Go Project Structure

This project follows the [Standard Go Project Layout](https://github.com/golang-standards/project-layout).

## Core Directories

### `/cmd`
Main applications for this project.
*   Each subdirectory here represents a binary to be built.
*   Example: `cmd/api/main.go` is the entry point for the Backend API.
*   There should be very little code here (mostly wiring and initialization).

### `/internal`
Private application and library code. This is the code you don't want others importing in their applications.
*   **`/internal/handlers`**: HTTP handlers (controllers).
*   **`/internal/core`**: Domain logic and business rules (Hexagonal/Clean Architecture).
*   **`/internal/repository`**: Database implementations.

### `/pkg` (Optional)
Library code that is safe for external applications to use. (Currently not used).

## Other Directories

### `/api`
OpenAPI/Swagger specifications, JSON schema files, protocol definition files.

### `/web`
Frontend assets, Single Page Applications (Vue.js), and static files.

### `/configs`
Configuration file templates or default configs.

### `/test`
Additional external test applications and test data.

### `/deployments` or `/infra`
IaaS, PaaS, systemd, and container orchestration configurations (Docker, Terraform).