# SQL-Sage

> AI-powered static analysis platform for Oracle PL/SQL and SQL code.

SQL-Sage helps developers, architects, DBAs, auditors, and modernization teams understand legacy database code by generating human-readable explanations, execution flows, dependency maps, business rules, risk assessments, and modernization recommendations.

**SQL-Sage never executes SQL.** It is strictly a static analysis and documentation platform.

---

## Features

- **AI Analysis Engine** — LLM-powered explanations of procedures, functions, packages, triggers, views, and SQL scripts.
- **Provider Agnostic** — Works with any OpenAI-compatible endpoint (OpenAI, OpenRouter, Groq, DeepSeek, vLLM, LM Studio, Ollama, LiteLLM, internal gateways).
- **Secure by Design** — Encrypted API keys, JWT authentication, role-based access.
- **Workspace** — Monaco editor with syntax highlighting and paste-to-analyze flow.
- **Visualizations** — Auto-generated Mermaid flowcharts, dependency tables, risk distribution charts.
- **History & Reports** — Searchable analysis history with Markdown / HTML / PDF export.
- **Dashboard** — Analytics, risk trends, provider usage, and high-risk findings.
- **Dark / Light Mode** — Professional enterprise UI.

## Tech Stack

| Layer       | Technology                                       |
| ----------- | ------------------------------------------------ |
| Frontend    | Next.js 15, TypeScript, Tailwind, shadcn/ui      |
| Editor      | Monaco Editor, Mermaid.js                        |
| Backend     | Go 1.24+, Gin, GORM                              |
| LLM Client  | OpenAI SDK compatible client                     |
| Database    | PostgreSQL 17                                    |
| Infra       | Docker, Docker Compose, Nginx                    |

## Project Structure

```
sql-sage/
├── backend/            # Go API server
├── frontend/           # Next.js application
├── docker-compose.yml
├── ARCHITECTURE.md
├── API.md
├── CONTRIBUTING.md
├── SECURITY.md
└── DEPLOYMENT.md
```

## Quick Start

See [DEPLOYMENT.md](./DEPLOYMENT.md) for full installation and deployment instructions.

## Documentation

- [ARCHITECTURE.md](./ARCHITECTURE.md) — System architecture & module design.
- [API.md](./API.md) — REST API reference.
- [CONTRIBUTING.md](./CONTRIBUTING.md) — Development workflow.
- [SECURITY.md](./SECURITY.md) — Security policy.
- [DEPLOYMENT.md](./DEPLOYMENT.md) — Deployment guide.

## License

See [LICENSE](./LICENSE).
