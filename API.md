# API Reference

Base URL: `/api/v1`

All endpoints (except `/auth/login`, `/auth/refresh`, `/auth/register`) require a Bearer token.

```
Authorization: Bearer <access_token>
```

## Conventions

- JSON in, JSON out.
- Timestamps are RFC3339 UTC.
- Errors return:

```json
{
  "error": {
    "code": "string",
    "message": "string",
    "details": {}
  }
}
```

## Authentication

### POST /auth/register

```json
{ "name": "Jane", "email": "j@x.com", "password": "string" }
```

### POST /auth/login

```json
{ "email": "user@example.com", "password": "string" }
```

Response `200`:

```json
{ "accessToken": "string", "refreshToken": "string", "expiresIn": 900 }
```

### POST /auth/refresh

```json
{ "refreshToken": "string" }
```

### POST /auth/logout

```json
{ "refreshToken": "string" }
```

### GET /auth/me

Returns the authenticated user.

### PATCH /users/me

```json
{ "name": "string" }
```

## Providers

### GET /providers
List providers for the current user.

### POST /providers
```json
{
  "name": "OpenAI",
  "baseUrl": "https://api.openai.com/v1",
  "apiKey": "sk-...",
  "modelName": "gpt-4o-mini",
  "temperature": 0.2,
  "maxTokens": 2048,
  "isDefault": true
}
```

### GET /providers/{id}
### PUT /providers/{id}
### DELETE /providers/{id}
### POST /providers/{id}/test

Response:
```json
{ "ok": true, "message": "ok", "latencyMs": 312 }
```

### POST /providers/{id}/default

## Analyses

### POST /analyses

```json
{
  "objectName": "close_account",
  "objectType": "procedure",
  "sourceCode": "...",
  "providerId": "optional"
}
```

Returns the analysis with `result` containing the structured JSON.

### GET /analyses?q=&objectType=&risk=&page=&pageSize=

### GET /analyses/{id}
### DELETE /analyses/{id}

### GET /analyses/{id}/report?format=md|html|pdf

Streams Markdown or HTML. `pdf` currently returns 501.

## Dashboard

### GET /dashboard/summary
### GET /dashboard/trend
### GET /dashboard/risk-distribution
### GET /dashboard/object-types

## Health

### GET /healthz — liveness
### GET /readyz — DB readiness
