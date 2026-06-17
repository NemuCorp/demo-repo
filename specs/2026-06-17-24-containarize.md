# Containarize

**Issue:** [NemuCorp/demo-repo#24](https://github.com/NemuCorp/demo-repo/issues/24)

## Problem

The demo-repo project currently requires manual setup steps to run: a locally installed PostgreSQL instance, `go run .` for the server, and `npm start` for the client. There is no Docker or container support of any kind. A developer cloning the repo must install Go, Node.js, and PostgreSQL, then run three separate processes. The goal is to collapse this into a single `docker compose up` command that brings up the full stack.

## Goals

- Provide a `docker-compose.yml` at the repo root so that `docker compose up` launches the entire application.
- Containerize the Go/Gin server with a multi-stage Docker build (compile then minimal runtime image).
- Containerize the React client: build static assets with Node, then serve them via nginx with API reverse-proxying to the server container.
- Include a PostgreSQL service in the compose stack with persistent volume storage.
- Run database migrations automatically on server startup (no separate manual step).
- Wire service dependencies so PostgreSQL is healthy before the server starts.
- Keep development convenience: support live-reload / hot-reload patterns for future compose profiles.

## Approach

### File Layout

```
demo-repo/
  docker-compose.yml          # Compose definition: postgres, server, client
  server/
    Dockerfile                # Multi-stage: Go build → distroless/scratch runtime
    .dockerignore             # Exclude local binaries, IDE files
  client/
    Dockerfile                # Multi-stage: Node build → nginx runtime
    .dockerignore             # Exclude node_modules, build artifacts
    nginx.conf                # Nginx config: serve static files + proxy /api to server
```

### docker-compose.yml Design

Three services:

| Service | Image | Port | Env Vars |
|---------|-------|------|----------|
| `postgres` | `postgres:16-alpine` | `5432` (internal only) | `POSTGRES_USER=postgres`, `POSTGRES_PASSWORD=postgres`, `POSTGRES_DB=demorepo` |
| `server` | Built from `server/Dockerfile` | `8080` (internal only) | `DATABASE_URL=postgres://postgres:postgres@postgres:5432/demorepo?sslmode=disable`, `PORT=8080` |
| `client` | Built from `client/Dockerfile` | `80` → mapped to host `3000` | None (nginx proxies to `server:8080`) |

- `postgres` uses a named volume (`pgdata`) to persist data across restarts.
- `server` depends on `postgres` with a health check (`pg_isready`) so PostgreSQL accepts connections before the server starts.
- `server` runs `go run . up` (migrations) then starts the HTTP server. This is handled by an entrypoint script or by running a subcommand in the compose `command` field before the main process.
- `client` depends on `server` and proxies `/api/*` requests to `http://server:8080` via nginx.

### Server Dockerfile (server/Dockerfile)

```
# Stage 1: build
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o server .

# Stage 2: runtime
FROM alpine:3.19
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/server .
COPY db/migrations/ ./db/migrations/
CMD ["./server"]
```

- Multi-stage build to keep the runtime image small.
- Migration SQL files are copied into the runtime image so the server binary can execute them from disk.
- The `CMD` runs the server binary, which detects no subcommand argument and starts the HTTP listener. The migration step is handled by the compose-level `command` override or an entrypoint script (see below).

### Migration Strategy

The server binary already supports a CLI subcommand `up` that runs pending migrations. In the compose setup, the `server` service uses a `command` that first runs migrations, then starts the HTTP server:

```yaml
server:
  command: sh -c "./server up && ./server"
```

This ensures the database schema is created before the HTTP server begins accepting requests.

### Client Dockerfile (client/Dockerfile)

```
# Stage 1: build
FROM node:18-alpine AS builder
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
RUN npm run build

# Stage 2: serve
FROM nginx:1.25-alpine
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/build /usr/share/nginx/html
EXPOSE 80
```

- `npm ci` for reproducible installs (respects `package-lock.json`).
- `npm run build` creates optimized static assets in `client/build/`.
- The nginx layer is minimal; `nginx.conf` serves the static files and reverse-proxies `/api/` to `http://server:8080`.

### nginx.conf (client/nginx.conf)

```nginx
server {
    listen 80;
    server_name localhost;

    location / {
        root /usr/share/nginx/html;
        index index.html;
        try_files $uri $uri/ /index.html;
    }

    location /api/ {
        proxy_pass http://server:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

- The `try_files` fallback to `/index.html` supports client-side routing (React Router).
- `/api/` is forwarded to the Go server, replacing the CRA `"proxy"` field that only works in dev mode.

### Health Checks

```yaml
postgres:
  healthcheck:
    test: ["CMD-SHELL", "pg_isready -U postgres"]
    interval: 5s
    timeout: 5s
    retries: 5

server:
  healthcheck:
    test: ["CMD", "wget", "--spider", "-q", "http://localhost:8080/api/products"]
    interval: 10s
    timeout: 5s
    retries: 3

client:
  healthcheck:
    test: ["CMD", "wget", "--spider", "-q", "http://localhost:80"]
    interval: 10s
    timeout: 5s
    retries: 3
```

### Development Convenience (Future Consideration)

The initial compose file targets a production-like setup (nginx serving static files). For local development with hot reload, a future compose profile (`dev`) could mount source directories and run `npm start` (CRA dev server) and `go run .` with live-reload tooling. This is out of scope for the initial implementation but the compose file structure should not preclude it.

## Risks

| Risk | Mitigation |
|------|------------|
| Migration re-run on every container start | The existing migration runner tracks applied migrations via a `schema_migrations` table; the `up` subcommand is idempotent — it only applies pending migrations. |
| Client `package.json` proxy field becomes stale or misleading | Document in the spec / README that the proxy field is a dev-only convenience; production routing uses nginx. The field is harmless if left in place. |
| Nginx image size or CVEs | Use pinned `nginx:1.25-alpine` image and update periodically. The compose file uses explicit tags, not `latest`. |
| PostgreSQL data loss on `docker compose down` | The named volume (`pgdata`) persists across `down` commands. Explicit `docker compose down -v` is required to destroy data. |
| Port conflicts (host port 3000 for client) | Map client to host port `3000` by default but document how to override with `CLIENT_PORT` env var or a `.env` file. |
| `package-lock.json` missing or out of date | The Docker build uses `npm ci` which requires a valid lockfile. CI should verify the lockfile is present and consistent. |

## Validation

- [ ] `docker compose up` starts all three services (postgres, server, client).
- [ ] `docker compose ps` shows all services as healthy.
- [ ] Visiting `http://localhost:3000` serves the React application.
- [ ] API calls to `/api/products` return product data from the database.
- [ ] User registration and login work end-to-end (client ↔ nginx ↔ server ↔ postgres).
- [ ] `docker compose down && docker compose up` reuses the persisted database volume (data survives restart).
- [ ] Running `docker compose up` twice does not error due to duplicate migration attempts.
- [ ] No secrets or credentials are hardcoded in Dockerfiles (only defaults in compose file).
