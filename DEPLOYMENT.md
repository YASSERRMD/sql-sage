# Deployment

SQL-Sage ships with Docker images for the backend, frontend, and a Postgres 17 database.

## Prerequisites

- Docker 24+
- Docker Compose v2

## Local (development)

```bash
cp .env.example .env
docker compose up -d
```

Services:

- Frontend: http://localhost:3000
- Backend:  http://localhost:8080
- Postgres: localhost:5432

## Environment Variables

See `.env.example` for the full list. Critical keys:

- `JWT_SECRET` — at least 32 random bytes.
- `DB_PASSWORD` — Postgres password.
- `ENCRYPTION_KEY` — 32-byte key for AES-GCM provider secrets.
- `ALLOWED_PROVIDER_HOSTS` — comma-separated host allowlist for SSRF defense.

## Production

1. Provision a managed PostgreSQL 17 instance.
2. Set strong secrets and rotate them regularly.
3. Build images: `docker compose -f docker-compose.prod.yml build`.
4. Run behind a reverse proxy (Nginx / Traefik) with TLS termination.
5. Configure backups for Postgres.

## Health Checks

- `GET /healthz` — liveness.
- `GET /readyz` — readiness (verifies DB).

## Migrations

Migrations run automatically on backend start. To run manually:

```bash
docker compose exec backend ./migrate up
```

## Backups

```bash
docker compose exec postgres pg_dump -U sage sage > backup.sql
```
