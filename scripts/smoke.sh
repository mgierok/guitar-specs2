#!/usr/bin/env bash
set -euo pipefail

cleanup() {
  docker compose down -v
}
trap cleanup EXIT

docker compose down -v || true
docker compose up -d --build

until docker compose exec -T db pg_isready -U postgres -d guitar_specs >/dev/null 2>&1; do
  sleep 1
done

docker compose exec -T db psql -U postgres -d guitar_specs -f /migrations/0001_init.up.sql

docker compose exec -T db psql -U postgres -d guitar_specs -c "TRUNCATE guitar_media, guitar_spec_value, spec_option, spec, guitar, brand RESTART IDENTITY CASCADE;"

docker compose exec -T backend /app/seed

curl -fsS http://localhost:8080/api/v1/health >/dev/null
curl -fsS http://localhost:8080/api/v1/guitars >/dev/null
curl -fsS http://localhost:8080/api/v1/guitars/prs-custom-24 >/dev/null
curl -fsS http://localhost:3000/guitars/prs-custom-24 >/dev/null
