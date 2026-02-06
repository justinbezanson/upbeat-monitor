# Upbeat

**Upbeat** is a modern, serverless uptime monitoring SaaS designed to be cost-effective, scalable, and reliable. It offers HTTP, Ping, and Port monitoring with real-time alerts and status pages.

## 🚀 Tech Stack

*   **Backend:** Go (AWS Lambda)
*   **Database:** Turso (Distributed SQLite)
*   **Frontend:** Vue 3, TypeScript, Tailwind CSS, Pinia (Hosted on Netlify)
*   **Marketing Site:** Astro (Hosted on Netlify)
*   **Infrastructure:** AWS EventBridge (Cron), AWS HTTP API Gateway
*   **Documentation:** VitePress

## 📂 Repository Structure

This project follows a monorepo structure:

*   `cmd/`: Go application entrypoints (API & Cron).
*   `internal/`: Shared Go code (Domain models, Repository implementation).
*   `web/`: Vue.js Single Page Application (The Dashboard).
*   `marketing/`: Astro-based public marketing website.
*   `docs/`: Project documentation and research (VitePress).
*   `infra/`: Infrastructure as Code (Terraform/Pulumi) and self-hosting configs.

## ✨ Features

*   **Multi-Protocol Monitoring:** HTTP(s), Ping, Port, and Keyword checks.
*   **Global Checks:** Distributed check runners via AWS Lambda.
*   **Status Pages:** Public and private status pages for your services.
*   **Alerting:** Email, SMS, and Webhook notifications.
*   **Team Collaboration:** Multi-user accounts with role-based access.

## 📖 Documentation

Detailed documentation regarding architecture, pricing, and implementation plans can be found in the `docs/` directory or by running the documentation site locally:

```bash
cd docs
npm install
npm run docs:dev
```

## 🛠️ Getting Started

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/yourusername/upbeat.git
    cd upbeat
    ```

2.  **Install Dependencies:**
    *   Backend: `go mod download`
    *   Frontend: `cd web && npm install`
    *   Marketing: `cd marketing && npm install`