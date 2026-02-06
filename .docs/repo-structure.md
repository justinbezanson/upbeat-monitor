# Repository Organization: Monorepo Strategy

## Recommendation: Single Monorepo
Instead of splitting the project into 4+ separate repositories, use a single repository (Monorepo).

### Why?
1.  **Go Code Sharing:** The **Backend API** and **Cron Runner** are both written in Go and share the exact same Database Models (`internal/core/domain`) and Database Logic (`internal/repository`). Splitting them would require managing a private Go module, which is unnecessary overhead.
2.  **Frontend Deployment:** Netlify natively supports "Base Directories," allowing you to deploy the **Marketing Site** and **SPA** from different folders in the same repo.
3.  **Context:** Your documentation, infrastructure code, and application code live together, making it easier to reference architecture decisions while coding.

## Directory Structure

```text
upbeat/
├── .github/
│   └── workflows/        # CI/CD Actions (Test Go, Deploy Docs)
├── cmd/                  # Go Application Entrypoints
│   ├── api/              # Lambda Handler for API Gateway
│   │   └── main.go
│   └── cron/             # Lambda Handler for EventBridge
│       └── main.go
├── internal/             # Shared Go Code (The Core)
│   ├── core/             # Domain Models (User, Monitor)
│   ├── repository/       # Turso/SQL implementations
│   └── platform/         # Shared utils (Logger, AWS SDK wrappers)
├── web/                  # Vue 3 SPA (The App)
│   ├── src/
│   ├── package.json
│   └── vite.config.ts
├── marketing/            # Astro Website (Public Site)
│   ├── src/
│   ├── package.json
│   └── astro.config.mjs
├── docs/                 # VitePress (Research & Plans)
│   ├── .vitepress/
│   └── architecture.md   # (Moved from root)
├── infra/                # Infrastructure as Code
│   ├── aws/              # Terraform/Pulumi for Lambda/EventBridge
│   └── home-server/      # Configs for Fizzy (Docker Compose, Systemd)
├── go.mod                # Single Go Module for the root
├── go.sum
└── README.md
```

## Deployment Configuration

### 1. Netlify (Frontend)
You will create **two** sites in the Netlify Dashboard, both linked to this same repository, but configured with different "Base directories".

*   **Site A (Marketing):** `monitor-saas.com`
    *   **Base directory:** `marketing`
    *   **Build command:** `npm run build`
    *   **Publish directory:** `dist`

*   **Site B (App):** `app.monitor-saas.com`
    *   **Base directory:** `web`
    *   **Build command:** `npm run build`
    *   **Publish directory:** `dist`

### 2. AWS Lambda (Backend)
Since Go compiles to a binary, your build script (GitHub Actions or local makefile) will target the specific `cmd` folders.

*   **API Build:** `GOOS=linux GOARCH=arm64 go build -o bootstrap ./cmd/api`
*   **Cron Build:** `GOOS=linux GOARCH=arm64 go build -o bootstrap ./cmd/cron`

### 3. Documentation
*   **Base directory:** `docs`
*   **Build command:** `npx vitepress build`

## Migration Steps

Since you currently have markdown files in the root, the first step of implementation would be:

1.  Initialize the git repo.
2.  Create the `docs/` folder.
3.  Move all `*.md` files (architecture, plans, etc.) into `docs/`.
4.  Initialize the Go module in the root: `go mod init github.com/yourname/upbeat`.