# Backend API Implementation Cards

## 1. Initialize Go Project & Architecture
**Description:** Initialize the Go module (`go mod init`) and set up the hexagonal architecture directory structure as defined in the docs (`cmd/`, `internal/core`, `internal/handler`, `internal/repository`, `internal/service`).
**Tags:** backend, setup, architecture

## 2. Turso Database Setup
**Description:** Implement the database connection logic for Turso (libSQL) in `internal/repository/turso`. Ensure configuration for `DATABASE_URL` and auth tokens is handled via environment variables.
**Tags:** backend, database, turso

## 3. Basic Lambda Handler & Router
**Description:** Create the Lambda entry point in `cmd/api/main.go`. Set up a lightweight router (e.g., Chi) and wrap it with `aws-lambda-go-api-proxy` to handle API Gateway events.
**Tags:** backend, aws, lambda

## 4. User Registration & Login Endpoints
**Description:** Implement `POST /auth/register` and `POST /auth/login`. Logic should include validating input and hashing passwords using bcrypt before storing them in the `users` table.
**Tags:** backend, auth, security

## 5. JWT & Cookie Implementation
**Description:** Implement JWT issuance upon successful login. The response should set a `Set-Cookie` header with the JWT, marked as `HttpOnly`, `Secure`, and `SameSite=Strict`.
**Tags:** backend, auth, security

## 6. Internal Auth Middleware
**Description:** Create middleware to protect `/internal/*` routes. It should parse the JWT from the cookie, validate the signature, and inject the user context into the request.
**Tags:** backend, auth, middleware

## 7. Monitor CRUD (Internal API)
**Description:** Implement the Service and Repository logic for `GET`, `POST`, and `DELETE` on `/internal/monitors`. This allows the SPA to manage monitors.
**Tags:** backend, monitors, crud

## 8. Enforce Plan Limits
**Description:** Add logic to the Monitor Creation service to check the user's account tier (Free, Solo, Team) and enforce limits (e.g., Max 50 monitors for Free tier) as defined in `features.md`.
**Tags:** backend, logic, billing

## 9. OAuth Client Management
**Description:** Add endpoints to the Internal API to allow users to generate and revoke `client_id` and `client_secret` pairs for accessing the Public API.
**Tags:** backend, oauth, public-api

## 10. OAuth Token Endpoint
**Description:** Implement `POST /oauth/token` following the Client Credentials Grant flow. Validate `client_id` and `client_secret` and return a short-lived Bearer token.
**Tags:** backend, oauth, security

## 11. Public API Middleware
**Description:** Create middleware for `/api/v1/*` routes that validates the Authorization Bearer token (JWT) issued by the OAuth endpoint.
**Tags:** backend, oauth, middleware

## 12. Public Monitor Endpoint (Read-Only)
**Description:** Expose `/api/v1/monitors` to return a JSON list of monitors and their current status. Ensure this endpoint is read-only and scoped to the authenticated OAuth client's account.
**Tags:** backend, public-api

## 13. EventBridge Integration Logic
**Description:** Ensure that creating or deleting a monitor updates the database state in a way that the Cron Runner (EventBridge + Lambda) picks up immediately. (Or explicitly trigger an EventBridge event if using a push model).
**Tags:** backend, aws, eventbridge