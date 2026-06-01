# Security Policy

## Our Commitments

- SQL-Sage is a **read-only** static analysis platform. It **never executes** SQL.
- User input is treated as **untrusted text**.
- All provider API keys are encrypted at rest.
- Authentication is enforced for every protected endpoint.

## Threat Model

| Threat            | Mitigation                                              |
| ----------------- | ------------------------------------------------------- |
| SQL Injection     | SQL is never parsed or executed; all input is plain text.|
| XSS               | React + strict CSP, no inline scripts, sanitized HTML.  |
| CSRF              | Bearer token auth, SameSite cookies, double-submit.     |
| SSRF              | Provider URL allowlist; block private/loopback ranges.  |
| Secret Leakage    | AES-GCM at-rest encryption, no keys in logs/responses.  |
| Brute Force       | Per-IP and per-user rate limits, account lockout.       |
| Token Theft       | Short-lived JWTs, rotated refresh tokens, revocation.  |

## Reporting a Vulnerability

Please email **arafath.yasser@gmail.com** with details. We will respond within 72 hours.

Do **not** file public GitHub issues for security vulnerabilities.

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 1.x     | :white_check_mark: |
