# Backend API Implementation Plan

## Overview
This document outlines the implementation strategy for the backend API of the server monitoring SaaS. The API is built using Go, hosted on AWS Lambda, and backed by Turso (libSQL). It serves two primary clients: the Vue.js SPA (Internal) and third-party integrations (Public API).

## 1. Architecture & Tech Stack

*   **Language:** Go (Golang)
*   **Compute:** AWS Lambda (Go runtime)
*   **API Gateway:** AWS HTTP API (v2) for routing and basic throttling.
*   **Database:** Turso (Distributed SQLite).
*   **Authentication:** Custom JWT for SPA, OAuth 2.0 (Client Credentials) for Public API.

## 2. API Interface Design

The API will be split into two distinct route groups to separate concerns and security models.

### A. Internal API (`/internal/v1`)
*   **Consumer:** Vue.js Single Page Application (hosted on Netlify).
*   **Auth:** Session-based via HTTP-Only, Secure Cookies (JWT).
*   **Features:** Full management capabilities (User settings, Billing, Monitor CRUD).

### B. Public API (`/api/v1`)
*   **Consumer:** External scripts, status page aggregators, custom dashboards.
*   **Auth:** OAuth 2.0 Bearer Tokens.
*   **Features:** Read-heavy access to Monitor status and metrics. Limited write capabilities.

## 3. Authentication & Security

### SPA Security (User Sessions)
*   **Login Flow:**
    *   `POST /auth/login` accepts email/password.
    *   Validates bcrypt hash against `users` table.
    *   Returns a `Set-Cookie` header containing a signed JWT (Access Token).
*   **CSRF Protection:**
    *   Cookie is `SameSite=Strict` and `HttpOnly`.
    *   Double-Submit Cookie pattern or custom header requirement (e.g., `X-Requested-With: XMLHttpRequest`) to prevent simple CSRF.

### Public API Security (OAuth 2.0)
*   **Flow:** Client Credentials Grant (RFC 6749).
*   **Mechanism:**
    1.  User generates `client_id` and `client_secret` in the SPA dashboard.
    2.  External client requests token: `POST /oauth/token`.
    3.  System validates credentials and returns a short-lived JWT Access Token.
    4.  Client sends `Authorization: Bearer <token>` for API requests.
*   **Scopes:** `monitors:read`, `monitors:write`.

## 4. Database Schema (Turso)

The database will handle relational data for users and monitors, and time-series data for checks.

```sql
-- Accounts (Billing/Subscription Entity)
CREATE TABLE accounts (
    id TEXT PRIMARY KEY,
    name TEXT, -- e.g., "Acme Corp" or "Personal"
    tier TEXT DEFAULT 'free', -- free, solo, team, enterprise
    created_at DATETIME
);

-- Users & Auth
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    account_id TEXT REFERENCES accounts(id),
    email TEXT UNIQUE,
    password_hash TEXT,
    role TEXT DEFAULT 'owner', -- owner, admin, viewer
    created_at DATETIME
);

CREATE TABLE oauth_clients (
    client_id TEXT PRIMARY KEY,
    account_id TEXT REFERENCES accounts(id),
    client_secret_hash TEXT,
    name TEXT,
    created_at DATETIME
);

-- Monitoring Core
CREATE TABLE monitors (
    id TEXT PRIMARY KEY,
    account_id TEXT REFERENCES accounts(id),
    friendly_name TEXT,
    type TEXT, -- http, ping, port, keyword
    url TEXT,
    interval_seconds INTEGER, -- 300 (Free), 60 (Solo/Team), 30 (Enterprise)
    status TEXT, -- up, down, paused
    last_checked_at DATETIME
);

-- Time Series Data (Potential for separate table per month or optimization)
CREATE TABLE checks (
    id TEXT PRIMARY KEY,
    monitor_id TEXT,
    status_code INTEGER,
    latency_ms INTEGER,
    success BOOLEAN,
    created_at DATETIME
);
```

## 5. Go Application Structure

The application will follow a hexagonal architecture to keep the core logic independent of AWS Lambda.

```text
/
├── cmd/
│   └── api/
│       └── main.go           # Lambda entry point, router setup
├── internal/
│   ├── core/
│   │   ├── domain/           # Structs: User, Monitor, Check
│   │   └── ports/            # Interfaces: MonitorRepository, AuthService
│   ├── handler/
│   │   ├── http/             # HTTP Handlers (Gin/Chi or stdlib)
│   │   └── middleware/       # Auth, Logging, CORS
│   ├── service/              # Business Logic (CRUD, Validation)
│   └── repository/
│       └── turso/            # SQL implementations
└── go.mod
```

## 6. Implementation Plan

### Phase 1: Foundation
1.  Initialize Go module.
2.  Set up Turso database connection logic.
3.  Create basic Lambda handler with a router (e.g., `chi` or `aws-lambda-go-api-proxy`).

### Phase 2: User Management & Internal Auth
1.  Implement `POST /auth/register` and `POST /auth/login`.
2.  Implement JWT issuance and Cookie setting.
3.  Create Auth Middleware to protect `/internal/*` routes.

### Phase 3: Monitor Management (Core Feature)
1.  Implement CRUD for Monitors (`GET`, `POST`, `DELETE` /internal/monitors).
2.  Enforce Tier limits (e.g., Max 50 monitors for Free tier) based on `features.md`.

### Phase 4: Public API & OAuth
1.  Implement `oauth_clients` management in Internal API.
2.  Implement `POST /oauth/token` endpoint.
3.  Create OAuth Middleware to validate Bearer tokens.
4.  Expose `/api/v1/monitors` with read-only access.

### Phase 5: Integration
1.  Ensure API can trigger or configure the EventBridge rules (via AWS SDK) for the check runner.