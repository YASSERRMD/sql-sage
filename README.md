<div align="center">

# SQL-Sage

**AI-powered static analysis platform for Oracle PL/SQL and SQL code.**

[![License](https://img.shields.io/badge/license-see%20LICENSE-blue.svg)](./LICENSE)
[![Backend](https://img.shields.io/badge/backend-Go%201.24-00ADD8.svg)](https://go.dev)
[![Frontend](https://img.shields.io/badge/frontend-Next.js%2015-000000.svg)](https://nextjs.org)
[![Database](https://img.shields.io/badge/database-PostgreSQL%2017-336791.svg)](https://www.postgresql.org)

</div>

---


## Overview

SQL-Sage helps developers, architects, DBAs, auditors, and modernization teams
understand legacy database code. It generates human-readable explanations,
execution flows, dependency maps, business rules, risk assessments, and
modernization recommendations from procedures, functions, packages, triggers,
views, and raw SQL scripts.

> **SQL-Sage never executes SQL.** It is strictly a static analysis and
> documentation platform — your code is read, explained, and visualized, but
> never run against a database.

## Key Features

- **AI Analysis Engine** — LLM-powered explanations of procedures, functions, packages, triggers, views, and SQL scripts, with structured output validated against a JSON schema.
- **Provider Agnostic** — Works with any OpenAI-compatible endpoint: OpenAI, OpenRouter, Groq, DeepSeek, vLLM, LM Studio, Ollama, LiteLLM, or an internal gateway.
- **Secure by Design** — AES-GCM encrypted API keys, JWT authentication with refresh tokens, role-based access control, per-IP rate limiting, and an SSRF allowlist for provider URLs.
- **Analysis Workspace** — Monaco-based SQL editor with syntax highlighting and a paste-to-analyze flow.
- **Visualizations** — Auto-generated Mermaid flowcharts, dependency tables, and risk distribution charts.
- **History & Reports** — Searchable analysis history with Markdown and HTML export.
- **Dashboard** — Summary statistics, recent analyses, object-type breakdowns, risk distribution, and trends over time.
- **Modern UI** — Responsive enterprise interface with dark and light modes.

## Architecture

```
┌─────────────────┐      HTTPS       ┌──────────────────┐      ┌──────────────┐
│   Next.js 15    │ ───────────────▶ │   Go API (Gin)   │ ───▶ │ PostgreSQL 17│
│  (App Router)   │ ◀─────────────── │  REST + JWT auth │ ◀─── │   via GORM   │
└─────────────────┘                  └────────┬─────────┘      └──────────────┘
                                              │
                                              │ OpenAI-compatible API
                                              ▼
                                     ┌──────────────────┐
                                     │   LLM Provider   │
                                     │ (SSRF-allowlisted)│
                                     └──────────────────┘
```

See [ARCHITECTURE.md](./ARCHITECTURE.md) for module-level design.

## Tech Stack

| Layer       | Technology                                       |
| ----------- | ------------------------------------------------ |
| Frontend    | Next.js 15, TypeScript, Tailwind CSS, shadcn/ui  |
| Editor      | Monaco Editor                                    |
| Backend     | Go 1.24+, Gin, GORM                              |
| LLM Client  | OpenAI-compatible chat completions client        |
| Database    | PostgreSQL 17                                    |
| Infra       | Docker, Docker Compose, Nginx reverse proxy      |

## Project Structure

```
sql-sage/
├── backend/                # Go API server
│   ├── cmd/server/         # Application entrypoint
│   ├── cmd/seed/           # Admin/seed utility
│   ├── internal/           # Domain logic (analysis, auth, provider, report, ...)
│   ├── pkg/                # Reusable packages (crypto, llm, httpx, logger)
│   └── migrations/         # Database migrations
├── frontend/               # Next.js application
│   └── app/                # App Router routes ((auth), (app)/workspace, history, providers)
├── infra/nginx/            # Reverse proxy configuration
├── docker-compose.yml
├── ARCHITECTURE.md
├── API.md
├── CONTRIBUTING.md
├── SECURITY.md
└── DEPLOYMENT.md
```

## Quick Start

**Prerequisites:** Docker and Docker Compose.

```bash
# 1. Configure environment
cp .env.example .env
# Generate strong secrets:
#   JWT_SECRET     -> openssl rand -base64 48
#   ENCRYPTION_KEY -> openssl rand -base64 32   (must decode to 32 bytes)

# 2. Launch the stack
docker compose up -d
```

| Service  | URL                     |
| -------- | ----------------------- |
| Frontend | http://localhost:3000   |
| Backend  | http://localhost:8080   |

Register the first account, add an OpenAI-compatible provider with your API key,
then paste SQL into the workspace to run your first analysis. See
[DEPLOYMENT.md](./DEPLOYMENT.md) for production deployment and seeding.

## Configuration

Key environment variables (full list in [`.env.example`](./.env.example)):

| Variable                 | Description                                                        |
| ------------------------ | ------------------------------------------------------------------ |
| `JWT_SECRET`             | Secret used to sign JWTs (use a long random value).                |
| `ENCRYPTION_KEY`         | 32-byte key for AES-GCM encryption of provider API keys.           |
| `ALLOWED_PROVIDER_HOSTS` | Comma-separated SSRF allowlist of permitted LLM provider hosts.    |
| `RATE_LIMIT_PER_MIN`     | Per-IP request rate limit.                                         |
| `NEXT_PUBLIC_API_URL`    | Backend base URL exposed to the frontend.                          |

## API

The backend exposes a REST API under `/api/v1`. Selected endpoints:

| Method | Endpoint                          | Description                          |
| ------ | --------------------------------- | ------------------------------------ |
| POST   | `/auth/register`                  | Create an account                    |
| POST   | `/auth/login`                     | Authenticate and receive tokens      |
| POST   | `/auth/refresh`                   | Refresh access token                 |
| GET    | `/providers`                      | List configured LLM providers        |
| POST   | `/providers`                      | Add a provider                       |
| POST   | `/providers/{id}/test`            | Test provider connectivity           |
| POST   | `/analyses`                       | Run a new analysis                   |
| GET    | `/analyses`                       | List analysis history                |
| GET    | `/analyses/{id}/report`           | Export report (Markdown / HTML)      |
| GET    | `/dashboard/summary`              | Dashboard statistics                 |
| GET    | `/healthz`, `/readyz`             | Liveness and readiness probes        |

See [API.md](./API.md) for the complete reference.

## Documentation

| Document                               | Contents                              |
| -------------------------------------- | ------------------------------------- |
| [ARCHITECTURE.md](./ARCHITECTURE.md)   | System architecture & module design   |
| [API.md](./API.md)                     | Full REST API reference               |
| [CONTRIBUTING.md](./CONTRIBUTING.md)   | Development workflow & conventions     |
| [SECURITY.md](./SECURITY.md)           | Security policy & reporting            |
| [DEPLOYMENT.md](./DEPLOYMENT.md)       | Deployment & operations guide          |

## Security

SQL-Sage encrypts provider credentials at rest, authenticates every request,
rate-limits by IP, and restricts outbound LLM calls to an allowlist to prevent
SSRF. To report a vulnerability, please follow the process in
[SECURITY.md](./SECURITY.md).

## Contributing

Contributions are welcome. Please read [CONTRIBUTING.md](./CONTRIBUTING.md) for
the development workflow, coding standards, and pull request process.

## License

See [LICENSE](./LICENSE).
