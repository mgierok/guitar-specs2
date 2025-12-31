# API Reference

This document is the canonical API reference for Guitar-Specs. Keep it aligned with the current implementation.

## Overview
- Base URL (local): `http://localhost:8080/api/v1`
- Base URL (prod): `https://www.guitar-specs.com/api/v1`
- Content type: `application/json`
- Auth: none (public read-only)

## Error Format
Errors return JSON with a stable shape:
```json
{
  "code": "list_failed",
  "message": "Failed to list guitars",
  "details": null
}
```

## Endpoints

### GET /health
Health check for the API.

- Params: none
- Response: `{"status":"ok"}`

Example:
```bash
curl -sS http://localhost:8080/api/v1/health
```

### GET /guitars
List guitars. Pagination params are parsed and echoed in the response. Filtering and sorting are accepted but not applied yet.

Query params:
- `page` (int, default `1`)
- `pageSize` (int, default `20`)
- `sort` (string, default `name:asc`) — accepted, not applied yet
- `filter` (repeatable, `key:value`) — accepted, not applied yet

Response:
```json
{
  "items": [
    {
      "id": "uuid",
      "slug": "prs-custom-24",
      "name": "Custom 24",
      "brand": "PRS",
      "model": "Custom 24",
      "type": "electric",
      "year": 2022,
      "thumbnail": "https://..."
    }
  ],
  "total": 10,
  "page": 1,
  "pageSize": 20
}
```

Example:
```bash
curl -sS "http://localhost:8080/api/v1/guitars?page=1&pageSize=20"
```

### GET /guitars/{slug}
Fetch a single guitar by slug.

Path params:
- `slug` (string, required)

Response:
```json
{
  "id": "uuid",
  "slug": "prs-custom-24",
  "name": "Custom 24",
  "brand": "PRS",
  "model": "Custom 24",
  "type": "electric",
  "year": 2022,
  "description": "..."
  ,
  "specs": {
    "body": "Mahogany",
    "bridge": "PRS Tremolo"
  },
  "media": [
    {
      "kind": "image",
      "url": "https://..."
    }
  ]
}
```

Example:
```bash
curl -sS "http://localhost:8080/api/v1/guitars/prs-custom-24"
```

## Internal Packages (Backend)
Each package should have a single, clear responsibility. This section documents the intended role of each `backend/internal/*` package.

- `internal/api`: HTTP request/response helpers and common API utilities (params parsing, JSON response writing, error shapes).
- `internal/app`: application bootstrap (config loading, DB connect, migrations/seed, router/server wiring).
- `internal/config`: configuration loading and defaults (environment variables, config struct).
- `internal/db`: database connectivity helpers and SQL definitions (`queries.sql`, schema), plus connection utilities.
- `internal/db/sqlc`: generated SQLC code (queries, models). No hand-edited logic.
- `internal/handlers`: HTTP handlers that translate requests into domain/data calls and API responses.
- `internal/log`: logging setup and helpers (structured JSON logging).
- `internal/middleware`: HTTP middleware (request ID, logging, etc.).
- `internal/migrate`: migration execution helper (apply SQL file).
- `internal/seed`: seed dataset loader and DB apply logic.
