# Self-Hosting Plan: Fizzy (Bug Tracker)

## Overview
This document outlines the deployment strategy for [Fizzy](https://github.com/basecamp/fizzy) on a Linux home server.

*   **App:** Fizzy (Assumed Ruby on Rails stack, typical of Basecamp).
*   **Infrastructure:** Docker & Docker Compose.
*   **Reverse Proxy:** Nginx (Existing on host).
*   **Database:** SQLite or PostgreSQL (Containerized).

## 1. Prerequisites
Ensure your Linux home server has the following installed:
*   **Git:** To clone the repository.
*   **Docker & Docker Compose:** To run the application isolated.
*   **Nginx:** To serve the app on your local network.

## 2. Installation & Deployment

We will use Docker Compose to orchestrate the application and its dependencies (Redis/DB).

### A. Clone & Setup
```bash
cd /opt
sudo git clone https://github.com/basecamp/fizzy.git
cd fizzy
```

### B. Docker Compose Configuration
Create a `docker-compose.yml` file in the project root if one does not exist.

```yaml
services:
  web:
    build: .
    restart: always
    ports:
      - "3000:3000"
    env_file: .env
    environment:
      - RAILS_ENV=production
      - RAILS_LOG_TO_STDOUT=true
      - RAILS_SERVE_STATIC_FILES=true
      - PIDFILE=/rails/tmp/pids/server.pid
    volumes:
      - fizzy_storage:/rails/storage
      - fizzy_db:/rails/db
    depends_on:
      - redis
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:3000/up"]
      interval: 5s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    restart: always
    volumes:
      - fizzy_redis:/data

volumes:
  fizzy_storage:
  fizzy_db:
  fizzy_redis:
```

### C. Environment Configuration
Create a `.env` file based on the repo's example (usually `.env.example`).

```bash
# Generate a secret key (Required for production)
SECRET_KEY_BASE=$(openssl rand -hex 64)

# Database (If using SQLite, this is often default. If Postgres, add connection string)
DATABASE_URL=sqlite3:///rails/db/production.sqlite3
```

### D. First Run
```bash
# Build and start
docker compose up -d --build

# Run database migrations (Rails specific)
docker compose exec web bin/rails db:prepare

# (Optional) Seed initial data if mentioned in development.md
# docker compose exec web bin/rails db:seed
```

## 3. Nginx Configuration (Reverse Proxy)

Configure your existing Nginx instance to proxy traffic to the Docker container.

Create `/etc/nginx/sites-available/fizzy.local`:

```nginx
server {
    listen 80;
    server_name fizzy.local; # Add this to your local DNS or /etc/hosts

    location / {
        proxy_pass http://127.0.0.1:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # WebSocket support (Critical for Hotwire/Turbo used by Basecamp apps)
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
```

Enable the site:
```bash
sudo ln -s /etc/nginx/sites-available/fizzy.local /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

## 4. Maintenance

*   **Updates:** `git pull && docker compose up -d --build`
*   **Backups:** Backup the `/opt/fizzy/storage` directory and the SQLite database file (or Postgres dump).

## 5. Auto-start on Boot (Systemd)

To ensure Fizzy starts automatically when the server reboots, create a systemd service.

Create `/etc/systemd/system/fizzy.service`:

```ini
[Unit]
Description=Fizzy Docker Compose Application
Requires=docker.service
After=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/opt/fizzy
ExecStart=/usr/bin/docker compose up -d
ExecStop=/usr/bin/docker compose down
TimeoutStartSec=0

[Install]
WantedBy=multi-user.target
```

Enable and start the service:
```bash
sudo systemctl enable fizzy.service
sudo systemctl start fizzy.service
```
