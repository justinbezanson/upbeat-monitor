# Frontend SPA Implementation Plan

## Overview
The Single Page Application (SPA) is the control center for users. It interacts with the `internal/v1` API to manage monitors, view status reports, and handle billing.

## 1. Tech Stack & Dependencies

### Core
*   **Framework:** Vue 3 (Composition API)
*   **Build Tool:** Vite
*   **Language:** TypeScript
*   **State Management:** Pinia
*   **Routing:** Vue Router 4

### UI & Styling
*   **Styling:** Tailwind CSS
*   **Components:** shadcn-vue (Headless UI + Radix Vue based)
*   **Icons:** lucide-vue-next
*   **Charts:** vue-chartjs (Chart.js) or Unovis (for latency/uptime graphs)

### Recommended Additions
*   **Data Fetching:** **TanStack Query (Vue Query)**
    *   *Why?* Essential for handling server state, caching, and **auto-polling** monitor statuses every 30 seconds without manual `setInterval` spaghetti code.
*   **HTTP Client:** **Axios**
    *   *Why?* Robust interceptors for handling 401 (Unauthorized) redirects and CSRF.
*   **Validation:** **Zod** + **Vee-Validate (v4)**
    *   *Why?* Type-safe schema validation for login/register and monitor creation forms. Vee-Validate v4 is rewritten specifically for Vue 3 and the Composition API.
*   **Utilities:** **@vueuse/core**
    *   *Why?* Standard library for Vue composables (Dark mode, LocalStorage, Window resizing).
*   **Date Formatting:** **date-fns**
    *   *Why?* Lightweight date manipulation for "Last checked 5 mins ago".

## 2. Project Structure

```text
src/
├── assets/
├── components/
│   ├── ui/           # shadcn-vue components (Button, Input, Card)
│   ├── common/       # App-specific shared components (Navbar, Footer)
│   └── monitors/     # Monitor-specific components (MonitorList, UptimeGraph)
├── composables/      # Shared logic (useAuth, useTheme)
├── layouts/          # AppLayout (Sidebar), AuthLayout (Centered)
├── pages/            # Route views
│   ├── auth/         # Login.vue, Register.vue
│   ├── dashboard/    # Overview.vue
│   ├── monitors/     # Index.vue, Create.vue, Details.vue
│   └── settings/     # Account.vue, APIKeys.vue
├── router/           # Vue Router config
├── services/         # API wrappers (Axios instances)
├── stores/           # Pinia stores (auth, preferences)
└── utils/            # Helpers (date formatting, error handling)
```

## 3. Key Features & Implementation

### A. Authentication (Cookie-based)
*   **Login:** POST `/auth/login`.
*   **Session:** The browser handles the `HttpOnly` cookie.
*   **Axios Config:** `withCredentials: true` is mandatory.
*   **Auth Guard:** Vue Router middleware to check authentication status before accessing protected routes.

### B. Dashboard & Polling
*   **Goal:** Show live status of 50+ monitors.
*   **Implementation:**
    *   Use `useQuery` from TanStack Query.
    *   Key: `['monitors']`.
    *   `refetchInterval: 30000` (30 seconds).
    *   **Optimistic Updates:** When a user pauses a monitor, update UI immediately while sending API request.

### C. Monitor Details
*   **Charts:** Line chart for Latency (ms) over last 24h. Bar chart for Uptime %.
*   **Logs:** Virtualized list (if logs are extensive) showing recent checks.

### D. Settings & API Keys
*   **API Keys:** Interface to generate/revoke Client Credentials for the Public API.
*   **Billing:** Display current Tier and usage (e.g., "12/50 Monitors").

## 4. Routing Map

| Path | Component | Protected | Description |
| :--- | :--- | :--- | :--- |
| `/login` | Login.vue | No | User sign in |
| `/register` | Register.vue | No | New account |
| `/` | Dashboard.vue | Yes | Overview of all monitors |
| `/monitors/new` | CreateMonitor.vue | Yes | Form to add HTTP/Ping check |
| `/monitors/:id` | MonitorDetails.vue | Yes | Charts, Logs, Edit settings |
| `/status-pages` | StatusPages.vue | Yes | Manage public status pages |
| `/settings` | Settings.vue | Yes | Profile, Password, API Keys |

## 5. State Management (Pinia)

We will keep client-state minimal, relying on Vue Query for server-state.

*   **AuthStore:** `isAuthenticated`, `user` (email, tier), `login()`, `logout()`.
*   **UIStore:** Sidebar state (collapsed/expanded), Dark mode preference.
