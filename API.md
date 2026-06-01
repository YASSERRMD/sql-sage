# API Reference

Base URL: `/api/v1`

All endpoints (except `/auth/login` and `/auth/refresh`) require a Bearer token.

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

### POST /auth/login

Request:

```json
{
  "email": "user@example.com",
  "password": "string"
}
```

Response `200`:

```json
{
  "accessToken": "string",
  "refreshToken": "string",
  "expiresIn": 900
}
```

### POST /auth/refresh

Request:

```json
{ "refreshToken": "string" }
```

Response `200`: same as login.

### POST /auth/logout

Invalidates the supplied refresh token.

### GET /auth/me

Returns the authenticated user profile.

## Providers

### GET /providers
### POST /providers
### GET /providers/{id}
### PUT /providers/{id}
### DELETE /providers/{id}
### POST /providers/{id}/test
### POST /providers/{id}/default

## Analyses

### POST /analyses
### GET /analyses
### GET /analyses/{id}
### DELETE /analyses/{id}

## Reports

### GET /analyses/{id}/report?format=md|html|pdf

## Dashboard

### GET /dashboard/summary
### GET /dashboard/trend
### GET /dashboard/risk-distribution

## Health

### GET /healthz
### GET /readyz
