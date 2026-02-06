# Cron & Check Runner Implementation

## Overview
The "Cron" system is the heartbeat of the monitoring SaaS. It is responsible for scheduling checks, executing them (HTTP, Ping, Port), recording results to Turso, and triggering alerts.

## 1. Architecture

*   **Scheduler:** AWS EventBridge Scheduler.
*   **Compute:** AWS Lambda (Go runtime).
*   **Database:** Turso (Shared with Backend API).
*   **Execution Model:** "Tick" based polling.

### The "Tick" Strategy
Instead of managing individual EventBridge rules for every user monitor (which hits quotas and is hard to manage), we use a single high-frequency trigger.

1.  **Trigger:** EventBridge invokes the `CheckRunner` Lambda every **30 seconds**.
2.  **Dispatcher:** The Lambda queries the `monitors` table for items where `next_check_at <= NOW()`.
3.  **Execution:** The Lambda spins up Go routines to perform the checks in parallel.
4.  **Completion:** Results are batched and written to the `checks` table; `next_check_at` is updated.

## 2. Technology Choice: Go
Go is explicitly chosen for the Check Runner for the following reasons:
*   **Concurrency:** A single Lambda (1GB RAM) can easily handle hundreds of concurrent HTTP requests using `goroutines` and `channels`.
*   **Cold Starts:** Go compiles to a static binary, offering faster cold starts than Java or .NET, which is crucial for a function running every 30 seconds.
*   **Code Reuse:** We can import the `domain` and `repository` packages from the `backend-api` module to share database models and logic.

## 3. Implementation Details

### Data Model Updates
To support efficient polling, the `monitors` table requires indexing on the scheduling column.
*   **Column:** `next_check_at` (DATETIME) - Calculated as `last_checked_at + interval_seconds`.
*   **Index:** `CREATE INDEX idx_monitors_next_check ON monitors(next_check_at);`

### Lambda Logic (Go)
The handler will follow a fan-out/fan-in pattern:

1.  **Fetch:** `SELECT * FROM monitors WHERE next_check_at <= ? LIMIT 500`
2.  **Fan-Out:** Iterate over monitors and launch a goroutine for each.
3.  **Check:**
    *   **HTTP:** `net/http` client with strict timeouts.
    *   **Ping:** ICMP (requires Lambda permissions) or unprivileged UDP ping.
    *   **Port:** `net.DialTimeout`.
4.  **Fan-In:** Collect `CheckResult` structs via a buffered channel.
5.  **Persist:** Perform a bulk `INSERT` into `checks` and bulk `UPDATE` on `monitors`.
6.  **Alert:** If `current_status != previous_status`, trigger the Alerting Service (separate Go struct or Lambda).

## 4. Security Concerns & Mitigations

Since this service acts as a proxy making requests to user-defined URLs, it is highly susceptible to **SSRF (Server-Side Request Forgery)**.

### A. SSRF Protection (Critical)
Users must not be able to use your infrastructure to scan internal AWS resources or attack third parties.
*   **Block Private IPs:** The HTTP client must validate the resolved IP address *before* connecting.
    *   Deny: `127.0.0.0/8`, `10.0.0.0/8`, `172.16.0.0/12`, `192.168.0.0/16`.
    *   Deny: `169.254.169.254` (AWS Instance Metadata Service).
*   **DNS Rebinding:** Resolve the IP once, validate it, and then dial that specific IP (setting the `Host` header manually) to prevent TOCTOU (Time-of-check to time-of-use) attacks.

### B. Resource Exhaustion
*   **Timeouts:** Hard timeout of 10 seconds per check.
*   **Response Limits:** Do not read the entire body. Use `io.LimitReader` to read only the first 10KB (enough for keyword checks) to prevent memory exhaustion from large files.
*   **Redirects:** Limit max redirects to 3 to prevent infinite loops.

### C. Egress Control
*   **VPC:** If the Lambda runs in a VPC, use Network ACLs to deny outbound traffic to internal subnets.

## 5. Scaling Plan

As the user base grows, a single Lambda may not process all checks within 30 seconds.

*   **Stage 1 (Current):** Single Lambda handles all checks.
*   **Stage 2 (Sharding):**
    *   EventBridge triggers a **Dispatcher Lambda**.
    *   Dispatcher queries DB and pushes batches of Monitor IDs to **AWS SQS**.
    *   **Worker Lambdas** consume SQS messages and run checks.
    *   *Benefit:* Infinite horizontal scaling.