# Contributing

Thank you for your interest in contributing to SQL-Sage.

## Development Rules

1. **Atomic Commits** — One logical change per commit.
2. **Branching** — Work in a feature branch and open a PR.
3. **Tests** — Add or update tests for any behavior change.
4. **Linting** — Code must pass `go vet`, `golangci-lint`, `eslint`, and `tsc --noEmit` before merge.
5. **No Secrets** — Never commit API keys, tokens, or credentials.

## Commit Convention

We follow [Conventional Commits](https://www.conventionalcommits.org/).

```
<type>(<scope>): <short description>

types: chore, feat, fix, refactor, test, docs, build, ci
```

Examples:
- `feat(backend): add provider test endpoint`
- `fix(frontend): handle 401 on token expiry`
- `docs: update deployment guide`

## Local Setup

### Backend

```bash
cd backend
go mod tidy
go run ./cmd/server
```

### Frontend

```bash
cd frontend
pnpm install
pnpm dev
```

## Pull Request Checklist

- [ ] Tests added / updated
- [ ] Lint passes
- [ ] Build passes
- [ ] Documentation updated (if applicable)
- [ ] No secrets introduced

## Code of Conduct

Be respectful. Focus on the work. Assume good intent.
