# Deployment

SQL-Sage ships as a multi-service Docker stack: frontend (Next.js), backend (Go), Postgres 17, and an optional Nginx reverse proxy.

## Services

| Service   | Port | Image                                       |
| --------- | ---- | ------------------------------------------- |
| Frontend  | 3000 | Custom build (Node 22 + Next.js standalone) |
| Backend   | 8080 | Custom build (Go 1.24 + Alpine)             |
| Postgres  | 5432 | `postgres:17-alpine`                        |
| Nginx     | 80/443 | `nginx:1.27-alpine` (prod profile)        |

## Local Development (no Docker)

### Backend

```bash
cd backend
cp .env.example .env
go run ./cmd/server
```

The server listens on `:8080` and runs migrations on startup.

### Frontend

```bash
cd frontend
npm install --legacy-peer-deps
cp .env.example .env.local  # set NEXT_PUBLIC_API_URL=http://localhost:8080
npm run dev
```

The app is at `http://localhost:3000`.

## Local (Docker)

```bash
cp .env.example .env
# edit JWT_SECRET, ENCRYPTION_KEY (use `openssl rand -base64 32`)
docker compose up -d
```

Services:

- Frontend: http://localhost:3000
- Backend:  http://localhost:8080
- Postgres: localhost:5432

Seed an admin:

```bash
docker compose exec backend /app/seed
```

## Production

```bash
docker compose --profile prod up -d
```

This adds Nginx as a TLS-terminating reverse proxy listening on `80` and `443`.

Place TLS certs in `infra/nginx/certs/`:

- `fullchain.pem`
- `privkey.pem`

See `infra/nginx/certs/README.md` for self-signed dev instructions.

## Environment Variables

| Key                       | Description                                        |
| ------------------------- | -------------------------------------------------- |
| `APP_ENV`                 | `development` or `production`                      |
| `HTTP_PORT`               | Backend listen port (default 8080)                 |
| `DB_*`                    | Postgres connection                                |
| `JWT_SECRET`              | At least 32 random characters                      |
| `ENCRYPTION_KEY`          | Exactly 32 bytes (use `openssl rand -base64 32` after base64-decode) |
| `ALLOWED_PROVIDER_HOSTS`  | Comma-separated SSRF allowlist for provider URLs   |
| `RATE_LIMIT_PER_MIN`      | Per-IP request cap                                 |
| `LOG_LEVEL`               | `debug`/`info`/`warn`/`error`                      |
| `NEXT_PUBLIC_API_URL`     | Public URL of the backend                          |

## Health Checks

- Backend: `GET /healthz` (liveness) and `GET /readyz` (DB readiness)
- Frontend: `GET /login` returns 200

## Backups

```bash
docker compose exec postgres pg_dump -U sage sage > backup.sql
```

Restore:

```bash
cat backup.sql | docker compose exec -T postgres psql -U sage sage
```
