# Architecture

## High-Level System

```
┌──────────────┐      ┌──────────────┐      ┌──────────────┐
│   Frontend   │ ───> │  REST API    │ ───> │  Go Backend  │
│   Next.js    │      │  (Gin)       │      │  (Services)  │
└──────────────┘      └──────────────┘      └──────────────┘
                                                   │
                                                   ▼
                                          ┌────────────────┐
                                          │  PostgreSQL 17 │
                                          └────────────────┘
                                                   │
                                                   ▼
                                     ┌────────────────────────┐
                                     │  LLM Provider Layer    │
                                     │  (OpenAI-compatible)   │
                                     └────────────────────────┘
                                                   │
                                                   ▼
                                       Configured AI Endpoint
```

## Module Map

| Module                       | Responsibility                                       |
| ---------------------------- | ---------------------------------------------------- |
| `internal/auth`              | Login, logout, refresh, JWT issuance & validation.   |
| `internal/provider`          | CRUD + test for OpenAI-compatible providers.         |
| `internal/analysis`          | LLM analysis pipeline, validation, retry, repair.    |
| `internal/report`            | Markdown / HTML / PDF rendering of analysis output.  |
| `internal/database`          | GORM initialization, migrations, connection pool.    |
| `internal/models`            | Domain entities & DTOs.                              |
| `internal/services`          | Application-level orchestration.                     |
| `internal/repositories`      | Persistence abstractions over GORM.                  |
| `internal/middleware`        | Auth, RBAC, rate limit, request ID, recovery.        |
| `internal/config`            | Environment & secrets loading.                       |
| `internal/api`               | Gin route registration & handlers.                   |

## Frontend Module Map

| Path                            | Purpose                                       |
| ------------------------------- | --------------------------------------------- |
| `app/`                          | Next.js App Router pages.                     |
| `components/editor/`            | Monaco wrapper and SQL paste components.      |
| `components/analysis/`          | Result views, diagrams, tables, badges.       |
| `components/provider/`          | Provider management forms & cards.            |
| `components/history/`           | History list, filters, detail viewer.         |
| `components/reports/`           | Report rendering & export.                    |
| `lib/`                          | Utilities, axios client, formatters.          |
| `hooks/`                        | React Query hooks.                            |
| `services/`                     | Typed API service modules.                    |
| `types/`                        | Shared TypeScript types.                      |

## Data Flow — Analysis

1. User pastes PL/SQL in the workspace and clicks **Analyze**.
2. Frontend POSTs the code + object type to `POST /api/v1/analyses`.
3. Backend resolves the user's default provider.
4. Backend builds a strict system prompt enforcing a fixed JSON schema.
5. Backend calls the LLM provider.
6. Backend validates the response (JSON shape + Mermaid syntax).
7. If invalid, retry with repair instructions; otherwise persist.
8. Response returned to frontend which renders the report.

## Security Boundaries

- API keys are encrypted at rest using AES-GCM with a server-side key.
- API keys are never returned to the frontend — only masked previews.
- LLM provider endpoints are validated against SSRF (no localhost / private ranges by default).
- All input is treated as **text only**; SQL is never parsed or executed.
- Rate limits applied per-IP and per-user.
- JWT tokens are short-lived; refresh tokens are rotated.

## Future-Ready Extensions

- **Redis** for caching analysis results and rate limits.
- **Qdrant** for semantic search across analyses.
- **Keycloak** as an alternative IdP.
- **Kubernetes** manifests with HPA and PDB.

## Database Tables (initial)

- `users` — credentials, role, status.
- `providers` — AI provider config (encrypted key).
- `analyses` — stored results.
- `refresh_tokens` — rotated refresh tokens.
